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

package generators

import (
	"io"
	"strings"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"

	"github.com/galaxyobe/gen/pkg/util"
	tpgenerator "github.com/galaxyobe/gen/third_party/gengo/generator"
	"github.com/galaxyobe/gen/third_party/gengo/parser"
)

type GenType struct {
	*types.Type
	AllowFields []string
}

func NewGenTypes(pkg *types.Package) (pkgEnabled bool, list GenTypes) {
	pkgAllowed := util.CheckTag(tagPackageName, pkg.Comments, util.Package)
	for _, t := range pkg.Types {
		ut := util.UnderlyingType(t)
		if ut.Kind != types.Struct {
			continue
		}
		comments := t.CommentLines
		set, enabled := util.GetTagBoolStatus(tagPackageName, comments)
		allowedFields := util.GetTagValues(tagSelectFieldsName, comments)
		if len(allowedFields) > 0 {
			set = true
			enabled = true
		}
		if (!pkgAllowed && !set) || !enabled {
			continue // ignore type
		}
		pkgEnabled = true
		list = append(list, &GenType{
			Type:        t,
			AllowFields: allowedFields,
		})
	}
	if len(list) == 0 {
		pkgEnabled = false
	}
	return
}

type GenTypes []*GenType

func (list GenTypes) allowed(t *types.Type) bool {
	for _, item := range list {
		if item.Name.Name == t.Name.Name && item.Name.Package == t.Name.Package {
			return true
		}
	}
	return false
}

// allowedField allowed type's field to generate Setter func.
// will be ignored field enabled status when it's in +gen:setter:fields allowed fields.
func (list GenTypes) allowedField(t *types.Type, m int) bool {
	for _, item := range list {
		if item.Name.Name == t.Name.Name && item.Name.Package == t.Name.Package {
			if len(t.Members) == 0 {
				return true
			}
			field := t.Members[m]
			set, enable := util.GetTagBoolStatus(tagFieldName, field.CommentLines)
			if set && !enable && len(item.AllowFields) == 0 {
				return false
			}
			if len(item.AllowFields) == 0 {
				return true
			}
			return util.Exist(item.AllowFields, field.Name)
		}
	}
	return false
}

type genSetter struct {
	generator.DefaultGen
	build         *parser.Builder
	targetPackage string
	boundingDirs  []string
	imports       namer.ImportTracker
	types         GenTypes
	packageTypes  util.PackageTypes
}

func NewGenSetter(build *parser.Builder, sanitizedName, targetPackage string, boundingDirs []string, types []*GenType, sourcePath string) generator.Generator {
	return &genSetter{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		build:         build,
		targetPackage: targetPackage,
		boundingDirs:  boundingDirs,
		imports:       generator.NewImportTracker(),
		types:         types,
		packageTypes:  util.NewPackageTypes(build),
	}
}

func (g *genSetter) Name() string {
	return "setter"
}

func (g *genSetter) Filter(c *generator.Context, t *types.Type) bool {
	ok := g.types.allowed(t)
	if !ok {
		klog.V(5).Infof("Ignore generate setter function for type %v", t)
	}
	return ok
}

func (g *genSetter) Namers(c *generator.Context) namer.NameSystems {
	// Have the raw namer for this file track what it imports.
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.targetPackage, g.imports),
	}
}

func (g *genSetter) Init(c *generator.Context, w io.Writer) error {
	return nil
}

func (g *genSetter) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	klog.V(5).Infof("Generating setter function for type %v", t)

	sw := tpgenerator.NewSnippetWriter(w, c, "", "")
	sw.AddFunc("slice", func(s string) string {
		if strings.HasPrefix(s, "[]") {
			return strings.ReplaceAll(s, "[]", "...")
		}
		return s
	})
	g.genSetFunc(sw, t)
	sw.Do("\n", nil)

	return sw.Error()
}

func (g *genSetter) isOtherPackage(pkg string) bool {
	if pkg == g.targetPackage {
		return false
	}
	if strings.HasSuffix(pkg, "\""+g.targetPackage+"\"") {
		return false
	}
	return true
}

func (g *genSetter) Imports(c *generator.Context) (imports []string) {
	var importLines []string
	for _, singleImport := range g.imports.ImportLines() {
		if g.isOtherPackage(singleImport) {
			importLines = append(importLines, singleImport)
		}
	}
	return importLines
}

func (g *genSetter) genSetFunc(sw *tpgenerator.SnippetWriter, t *types.Type) {
	receiver := strings.ToLower(t.Name.Name[:1])
	isExternalType := g.packageTypes.IsExternalType(t.Name.Package, t.Name.Name)
	var methodSet = util.NewMethodSet()
	var methodGen = util.NewMethodGenerate(util.GenName("Set", ""))

	for idx, m := range t.Members {
		if util.IsLower(m.Name) && isExternalType {
			continue
		}
		methods := util.GetTagValues(tagMethodName, m.CommentLines)
		methodSet.AddMethods("Set", methods, t.Members[idx])
		if !g.types.allowedField(t, idx) {
			continue
		}
		method := methodGen.GenName(m.Name)
		if _, ok := t.Methods[method]; ok {
			continue
		}
		args := generator.Args{
			"type":     t,
			"field":    m,
			"receiver": receiver,
			"method":   method,
		}
		sw.Do("func ({{.receiver}} *{{.type|public}}) {{.method}}(val {{.field.Type|raw|slice}}) *{{.type|public}} {\n", args)
		sw.Do("{{.receiver}}.{{.field.Name}} = val\n", args)
		sw.Do("return {{.receiver}}", args)
		sw.Do("}\n\n", nil)
	}
	// add exist methods
	methodGen.AddExistNames(func() []string {
		var existMethods []string
		for name := range t.Methods {
			existMethods = append(existMethods, name)
		}
		return existMethods
	}()...)
	// gen aggregate method
	for method, members := range methodSet {
		if ok := methodGen.ExistName(method); ok {
			klog.Fatalf("exist method: %s when generate aggregate method", method)
		}
		args := generator.Args{
			"type":     t,
			"receiver": receiver,
			"method":   method,
			"fields":   members,
		}
		sw.Do("func ({{.receiver}} *{{.type|public}}) {{.method}}({{ range $i, $field := .fields }} in{{$i}} {{$field.Type|raw}}, {{end}}) *{{.type|public}} {\n", args)
		sw.Do("{{ range $i, $field := .fields }} {{$.receiver}}.{{$field.Name}} = in{{$i}} \n{{ end }}", args)
		sw.Do("return {{.receiver}}", args)
		sw.Do("}\n\n", nil)
	}
}
