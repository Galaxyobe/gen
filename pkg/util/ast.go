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

func FindInt8OrUint8Type(files []*ast.File) (list []*ast.StructType) {
	for _, file := range files {
		for _, decl := range file.Decls {
			_ = decl
		}
	}
	return nil
}
