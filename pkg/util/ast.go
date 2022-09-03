package util

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func walkAstParseFile(files *[]*ast.File) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "_generated.go") {
			return nil
		}
		if d.IsDir() {
			walkAstParseFile(files)
			return nil
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.DeclarationErrors|parser.ParseComments)
		if err != nil {
			return err
		}
		*files = append(*files, f)
		return nil
	}
}

func AstParseDir(dir string) (files []*ast.File, err error) {
	err = filepath.WalkDir(dir, walkAstParseFile(&files))
	return
}

func getAstObjectName(o any) string {
	switch v := o.(type) {
	case *ast.Ident:
		return v.Name
	default:
		return ""
	}
}

func getFuncType(funcType *ast.FuncType) string {
	var buf bytes.Buffer
	buf.WriteString("func(")
	for idx, param := range funcType.Params.List {
		buf.WriteString(getAstObjectName(param.Type))
		if idx < len(funcType.Params.List)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(')')
	if len(funcType.Results.List) > 1 {
		buf.WriteByte('(')
	}
	for idx, param := range funcType.Results.List {
		buf.WriteString(getAstObjectName(param.Type))
		if idx < len(funcType.Params.List)-1 {
			buf.WriteByte(',')
		}
	}
	if len(funcType.Results.List) > 1 {
		buf.WriteByte(')')
	}
	return buf.String()
}

func getTypeInTypeSpec(typeSpec *ast.TypeSpec) (type_ string, name string) {
	name = typeSpec.Name.Name
	switch t := typeSpec.Type.(type) {
	case *ast.SelectorExpr:
		type_ = getAstObjectName(t.X) + "." + getAstObjectName(t.Sel)
	case *ast.Ident:
		type_ = t.Name
	case *ast.FuncType:
		// type_ = getFuncType(t)
		type_ = "func"
	case *ast.StructType:
		type_ = "struct"
	case *ast.StarExpr:
		type_ = getAstObjectName(t.X)
	default:
		return
	}
	return
}

func getTypesInGenDecl(genDecl *ast.GenDecl) map[string]string {
	var m = make(map[string]string)
	for _, spec := range genDecl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			type_, name := getTypeInTypeSpec(s)
			m[name] = type_
		}
	}
	return m
}

func FindTypeNames(files []*ast.File) (m map[string]string) {
	m = make(map[string]string)
	for _, file := range files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.GenDecl:
				result := getTypesInGenDecl(n)
				for k, v := range result {
					m[k] = v
				}
			}
			return true
		})
	}
	return
}
