package expressions

import (
	"github.com/RyanCopley/expression-parser/pkg/ast"
	"github.com/RyanCopley/expression-parser/pkg/env"
	"github.com/RyanCopley/expression-parser/pkg/errors"
	"github.com/RyanCopley/expression-parser/pkg/tokens"
	"github.com/RyanCopley/expression-parser/pkg/types"
)

// UnaryExpr represents a unary operation.
type UnaryExpr struct {
	Operator tokens.TokenType
	Expr     ast.Expression
	Line     int
	Column   int
}

func (u *UnaryExpr) Eval(ctx map[string]interface{}, env *env.Environment) (interface{}, error) {
	val, err := u.Expr.Eval(ctx, env)
	if err != nil {
		return nil, err
	}
	switch u.Operator {
	case tokens.TokenMinus:
		num, ok := types.ToFloat(val)
		if !ok {
			return nil, errors.NewSemanticError("unary '-' operator requires a numeric operand", u.Line, u.Column)
		}
		if types.IsInt(val) {
			return int64(-num), nil
		}
		return -num, nil
	case tokens.TokenNot:
		b, ok := val.(bool)
		if !ok {
			return nil, errors.NewSemanticError("NOT operator requires a boolean operand", u.Line, u.Column)
		}
		return !b, nil
	default:
		return nil, errors.NewUnknownOperatorError("unknown unary operator", u.Line, u.Column)
	}
}

func (u *UnaryExpr) Pos() (int, int) {
	return u.Line, u.Column
}
func (u *UnaryExpr) String() string {
	exprStr := u.Expr.String()

	// Convert operator token to its DSL string form.
	var opStr string
	switch u.Operator {
	case tokens.TokenMinus:
		opStr = "-"
	case tokens.TokenNot:
		opStr = "NOT"
	default:
		// Fallback if needed; your tokens may or may not have a .String()
		opStr = tokens.FixedTokenLiterals[u.Operator]
	}

	// Apply operator color if enabled.
	if ColorEnabled {
		opStr = OperatorColor + opStr + ColorReset
	}

	// For a minus operator, we typically do "-(expr)" if expression is more complex,
	// or just "-expr" if it's a single literal or variable. For simplicity:
	if u.Operator == tokens.TokenMinus {
		// E.g. "-(x + y)" or "-3"
		return opStr + exprStr
	}

	// For NOT, we often do "NOT expr" with a space.
	return opStr + " " + exprStr
}
