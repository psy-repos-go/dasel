package execution_test

import (
	"context"
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func TestVariable(t *testing.T) {
	t.Run("use passed variable", testCase{
		s:   `$x + 1`,
		out: model.NewIntValue(11),
		opts: []execution.ExecuteOptionFn{
			execution.WithVariable("x", model.NewIntValue(10)),
		},
	}.run)

	t.Run("variable in expression", testCase{
		s:   `$x + $y`,
		out: model.NewIntValue(30),
		opts: []execution.ExecuteOptionFn{
			execution.WithVariable("x", model.NewIntValue(10)),
			execution.WithVariable("y", model.NewIntValue(20)),
		},
	}.run)

	t.Run("undefined variable errors", func(t *testing.T) {
		in := model.NewNullValue()
		_, err := execution.ExecuteSelector(context.Background(), `$undefined_var_xyz`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for undefined variable")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("env variable", func(t *testing.T) {
		t.Setenv("DASEL_TEST_VAR_ABC", "from_env")
		in := model.NewNullValue()
		res, err := execution.ExecuteSelector(context.Background(), `$DASEL_TEST_VAR_ABC`, in, execution.NewOptions())
		if err != nil {
			t.Fatal(err)
		}
		got, err := res.StringValue()
		if err != nil {
			t.Fatal(err)
		}
		if got != "from_env" {
			t.Errorf("expected 'from_env', got %s", got)
		}
	})
}
