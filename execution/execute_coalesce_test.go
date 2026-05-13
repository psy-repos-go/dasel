package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestCoalesce(t *testing.T) {
	t.Run("first non-null", testCase{
		inFn: func() *model.Value {
			return model.NewValue(orderedmap.NewMap().Set("a", "hello"))
		},
		s:   `a ?? "default"`,
		out: model.NewStringValue("hello"),
	}.run)

	t.Run("fallback on missing", testCase{
		inFn: func() *model.Value {
			return model.NewValue(orderedmap.NewMap().Set("a", "hello"))
		},
		s:   `b ?? "default"`,
		out: model.NewStringValue("default"),
	}.run)

	t.Run("chained coalesce", testCase{
		inFn: func() *model.Value {
			return model.NewValue(orderedmap.NewMap().Set("c", "found"))
		},
		s:   `a ?? b ?? c ?? "none"`,
		out: model.NewStringValue("found"),
	}.run)
}
