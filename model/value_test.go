package model_test

import (
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestType_String(t *testing.T) {
	run := func(ty model.Type, exp string) func(*testing.T) {
		return func(t *testing.T) {
			got := ty.String()
			if got != exp {
				t.Errorf("expected %s, got %s", exp, got)
			}
		}
	}
	t.Run("string", run(model.TypeString, "string"))
	t.Run("int", run(model.TypeInt, "int"))
	t.Run("float", run(model.TypeFloat, "float"))
	t.Run("bool", run(model.TypeBool, "bool"))
	t.Run("map", run(model.TypeMap, "map"))
	t.Run("slice", run(model.TypeSlice, "array"))
	t.Run("unknown", run(model.TypeUnknown, "unknown"))
	t.Run("null", run(model.TypeNull, "null"))
}

func TestValue_Len(t *testing.T) {
	run := func(v *model.Value, exp int) func(*testing.T) {
		return func(t *testing.T) {
			got, err := v.Len()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if got != exp {
				t.Errorf("expected %d, got %d", exp, got)
			}
		}
	}
	t.Run("string", func(t *testing.T) {
		t.Run("empty", run(model.NewStringValue(""), 0))
		t.Run("non-empty", run(model.NewStringValue("hello"), 5))
	})
	t.Run("slice", func(t *testing.T) {
		t.Run("empty", run(model.NewSliceValue(), 0))
		t.Run("non-empty", run(model.NewValue([]any{1, 2, 3}), 3))
	})
	t.Run("map", func(t *testing.T) {
		t.Run("empty", run(model.NewMapValue(), 0))
		t.Run("non-empty", run(model.NewValue(map[string]any{"one": 1, "two": 2, "three": 3}), 3))
	})
}

func TestValue_IsScalar(t *testing.T) {
	run := func(v *model.Value, exp bool) func(*testing.T) {
		return func(t *testing.T) {
			got := v.IsScalar()
			if got != exp {
				t.Errorf("expected %v, got %v", exp, got)
			}
		}
	}
	t.Run("string", run(model.NewStringValue("foo"), true))
	t.Run("bool", run(model.NewBoolValue(true), true))
	t.Run("int", run(model.NewIntValue(1), true))
	t.Run("float", run(model.NewFloatValue(1.0), true))
	t.Run("null", run(model.NewNullValue(), true))
	t.Run("map", run(model.NewMapValue(), false))
	t.Run("slice", run(model.NewSliceValue(), false))

	t.Run("nested", func(t *testing.T) {
		t.Run("nested string", run(model.NewNestedValue(model.NewStringValue("foo")), true))
		t.Run("nested bool", run(model.NewNestedValue(model.NewBoolValue(true)), true))
		t.Run("nested int", run(model.NewNestedValue(model.NewIntValue(1)), true))
		t.Run("nested float", run(model.NewNestedValue(model.NewFloatValue(1.0)), true))
		t.Run("nested null", run(model.NewNestedValue(model.NewNullValue()), true))
		t.Run("nested map", run(model.NewNestedValue(model.NewMapValue()), false))
		t.Run("nested slice", run(model.NewNestedValue(model.NewSliceValue()), false))

		t.Run("double nested string", run(model.NewNestedValue(model.NewNestedValue(model.NewStringValue("foo"))), true))
	})
}

func TestErrorTypes(t *testing.T) {
	t.Run("MapKeyNotFound", func(t *testing.T) {
		err := model.MapKeyNotFound{Key: "foo"}
		if !strings.Contains(err.Error(), "foo") {
			t.Errorf("expected error to contain 'foo', got %s", err.Error())
		}
	})
	t.Run("SliceIndexOutOfRange", func(t *testing.T) {
		err := model.SliceIndexOutOfRange{Index: 5}
		if !strings.Contains(err.Error(), "5") {
			t.Errorf("expected error to contain '5', got %s", err.Error())
		}
	})
	t.Run("ErrIncompatibleTypes", func(t *testing.T) {
		err := model.ErrIncompatibleTypes{
			A: model.NewStringValue("hello"),
			B: model.NewIntValue(42),
		}
		msg := err.Error()
		if !strings.Contains(msg, "string") || !strings.Contains(msg, "int") {
			t.Errorf("expected error to contain both types, got %s", msg)
		}
	})
	t.Run("ErrUnexpectedType", func(t *testing.T) {
		err := model.ErrUnexpectedType{
			Expected: model.TypeString,
			Actual:   model.TypeInt,
		}
		msg := err.Error()
		if !strings.Contains(msg, "string") || !strings.Contains(msg, "int") {
			t.Errorf("expected error to contain expected and actual types, got %s", msg)
		}
	})
	t.Run("ErrUnexpectedTypes", func(t *testing.T) {
		err := model.ErrUnexpectedTypes{
			Expected: []model.Type{model.TypeSlice, model.TypeMap},
			Actual:   model.TypeInt,
		}
		msg := err.Error()
		if !strings.Contains(msg, "array") || !strings.Contains(msg, "map") || !strings.Contains(msg, "int") {
			t.Errorf("expected error to contain expected list and actual, got %s", msg)
		}
	})
}

func TestErrCouldNotUnpackToType(t *testing.T) {
	err := model.ErrCouldNotUnpackToType{}
	got := err.Error()
	if got == "" {
		t.Error("expected non-empty error message")
	}
}

func TestValue_Len_Error(t *testing.T) {
	_, err := model.NewIntValue(42).Len()
	if err == nil {
		t.Fatal("expected error for Len on int")
	}
}

func TestValue_Copy(t *testing.T) {
	t.Run("map copy is independent", func(t *testing.T) {
		orig := model.NewValue(orderedmap.NewMap().Set("a", "one"))
		cp, err := orig.Copy()
		if err != nil {
			t.Fatal(err)
		}
		// Modify original
		if err := orig.SetMapKey("a", model.NewStringValue("modified")); err != nil {
			t.Fatal(err)
		}
		// Copy should be unchanged
		val, err := cp.GetMapKey("a")
		if err != nil {
			t.Fatal(err)
		}
		got, err := val.StringValue()
		if err != nil {
			t.Fatal(err)
		}
		if got != "one" {
			t.Errorf("expected copy to be 'one', got %s", got)
		}
	})
	t.Run("non-map returns error", func(t *testing.T) {
		_, err := model.NewIntValue(1).Copy()
		if err == nil {
			t.Fatal("expected error for Copy on non-map")
		}
	})
}

func TestValues_ToSliceValue(t *testing.T) {
	vals := model.Values{model.NewIntValue(1), model.NewStringValue("hello")}
	sv, err := vals.ToSliceValue()
	if err != nil {
		t.Fatal(err)
	}
	l, err := sv.Len()
	if err != nil {
		t.Fatal(err)
	}
	if l != 2 {
		t.Errorf("expected length 2, got %d", l)
	}
}

func TestValue_Interface(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := model.NewStringValue("hello")
		got := v.Interface()
		if got == nil {
			t.Fatal("expected non-nil interface")
		}
	})
	t.Run("int", func(t *testing.T) {
		v := model.NewIntValue(42)
		got := v.Interface()
		if got == nil {
			t.Fatal("expected non-nil interface")
		}
	})
	t.Run("null", func(t *testing.T) {
		v := model.NewNullValue()
		got := v.Interface()
		if got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})
}

func TestValue_EqualTypeValue(t *testing.T) {
	t.Run("same strings", func(t *testing.T) {
		a := model.NewStringValue("hello")
		b := model.NewStringValue("hello")
		eq, err := a.EqualTypeValue(b)
		if err != nil {
			t.Fatal(err)
		}
		if !eq {
			t.Error("expected equal")
		}
	})
	t.Run("different strings", func(t *testing.T) {
		a := model.NewStringValue("hello")
		b := model.NewStringValue("world")
		eq, err := a.EqualTypeValue(b)
		if err != nil {
			t.Fatal(err)
		}
		if eq {
			t.Error("expected not equal")
		}
	})
	t.Run("different types", func(t *testing.T) {
		a := model.NewStringValue("1")
		b := model.NewIntValue(1)
		eq, err := a.EqualTypeValue(b)
		if err != nil {
			t.Fatal(err)
		}
		if eq {
			t.Error("expected not equal for different types")
		}
	})
	t.Run("same ints", func(t *testing.T) {
		eq, err := model.NewIntValue(42).EqualTypeValue(model.NewIntValue(42))
		if err != nil {
			t.Fatal(err)
		}
		if !eq {
			t.Error("expected equal")
		}
	})
	t.Run("different ints", func(t *testing.T) {
		eq, err := model.NewIntValue(1).EqualTypeValue(model.NewIntValue(2))
		if err != nil {
			t.Fatal(err)
		}
		if eq {
			t.Error("expected not equal")
		}
	})
	t.Run("same floats", func(t *testing.T) {
		eq, err := model.NewFloatValue(3.14).EqualTypeValue(model.NewFloatValue(3.14))
		if err != nil {
			t.Fatal(err)
		}
		if !eq {
			t.Error("expected equal")
		}
	})
	t.Run("same bools", func(t *testing.T) {
		eq, err := model.NewBoolValue(true).EqualTypeValue(model.NewBoolValue(true))
		if err != nil {
			t.Fatal(err)
		}
		if !eq {
			t.Error("expected equal")
		}
	})
	t.Run("both null", func(t *testing.T) {
		eq, err := model.NewNullValue().EqualTypeValue(model.NewNullValue())
		if err != nil {
			t.Fatal(err)
		}
		if !eq {
			t.Error("expected equal")
		}
	})
	t.Run("same slices", func(t *testing.T) {
		a := model.NewSliceValue()
		_ = a.Append(model.NewIntValue(1))
		b := model.NewSliceValue()
		_ = b.Append(model.NewIntValue(1))
		eq, err := a.EqualTypeValue(b)
		if err != nil {
			t.Fatal(err)
		}
		if !eq {
			t.Error("expected equal")
		}
	})
	t.Run("different length slices", func(t *testing.T) {
		a := model.NewSliceValue()
		_ = a.Append(model.NewIntValue(1))
		b := model.NewSliceValue()
		eq, err := a.EqualTypeValue(b)
		if err != nil {
			t.Fatal(err)
		}
		if eq {
			t.Error("expected not equal")
		}
	})
	t.Run("same maps", func(t *testing.T) {
		a := model.NewMapValue()
		_ = a.SetMapKey("k", model.NewStringValue("v"))
		b := model.NewMapValue()
		_ = b.SetMapKey("k", model.NewStringValue("v"))
		eq, err := a.EqualTypeValue(b)
		if err != nil {
			t.Fatal(err)
		}
		if !eq {
			t.Error("expected equal")
		}
	})
	t.Run("different maps", func(t *testing.T) {
		a := model.NewMapValue()
		_ = a.SetMapKey("k", model.NewStringValue("v"))
		b := model.NewMapValue()
		_ = b.SetMapKey("k", model.NewStringValue("w"))
		eq, err := a.EqualTypeValue(b)
		if err != nil {
			t.Fatal(err)
		}
		if eq {
			t.Error("expected not equal")
		}
	})
}

func TestValue_GoValue_AllTypes(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		v := model.NewStringValue("hello")
		got, err := v.GoValue()
		if err != nil {
			t.Fatal(err)
		}
		if got != "hello" {
			t.Errorf("expected 'hello', got %v", got)
		}
	})
	t.Run("int", func(t *testing.T) {
		v := model.NewIntValue(42)
		got, err := v.GoValue()
		if err != nil {
			t.Fatal(err)
		}
		if got != int64(42) {
			t.Errorf("expected 42, got %v", got)
		}
	})
	t.Run("float", func(t *testing.T) {
		v := model.NewFloatValue(3.14)
		got, err := v.GoValue()
		if err != nil {
			t.Fatal(err)
		}
		if got != 3.14 {
			t.Errorf("expected 3.14, got %v", got)
		}
	})
	t.Run("bool", func(t *testing.T) {
		v := model.NewBoolValue(true)
		got, err := v.GoValue()
		if err != nil {
			t.Fatal(err)
		}
		if got != true {
			t.Errorf("expected true, got %v", got)
		}
	})
	t.Run("null", func(t *testing.T) {
		v := model.NewNullValue()
		got, err := v.GoValue()
		if err != nil {
			t.Fatal(err)
		}
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("map", func(t *testing.T) {
		v := model.NewMapValue()
		_ = v.SetMapKey("key", model.NewStringValue("val"))
		got, err := v.GoValue()
		if err != nil {
			t.Fatal(err)
		}
		m, ok := got.(map[string]any)
		if !ok {
			t.Fatalf("expected map, got %T", got)
		}
		if m["key"] != "val" {
			t.Errorf("expected 'val', got %v", m["key"])
		}
	})
	t.Run("slice", func(t *testing.T) {
		v := model.NewSliceValue()
		_ = v.Append(model.NewIntValue(1))
		_ = v.Append(model.NewIntValue(2))
		got, err := v.GoValue()
		if err != nil {
			t.Fatal(err)
		}
		s, ok := got.([]any)
		if !ok {
			t.Fatalf("expected slice, got %T", got)
		}
		if len(s) != 2 {
			t.Errorf("expected length 2, got %d", len(s))
		}
	})
}

func TestValue_Set_SetFn(t *testing.T) {
	// When we get a map key, the returned value has a setFn that allows setting it back.
	m := model.NewValue(orderedmap.NewMap().Set("key", "original"))
	val, err := m.GetMapKey("key")
	if err != nil {
		t.Fatal(err)
	}
	if err := val.Set(model.NewStringValue("modified")); err != nil {
		t.Fatal(err)
	}
	// Verify the map was updated
	updated, err := m.GetMapKey("key")
	if err != nil {
		t.Fatal(err)
	}
	got, err := updated.StringValue()
	if err != nil {
		t.Fatal(err)
	}
	if got != "modified" {
		t.Errorf("expected 'modified', got %s", got)
	}
}

func TestValue_Set_Nested(t *testing.T) {
	// Test setting a value via a nested dasel value (exercises the isDaselValue branch in Set)
	inner := model.NewStringValue("original")
	outer := model.NewNestedValue(inner)
	if err := outer.Set(model.NewStringValue("updated")); err != nil {
		t.Fatal(err)
	}
}

func TestValue_Kind(t *testing.T) {
	// Exercise the Kind method
	v := model.NewStringValue("hello")
	k := v.Kind()
	if k.String() == "" {
		t.Error("expected non-empty kind")
	}
}

func TestValue_Type_AllTypes(t *testing.T) {
	tests := []struct {
		name string
		val  *model.Value
		typ  model.Type
	}{
		{"string", model.NewStringValue("hello"), model.TypeString},
		{"int", model.NewIntValue(1), model.TypeInt},
		{"float", model.NewFloatValue(1.0), model.TypeFloat},
		{"bool", model.NewBoolValue(true), model.TypeBool},
		{"map", model.NewMapValue(), model.TypeMap},
		{"slice", model.NewSliceValue(), model.TypeSlice},
		{"null", model.NewNullValue(), model.TypeNull},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.Type(); got != tt.typ {
				t.Errorf("expected %s, got %s", tt.typ, got)
			}
		})
	}
}

func TestValue_MapLen(t *testing.T) {
	m := model.NewValue(orderedmap.NewMap().Set("a", 1).Set("b", 2).Set("c", 3))
	l, err := m.MapLen()
	if err != nil {
		t.Fatal(err)
	}
	if l != 3 {
		t.Errorf("expected 3, got %d", l)
	}
}

func TestValue_SliceLen_Error(t *testing.T) {
	_, err := model.NewIntValue(42).SliceLen()
	if err == nil {
		t.Fatal("expected error for SliceLen on int")
	}
}

func TestValue_SetSliceIndex_Error(t *testing.T) {
	s := model.NewSliceValue()
	_ = s.Append(model.NewIntValue(1))
	err := s.SetSliceIndex(5, model.NewIntValue(2))
	if err == nil {
		t.Fatal("expected error for out of range SetSliceIndex")
	}
}

func TestValue_String(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		got := model.NewStringValue("hello").String()
		if !strings.Contains(got, "hello") {
			t.Errorf("expected output to contain 'hello', got %s", got)
		}
	})
	t.Run("int", func(t *testing.T) {
		got := model.NewIntValue(42).String()
		if !strings.Contains(got, "42") {
			t.Errorf("expected output to contain '42', got %s", got)
		}
	})
	t.Run("float", func(t *testing.T) {
		got := model.NewFloatValue(3.14).String()
		if !strings.Contains(got, "3.14") {
			t.Errorf("expected output to contain '3.14', got %s", got)
		}
	})
	t.Run("bool", func(t *testing.T) {
		got := model.NewBoolValue(true).String()
		if !strings.Contains(got, "true") {
			t.Errorf("expected output to contain 'true', got %s", got)
		}
	})
	t.Run("null", func(t *testing.T) {
		got := model.NewNullValue().String()
		if !strings.Contains(got, "null") {
			t.Errorf("expected output to contain 'null', got %s", got)
		}
	})
	t.Run("slice", func(t *testing.T) {
		s := model.NewSliceValue()
		_ = s.Append(model.NewIntValue(1))
		got := s.String()
		if !strings.Contains(got, "array") {
			t.Errorf("expected output to contain 'array', got %s", got)
		}
	})
	t.Run("map", func(t *testing.T) {
		m := model.NewValue(orderedmap.NewMap().Set("key", "val"))
		got := m.String()
		if !strings.Contains(got, "key") {
			t.Errorf("expected output to contain 'key', got %s", got)
		}
	})
}
