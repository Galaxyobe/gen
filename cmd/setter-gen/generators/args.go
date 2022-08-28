package generators

import (
	"github.com/spf13/pflag"
)

// CustomArgs is used tby the go2idl framework to pass args specific to this
// generator.
type CustomArgs struct {
	BoundingDirs    []string // Only deal with types rooted under these dirs.
	TrimPackagePath string   // If specified, trim the path from PackagePath before writing files.
}

func (a *CustomArgs) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&a.BoundingDirs, "bounding-dirs", a.BoundingDirs,
		"Comma-separated list of import paths which bound the types for which deep-copies will be generated.")
	fs.StringVar(&a.TrimPackagePath, "trim-package-path", a.TrimPackagePath,
		"If set, trim the specified path from PackagePath when generating files.")
}
