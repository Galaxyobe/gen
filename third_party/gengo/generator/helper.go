package generator

import (
	"text/template"
)

func (s *SnippetWriter) AddFunc(name string, f interface{}) *SnippetWriter {
	s.funcMap[name] = f
	return s
}

func (s *SnippetWriter) AddFuncMap(m template.FuncMap) *SnippetWriter {
	for k, v := range m {
		s.funcMap[k] = v
	}
	return s
}
