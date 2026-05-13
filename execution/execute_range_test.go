package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestRange(t *testing.T) {
	t.Run("string range", testCase{
		in:  model.NewStringValue("hello"),
		s:   `$this[1:3]`,
		out: model.NewStringValue("ell"),
	}.run)

	t.Run("slice range", testCase{
		s: `[10,20,30,40,50][1:3]`,
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(20))
			_ = s.Append(model.NewIntValue(30))
			_ = s.Append(model.NewIntValue(40))
			return s
		},
	}.run)

	t.Run("slice range start only", testCase{
		s: `[10,20,30,40,50][3:]`,
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(40))
			_ = s.Append(model.NewIntValue(50))
			return s
		},
	}.run)
}
