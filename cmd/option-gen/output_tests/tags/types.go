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

// +gen:option=true
// type Builtins builtins.Builtins

// +gen:option=false
// type Builtins2 builtins.Builtins

// +gen:option=true
// type Structs struct {
// 	b      byte
// 	String string
// 	Builtins
// 	Builtins1 Builtins
// 	Builtins2 *Builtins
// 	// +option=false
// 	Builtins3 *Builtins
// }

// +gen:option:fields=b,s,String,Builtins,Builtins2
// +gen:option:suffix
// type Structs2 struct {
// 	b byte
// 	u uint8
// 	// +option=true
// 	i int8 // will be ignored
// 	// +option=false
// 	s      string // will be ignored
// 	String string
// 	Builtins
// 	Builtins1 Builtins
// 	Builtins2 *Builtins
// 	// +option=false
// 	Builtins3 *Builtins
// }

// +gen:option:name=NameStructs3
type Structs3 struct {
	// +option=false
	u8 uint8
	// +option=true
	i8 int8
	s  string
	b  byte
	U8 uint
	I8 int8
}
