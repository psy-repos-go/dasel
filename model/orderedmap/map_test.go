package orderedmap_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestNewMap(t *testing.T) {
	m := orderedmap.NewMap()
	if m.Len() != 0 {
		t.Fatalf("expected Len 0, got %d", m.Len())
	}
	if len(m.Keys()) != 0 {
		t.Fatal("expected empty keys")
	}
}

func TestSetGet(t *testing.T) {
	m := orderedmap.NewMap()
	m.Set("a", 1)

	v, ok := m.Get("a")
	if !ok {
		t.Fatal("expected key a to exist")
	}
	if v != 1 {
		t.Fatalf("expected 1, got %v", v)
	}

	// Missing key.
	_, ok = m.Get("missing")
	if ok {
		t.Fatal("expected key missing to not exist")
	}
}

func TestSetOverwrite(t *testing.T) {
	m := orderedmap.NewMap()
	m.Set("a", 1)
	m.Set("a", 2)

	v, _ := m.Get("a")
	if v != 2 {
		t.Fatalf("expected 2, got %v", v)
	}
	// Overwrite must not duplicate the key.
	if m.Len() != 1 {
		t.Fatalf("expected Len 1, got %d", m.Len())
	}
}

func TestDelete(t *testing.T) {
	m := orderedmap.NewMap()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Delete("a")

	_, ok := m.Get("a")
	if ok {
		t.Fatal("expected key a to be deleted")
	}
	if m.Len() != 1 {
		t.Fatalf("expected Len 1, got %d", m.Len())
	}
	keys := m.Keys()
	if len(keys) != 1 || keys[0] != "b" {
		t.Fatalf("expected keys [b], got %v", keys)
	}
}

func TestDeleteNonExistent(t *testing.T) {
	m := orderedmap.NewMap()
	m.Set("a", 1)
	m.Delete("nope") // no-op
	if m.Len() != 1 {
		t.Fatalf("expected Len 1, got %d", m.Len())
	}
}

func TestLen(t *testing.T) {
	m := orderedmap.NewMap()
	if m.Len() != 0 {
		t.Fatalf("expected 0, got %d", m.Len())
	}
	m.Set("a", 1)
	m.Set("b", 2)
	if m.Len() != 2 {
		t.Fatalf("expected 2, got %d", m.Len())
	}
	m.Delete("a")
	if m.Len() != 1 {
		t.Fatalf("expected 1, got %d", m.Len())
	}
}

func TestKeysOrder(t *testing.T) {
	m := orderedmap.NewMap()
	m.Set("c", 3)
	m.Set("a", 1)
	m.Set("b", 2)

	keys := m.Keys()
	expected := []string{"c", "a", "b"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, keys)
	}
	for i, k := range expected {
		if keys[i] != k {
			t.Fatalf("expected key %d to be %s, got %s", i, k, keys[i])
		}
	}
}

func TestKeyValues(t *testing.T) {
	m := orderedmap.NewMap()
	m.Set("x", 10)
	m.Set("y", 20)

	kvs := m.KeyValues()
	if len(kvs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(kvs))
	}
	if kvs[0].Key != "x" || kvs[0].Value != 10 {
		t.Fatalf("unexpected first pair: %v", kvs[0])
	}
	if kvs[1].Key != "y" || kvs[1].Value != 20 {
		t.Fatalf("unexpected second pair: %v", kvs[1])
	}
}

func TestEqual(t *testing.T) {
	a := orderedmap.NewMap()
	a.Set("x", 1)
	a.Set("y", 2)

	b := orderedmap.NewMap()
	b.Set("x", 1)
	b.Set("y", 2)

	if !a.Equal(b) {
		t.Fatal("expected maps to be equal")
	}

	// Different value.
	c := orderedmap.NewMap()
	c.Set("x", 1)
	c.Set("y", 99)
	if a.Equal(c) {
		t.Fatal("expected maps to not be equal (different value)")
	}

	// Different key order.
	d := orderedmap.NewMap()
	d.Set("y", 2)
	d.Set("x", 1)
	if a.Equal(d) {
		t.Fatal("expected maps to not be equal (different order)")
	}

	// Different length.
	e := orderedmap.NewMap()
	e.Set("x", 1)
	if a.Equal(e) {
		t.Fatal("expected maps to not be equal (different length)")
	}
}

func TestFromMap(t *testing.T) {
	src := map[string]any{"a": 1, "b": 2, "c": 3}
	m := orderedmap.FromMap(src)

	if m.Len() != 3 {
		t.Fatalf("expected Len 3, got %d", m.Len())
	}
	for k, v := range src {
		got, ok := m.Get(k)
		if !ok {
			t.Fatalf("expected key %s to exist", k)
		}
		if got != v {
			t.Fatalf("expected %v for key %s, got %v", v, k, got)
		}
	}
}

func TestUnorderedData(t *testing.T) {
	m := orderedmap.NewMap()
	m.Set("a", 1)
	m.Set("b", 2)

	data := m.UnorderedData()
	if len(data) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(data))
	}
	if data["a"] != 1 || data["b"] != 2 {
		t.Fatalf("unexpected data: %v", data)
	}
}
