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

	Builtins = &types.Package{
		Types: map[string]*types.Type{
			"int8":  Int8,
			"uint8": Uint8,
		},
		Imports: map[string]*types.Package{},
		Path:    "",
		Name:    "",
	}
)
