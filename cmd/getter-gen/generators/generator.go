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

// allowedField allowed type's field to generate Getter func.
// will be ignored field enabled status when it's in +gen:getter:fields allowed fields.
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

type genGetter struct {
	generator.DefaultGen
	build         *parser.Builder
	targetPackage string
	boundingDirs  []string
	imports       namer.ImportTracker
	types         GenTypes
	packageTypes  util.PackageTypes
}

func NewGenGetter(build *parser.Builder, sanitizedName, targetPackage string, boundingDirs []string, types []*GenType, sourcePath string) generator.Generator {
	return &genGetter{
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

func (g *genGetter) Name() string {
	return "getter"
}

func (g *genGetter) Filter(c *generator.Context, t *types.Type) bool {
	ok := g.types.allowed(t)
	if !ok {
		klog.V(5).Infof("Ignore generate getter function for type %v", t)
	}
	return ok
}

func (g *genGetter) Namers(c *generator.Context) namer.NameSystems {
	// Have the raw namer for this file track what it imports.
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.targetPackage, g.imports),
	}
}

func (g *genGetter) Init(c *generator.Context, w io.Writer) error {
	return nil
}

func (g *genGetter) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	klog.V(5).Infof("Generating getter function for type %v", t)

	sw := generator.NewSnippetWriter(w, c, "", "")
	g.genGetFunc(sw, t)
	sw.Do("\n", nil)

	return sw.Error()
}

func (g *genGetter) isOtherPackage(pkg string) bool {
	if pkg == g.targetPackage {
		return false
	}
	if strings.HasSuffix(pkg, "\""+g.targetPackage+"\"") {
		return false
	}
	return true
}

func (g *genGetter) Imports(c *generator.Context) (imports []string) {
	var importLines []string
	for _, singleImport := range g.imports.ImportLines() {
		if g.isOtherPackage(singleImport) {
			importLines = append(importLines, singleImport)
		}
	}
	return importLines
}

func (g *genGetter) genGetFunc(sw *generator.SnippetWriter, t *types.Type) {
	receiver := strings.ToLower(t.Name.Name[:1])
	isExternalType := g.packageTypes.IsExternalType(t.Name.Package, t.Name.Name)
	var methodSet = util.NewMethodSet()
	var methodGen = util.NewMethodGenerate(util.GenName("Get", ""))

	for idx, m := range t.Members {
		if util.IsLower(m.Name) && isExternalType {
			continue
		}
		methods := util.GetTagValues(tagMethodName, m.CommentLines)
		methodSet.AddMethods("Get", methods, t.Members[idx])
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
		sw.Do("func ({{.receiver}} *{{.type|public}}) {{.method}} () {{.field.Type|raw}} {\n", args)
		sw.Do("return {{.receiver}}.{{.field.Name}}\n", args)
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
		sw.Do("func ({{.receiver}} *{{.type|public}}) {{.method}} () ({{ range $i, $field := .fields }} in{{$i}} {{$field.Type|raw}}, {{end}}) {\n", args)
		sw.Do("{{ range $i, $field := .fields }} in{{$i}} = {{$.receiver}}.{{$field.Name}}\n {{ end }}", args)
		sw.Do("return", nil)
		sw.Do("}\n\n", nil)
	}
}
