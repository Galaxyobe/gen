package util

import (
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

func getSpecKindFieldsInGenDecl(genDecl *ast.GenDecl, kind string) map[string][]string {
	var m = make(map[string][]string)
	for _, spec := range genDecl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			name, list := getSpecKindFieldsInTypeSpec(s, kind)
			if len(list) == 0 {
				continue
			}
			exists := m[name]
			if exists == nil {
				m[name] = list
			} else {
				exists = append(exists, list...)
				m[name] = exists
			}
		}
	}
	return m
}

func getSpecKindFieldsInTypeSpec(typeSpec *ast.TypeSpec, kind string) (name string, fields []string) {
	name = typeSpec.Name.Name
	switch t := typeSpec.Type.(type) {
	case *ast.StructType:
		fields = append(fields, getSpecKindFieldsInStructType(t, kind)...)
	case *ast.Ident:
		fields = append(fields, getSpecKindFieldsInIdent(t, kind)...)
	}
	return
}

func getSpecKindFieldsInStructType(structType *ast.StructType, kind string) (fields []string) {
	for _, field := range structType.Fields.List {
		switch f := field.Type.(type) {
		case *ast.Ident:
			if f.Name == kind {
				if len(field.Names) == 0 {
					continue
				}
				fields = append(fields, field.Names[0].Name)
			}
		}
	}
	return
}

func getSpecKindFieldsInIdent(ident *ast.Ident, kind string) (fields []string) {
	if ident.Obj == nil {
		return
	}
	switch decl := ident.Obj.Decl.(type) {
	case *ast.TypeSpec:
		_, fields = getSpecKindFieldsInTypeSpec(decl, kind)
	}
	return
}

func FindStructSpecKindFields(files []*ast.File, kind string) map[string][]string {
	var m = make(map[string][]string)
	for _, file := range files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.GenDecl:
				result := getSpecKindFieldsInGenDecl(n, kind)
				for k, v := range result {
					m[k] = v
				}
			}
			return true
		})
	}
	return m
}

func FindUint8Type(files []*ast.File) map[string][]string {
	return FindStructSpecKindFields(files, "uint8")
}

func FindInt8Type(files []*ast.File) map[string][]string {
	return FindStructSpecKindFields(files, "int8")
}
