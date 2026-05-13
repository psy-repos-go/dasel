package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/execution"
)

func TestOptions(t *testing.T) {
	t.Run("WithFuncs", func(t *testing.T) {
		fc := execution.DefaultFuncCollection.Copy()
		opts := execution.NewOptions(execution.WithFuncs(fc))
		if _, ok := opts.Funcs.Get("len"); !ok {
			t.Error("expected funcs to contain 'len'")
		}
	})

	t.Run("WithoutUnstable", func(t *testing.T) {
		opts := execution.NewOptions(execution.WithUnstable(), execution.WithoutUnstable())
		if opts.Unstable {
			t.Error("expected unstable to be false")
		}
	})
}

func TestFuncCollection(t *testing.T) {
	t.Run("Copy", func(t *testing.T) {
		fc := execution.DefaultFuncCollection.Copy()
		if _, ok := fc.Get("len"); !ok {
			t.Error("expected copy to contain 'len'")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		fc := execution.DefaultFuncCollection.Copy()
		fc.Delete("len")
		if _, ok := fc.Get("len"); ok {
			t.Error("expected len to be deleted")
		}
	})
}
