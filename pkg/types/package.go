package types

import (
	"go/ast"
	"strings"

	"k8s.io/klog/v2"

	"github.com/galaxyobe/gen/pkg/util"
	"github.com/galaxyobe/gen/third_party/gengo/parser"
)

type PackageTypes []*PackageType

type PackageType struct {
	Name        string              // pkg name
	Types       map[string]string   // key: type's name value: type
	Int8Fields  map[string][]string // type's fields
	Uint8Fields map[string][]string // type's fields
}

func IsExternalType(type_ string) bool {
	return strings.Contains(type_, ".")
}

func (p PackageType) IsExternalType(name string) bool {
	for name_, type_ := range p.Types {
		if name_ == name {
			return IsExternalType(type_)
		}
	}
	return false
}

func (p PackageType) GetType(name string) string {
	for name_, type_ := range p.Types {
		if name_ == name {
			return type_
		}
	}
	return ""
}

func (p PackageType) GetUint8Fields(name string) []string {
	for name_, fields := range p.Uint8Fields {
		if name_ == name {
			return fields
		}
	}
	return nil
}

func (p PackageType) GetInt8Fields(name string) []string {
	for name_, fields := range p.Int8Fields {
		if name_ == name {
			return fields
		}
	}
	return nil
}

func NewPackageTypes(build *parser.Builder) (ret PackageTypes) {
	for pkg, list := range build.GetParsedFiles() {
		var files []*ast.File
		for _, item := range list {
			files = append(files, item.File)
		}
		ret = append(ret, &PackageType{
			Name:        pkg,
			Types:       util.FindTypeNames(files),
			Int8Fields:  util.FindInt8Type(files),
			Uint8Fields: util.FindUint8Type(files),
		})
	}
	return
}

func (list PackageTypes) Find(pkg string) *PackageType {
	for _, item := range list {
		switch {
		case item.Name == pkg:
			return item
		case strings.HasSuffix(item.Name, pkg):
			return item
		}
	}
	return nil
}

func (list PackageTypes) IsExternalType(pkg, name string) bool {
	p := list.Find(pkg)
	if p == nil {
		return false
	}
	return p.IsExternalType(name)
}

func (list PackageTypes) GetUint8Fields(pkg, name string) []string {
	p := list.Find(pkg)
	if p == nil {
		return nil
	}
	type_ := p.GetType(name)
	if !IsExternalType(type_) {
		return p.GetUint8Fields(name)
	}
	items := strings.Split(type_, ".")
	if len(items) != 2 {
		klog.Fatalf("unexpect external type: %s", type_)
	}
	return list.GetUint8Fields(items[0], items[1])
}

func (list PackageTypes) GetInt8Fields(pkg, name string) []string {
	p := list.Find(pkg)
	if p == nil {
		return nil
	}
	type_ := p.GetType(name)
	if !IsExternalType(type_) {
		return p.GetInt8Fields(name)
	}
	items := strings.Split(type_, ".")
	if len(items) != 2 {
		klog.Fatalf("unexpect external type: %s", type_)
	}
	return list.GetInt8Fields(items[0], items[1])
}
