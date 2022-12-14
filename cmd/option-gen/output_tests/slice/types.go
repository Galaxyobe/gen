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

package slice

import (
	"github.com/galaxyobe/gen/cmd/option-gen/output_tests/builtins"
)

// +gen:option=true
type B struct{}

// +gen:option=true
type Slice struct {
	i8S        []int8
	I8S        []*int8
	I8pS       *[]int8
	u8S        []uint8
	U8S        []*uint8
	U8pS       *[]uint8
	sS         []string
	SS         []*string
	SpS        *[]string
	bS         []byte
	BS         []*byte
	BpS        *[]byte
	BBS        []B
	builtinsS  []builtins.Builtins
	BuiltinsS  []*builtins.Builtins
	BuiltinsPs *[]builtins.Builtins
}

// +gen:option=true
type AliasSlice builtins.Slice
