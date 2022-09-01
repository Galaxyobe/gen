package types

import (
	"k8s.io/gengo/types"
)

var (
	Uint8 = &types.Type{
		Name: types.Name{Name: "uint8"},
		Kind: types.Builtin,
	}
	Int8 = &types.Type{
		Name: types.Name{Name: "int8"},
		Kind: types.Builtin,
	}
)

func ReplacePackageTypes(m map[string]*types.Type) {
	m["uint8"] = Uint8
	m["int8"] = Int8
}
