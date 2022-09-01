package generators

import (
	"io"
	"strings"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"

	pkgtypes "github.com/galaxyobe/gen/pkg/types"
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
		comments := t.SecondClosestCommentLines
		comments = append(comments, t.CommentLines...)
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

type MethodSet map[string][]types.Member // key: method name value: types.Member

func NewMethodSet() MethodSet {
	return make(MethodSet)
}

func (m MethodSet) AddMethod(method string, member types.Member) {
	list, ok := m[method]
	if !ok {
		m[method] = []types.Member{member}
		return
	}
	list = append(list, member)
	m[method] = list
}

func (m MethodSet) AddMethods(methods []string, member types.Member) {
	for _, method := range methods {
		m.AddMethod(method, member)
	}
}

type genSetter struct {
	generator.DefaultGen
	build         *parser.Builder
	targetPackage string
	boundingDirs  []string
	imports       namer.ImportTracker
	types         GenTypes
	packageTypes  pkgtypes.PackageTypes
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
		packageTypes:  pkgtypes.NewPackageTypes(build),
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
	g.updateTypeMembers(t)
	sw := generator.NewSnippetWriter(w, c, "", "")
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

func (g *genSetter) updateTypeMembers(t *types.Type) {
	uint8Fields := g.packageTypes.GetUint8Fields(t.Name.Package, t.Name.Name)
	int8Fields := g.packageTypes.GetInt8Fields(t.Name.Package, t.Name.Name)

	for i, m := range t.Members {
		if util.Exist(uint8Fields, m.Name) {
			t.Members[i].Type = pkgtypes.Uint8
		} else if util.Exist(int8Fields, m.Name) {
			t.Members[i].Type = pkgtypes.Int8
		}
	}
}

func (g *genSetter) genSetFunc(sw *generator.SnippetWriter, t *types.Type) {
	receiver := strings.ToLower(t.Name.Name[:1])
	isExternalType := g.packageTypes.IsExternalType(t.Name.Package, t.Name.Name)
	var methodSet = NewMethodSet()
	var genMethodSet = make(map[string]struct{})
	for idx, m := range t.Members {
		if isExternalType && util.IsLower(m.Name) {
			continue
		}
		methods := util.GetTagValues(tagMethodName, m.CommentLines)
		methodSet.AddMethods(methods, t.Members[idx])
		if !g.types.allowedField(t, idx) {
			continue
		}
		method := "Set" + m.Name
		if _, ok := t.Methods[method]; ok {
			continue
		}
		genMethodSet[method] = struct{}{}
		args := generator.Args{
			"type":     t,
			"field":    m,
			"receiver": receiver,
			"method":   method,
		}
		sw.Do("func ({{.receiver}} *{{.type|public}}) {{.method}}(val {{.field.Type|raw}}) *{{.type|public}} {\n", args)
		sw.Do("{{.receiver}}.{{.field.Name}} = val\n", args)
		sw.Do("return {{.receiver}}", args)
		sw.Do("}\n\n", nil)
	}
	for name := range t.Methods {
		genMethodSet[name] = struct{}{}
	}
	for method, members := range methodSet {
		if !strings.HasPrefix(method, "Set") {
			method = "Set" + method
		}
		if _, ok := genMethodSet[method]; ok {
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
