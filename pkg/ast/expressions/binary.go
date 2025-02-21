package expressions

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/ast"
	"github.com/RyanCopley/expression-parser/pkg/env"
	"github.com/RyanCopley/expression-parser/pkg/errors"
	"github.com/RyanCopley/expression-parser/pkg/tokens"
	"github.com/RyanCopley/expression-parser/pkg/types"
)

// BinaryExpr represents a binary operation.
type BinaryExpr struct {
	Left     ast.Expression
	Operator tokens.TokenType
	Right    ast.Expression
	Line     int
	Column   int
}

func (b *BinaryExpr) Eval(ctx map[string]interface{}, env *env.Environment) (interface{}, error) {
	switch b.Operator {
	case tokens.TokenAnd:
		// Short-circuit: evaluate left operand first.
		leftVal, err := b.Left.Eval(ctx, env)
		if err != nil {
			return nil, err
		}
		lb, ok := leftVal.(bool)
		if !ok {
			return nil, errors.NewSemanticError("AND operator requires boolean operand", b.Line, b.Column)
		}
		if !lb {
			return false, nil
		}
		rightVal, err := b.Right.Eval(ctx, env)
		if err != nil {
			return nil, err
		}
		rb, ok := rightVal.(bool)
		if !ok {
			return nil, errors.NewSemanticError("AND operator requires boolean operand", b.Line, b.Column)
		}
		return rb, nil

	case tokens.TokenOr:
		// Short-circuit: evaluate left operand first.
		leftVal, err := b.Left.Eval(ctx, env)
		if err != nil {
			return nil, err
		}
		lb, ok := leftVal.(bool)
		if !ok {
			return nil, errors.NewSemanticError("OR operator requires boolean operand", b.Line, b.Column)
		}
		if lb {
			return true, nil
		}
		rightVal, err := b.Right.Eval(ctx, env)
		if err != nil {
			return nil, err
		}
		rb, ok := rightVal.(bool)
		if !ok {
			return nil, errors.NewSemanticError("OR operator requires boolean operand", b.Line, b.Column)
		}
		return rb, nil

	default:
		// Evaluate both operands for other operators.
		leftVal, err := b.Left.Eval(ctx, env)
		if err != nil {
			return nil, err
		}
		rightVal, err := b.Right.Eval(ctx, env)
		if err != nil {
			return nil, err
		}
		switch b.Operator {
		case tokens.TokenPlus:
			ln, lok := types.ToFloat(leftVal)
			rn, rok := types.ToFloat(rightVal)
			if !lok || !rok {
				return nil, errors.NewSemanticError("'+' operator used on non‑numeric type", b.Line, b.Column)
			}
			if types.IsInt(leftVal) != types.IsInt(rightVal) {
				return nil, errors.NewSemanticError("Mixed numeric types require explicit conversion", b.Line, b.Column)
			}
			if types.IsInt(leftVal) {
				return int64(ln + rn), nil
			}
			return ln + rn, nil

		case tokens.TokenMinus:
			ln, lok := types.ToFloat(leftVal)
			rn, rok := types.ToFloat(rightVal)
			if !lok || !rok {
				return nil, errors.NewSemanticError("'-' operator used on non‑numeric type", b.Line, b.Column)
			}
			if types.IsInt(leftVal) != types.IsInt(rightVal) {
				return nil, errors.NewSemanticError("Mixed numeric types require explicit conversion", b.Line, b.Column)
			}
			if types.IsInt(leftVal) {
				return int64(ln - rn), nil
			}
			return ln - rn, nil

		case tokens.TokenMultiply:
			ln, lok := types.ToFloat(leftVal)
			rn, rok := types.ToFloat(rightVal)
			if !lok || !rok {
				return nil, errors.NewSemanticError("'*' operator used on non‑numeric type", b.Line, b.Column)
			}
			if types.IsInt(leftVal) != types.IsInt(rightVal) {
				return nil, errors.NewSemanticError("Mixed numeric types require explicit conversion", b.Line, b.Column)
			}
			if types.IsInt(leftVal) {
				return int64(ln * rn), nil
			}
			return ln * rn, nil

		case tokens.TokenDivide:
			ln, lok := types.ToFloat(leftVal)
			rn, rok := types.ToFloat(rightVal)
			if !lok || !rok {
				return nil, errors.NewSemanticError("'/' operator used on non‑numeric type", b.Line, b.Column)
			}
			if rn == 0 {
				return nil, errors.NewDivideByZeroError("division by zero", b.Line, b.Column)
			}
			if types.IsInt(leftVal) != types.IsInt(rightVal) {
				return nil, errors.NewSemanticError("Mixed numeric types require explicit conversion", b.Line, b.Column)
			}
			if types.IsInt(leftVal) {
				return int64(ln / rn), nil
			}
			return ln / rn, nil

		case tokens.TokenLt:
			return types.Compare(leftVal, rightVal, "<", b.Line, b.Column)
		case tokens.TokenGt:
			return types.Compare(leftVal, rightVal, ">", b.Line, b.Column)
		case tokens.TokenLte:
			return types.Compare(leftVal, rightVal, "<=", b.Line, b.Column)
		case tokens.TokenGte:
			return types.Compare(leftVal, rightVal, ">=", b.Line, b.Column)
		case tokens.TokenEq:
			return types.Equals(leftVal, rightVal), nil
		case tokens.TokenNeq:
			return !types.Equals(leftVal, rightVal), nil
		}
	}
	return nil, errors.NewUnknownOperatorError("unknown binary operator", b.Line, b.Column)
}

func (b *BinaryExpr) Pos() (int, int) {
	return b.Line, b.Column
}

func (b *BinaryExpr) String() string {
	leftStr := b.Left.String()
	rightStr := b.Right.String()
	opStr := tokens.FixedTokenLiterals[b.Operator]
	if ColorEnabled {
		opStr = OperatorColor + opStr + ColorReset
	}
	return fmt.Sprintf("%s %s %s", leftStr, opStr, rightStr)
}
