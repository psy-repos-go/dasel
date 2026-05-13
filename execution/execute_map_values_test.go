package execution_test

import (
	"context"
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func TestMapValues(t *testing.T) {
	t.Run("multiply int values", testCase{
		inFn: func() *model.Value {
			res := model.NewMapValue()
			_ = res.SetMapKey("a", model.NewIntValue(1))
			_ = res.SetMapKey("b", model.NewIntValue(2))
			_ = res.SetMapKey("c", model.NewIntValue(3))
			return res
		},
		s: `mapValues($this * 2)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()
			_ = res.SetMapKey("a", model.NewIntValue(2))
			_ = res.SetMapKey("b", model.NewIntValue(4))
			_ = res.SetMapKey("c", model.NewIntValue(6))
			return res
		},
	}.run)

	t.Run("add to int values", testCase{
		inFn: func() *model.Value {
			res := model.NewMapValue()
			_ = res.SetMapKey("x", model.NewIntValue(10))
			_ = res.SetMapKey("y", model.NewIntValue(20))
			return res
		},
		s: `mapValues($this + 1)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()
			_ = res.SetMapKey("x", model.NewIntValue(11))
			_ = res.SetMapKey("y", model.NewIntValue(21))
			return res
		},
	}.run)

	t.Run("boolean expression", testCase{
		inFn: func() *model.Value {
			res := model.NewMapValue()
			_ = res.SetMapKey("a", model.NewIntValue(3))
			_ = res.SetMapKey("b", model.NewIntValue(8))
			return res
		},
		s: `mapValues($this > 5)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()
			_ = res.SetMapKey("a", model.NewBoolValue(false))
			_ = res.SetMapKey("b", model.NewBoolValue(true))
			return res
		},
	}.run)

	t.Run("empty map", testCase{
		inFn: func() *model.Value {
			return model.NewMapValue()
		},
		s: `mapValues($this * 2)`,
		outFn: func() *model.Value {
			return model.NewMapValue()
		},
	}.run)

	t.Run("non-map input errors", func(t *testing.T) {
		in := model.NewStringValue("not a map")
		_, err := execution.ExecuteSelector(context.Background(), `mapValues($this * 2)`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for non-map input")
		}
		if !strings.Contains(err.Error(), "cannot mapValues over non-map") {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("mapValues with $key", testCase{
		inFn: func() *model.Value {
			res := model.NewMapValue()
			_ = res.SetMapKey("a", model.NewIntValue(1))
			_ = res.SetMapKey("b", model.NewIntValue(2))
			return res
		},
		s: `mapValues($key)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()
			_ = res.SetMapKey("a", model.NewStringValue("a"))
			_ = res.SetMapKey("b", model.NewStringValue("b"))
			return res
		},
	}.run)

	t.Run("chaining after mapValues", testCase{
		inFn: func() *model.Value {
			res := model.NewMapValue()
			_ = res.SetMapKey("a", model.NewIntValue(1))
			_ = res.SetMapKey("b", model.NewIntValue(2))
			return res
		},
		s:   `mapValues($this + 10).a`,
		out: model.NewIntValue(11),
	}.run)
}
