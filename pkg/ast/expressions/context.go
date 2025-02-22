package expressions

import (
	"fmt"
	"github.com/SpecDrivenDesign/lql/pkg/ast"
	"github.com/SpecDrivenDesign/lql/pkg/env"
	"github.com/SpecDrivenDesign/lql/pkg/errors"
)

// ContextExpr represents a context reference (e.g. $identifier or $[expression]).
type ContextExpr struct {
	Ident     *IdentifierExpr
	Subscript ast.Expression
	Line      int
	Column    int
}

func (c *ContextExpr) Eval(ctx map[string]interface{}, env *env.Environment) (interface{}, error) {
	if c.Ident != nil {
		if val, ok := ctx[c.Ident.Name]; ok {
			return val, nil
		}
		return nil, errors.NewReferenceError(fmt.Sprintf("field '%s' not found", c.Ident.Name), c.Ident.Line, c.Ident.Column)
	}
	return ctx, nil
}

func (c *ContextExpr) Pos() (int, int) {
	return c.Line, c.Column
}

func (c *ContextExpr) String() string {
	// If there's an identifier, we produce something like "$myField".
	// If there's a subscript expression, we produce something like "$[someExpr]".
	// If both are nil, it's just "$".

	// Base "$" symbol (maybe colored if ColorEnabled).
	dollar := "$"
	if ColorEnabled {
		dollar = PunctuationColor + "$" + ColorReset
	}

	// If we have an identifier, we build "$ident".
	if c.Ident != nil {
		identName := c.Ident.Name
		if ColorEnabled {
			identName = ContextColor + identName + ColorReset
		}
		return dollar + identName
	}

	// If we have a subscript expression, build "$[ expression ]".
	if c.Subscript != nil {
		openBracket := "["
		closeBracket := "]"

		if ColorEnabled {
			openBracket = PunctuationColor + "[" + ColorReset
			closeBracket = PunctuationColor + "]" + ColorReset
		}

		subscriptStr := c.Subscript.String()
		return dollar + openBracket + subscriptStr + closeBracket
	}

	// Otherwise, it's just "$" referencing the entire context.
	return dollar
}
