package util

import (
	"reflect"
	"strings"

	"k8s.io/gengo/types"
)

const (
	Package = "package"
	True    = "true"
	False   = "false"
)

func CheckTag(tag string, comments []string, require ...string) bool {
	values := types.ExtractCommentTags("+", comments)[tag]
	if len(require) == 0 {
		return len(values) == 1 && values[0] == ""
	}
	if values == nil {
		return false
	}
	return reflect.DeepEqual(values, require)
}

func GetTagBoolStatus(tag string, comments []string) (set bool, enabled bool) {
	values := types.ExtractCommentTags("+", comments)[tag]
	switch len(values) {
	case 0:
		enabled = true
		return
	default:
		set = true
		enabled = strings.Split(values[0], ",")[0] == True
		return
	}
}

func GetTagValues(tag string, comments []string) []string {
	values := types.ExtractCommentTags("+", comments)[tag]
	if len(values) == 0 {
		return nil
	}
	var result []string
	for _, item := range values {
		result = append(result, strings.Split(item, ",")...)
	}
	return result
}

func GetTagValuesStatus(tag string, comments []string) (bool, []string) {
	values := types.ExtractCommentTags("+", comments)[tag]
	if len(values) == 0 {
		return false, nil
	}
	var result []string
	for _, item := range values {
		result = append(result, strings.Split(item, ",")...)
	}
	return true, result
}

func GetTagValueStatus(tag string, comments []string) (bool, string) {
	set, values := GetTagValuesStatus(tag, comments)
	if len(values) > 0 {
		return set, values[0]
	}
	return false, ""
}
