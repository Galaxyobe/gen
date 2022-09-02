package parser

import (
	"go/ast"

	"k8s.io/gengo/types"

	tptypes "github.com/galaxyobe/gen/third_party/gengo/types"
)

func (b *Builder) GetAllParsedFiles() (list []*ast.File) {
	for _, files := range b.parsed {
		for _, file := range files {
			list = append(list, file.file)
		}
	}
	return
}

// ParsedFile is for tracking files with name
type ParsedFile struct {
	Path string
	File *ast.File
}

func (b *Builder) GetParsedFiles() map[string][]ParsedFile {
	var m = make(map[string][]ParsedFile, len(b.parsed))
	for pkg, files := range b.parsed {
		var list = make([]ParsedFile, 0, len(files))
		for _, file := range files {
			list = append(list, ParsedFile{
				Path: file.name,
				File: file.file,
			})
		}
		m[string(pkg)] = list
	}
	return m
}

func (b *Builder) ReplaceUniverse(u types.Universe) error {
	universe, err := b.FindTypes()
	if err != nil {
		return err
	}
	tptypes.ReplaceUniverse(u, universe)
	return nil
}
