package execution_test

import (
	"context"
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func TestErrors(t *testing.T) {
	t.Run("unknown function", func(t *testing.T) {
		in := model.NewNullValue()
		_, err := execution.ExecuteSelector(context.Background(), `unknownFunc123()`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for unknown function")
		}
		if !strings.Contains(err.Error(), "unknown function") {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("empty selector returns input", func(t *testing.T) {
		in := model.NewStringValue("hello")
		res, err := execution.ExecuteSelector(context.Background(), ``, in, execution.NewOptions())
		if err != nil {
			t.Fatal(err)
		}
		got, err := res.StringValue()
		if err != nil {
			t.Fatal(err)
		}
		if got != "hello" {
			t.Errorf("expected 'hello', got %s", got)
		}
	})

	t.Run("filter on non-array", func(t *testing.T) {
		in := model.NewStringValue("not a slice")
		_, err := execution.ExecuteSelector(context.Background(), `filter($this == "x")`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for filter on non-array")
		}
	})

	t.Run("each on non-array", func(t *testing.T) {
		in := model.NewStringValue("not a slice")
		_, err := execution.ExecuteSelector(context.Background(), `each($this = $this)`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for each on non-array")
		}
	})

	t.Run("sortBy on non-array", func(t *testing.T) {
		in := model.NewStringValue("not a slice")
		_, err := execution.ExecuteSelector(context.Background(), `sortBy($this)`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for sortBy on non-array")
		}
	})

	t.Run("map on non-array", func(t *testing.T) {
		in := model.NewStringValue("not a slice")
		_, err := execution.ExecuteSelector(context.Background(), `map($this)`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for map on non-array")
		}
	})

	t.Run("spread on non-spreadable", func(t *testing.T) {
		in := model.NewIntValue(42)
		_, err := execution.ExecuteSelector(context.Background(), `$this...`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for spread on int")
		}
	})

	t.Run("all on non-array", func(t *testing.T) {
		in := model.NewStringValue("not a slice")
		_, err := execution.ExecuteSelector(context.Background(), `all($this == "x")`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for all on non-array")
		}
	})

	t.Run("any on non-array", func(t *testing.T) {
		in := model.NewStringValue("not a slice")
		_, err := execution.ExecuteSelector(context.Background(), `any($this == "x")`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for any on non-array")
		}
	})

	t.Run("count on non-array", func(t *testing.T) {
		in := model.NewStringValue("not a slice")
		_, err := execution.ExecuteSelector(context.Background(), `count($this == "x")`, in, execution.NewOptions())
		if err == nil {
			t.Fatal("expected error for count on non-array")
		}
	})
}
