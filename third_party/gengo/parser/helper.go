package parser

import (
	"go/ast"
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
