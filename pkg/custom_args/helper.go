package custom_args

import (
	"k8s.io/gengo/args"
)

func GetCustomArgs(args *args.GeneratorArgs) *CustomArgs {
	customArgs, ok := args.CustomArgs.(*CustomArgs)
	if !ok {
		return NewCustomArgs(args)
	}
	return customArgs
}
