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

package builtins

type (
	Uint8   uint8
	Uint    uint
	Int8    int8
	Int     int
	StrFunc func(string) error
	I8Func  func(int8) error
	U8Func  func(uint8) error
)

type Builtins struct {
	U8       Uint8
	I8       Int8
	U        Uint
	I        Int
	StrFunc  StrFunc
	I8Func   I8Func
	U8Func   U8Func
	Int8     int8
	Uint8    uint8
	a        any
	A        any
	b        byte
	i        int
	u        uint
	Bool     bool
	Byte     byte
	Int16    int16
	Int32    int32
	Int64    int64
	Uint16   uint16
	Uint32   uint32
	Uint64   uint64
	Float32  float32
	Float64  float64
	String   string
	Bytes    []byte
	BoolP    *bool
	ByteP    *byte
	IntP     *int
	UintP    *uint
	BytesP   *[]byte
	Float64P *float64
	StringP  *string
}

func (a *Alias) SetA(v interface{}) *Alias {
	a.A = v
	return a
}

type Alias Builtins
type Alias2 *Builtins
type Alias3 = Builtins
type AliasString string
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
	builtinsS  []Builtins
	BuiltinsS  []*Builtins
	BuiltinsPs *[]Builtins
}
