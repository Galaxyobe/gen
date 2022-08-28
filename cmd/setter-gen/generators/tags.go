package generators

import (
	"reflect"
	"strings"

	"k8s.io/gengo/types"
	"k8s.io/klog/v2"
)

const (
	tagName = "gen:setter"
)

// Known values for the comment tag.
const tagValuePackage = "package"

func extractTag(comments []string) []string {
	return types.ExtractCommentTags("+", comments)[tagName]
}

func checkTag(comments []string, require ...string) bool {
	values := types.ExtractCommentTags("+", comments)[tagName]
	if len(require) == 0 {
		return len(values) == 1 && values[0] == ""
	}
	return reflect.DeepEqual(values, require)
}

// enabledTagValue holds parameters from a tagName tag.
type enabledTagValue struct {
	value string
}

func extractEnabledTypeTag(t *types.Type) *enabledTagValue {
	comments := append(append([]string{}, t.SecondClosestCommentLines...), t.CommentLines...)
	return extractEnabledTag(comments)
}

func extractEnabledTag(comments []string) *enabledTagValue {
	tagVals := types.ExtractCommentTags("+", comments)[tagName]
	if tagVals == nil {
		// No match for the tag.
		return nil
	}
	// If there are multiple values, abort.
	if len(tagVals) > 1 {
		klog.Fatalf("Found %d %s tags: %q", len(tagVals), tagName, tagVals)
	}

	// If we got here we are returning something.
	tag := &enabledTagValue{}

	// Get the primary value.
	parts := strings.Split(tagVals[0], ",")
	if len(parts) >= 1 {
		tag.value = parts[0]
	}

	return tag
}
