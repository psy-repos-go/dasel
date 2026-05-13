package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestEach(t *testing.T) {
	t.Run("all true", testCase{
		s: "[1,2,3].each($this = $this + 1)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
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

	t.Run("multiply each element", testCase{
		s: "[2,3,4].each($this = $this * 2)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(4))
			_ = s.Append(model.NewIntValue(6))
			_ = s.Append(model.NewIntValue(8))
			return s
		},
	}.run)
}
