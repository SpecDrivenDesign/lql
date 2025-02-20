package expressions

import (
	"github.com/RyanCopley/expression-parser/pkg/ast"
	"github.com/RyanCopley/expression-parser/pkg/env"
	"strings"
)

// ObjectLiteralExpr represents an object literal.
type ObjectLiteralExpr struct {
	Fields map[string]ast.Expression
	Line   int
	Column int
}

func (o *ObjectLiteralExpr) Eval(ctx map[string]interface{}, env *env.Environment) (interface{}, error) {
	result := make(map[string]interface{})
	for key, expr := range o.Fields {
		val, err := expr.Eval(ctx, env)
		if err != nil {
			return nil, err
		}
		result[key] = val
	}
	return result, nil
}

func (o *ObjectLiteralExpr) Pos() (int, int) {
	return o.Line, o.Column
}
func (o *ObjectLiteralExpr) String() string {
	var sb strings.Builder

	// Basic punctuation
	openBrace := "{"
	closeBrace := "}"
	colon := ": "
	comma := ", "

	// If color is enabled, wrap punctuation in ANSI color codes
	if ColorEnabled {
		openBrace = PunctuationColor + "{" + ColorReset
		closeBrace = PunctuationColor + "}" + ColorReset
		colon = PunctuationColor + ":" + ColorReset + " "
		comma = PunctuationColor + "," + ColorReset + " "
	}

	sb.WriteString(openBrace)

	i := 0
	for key, expr := range o.Fields {
		// Insert commas between fields
		if i > 0 {
			sb.WriteString(comma)
		}

		// Decide how to print the key: If it's a valid identifier or not.
		// For simplicity, always quote the key here. You could do a check if you want.
		quotedKey := `"` + key + `"`
		if ColorEnabled {
			// Color the key as an identifier or as a string—your choice.
			// We'll treat it like a string literal for consistency.
			quotedKey = StringColor + quotedKey + ColorReset
		}

		sb.WriteString(quotedKey)
		sb.WriteString(colon)

		// The expression value
		sb.WriteString(expr.String())

		i++
	}

	sb.WriteString(closeBrace)
	return sb.String()
}
