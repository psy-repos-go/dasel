package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestRecursiveDescent(t *testing.T) {
	t.Run("find key at any depth", testCase{
		inFn: func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("a", orderedmap.NewMap().
						Set("x", "found1").
						Set("b", orderedmap.NewMap().
							Set("x", "found2"))),
			)
		},
		s: `..x`,
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewStringValue("found1"))
			_ = s.Append(model.NewStringValue("found2"))
			return s
		},
	}.run)

	t.Run("wildcard collects all scalars", testCase{
		inFn: func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("a", "hello").
					Set("b", orderedmap.NewMap().
						Set("c", "world")),
			)
		},
		s: `..*`,
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewStringValue("hello"))
			_ = s.Append(model.NewStringValue("world"))
			return s
		},
	}.run)

	t.Run("find key in nested slices containing maps", testCase{
		inFn: func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("items", []any{
						orderedmap.NewMap().Set("x", "first"),
						orderedmap.NewMap().Set("x", "second"),
					}),
			)
		},
		s: `..x`,
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewStringValue("first"))
			_ = s.Append(model.NewStringValue("second"))
			return s
		},
	}.run)

	t.Run("wildcard with $key available", testCase{
		inFn: func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("a", "hello").
					Set("b", "world"),
			)
		},
		s: `..*`,
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewStringValue("hello"))
			_ = s.Append(model.NewStringValue("world"))
			return s
		},
	}.run)

	t.Run("key not found returns empty", testCase{
		inFn: func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("a", orderedmap.NewMap().
						Set("b", "value")),
			)
		},
		s: `..z`,
		outFn: func() *model.Value {
			return model.NewSliceValue()
		},
	}.run)
}
