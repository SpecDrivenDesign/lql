package libraries

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/param"
	"github.com/RyanCopley/expression-parser/pkg/types"

	"github.com/RyanCopley/expression-parser/pkg/errors"
)

// CondLib implements conditional library functions.
type CondLib struct{}

func NewCondLib() *CondLib {
	return &CondLib{}
}

func (c *CondLib) Call(functionName string, args []param.Arg, line, col, parenLine, parenCol int) (interface{}, error) {
	switch functionName {
	case "ifExpr":
		if len(args) != 3 {
			return nil, errors.NewParameterError("cond.ifExpr requires 3 arguments", line, col)
		}
		arg0 := args[0]
		condVal, ok := arg0.Value.(bool)
		if !ok {
			if arg0.Value == nil {
				condVal = false
			} else {
				return nil, errors.NewTypeError("cond.ifExpr: first argument must be boolean", arg0.Line, arg0.Column)
			}
		}
		if condVal {
			return args[1].Value, nil
		}
		return args[2].Value, nil

	case "coalesce":
		if len(args) < 1 {
			return nil, errors.NewParameterError("cond.coalesce requires at least 1 argument", parenLine, parenCol)
		}
		for _, arg := range args {
			if arg.Value != nil {
				return arg.Value, nil
			}
		}
		return nil, errors.NewFunctionCallError("cond.coalesce: all arguments are null", args[0].Line, args[0].Column)

	case "isFieldPresent":
		if len(args) != 2 {
			return nil, errors.NewParameterError("cond.isFieldPresent requires 2 arguments", line, col)
		}
		arg0 := args[0]
		obj, ok := types.ConvertToStringMap(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("cond.isFieldPresent: first argument must be an object", arg0.Line, arg0.Column)
		}
		arg1 := args[1]
		fieldPath, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("cond.isFieldPresent: second argument must be a string", arg1.Line, arg1.Column)
		}
		_, exists := obj[fieldPath]
		return exists, nil

	default:
		return nil, errors.NewFunctionCallError(fmt.Sprintf("unknown cond function '%s'", functionName), 0, 0)
	}
}
