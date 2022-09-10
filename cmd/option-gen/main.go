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

// option-gen is a tool for auto-generating with struct field functions.
//
// option-gen can be generated for individual types, rather than the
// entire package by specifying a comment on the type or filed definition of the form:
//
//	// +gen:option=true
//
// You can not participate in the generation by setting the field a
// comment on the type definition of the form:
//
//	// +gen:option=true
//	type Struct struct {
//		// +option=false
//		Field string
//	}
//
// You can also specify the field name that needs to be generated a
// comment on the type definition of the form:
//
//	// +gen:option:fields=field1,field2,field3
//
// Option names and suffixes can be customized by setting the field a
// comment on the type definition of the form:
//
//	// +gen:option:name=StructOptionName
//
//	// +gen:option:suffix=StructOptionNameSuffix
package main

import (
	"k8s.io/gengo/args"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/galaxyobe/gen/cmd/option-gen/generators"
	custom_args "github.com/galaxyobe/gen/pkg/custom_args"
	"github.com/galaxyobe/gen/pkg/util"
)

// main option-gen at project root, program flag:
// -v 7 -i github.com/galaxyobe/gen/cmd/option-gen/output_tests/... --trim-path-prefix github.com/galaxyobe/gen -o .
func main() {
	klog.InitFlags(nil)
	arguments := args.Default()

	// Override defaults.
	arguments.GoHeaderFilePath = util.BoilerplatePath()
	arguments.OutputFileBaseName = "option_generated"

	// Custom custom_args.
	customArgs := custom_args.NewCustomArgs(arguments)
	customArgs.AddFlags(pflag.CommandLine)

	// Validate checks the given arguments.
	if len(arguments.OutputFileBaseName) == 0 {
		klog.Fatalf("output file base name cannot be empty")
	}

	// Run it.
	if err := arguments.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}
	klog.V(2).Info("Completed successfully.")
}
