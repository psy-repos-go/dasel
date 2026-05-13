package execution_test

import (
	"context"
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func TestGroupBy(t *testing.T) {
	t.Run("group by $this strings", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewStringValue("a"))
			_ = res.Append(model.NewStringValue("b"))
			_ = res.Append(model.NewStringValue("a"))
			_ = res.Append(model.NewStringValue("c"))
			_ = res.Append(model.NewStringValue("b"))
			return res
		},
		s: `groupBy($this)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()

			groupA := model.NewSliceValue()
			_ = groupA.Append(model.NewStringValue("a"))
			_ = groupA.Append(model.NewStringValue("a"))
			_ = res.SetMapKey("a", groupA)

			groupB := model.NewSliceValue()
			_ = groupB.Append(model.NewStringValue("b"))
			_ = groupB.Append(model.NewStringValue("b"))
			_ = res.SetMapKey("b", groupB)

			groupC := model.NewSliceValue()
			_ = groupC.Append(model.NewStringValue("c"))
			_ = res.SetMapKey("c", groupC)

			return res
		},
	}.run)

	t.Run("group by $this ints", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewIntValue(1))
			_ = res.Append(model.NewIntValue(2))
			_ = res.Append(model.NewIntValue(1))
			_ = res.Append(model.NewIntValue(2))
			_ = res.Append(model.NewIntValue(3))
			return res
		},
		s: `groupBy($this)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()

			group1 := model.NewSliceValue()
			_ = group1.Append(model.NewIntValue(1))
			_ = group1.Append(model.NewIntValue(1))
			_ = res.SetMapKey("1", group1)

			group2 := model.NewSliceValue()
			_ = group2.Append(model.NewIntValue(2))
			_ = group2.Append(model.NewIntValue(2))
			_ = res.SetMapKey("2", group2)

			group3 := model.NewSliceValue()
			_ = group3.Append(model.NewIntValue(3))
			_ = res.SetMapKey("3", group3)

			return res
		},
	}.run)

	t.Run("group by field name", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()

			alice := model.NewMapValue()
			_ = alice.SetMapKey("name", model.NewStringValue("Alice"))
			_ = alice.SetMapKey("dept", model.NewStringValue("eng"))
			_ = res.Append(alice)

			bob := model.NewMapValue()
			_ = bob.SetMapKey("name", model.NewStringValue("Bob"))
			_ = bob.SetMapKey("dept", model.NewStringValue("eng"))
			_ = res.Append(bob)

			carol := model.NewMapValue()
			_ = carol.SetMapKey("name", model.NewStringValue("Carol"))
			_ = carol.SetMapKey("dept", model.NewStringValue("sales"))
			_ = res.Append(carol)

			return res
		},
		s: `groupBy(dept)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()

			alice := model.NewMapValue()
			_ = alice.SetMapKey("name", model.NewStringValue("Alice"))
			_ = alice.SetMapKey("dept", model.NewStringValue("eng"))

			bob := model.NewMapValue()
			_ = bob.SetMapKey("name", model.NewStringValue("Bob"))
			_ = bob.SetMapKey("dept", model.NewStringValue("eng"))

			carol := model.NewMapValue()
			_ = carol.SetMapKey("name", model.NewStringValue("Carol"))
			_ = carol.SetMapKey("dept", model.NewStringValue("sales"))

			engGroup := model.NewSliceValue()
			_ = engGroup.Append(alice)
			_ = engGroup.Append(bob)
			_ = res.SetMapKey("eng", engGroup)

			salesGroup := model.NewSliceValue()
			_ = salesGroup.Append(carol)
			_ = res.SetMapKey("sales", salesGroup)

			return res
		},
	}.run)

	t.Run("group by boolean expression", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewIntValue(3))
			_ = res.Append(model.NewIntValue(7))
			_ = res.Append(model.NewIntValue(1))
			_ = res.Append(model.NewIntValue(8))
			_ = res.Append(model.NewIntValue(5))
			return res
		},
		s: `groupBy($this > 5)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()

			falseGroup := model.NewSliceValue()
			_ = falseGroup.Append(model.NewIntValue(3))
			_ = falseGroup.Append(model.NewIntValue(1))
			_ = falseGroup.Append(model.NewIntValue(5))
			_ = res.SetMapKey("false", falseGroup)

			trueGroup := model.NewSliceValue()
			_ = trueGroup.Append(model.NewIntValue(7))
			_ = trueGroup.Append(model.NewIntValue(8))
			_ = res.SetMapKey("true", trueGroup)

			return res
		},
	}.run)

	t.Run("empty array", testCase{
		inFn: func() *model.Value {
			return model.NewSliceValue()
		},
		s: `groupBy($this)`,
		outFn: func() *model.Value {
			return model.NewMapValue()
		},
	}.run)

	t.Run("group by $this floats", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewFloatValue(1.5))
			_ = res.Append(model.NewFloatValue(2.5))
			_ = res.Append(model.NewFloatValue(1.5))
			return res
		},
		s: `groupBy($this)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()

			group1 := model.NewSliceValue()
			_ = group1.Append(model.NewFloatValue(1.5))
			_ = group1.Append(model.NewFloatValue(1.5))
			_ = res.SetMapKey("1.5", group1)

			group2 := model.NewSliceValue()
			_ = group2.Append(model.NewFloatValue(2.5))
			_ = res.SetMapKey("2.5", group2)

			return res
		},
	}.run)

	t.Run("single element", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewStringValue("only"))
			return res
		},
		s: `groupBy($this)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()
			group := model.NewSliceValue()
			_ = group.Append(model.NewStringValue("only"))
			_ = res.SetMapKey("only", group)
			return res
		},
	}.run)

	t.Run("chaining after groupBy", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()

			alice := model.NewMapValue()
			_ = alice.SetMapKey("name", model.NewStringValue("Alice"))
			_ = alice.SetMapKey("dept", model.NewStringValue("eng"))
			_ = res.Append(alice)

			bob := model.NewMapValue()
			_ = bob.SetMapKey("name", model.NewStringValue("Bob"))
			_ = bob.SetMapKey("dept", model.NewStringValue("eng"))
			_ = res.Append(bob)

			carol := model.NewMapValue()
			_ = carol.SetMapKey("name", model.NewStringValue("Carol"))
			_ = carol.SetMapKey("dept", model.NewStringValue("sales"))
			_ = res.Append(carol)

			return res
		},
		s: `groupBy(dept).eng`,
		outFn: func() *model.Value {
			alice := model.NewMapValue()
			_ = alice.SetMapKey("name", model.NewStringValue("Alice"))
			_ = alice.SetMapKey("dept", model.NewStringValue("eng"))

			bob := model.NewMapValue()
			_ = bob.SetMapKey("name", model.NewStringValue("Bob"))
			_ = bob.SetMapKey("dept", model.NewStringValue("eng"))

			res := model.NewSliceValue()
			_ = res.Append(alice)
			_ = res.Append(bob)
			return res
		},
	}.run)

	t.Run("group by $key", testCase{
		inFn: func() *model.Value {
			res := model.NewSliceValue()
			_ = res.Append(model.NewStringValue("a"))
			_ = res.Append(model.NewStringValue("b"))
			_ = res.Append(model.NewStringValue("c"))
			return res
		},
		s: `groupBy($key)`,
		outFn: func() *model.Value {
			res := model.NewMapValue()
			s0 := model.NewSliceValue()
			_ = s0.Append(model.NewStringValue("a"))
			_ = res.SetMapKey("0", s0)
			s1 := model.NewSliceValue()
			_ = s1.Append(model.NewStringValue("b"))
			_ = res.SetMapKey("1", s1)
			s2 := model.NewSliceValue()
			_ = s2.Append(model.NewStringValue("c"))
			_ = res.SetMapKey("2", s2)
			return res
		},
	}.run)

	t.Run("non-slice input errors", func(t *testing.T) {
		in := model.NewStringValue("not a slice")
		_, err := execution.ExecuteSelector(context.Background(), `groupBy($this)`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for non-slice input")
		}
		if !strings.Contains(err.Error(), "cannot group by on non-slice data") {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}
