/*
 Copyright 2022 Galaxyobe.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package util

import (
	"go/ast"
	"strings"

	"github.com/galaxyobe/gen/third_party/gengo/parser"
)

type PackageTypes []*PackageType

type PackageType struct {
	Name  string            // pkg name
	Types map[string]string // key: type's name value: type
}

func (p PackageType) IsExternalType(name string) bool {
	for name_, type_ := range p.Types {
		if name_ == name {
			return strings.Contains(type_, ".")
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

func NewPackageTypes(build *parser.Builder) (ret PackageTypes) {
	for pkg, list := range build.GetParsedFiles() {
		var files []*ast.File
		for _, item := range list {
			files = append(files, item.File)
		}
		ret = append(ret, &PackageType{
			Name:  pkg,
			Types: FindTypeNames(files),
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
