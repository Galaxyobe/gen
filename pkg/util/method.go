package util

import (
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
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

type MethodGenerate struct {
	set      map[string]int
	gen      func(name string) string
	handlers []func(string) string
}

func NewMethodGenerate(fn func(name string) string, handler ...func(string) string) MethodGenerate {
	if len(handler) == 0 {
		handler = append(handler, strcase.ToCamel)
	}
	return MethodGenerate{
		set:      make(map[string]int),
		gen:      fn,
		handlers: handler,
	}
}

func (m *MethodGenerate) handle(name string) string {
	for _, f := range m.handlers {
		name = f(name)
	}
	return name
}

func (m *MethodGenerate) AddExistName(name string) (bool, int) {
	count, ok := m.set[name]
	if ok {
		count += 1
	} else {
		count = 1
	}
	m.set[name] = count
	return ok, count
}

func (m *MethodGenerate) AddExistNames(name ...string) *MethodGenerate {
	for _, item := range name {
		m.AddExistName(item)
	}
	return m
}

func (m *MethodGenerate) ExistName(name string) bool {
	_, ok := m.set[name]
	return ok
}

func (m *MethodGenerate) GenName(name string) string {
	name = m.handle(name)
	method := m.gen(name)
	if ok, count := m.AddExistName(method); ok {
		method = m.gen(name + strconv.Itoa(count))
	}
	return method
}

func GenName(prefix, suffix string) func(name string) string {
	return func(name string) string {
		return prefix + name + suffix
	}
}
