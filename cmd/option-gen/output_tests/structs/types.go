/*
 Copyright 2022 Galaxyobe.

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

package structs

import (
	"github.com/galaxyobe/gen/cmd/option-gen/output_tests/builtins"
)

// +gen:option=true
type Alias builtins.Builtins

func (a *Alias) SetA(v interface{}) *Alias {
	a.A = v
	return a
}

// +gen:option=true
type Alias2 builtins.Alias2

// +gen:option=true
type Alias3 = builtins.Alias3

// +gen:option=true
type Age int

// +gen:option=true
type User struct {
	Name string
	Age  Age
}

// +gen:option=true
type UserInfo User

// +gen:option=true
type Structs struct {
	b      byte
	u      uint8
	i      int8
	String string
	builtins.Builtins
	Builtins1 builtins.Builtins
	Builtins2 *builtins.Builtins
	User
	User2 User
	User3 *User
	Age   Age
}
