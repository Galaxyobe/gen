package types

import (
	"k8s.io/gengo/types"
)

func ModifyType(u types.Universe, typeName string, t *types.Type) *types.Type {
	p := u.Package(t.Name.Package)
	if p.Path != "" {
		return t
	}
	// modify standard builtin types!
	if t, ok := Builtins.Types[typeName]; ok {
		p.Types[typeName] = t
		return t
	}
	return t
}

// ReplaceUniverse replace u2 to u1 Universe
func ReplaceUniverse(u1, u2 types.Universe) {
	for k2, v2 := range u2 {
		if v2 == nil {
			continue
		}
		_, ok := u1[k2]
		if !ok {
			continue
		}
		u1[k2] = v2
	}
}
