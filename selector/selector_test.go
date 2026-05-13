package selector_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/selector"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func TestParse_SimpleProperty(t *testing.T) {
	expr, err := selector.Parse("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ast.IsType[ast.PropertyExpr](expr) {
		t.Fatalf("expected PropertyExpr, got %T", expr)
	}
}

func TestParse_ChainedSelector(t *testing.T) {
	expr, err := selector.Parse("foo.bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ast.IsType[ast.ChainedExpr](expr) {
		t.Fatalf("expected ChainedExpr, got %T", expr)
	}
}

func TestParse_ParserError(t *testing.T) {
	_, err := selector.Parse("if")
	if err == nil {
		t.Fatal("expected error for invalid selector")
	}
}

func TestParse_LexerError(t *testing.T) {
	_, err := selector.Parse(`"hello`)
	if err == nil {
		t.Fatal("expected error for unterminated string")
	}
}
