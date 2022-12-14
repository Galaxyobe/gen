//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by ___deepcopy_gen. DO NOT EDIT.

package pointer

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ttest) DeepCopyInto(out *Ttest) {
	*out = *in
	if in.Builtin != nil {
		in, out := &in.Builtin, &out.Builtin
		*out = new(string)
		**out = **in
	}
	if in.Ptr != nil {
		in, out := &in.Ptr, &out.Ptr
		*out = new(*string)
		if **in != nil {
			in, out := *in, *out
			*out = new(string)
			**out = **in
		}
	}
	if in.Map != nil {
		in, out := &in.Map, &out.Map
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
	if in.Slice != nil {
		in, out := &in.Slice, &out.Slice
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	if in.MapPtr != nil {
		in, out := &in.MapPtr, &out.MapPtr
		*out = new(*map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = new(map[string]string)
			if **in != nil {
				in, out := *in, *out
				*out = make(map[string]string, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
		}
	}
	if in.SlicePtr != nil {
		in, out := &in.SlicePtr, &out.SlicePtr
		*out = new(*[]string)
		if **in != nil {
			in, out := *in, *out
			*out = new([]string)
			if **in != nil {
				in, out := *in, *out
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
		}
	}
	if in.Struct != nil {
		in, out := &in.Struct, &out.Struct
		*out = new(Ttest)
		(*in).DeepCopyInto(*out)
	}
	if in.StructPtr != nil {
		in, out := &in.StructPtr, &out.StructPtr
		*out = new(*Ttest)
		if **in != nil {
			in, out := *in, *out
			*out = new(Ttest)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ttest.
func (in *Ttest) DeepCopy() *Ttest {
	if in == nil {
		return nil
	}
	out := new(Ttest)
	in.DeepCopyInto(out)
	return out
}
