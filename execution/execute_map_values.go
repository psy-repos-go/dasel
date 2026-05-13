package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func mapValuesExprExecutor(e ast.MapValuesExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "mapValuesExpr")
		if !data.IsMap() {
			return nil, fmt.Errorf("cannot mapValues over non-map")
		}

		res := model.NewMapValue()

		if err := data.RangeMap(func(key string, value *model.Value) error {
			transformed, err := ExecuteAST(ctx, e.Expr, value, options)
			if err != nil {
				return fmt.Errorf("error evaluating mapValues expr for key %q: %w", key, err)
			}
			return res.SetMapKey(key, transformed)
		}); err != nil {
			return nil, fmt.Errorf("error ranging over map: %w", err)
		}

		return res, nil
	}, nil
}
