package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestSearch(t *testing.T) {
	t.Run("search has key in nested map", testCase{
		inFn: func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("a", orderedmap.NewMap().
						Set("x", 1)).
					Set("b", orderedmap.NewMap().
						Set("y", 2)),
			)
		},
		s: `search(has("x"))`,
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewValue(orderedmap.NewMap().Set("x", 1)))
			return s
		},
	}.run)

	t.Run("search on array of objects", testCase{
		inFn: func() *model.Value {
			return model.NewValue([]any{
				orderedmap.NewMap().Set("name", "alice"),
				orderedmap.NewMap().Set("name", "bob"),
				orderedmap.NewMap().Set("name", "alice"),
			})
		},
		s: `search(name == "alice")`,
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewValue(orderedmap.NewMap().Set("name", "alice")))
			_ = s.Append(model.NewValue(orderedmap.NewMap().Set("name", "alice")))
			return s
		},
	}.run)

	t.Run("no matches returns empty", testCase{
		inFn: func() *model.Value {
			return model.NewValue([]any{
				orderedmap.NewMap().Set("name", "alice"),
				orderedmap.NewMap().Set("name", "bob"),
			})
		},
		s: `search(name == "charlie")`,
		outFn: func() *model.Value {
			return model.NewSliceValue()
		},
	}.run)
}
