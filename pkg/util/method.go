package util

import (
	"strings"

	"k8s.io/gengo/types"
)

type MethodSet map[string][]types.Member // key: method name, value: types.Member

func NewMethodSet() MethodSet {
	return make(MethodSet)
}

func (m MethodSet) AddMethod(prefix, method string, member types.Member) {
	if !strings.HasPrefix(method, prefix) {
		method = prefix + method
	}
	list, ok := m[method]
	if !ok {
		m[method] = []types.Member{member}
		return
	}
	list = append(list, member)
	m[method] = list
}

func (m MethodSet) AddMethods(prefix string, methods []string, member types.Member) {
	for _, method := range methods {
		m.AddMethod(prefix, method, member)
	}
}
