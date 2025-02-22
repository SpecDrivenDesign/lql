package expressions

import (
	"fmt"
	"github.com/SpecDrivenDesign/lql/pkg/ast"
	"github.com/SpecDrivenDesign/lql/pkg/env"
	"strings"
)

// ArrayLiteralExpr represents an array literal.
type ArrayLiteralExpr struct {
	Elements []ast.Expression
	Line     int
	Column   int
}

func (a *ArrayLiteralExpr) Eval(ctx map[string]interface{}, env *env.Environment) (interface{}, error) {
	var result []interface{}
	for _, expr := range a.Elements {
		val, err := expr.Eval(ctx, env)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}

func (a *ArrayLiteralExpr) Pos() (int, int) {
	return a.Line, a.Column
}

func (a *ArrayLiteralExpr) String() string {
	var sb strings.Builder

	// Default punctuation strings (uncolored).
	openBracket := "["
	closeBracket := "]"
	comma := ", "

	// If color is enabled, override with colored brackets/commas.
	if ColorEnabled {
		openBracket = fmt.Sprintf("%s[%s", PunctuationColor, ColorReset)
		closeBracket = fmt.Sprintf("%s]%s", PunctuationColor, ColorReset)
		comma = fmt.Sprintf("%s,%s ", PunctuationColor, ColorReset)
	}

	sb.WriteString(openBracket)

	for i, elem := range a.Elements {
		if i > 0 {
			sb.WriteString(comma)
		}
		sb.WriteString(elem.String())
	}

	sb.WriteString(closeBracket)
	return sb.String()
}
