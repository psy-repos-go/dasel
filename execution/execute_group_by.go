package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func groupByExprExecutor(e ast.GroupByExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "groupByExpr")
		if !data.IsSlice() {
			return nil, fmt.Errorf("cannot group by on non-slice data")
		}

		res := model.NewMapValue()

		if err := data.RangeSlice(func(i int, item *model.Value) error {
			keyVal, err := ExecuteAST(ctx, e.Expr, item, options)
			if err != nil {
				return err
			}

			key, err := valueToString(keyVal)
			if err != nil {
				return fmt.Errorf("cannot use %s as group key: %w", keyVal.Type(), err)
			}

			exists, err := res.MapKeyExists(key)
			if err != nil {
				return err
			}

			if exists {
				group, err := res.GetMapKey(key)
				if err != nil {
					return err
				}
				if err := group.Append(item); err != nil {
					return err
				}
			} else {
				group := model.NewSliceValue()
				if err := group.Append(item); err != nil {
					return err
				}
				if err := res.SetMapKey(key, group); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return nil, fmt.Errorf("error ranging over slice: %w", err)
		}

		return res, nil
	}, nil
}

func valueToString(v *model.Value) (string, error) {
	switch v.Type() {
	case model.TypeString:
		return v.StringValue()
	case model.TypeInt:
		i, err := v.IntValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case model.TypeFloat:
		f, err := v.FloatValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%g", f), nil
	case model.TypeBool:
		b, err := v.BoolValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", b), nil
	default:
		return "", fmt.Errorf("cannot convert %s to string for use as group key", v.Type())
	}
}
