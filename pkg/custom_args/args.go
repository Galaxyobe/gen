package custom_args

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"k8s.io/gengo/args"

	"github.com/galaxyobe/gen/third_party/gengo/parser"
)

// CustomArgs is used tby the go2idl framework to pass custom_args specific to this
// generator.
type CustomArgs struct {
	*args.GeneratorArgs
	BoundingDirs    []string // Only deal with types rooted under these dirs.
	TrimPackagePath string   // If specified, trim the path from PackagePath before writing files.
}

func NewCustomArgs(args *args.GeneratorArgs) *CustomArgs {
	customArgs := &CustomArgs{
		GeneratorArgs: args,
	}
	args.CustomArgs = customArgs
	return customArgs
}

func (a *CustomArgs) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&a.BoundingDirs, "bounding-dirs", a.BoundingDirs,
		"Comma-separated list of import paths which bound the types for which deep-copies will be generated.")
	fs.StringVar(&a.TrimPackagePath, "trim-package-path", a.TrimPackagePath,
		"If set, trim the specified path from PackagePath when generating files.")
}

// NewBuilder makes a new parser.Builder and populates it with the input
// directories.
func (a *CustomArgs) NewBuilder() (*parser.Builder, error) {
	b := parser.New()

	// flag for including *_test.go
	b.IncludeTestFiles = a.IncludeTestFiles

	// Ignore all auto-generated files.
	b.AddBuildTags(a.GeneratedBuildTag)

	for _, d := range a.InputDirs {
		var err error
		if strings.HasSuffix(d, "/...") {
			err = b.AddDirRecursive(strings.TrimSuffix(d, "/..."))
		} else {
			err = b.AddDir(d)
		}
		if err != nil {
			return nil, fmt.Errorf("unable to add directory %q: %v", d, err)
		}
	}
	return b, nil
}
