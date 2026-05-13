package model_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestValue_IsNull(t *testing.T) {
	v := model.NewNullValue()
	if !v.IsNull() {
		t.Fatalf("expected value to be null")
	}
}

func TestStringIndexRange(t *testing.T) {
	t.Run("positive range", func(t *testing.T) {
		v := model.NewStringValue("hello")
		got, err := v.StringIndexRange(1, 3)
		if err != nil {
			t.Fatal(err)
		}
		s, _ := got.StringValue()
		if s != "ell" {
			t.Errorf("expected 'ell', got %q", s)
		}
	})
	t.Run("negative indices", func(t *testing.T) {
		v := model.NewStringValue("hello")
		got, err := v.StringIndexRange(-3, -1)
		if err != nil {
			t.Fatal(err)
		}
		s, _ := got.StringValue()
		if s != "llo" {
			t.Errorf("expected 'llo', got %q", s)
		}
	})
	t.Run("reverse range", func(t *testing.T) {
		v := model.NewStringValue("hello")
		got, err := v.StringIndexRange(3, 1)
		if err != nil {
			t.Fatal(err)
		}
		s, _ := got.StringValue()
		if s != "lle" {
			t.Errorf("expected 'lle', got %q", s)
		}
	})
}
