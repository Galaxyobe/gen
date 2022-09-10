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

// +gen:deepcopy=package

// This is a test package.
package pointer

type Ttest struct {
	Builtin   *string
	Ptr       **string
	Map       *map[string]string
	Slice     *[]string
	MapPtr    **map[string]string
	SlicePtr  **[]string
	Struct    *Ttest
	StructPtr **Ttest
}
