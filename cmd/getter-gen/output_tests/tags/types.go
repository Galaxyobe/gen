/*
Copyright 2022 The Gen Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tags

import (
	"github.com/galaxyobe/gen/cmd/getter-gen/output_tests/builtins"
)

// +gen:getter=true
type Builtins builtins.Builtins

// +gen:getter=false
type Builtins2 builtins.Builtins

// +gen:getter=true
type Structs struct {
	b      byte
	String string
	Builtins
	Builtins1 Builtins
	Builtins2 *Builtins
	// +getter=false
	Builtins3 *Builtins
}

// +gen:getter:fields=b,s,String,Builtins,Builtins2
type Structs2 struct {
	b byte
	u uint8
	// +getter=true
	i int8 // will be ignored
	// +getter=false
	s      string // will be ignored
	String string
	Builtins
	Builtins1 Builtins
	Builtins2 *Builtins
	// +getter=false
	Builtins3 *Builtins
}

func (s *Structs2) Gets() string {
	return s.s
}

// +gen:getter=true
type Structs3 struct {
	// +getter:method=Get8Bits
	// +getter=false
	u8 uint8
	// +getter=true
	// +getter:method=8Bits
	i8 int8
	s  string
	// +getter:method=GetByte
	b byte
}
