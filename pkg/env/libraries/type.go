package libraries

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/param"
	"github.com/RyanCopley/expression-parser/pkg/types"
	"strconv"
	"strings"

	"github.com/RyanCopley/expression-parser/pkg/errors"
)

// TypeLib implements type conversion and type-checking functions.
type TypeLib struct{}

func NewTypeLib() *TypeLib {
	return &TypeLib{}
}

func (t *TypeLib) Call(functionName string, args []param.Arg, line, col, _, _ int) (interface{}, error) {
	switch functionName {
	case "string":
		if len(args) != 1 {
			return nil, errors.NewParameterError("type.string requires 1 argument", line, col)
		}
		arg0 := args[0]
		if arg0.Value == nil {
			return "", nil
		}
		return fmt.Sprintf("%v", arg0.Value), nil

	case "int":
		if len(args) != 1 {
			return nil, errors.NewParameterError("type.int requires 1 argument", line, col)
		}
		arg0 := args[0]
		if arg0.Value == nil {
			return int64(0), nil
		}
		switch v := arg0.Value.(type) {
		case string:
			s := strings.TrimSpace(v)
			if len(s) >= 2 {
				if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
					s = s[1 : len(s)-1]
				}
			}
			if i, err := strconv.ParseInt(s, 10, 64); err == nil {
				return i, nil
			}
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				return int64(f), nil
			}
			return nil, errors.NewFunctionCallError(fmt.Sprintf("type.int: string '%s' cannot be converted to int", v), arg0.Line, arg0.Column)
		default:
			num, ok := types.ToFloat(arg0.Value)
			if !ok {
				return nil, errors.NewTypeError("type.int: argument cannot be converted to int", arg0.Line, arg0.Column)
			}
			return int64(num), nil
		}

	case "float":
		if len(args) != 1 {
			return nil, errors.NewParameterError("type.float requires 1 argument", line, col)
		}
		arg0 := args[0]
		if arg0.Value == nil {
			return 0.0, nil
		}
		switch v := arg0.Value.(type) {
		case string:
			s := strings.TrimSpace(v)
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				return f, nil
			} else {
				return nil, errors.NewFunctionCallError(fmt.Sprintf("type.float: string '%s' cannot be converted to float", v), arg0.Line, arg0.Column)
			}
		default:
			num, ok := types.ToFloat(arg0.Value)
			if !ok {
				return nil, errors.NewTypeError("type.float: argument cannot be converted to float", arg0.Line, arg0.Column)
			}
			return num, nil
		}

	case "isNumber":
		if len(args) != 1 {
			return nil, errors.NewParameterError("type.isNumber requires 1 argument", line, col)
		}
		arg0 := args[0]
		switch v := arg0.Value.(type) {
		case string:
			_, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			return err == nil, nil
		default:
			_, ok := types.ToFloat(arg0.Value)
			return ok, nil
		}

	case "isString":
		if len(args) != 1 {
			return nil, errors.NewParameterError("type.isString requires 1 argument", line, col)
		}
		_, ok := args[0].Value.(string)
		return ok, nil

	case "isBoolean":
		if len(args) != 1 {
			return nil, errors.NewParameterError("type.isBoolean requires 1 argument", line, col)
		}
		_, ok := args[0].Value.(bool)
		return ok, nil

	case "isArray":
		if len(args) != 1 {
			return nil, errors.NewParameterError("type.isArray requires 1 argument", line, col)
		}
		_, ok := types.ConvertToInterfaceSlice(args[0].Value)
		return ok, nil

	case "isObject":
		if len(args) != 1 {
			return nil, errors.NewParameterError("type.isObject requires 1 argument", line, col)
		}
		_, ok := types.ConvertToStringMap(args[0].Value)
		return ok, nil

	case "isNull":
		if len(args) != 1 {
			return nil, errors.NewParameterError("type.isNull requires 1 argument", line, col)
		}
		return args[0].Value == nil, nil

	default:
		return nil, errors.NewFunctionCallError(fmt.Sprintf("unknown type function '%s'", functionName), 0, 0)
	}
}
