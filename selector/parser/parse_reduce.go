package parser

import (
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func parseReduce(p *Parser) (ast.Expr, error) {
	if err := p.expect(lexer.Reduce); err != nil {
		return nil, err
	}
	p.advance()

	if err := p.expect(lexer.OpenParen); err != nil {
		return nil, err
	}
	p.advance()

	// Parse first expression: expr (break on Comma)
	expr, err := p.parseExpressions(
		lexer.TokenKinds(lexer.Comma),
		nil,
		true,
		bpDefault,
		false,
	)
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.Comma); err != nil {
		return nil, err
	}
	p.advance()

	// Parse second expression: init (break on Comma)
	init, err := p.parseExpressions(
		lexer.TokenKinds(lexer.Comma),
		nil,
		true,
		bpDefault,
		false,
	)
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.Comma); err != nil {
		return nil, err
	}
	p.advance()

	// Parse third expression: update (break on CloseParen)
	update, err := p.parseExpressions(
		lexer.TokenKinds(lexer.CloseParen),
		nil,
		true,
		bpDefault,
		false,
	)
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.CloseParen); err != nil {
		return nil, err
	}
	p.advance()

	reduceExpr := ast.ReduceExpr{
		Expr:   expr,
		Init:   init,
		Update: update,
	}

	res, err := parseFollowingSymbol(p, reduceExpr)
	if err != nil {
		return nil, err
	}

	return res, nil
}
