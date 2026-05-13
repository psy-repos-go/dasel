package execution_test

import (
	"context"
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func TestReduce(t *testing.T) {
	t.Run("sum ints", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewIntValue(1))
			_ = res.Append(model.NewIntValue(2))
			_ = res.Append(model.NewIntValue(3))
			_ = res.Append(model.NewIntValue(4))
			return res
		},
		s: `reduce($this, 0, $acc + $this)`,
		out: model.NewIntValue(10),
	}.run)

	t.Run("concatenate strings", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewStringValue("a"))
			_ = res.Append(model.NewStringValue("b"))
			_ = res.Append(model.NewStringValue("c"))
			return res
		},
		s:   `reduce($this, "", $acc + $this)`,
		out: model.NewStringValue("abc"),
	}.run)

	t.Run("reduce by field", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			p1 := model.NewMapValue()
			_ = p1.SetMapKey("name", model.NewStringValue("Alice"))
			_ = p1.SetMapKey("age", model.NewIntValue(30))
			_ = res.Append(p1)
			p2 := model.NewMapValue()
			_ = p2.SetMapKey("name", model.NewStringValue("Bob"))
			_ = p2.SetMapKey("age", model.NewIntValue(25))
			_ = res.Append(p2)
			return res
		},
		s:   `reduce(age, 0, $acc + $this)`,
		out: model.NewIntValue(55),
	}.run)

	t.Run("empty array returns init", testCase{
		inFn: func() *model.Value {
			return model.NewSliceValue()
		},
		s:   `reduce($this, 0, $acc + $this)`,
		out: model.NewIntValue(0),
	}.run)

	t.Run("non-slice input errors", func(t *testing.T) {
		in := model.NewStringValue("not a slice")
		_, err := execution.ExecuteSelector(context.Background(), `reduce($this, 0, $acc + $this)`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for non-slice input")
		}
		if !strings.Contains(err.Error(), "cannot reduce over non-array") {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("chaining after reduce", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewIntValue(1))
			_ = res.Append(model.NewIntValue(2))
			_ = res.Append(model.NewIntValue(3))
			return res
		},
		s:   `reduce($this, 0, $acc + $this) > 5`,
		out: model.NewBoolValue(true),
	}.run)

	t.Run("product of ints", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewIntValue(2))
			_ = res.Append(model.NewIntValue(3))
			_ = res.Append(model.NewIntValue(4))
			return res
		},
		s:   `reduce($this, 1, $acc * $this)`,
		out: model.NewIntValue(24),
	}.run)

	t.Run("single element", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewIntValue(42))
			return res
		},
		s:   `reduce($this, 0, $acc + $this)`,
		out: model.NewIntValue(42),
	}.run)
}
