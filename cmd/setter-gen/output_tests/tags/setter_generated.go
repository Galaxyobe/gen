//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Gen Authors.

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

// Code generated by ___setter_gen. DO NOT EDIT.

package tags

import (
	builtins "github.com/galaxyobe/gen/cmd/setter-gen/output_tests/builtins"
)

func (b *Builtins) SetU8(val builtins.Uint8) *Builtins {
	b.U8 = val
	return b
}

func (b *Builtins) SetI8(val builtins.Int8) *Builtins {
	b.I8 = val
	return b
}

func (b *Builtins) SetU(val builtins.Uint) *Builtins {
	b.U = val
	return b
}

func (b *Builtins) SetI(val builtins.Int) *Builtins {
	b.I = val
	return b
}

func (b *Builtins) SetStrFunc(val builtins.StrFunc) *Builtins {
	b.StrFunc = val
	return b
}

func (b *Builtins) SetI8Func(val builtins.I8Func) *Builtins {
	b.I8Func = val
	return b
}

func (b *Builtins) SetU8Func(val builtins.U8Func) *Builtins {
	b.U8Func = val
	return b
}

func (b *Builtins) SetInt8(val int8) *Builtins {
	b.Int8 = val
	return b
}

func (b *Builtins) SetUint8(val uint8) *Builtins {
	b.Uint8 = val
	return b
}

func (b *Builtins) SetA(val interface{}) *Builtins {
	b.A = val
	return b
}

func (b *Builtins) SetBool(val bool) *Builtins {
	b.Bool = val
	return b
}

func (b *Builtins) SetByte(val byte) *Builtins {
	b.Byte = val
	return b
}

func (b *Builtins) SetInt16(val int16) *Builtins {
	b.Int16 = val
	return b
}

func (b *Builtins) SetInt32(val int32) *Builtins {
	b.Int32 = val
	return b
}

func (b *Builtins) SetInt64(val int64) *Builtins {
	b.Int64 = val
	return b
}

func (b *Builtins) SetUint16(val uint16) *Builtins {
	b.Uint16 = val
	return b
}

func (b *Builtins) SetUint32(val uint32) *Builtins {
	b.Uint32 = val
	return b
}

func (b *Builtins) SetUint64(val uint64) *Builtins {
	b.Uint64 = val
	return b
}

func (b *Builtins) SetFloat32(val float32) *Builtins {
	b.Float32 = val
	return b
}

func (b *Builtins) SetFloat64(val float64) *Builtins {
	b.Float64 = val
	return b
}

func (b *Builtins) SetString(val string) *Builtins {
	b.String = val
	return b
}

func (b *Builtins) SetBytes(val []byte) *Builtins {
	b.Bytes = val
	return b
}

func (b *Builtins) SetBoolP(val *bool) *Builtins {
	b.BoolP = val
	return b
}

func (b *Builtins) SetByteP(val *byte) *Builtins {
	b.ByteP = val
	return b
}

func (b *Builtins) SetIntP(val *int) *Builtins {
	b.IntP = val
	return b
}

func (b *Builtins) SetUintP(val *uint) *Builtins {
	b.UintP = val
	return b
}

func (b *Builtins) SetBytesP(val *[]byte) *Builtins {
	b.BytesP = val
	return b
}

func (b *Builtins) SetFloat64P(val *float64) *Builtins {
	b.Float64P = val
	return b
}

func (b *Builtins) SetStringP(val *string) *Builtins {
	b.StringP = val
	return b
}

func (s *Structs) Setb(val byte) *Structs {
	s.b = val
	return s
}

func (s *Structs) SetString(val string) *Structs {
	s.String = val
	return s
}

func (s *Structs) SetBuiltins(val Builtins) *Structs {
	s.Builtins = val
	return s
}

func (s *Structs) SetBuiltins1(val Builtins) *Structs {
	s.Builtins1 = val
	return s
}

func (s *Structs) SetBuiltins2(val *Builtins) *Structs {
	s.Builtins2 = val
	return s
}

func (s *Structs2) Setb(val byte) *Structs2 {
	s.b = val
	return s
}

func (s *Structs2) Sets(val string) *Structs2 {
	s.s = val
	return s
}

func (s *Structs2) SetString(val string) *Structs2 {
	s.String = val
	return s
}

func (s *Structs2) SetBuiltins(val Builtins) *Structs2 {
	s.Builtins = val
	return s
}

func (s *Structs2) SetBuiltins2(val *Builtins) *Structs2 {
	s.Builtins2 = val
	return s
}

func (s *Structs3) Seti8(val int8) *Structs3 {
	s.i8 = val
	return s
}

func (s *Structs3) Sets(val string) *Structs3 {
	s.s = val
	return s
}

func (s *Structs3) Setb(val byte) *Structs3 {
	s.b = val
	return s
}

func (s *Structs3) Set8Bits(in0 uint8, in1 int8) *Structs3 {
	s.u8 = in0
	s.i8 = in1
	return s
}
