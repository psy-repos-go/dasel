package execution_test

import (
	"context"
	"testing"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func TestKeyVar(t *testing.T) {
	intSlice := func() *model.Value {
		s := model.NewSliceValue()
		_ = s.Append(model.NewIntValue(10))
		_ = s.Append(model.NewIntValue(20))
		_ = s.Append(model.NewIntValue(30))
		return s
	}

	t.Run("filter with $key", testCase{
		inFn: intSlice,
		s:    "filter($key >= 1)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(20))
			_ = s.Append(model.NewIntValue(30))
			return s
		},
	}.run)

	t.Run("map with $key", testCase{
		inFn: intSlice,
		s:    "map($key)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(0))
			_ = s.Append(model.NewIntValue(1))
			_ = s.Append(model.NewIntValue(2))
			return s
		},
	}.run)

	t.Run("mapValues with $key", testCase{
		inFn: func() *model.Value {
			m := model.NewMapValue()
			_ = m.SetMapKey("a", model.NewIntValue(1))
			_ = m.SetMapKey("b", model.NewIntValue(2))
			return m
		},
		s: "mapValues($key)",
		outFn: func() *model.Value {
			m := model.NewMapValue()
			_ = m.SetMapKey("a", model.NewStringValue("a"))
			_ = m.SetMapKey("b", model.NewStringValue("b"))
			return m
		},
	}.run)

	t.Run("search with $key on map", testCase{
		inFn: func() *model.Value {
			m := model.NewMapValue()
			_ = m.SetMapKey("other", model.NewIntValue(1))
			_ = m.SetMapKey("target", model.NewIntValue(42))
			return m
		},
		s: `search($key == "target")`,
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(42))
			return s
		},
	}.run)

	t.Run("search with $key on slice", testCase{
		inFn: intSlice,
		s:    "search($key == 1)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(20))
			return s
		},
	}.run)

	t.Run("reduce with $key", testCase{
		inFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewStringValue("a"))
			_ = s.Append(model.NewStringValue("b"))
			_ = s.Append(model.NewStringValue("c"))
			return s
		},
		s:   "reduce($this, 0, $acc + $key)",
		out: model.NewIntValue(3), // 0+0+1+2
	}.run)

	t.Run("each with $key", testCase{
		inFn:        intSlice,
		s:           "each($key + $this)",
		compareRoot: true,
		outFn:       intSlice,
	}.run)

	t.Run("any with $key true", testCase{
		inFn: intSlice,
		s:    "any($key == 2)",
		out:  model.NewBoolValue(true),
	}.run)

	t.Run("any with $key false", testCase{
		inFn: intSlice,
		s:    "any($key == 5)",
		out:  model.NewBoolValue(false),
	}.run)

	t.Run("all with $key true", testCase{
		inFn: intSlice,
		s:    "all($key < 3)",
		out:  model.NewBoolValue(true),
	}.run)

	t.Run("all with $key false", testCase{
		inFn: intSlice,
		s:    "all($key < 2)",
		out:  model.NewBoolValue(false),
	}.run)

	t.Run("count with $key", testCase{
		inFn: intSlice,
		s:    "count($key > 0)",
		out:  model.NewIntValue(2),
	}.run)

	t.Run("sortBy with $key descending", testCase{
		inFn: intSlice,
		s:    "sortBy($key, desc)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(30))
			_ = s.Append(model.NewIntValue(20))
			_ = s.Append(model.NewIntValue(10))
			return s
		},
	}.run)

	t.Run("groupBy with $key", testCase{
		inFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewStringValue("a"))
			_ = s.Append(model.NewStringValue("b"))
			_ = s.Append(model.NewStringValue("c"))
			return s
		},
		s: "groupBy($key)",
		outFn: func() *model.Value {
			m := model.NewMapValue()
			s0 := model.NewSliceValue()
			_ = s0.Append(model.NewStringValue("a"))
			_ = m.SetMapKey("0", s0)
			s1 := model.NewSliceValue()
			_ = s1.Append(model.NewStringValue("b"))
			_ = m.SetMapKey("1", s1)
			s2 := model.NewSliceValue()
			_ = s2.Append(model.NewStringValue("c"))
			_ = m.SetMapKey("2", s2)
			return m
		},
	}.run)

	t.Run("$key does not leak after map", func(t *testing.T) {
		in := model.NewSliceValue()
		_ = in.Append(model.NewIntValue(1))
		_ = in.Append(model.NewIntValue(2))

		opts := execution.NewOptions()
		// $key should not exist before
		if _, exists := opts.Vars["key"]; exists {
			t.Fatal("$key should not exist before iteration")
		}

		_, err := execution.ExecuteSelector(context.Background(), "map($key)", in, opts)
		if err != nil {
			t.Fatal(err)
		}

		// $key should not exist after
		if _, exists := opts.Vars["key"]; exists {
			t.Fatal("$key should not leak after iteration")
		}
	})

	t.Run("$key scoping restores previous value", func(t *testing.T) {
		in := model.NewSliceValue()
		_ = in.Append(model.NewIntValue(1))

		opts := execution.NewOptions(execution.WithVariable("key", model.NewStringValue("preserved")))

		_, err := execution.ExecuteSelector(context.Background(), "map($key)", in, opts)
		if err != nil {
			t.Fatal(err)
		}

		// $key should be restored to the previous value
		v, exists := opts.Vars["key"]
		if !exists {
			t.Fatal("$key should still exist after iteration")
		}
		s, err := v.StringValue()
		if err != nil {
			t.Fatal(err)
		}
		if s != "preserved" {
			t.Fatalf("expected $key to be 'preserved', got %q", s)
		}
	})

	t.Run("filter with $key and $this", testCase{
		inFn: intSlice,
		s:    "filter($key == 0 || $this == 30)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(10))
			_ = s.Append(model.NewIntValue(30))
			return s
		},
	}.run)

	t.Run("map $key + $this", testCase{
		inFn: intSlice,
		s:    "map($key + $this)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(10))  // 0 + 10
			_ = s.Append(model.NewIntValue(21))  // 1 + 20
			_ = s.Append(model.NewIntValue(32))  // 2 + 30
			return s
		},
	}.run)
}
