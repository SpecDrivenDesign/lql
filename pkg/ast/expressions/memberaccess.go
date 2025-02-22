package expressions

import (
	"fmt"
	"github.com/SpecDrivenDesign/lql/pkg/ast"
	"github.com/SpecDrivenDesign/lql/pkg/env"
	"github.com/SpecDrivenDesign/lql/pkg/errors"
	"github.com/SpecDrivenDesign/lql/pkg/types"
	"strings"
)

// MemberPart represents a part of a member access (either dot or bracket).
type MemberPart struct {
	Optional bool
	IsIndex  bool
	Key      string
	Expr     ast.Expression
	Line     int
	Column   int
}

// MemberAccessExpr represents member access (dot or bracket notation).
type MemberAccessExpr struct {
	Target      ast.Expression
	AccessParts []MemberPart
}

func (m *MemberAccessExpr) Eval(ctx map[string]interface{}, env *env.Environment) (interface{}, error) {
	val, err := m.Target.Eval(ctx, env)
	if err != nil {
		return nil, err
	}
	for _, part := range m.AccessParts {
		if val == nil && part.Optional {
			return nil, nil
		}
		if part.IsIndex {
			indexVal, err := part.Expr.Eval(ctx, env)
			if err != nil {
				return nil, err
			}
			if obj, ok := types.ConvertToStringMap(val); ok {
				var key string
				switch v := indexVal.(type) {
				case string:
					key = v
				default:
					key = fmt.Sprintf("%v", v)
				}
				if v, exists := obj[key]; exists {
					val = v
				} else {
					if part.Optional {
						return nil, nil
					}
					return nil, errors.NewReferenceError(fmt.Sprintf("field '%s' not found", key), part.Line, part.Column)
				}
			} else if arr, ok := types.ConvertToInterfaceSlice(val); ok {
				idx, ok := types.ToInt(indexVal)
				if !ok {
					return nil, errors.NewTypeError("array index must be numeric", part.Line, part.Column)
				}
				if idx < 0 || idx >= int64(len(arr)) {
					if part.Optional {
						return nil, nil
					}
					return nil, errors.NewArrayOutOfBoundsError("array index out of bounds", part.Line, part.Column)
				}
				val = arr[idx]
			} else {
				return nil, errors.NewTypeError("target is not an object or array", part.Line, part.Column)
			}
		} else {
			obj, ok := types.ConvertToStringMap(val)
			if !ok {
				return nil, errors.NewTypeError("dot access on non‑object", part.Line, part.Column)
			}
			if v, exists := obj[part.Key]; exists {
				val = v
			} else {
				if part.Optional {
					return nil, nil
				}
				return nil, errors.NewReferenceError(fmt.Sprintf("field '%s' not found", part.Key), part.Line, part.Column)
			}
		}
	}
	return val, nil
}

func (m *MemberAccessExpr) Pos() (int, int) {
	return m.Target.Pos()
}
func (m *MemberAccessExpr) String() string {
	var sb strings.Builder

	// Start with the string form of the target expression.
	sb.WriteString(m.Target.String())

	for _, part := range m.AccessParts {

		// Optional chaining operator ('?') if part.Optional == true
		if part.Optional {
			if ColorEnabled {
				sb.WriteString(PunctuationColor + "?" + ColorReset)
			} else {
				sb.WriteString("?")
			}
		}

		// Bracket vs. dot notation
		if part.IsIndex {
			// Build something like "[expr]" or "[0]" (colored if enabled)
			openBracket := "["
			closeBracket := "]"

			if ColorEnabled {
				openBracket = PunctuationColor + "[" + ColorReset
				closeBracket = PunctuationColor + "]" + ColorReset
			}
			sb.WriteString(openBracket)

			if part.Expr != nil {
				sb.WriteString(part.Expr.String())
			}
			sb.WriteString(closeBracket)
		} else {
			// Dot notation
			dot := "."
			if ColorEnabled {
				dot = PunctuationColor + "." + ColorReset
			}
			sb.WriteString(dot)

			keyStr := part.Key
			if ColorEnabled {
				keyStr = ContextColor + keyStr + ColorReset
			}
			sb.WriteString(keyStr)
		}
	}

	return sb.String()
}
