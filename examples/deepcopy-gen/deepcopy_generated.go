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

// Code generated by deepcopy-gen. DO NOT EDIT.

package deepcopy_gen

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Builtin) DeepCopyInto(out *Builtin) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Builtin.
func (in *Builtin) DeepCopy() *Builtin {
	if in == nil {
		return nil
	}
	out := new(Builtin)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Map) DeepCopyInto(out *Map) {
	{
		in := &in
		*out = make(Map, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Map.
func (in Map) DeepCopy() Map {
	if in == nil {
		return nil
	}
	out := new(Map)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Slice) DeepCopyInto(out *Slice) {
	{
		in := &in
		*out = make(Slice, len(*in))
		copy(*out, *in)
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Slice.
func (in Slice) DeepCopy() Slice {
	if in == nil {
		return nil
	}
	out := new(Slice)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Struct) DeepCopyInto(out *Struct) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Struct.
func (in *Struct) DeepCopy() *Struct {
	if in == nil {
		return nil
	}
	out := new(Struct)
	in.DeepCopyInto(out)
	return out
}
