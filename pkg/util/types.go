package util

import (
	"k8s.io/gengo/types"
)

func UnderlyingType(t *types.Type) *types.Type {
	for t.Kind == types.Alias {
		t = t.Underlying
	}
	return t
}

// IsReference return true for pointer, maps, slices and aliases of those.
func IsReference(t *types.Type) bool {
	if t.Kind == types.Pointer || t.Kind == types.Map || t.Kind == types.Slice {
		return true
	}
	return t.Kind == types.Alias && IsReference(UnderlyingType(t))
}
