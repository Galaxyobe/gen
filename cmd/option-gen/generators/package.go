/*
Copyright 2015 The Kubernetes Authors.

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

package generators

import (
	"fmt"
	"path/filepath"
	"strings"

	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/set-gen/sets"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"

	"github.com/galaxyobe/gen/pkg/custom_args"
	"github.com/galaxyobe/gen/pkg/util"
)

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public":  namer.NewPublicNamer(0),
		"private": namer.NewPrivateNamer(0),
		"raw":     namer.NewRawNamer("", nil),
	}
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}

func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		klog.Fatalf("Failed loading boilerplate: %v", err)
	}

	inputs := sets.NewString(context.Inputs...)
	packages := generator.Packages{}
	header := append([]byte(fmt.Sprintf("//go:build !%s\n// +build !%s\n\n", arguments.GeneratedBuildTag, arguments.GeneratedBuildTag)), boilerplate...)

	customArgs := custom_args.GetCustomArgs(arguments)
	if customArgs.BoundingDirs == nil {
		customArgs.BoundingDirs = context.Inputs
	}

	var boundingDirs []string
	for i := range customArgs.BoundingDirs {
		// Strip any trailing slashes - they are not exactly "correct" but
		// this is friendlier.
		boundingDirs = append(boundingDirs, strings.TrimRight(customArgs.BoundingDirs[i], "/"))
	}

	// new parse build to replace context universe for basic int8/uint8 types
	build, err := customArgs.NewBuilder()
	if err != nil {
		klog.Fatalf("Failed making a parser: %v", err)
	}
	// replace universe
	err = build.ReplaceUniverse(context.Universe)
	if err != nil {
		klog.Fatalf("Find types error: %v", err)
	}
	// reorder
	util.ReContextOrder(context, DefaultNameSystem())

	for i := range inputs {
		klog.V(5).Infof("Considering pkg %q", i)
		pkg := context.Universe[i]
		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}
		pkgEnabled, genTypes := NewGenTypes(pkg)
		if !pkgEnabled {
			continue
		}
		klog.V(3).Infof("Package %q needs generation", pkg.Path)
		path := pkg.Path
		if strings.HasPrefix(pkg.SourcePath, arguments.OutputBase) {
			expandedPath := strings.TrimPrefix(pkg.SourcePath, arguments.OutputBase)
			if strings.Contains(expandedPath, "/vendor/") {
				path = expandedPath
			}
		}
		if customArgs.TrimPackagePath != "" {
			path = strings.ReplaceAll(path, customArgs.TrimPackagePath, "")
			separator := string(filepath.Separator)
			if path != "" && strings.HasPrefix(path, separator) {
				path = path[1:]
			}
		}

		packages = append(packages,
			&generator.DefaultPackage{
				PackageName: pkg.Name,
				PackagePath: path,
				HeaderText:  header,
				GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
					return []generator.Generator{
						NewGenOption(
							build,
							arguments.OutputFileBaseName,
							pkg.Path,
							boundingDirs,
							genTypes,
							pkg.SourcePath,
						),
					}
				},
				FilterFunc: func(c *generator.Context, t *types.Type) bool {
					return t.Name.Package == pkg.Path
				},
			})
	}
	return packages
}
