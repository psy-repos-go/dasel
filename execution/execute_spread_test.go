package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestSpread(t *testing.T) {
	t.Run("build new array", testCase{
		s: "[[1,2,3]..., 4]",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			if err := s.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if err := s.Append(model.NewIntValue(2)); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if err := s.Append(model.NewIntValue(3)); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if err := s.Append(model.NewIntValue(4)); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			return s
		},
	}.run)

	t.Run("map spread into object", testCase{
		inFn: func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("a", "one").
					Set("b", "two"),
			)
		},
		s: `{..., "extra": "three"}`,
		outFn: func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("a", "one").
					Set("b", "two").
					Set("extra", "three"),
			)
		},
	}.run)

	t.Run("spread map values", testCase{
		inFn: func() *model.Value {
			return model.NewValue(
				orderedmap.NewMap().
					Set("x", "one").
					Set("y", "two"),
			)
		},
		s: `$this...`,
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			s.MarkAsSpread()
			_ = s.Append(model.NewStringValue("one"))
			_ = s.Append(model.NewStringValue("two"))
			return s
		},
	}.run)
}
