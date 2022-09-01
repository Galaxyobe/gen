package types

import (
	"k8s.io/gengo/types"
)

var (
	Int8 = &types.Type{
		Name: types.Name{Name: "int8"},
		Kind: types.Builtin,
	}
	Uint8 = &types.Type{
		Name: types.Name{Name: "uint8"},
		Kind: types.Builtin,
	}
)
