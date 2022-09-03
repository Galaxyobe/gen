package generators

import (
	"io"
	"strings"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"

	"github.com/galaxyobe/gen/pkg/util"

	"github.com/galaxyobe/gen/third_party/gengo/parser"
)

type GenType struct {
	*types.Type
	AllowFields  []string
	OptionName   string
	OptionSuffix string
}

func NewGenTypes(pkg *types.Package) (pkgEnabled bool, genTypes GenTypes) {
	for _, t := range pkg.Types {
		ut := util.UnderlyingType(t)
		switch ut.Kind {
		case types.Struct:
		case types.Func:
			genTypes.Functions = append(genTypes.Functions, ut.Name.Name)
		default:
			continue
		}
		comments := t.CommentLines
		setTag, enabled := util.GetTagBoolStatus(tagName, comments)
		setName, name := util.GetTagValueStatus(tagTypeName, comments)
		setSuffix, suffix := util.GetTagValueStatus(tagSuffixName, comments)
		allowedFields := util.GetTagValues(tagSelectFieldsName, comments)
		if len(allowedFields) > 0 || setName || setSuffix {
			setTag = true
			enabled = true
		}
		if !setTag || !enabled {
			continue // ignore type
		}
		optionName := getOptionName(t.Name.Name)
		if !setName {
			name = optionName
		}
		if !setSuffix {
			suffix = name
		}
		pkgEnabled = true
		genTypes.Types = append(genTypes.Types, &GenType{
			Type:         t,
			AllowFields:  allowedFields,
			OptionName:   name,
			OptionSuffix: suffix,
		})
	}
	if len(genTypes.Types) == 0 {
		pkgEnabled = false
		return
	}
	for name := range pkg.Functions {
		genTypes.Functions = append(genTypes.Functions, name)
	}
	return
}

type GenTypes struct {
	Types     []*GenType
	Functions []string
}

func (g GenTypes) find(t *types.Type) *GenType {
	for _, item := range g.Types {
		if item.Name.Name == t.Name.Name && item.Name.Package == t.Name.Package {
			return item
		}
	}
	return nil
}

func (g GenTypes) allowed(t *types.Type) bool {
	return g.find(t) != nil
}

// allowedField allowed type's field to generate Option func.
// will be ignored field enabled status when it's in +gen:option:fields allowed fields.
func (g GenTypes) allowedField(t *types.Type, m int) bool {
	for _, item := range g.Types {
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

func (g GenTypes) existMethod(funcName string) bool {
	for _, item := range g.Functions {
		if item == funcName {
			return true
		}
	}
	return false
}

type genOption struct {
	generator.DefaultGen
	build         *parser.Builder
	targetPackage string
	boundingDirs  []string
	imports       namer.ImportTracker
	genTypes      GenTypes
	packageTypes  util.PackageTypes
}

func NewGenOption(build *parser.Builder, sanitizedName, targetPackage string, boundingDirs []string, genTypes GenTypes, sourcePath string) generator.Generator {
	return &genOption{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		build:         build,
		targetPackage: targetPackage,
		boundingDirs:  boundingDirs,
		imports:       generator.NewImportTracker(),
		genTypes:      genTypes,
		packageTypes:  util.NewPackageTypes(build),
	}
}

func (g *genOption) Name() string {
	return "option"
}

func (g *genOption) Filter(c *generator.Context, t *types.Type) bool {
	ok := g.genTypes.allowed(t)
	if !ok {
		klog.V(5).Infof("Ignore generate option function for type %v", t)
	}
	return ok
}

func (g *genOption) Namers(c *generator.Context) namer.NameSystems {
	// Have the raw namer for this file track what it imports.
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.targetPackage, g.imports),
	}
}

func (g *genOption) Init(c *generator.Context, w io.Writer) error {
	return nil
}

func (g *genOption) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	klog.V(5).Infof("Generating option function for type %v", t)
	sw := generator.NewSnippetWriter(w, c, "", "")
	g.genOptionType(sw, t)
	g.genWithFieldFunc(sw, t)
	sw.Do("\n", nil)

	return sw.Error()
}

func (g *genOption) isOtherPackage(pkg string) bool {
	if pkg == g.targetPackage {
		return false
	}
	if strings.HasSuffix(pkg, "\""+g.targetPackage+"\"") {
		return false
	}
	return true
}

func (g *genOption) Imports(c *generator.Context) (imports []string) {
	var importLines []string
	for _, singleImport := range g.imports.ImportLines() {
		if g.isOtherPackage(singleImport) {
			importLines = append(importLines, singleImport)
		}
	}
	return importLines
}

func getOptionName(typeName string) string {
	if strings.HasSuffix(typeName, "Option") {
		return typeName
	}
	return typeName + "Option"
}

func (g *genOption) getOptionSuffixName(t *types.Type) string {
	suffix := getOptionName(t.Name.Name)
	if v := g.genTypes.find(t); v != nil {
		suffix = v.OptionSuffix
	}
	return suffix
}

func (g *genOption) getOptionTypeName(t *types.Type) string {
	name := getOptionName(t.Name.Name)
	if v := g.genTypes.find(t); v != nil {
		name = v.OptionName
	}
	return name
}

func (g *genOption) genOptionType(sw *generator.SnippetWriter, t *types.Type) {
	name := g.getOptionTypeName(t)
	if g.genTypes.existMethod(name) {
		return
	}
	args := generator.Args{
		"type": t,
		"name": name,
	}
	sw.Do("type {{.name}} func(* {{.type|public}})\n", args)
}

func (g *genOption) genWithFieldFunc(sw *generator.SnippetWriter, t *types.Type) {
	isExternalType := g.packageTypes.IsExternalType(t.Name.Package, t.Name.Name)
	option := g.getOptionTypeName(t)
	suffix := g.getOptionSuffixName(t)
	methodGen := util.NewMethodGenerate(util.GenName("With", suffix))
	for idx, m := range t.Members {
		if util.IsLower(m.Name) && isExternalType {
			continue
		}
		if !g.genTypes.allowedField(t, idx) {
			continue
		}
		method := methodGen.GenName(m.Name)
		if g.genTypes.existMethod(method) {
			continue
		}
		args := generator.Args{
			"type":   t,
			"option": option,
			"field":  m,
			"method": method,
		}
		sw.Do("func {{.method}} (val {{.field.Type|raw}}) {{.option}} {\n", args)
		sw.Do("return func(object * {{.type|public}}) {\n", args)
		sw.Do("object.{{.field.Name}} = val\n", args)
		sw.Do("}\n", nil)
		sw.Do("}\n\n", nil)
	}
}
