package generators

import (
	"io"
	"strings"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"

	"github.com/galaxyobe/gen/pkg/util"
)

type GenType struct {
	*types.Type
	AllowFields []string
}

func NewGenTypes(pkg *types.Package) (pkgEnabled bool, list GenTypes) {
	pkgAllowed := util.CheckTag(tagName, pkg.Comments, util.Package)
	for _, t := range pkg.Types {
		ut := util.UnderlyingType(t)
		if ut.Kind != types.Struct {
			continue
		}
		comments := t.SecondClosestCommentLines
		comments = append(comments, t.CommentLines...)
		set, enabled := util.GetTagBoolStatus(tagName, comments)
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

type genSetter struct {
	generator.DefaultGen
	targetPackage string
	boundingDirs  []string
	imports       namer.ImportTracker
	types         GenTypes
	int8s         map[string][]string
	uint8s        map[string][]string
}

func NewGenSetter(sanitizedName, targetPackage string, boundingDirs []string, types []*GenType, int8s, uint8s map[string][]string) generator.Generator {
	return &genSetter{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		targetPackage: targetPackage,
		boundingDirs:  boundingDirs,
		imports:       generator.NewImportTracker(),
		types:         types,
		int8s:         int8s,
		uint8s:        uint8s,
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

	sw := generator.NewSnippetWriter(w, c, "$", "$")
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

func (g *genSetter) convertFieldToInt8(typeName, fieldName string) bool {
	list := g.int8s[typeName]
	if len(list) == 0 {
		return false
	}
	for _, item := range list {
		if item == fieldName {
			return true
		}
	}
	return false
}

func (g *genSetter) convertFieldToUint8(typeName, fieldName string) bool {
	list := g.uint8s[typeName]
	if len(list) == 0 {
		return false
	}
	for _, item := range list {
		if item == fieldName {
			return true
		}
	}
	return false
}

func (g *genSetter) genSetFunc(sw *generator.SnippetWriter, t *types.Type) {
	receiver := strings.ToLower(t.Name.Name[:1])
	for _, m := range t.Members {
		method := "Set" + m.Name
		if _, ok := t.Methods[method]; ok {
			continue
		}
		args := generator.Args{
			"type":     t,
			"field":    m,
			"receiver": receiver,
			"method":   method,
			"byte":     m.Type.Name.Name == "byte",
		}
		if g.convertFieldToInt8(t.Name.Name, m.Name) {
			sw.Do("func ($.receiver$ *$.type|public$) $.method$(val int8) *$.type|public$ {\n", args)
		} else if g.convertFieldToUint8(t.Name.Name, m.Name) {
			sw.Do("func ($.receiver$ *$.type|public$) $.method$(val uint8) *$.type|public$ {\n", args)
		} else {
			sw.Do("func ($.receiver$ *$.type|public$) $.method$(val $.field.Type|raw$) *$.type|public$ {\n", args)
		}
		sw.Do("$.receiver$.$.field.Name$ = val\n", args)
		sw.Do("return $.receiver$", args)
		sw.Do("}\n\n", nil)
	}
}
