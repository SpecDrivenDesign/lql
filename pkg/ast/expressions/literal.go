package expressions

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/env"
)

// LiteralExpr represents a literal value.
type LiteralExpr struct {
	Value  interface{}
	Line   int
	Column int
}

func (l *LiteralExpr) Eval(ctx map[string]interface{}, env *env.Environment) (interface{}, error) {
	return l.Value, nil
}

func (l *LiteralExpr) Pos() (int, int) {
	return l.Line, l.Column
}
func (l *LiteralExpr) String() string {
	var s string

	switch v := l.Value.(type) {

	case string:
		// Enclose strings in quotes, then optionally color.
		s = `"` + v + `"`
		if ColorEnabled {
			s = StringColor + s + ColorReset
		}

	case bool:
		// Lowercase "true"/"false" to match DSL specs, then optionally color.
		if v {
			s = "true"
		} else {
			s = "false"
		}
		if ColorEnabled {
			s = BoolNullColor + s + ColorReset
		}

	case nil:
		// null literal.
		s = "null"
		if ColorEnabled {
			s = BoolNullColor + s + ColorReset
		}

	case int, int64, float64:
		// Numeric literal -> convert to string, optionally color.
		s = fmt.Sprintf("%v", v)
		if ColorEnabled {
			s = NumberColor + s + ColorReset
		}

	default:
		// Fallback: just stringify with fmt.
		s = fmt.Sprintf("%v", v)
		// Optionally color if you'd like,
		// but typically unknown types are not recognized by DSL specs.
	}

	return s
}
