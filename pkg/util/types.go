package util

import (
	"strings"

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

func IsLower(s string) bool {
	if s == "" {
		return false
	}
	return strings.ToLower(s[:1]) == s[:1]
}
