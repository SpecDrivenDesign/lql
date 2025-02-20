package libraries

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/param"
	"github.com/RyanCopley/expression-parser/pkg/types"
	"strings"

	"github.com/RyanCopley/expression-parser/pkg/errors"
)

// StringLib implements string manipulation functions.
type StringLib struct{}

func NewStringLib() *StringLib {
	return &StringLib{}
}

func (s *StringLib) Call(functionName string, args []param.Arg, line, col, parenLine, parenCol int) (interface{}, error) {
	switch functionName {
	case "concat":
		if len(args) < 1 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("string.concat requires at least 1 argument", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("string.concat requires at least 1 argument", lastArg.Line, lastArg.Column)
		}
		var sb strings.Builder
		for _, arg := range args {
			str, ok := arg.Value.(string)
			if !ok {
				return nil, errors.NewTypeError("string.concat: all arguments must be strings", arg.Line, arg.Column)
			}
			sb.WriteString(str)
		}
		return sb.String(), nil

	case "toLower":
		if len(args) != 1 {
			return nil, errors.NewParameterError("string.toLower requires 1 argument", line, col)
		}
		arg0 := args[0]
		str, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.toLower: argument must be string", arg0.Line, arg0.Column)
		}
		return strings.ToLower(str), nil

	case "toUpper":
		if len(args) != 1 {
			return nil, errors.NewParameterError("string.toUpper requires 1 argument", line, col)
		}
		arg0 := args[0]
		str, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.toUpper: argument must be string", arg0.Line, arg0.Column)
		}
		return strings.ToUpper(str), nil

	case "trim":
		if len(args) != 1 {
			return nil, errors.NewParameterError("string.trim requires 1 argument", line, col)
		}
		arg0 := args[0]
		str, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.trim: argument must be string", arg0.Line, arg0.Column)
		}
		return strings.TrimSpace(str), nil

	case "startsWith":
		if len(args) != 2 {
			return nil, errors.NewParameterError("string.startsWith requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arg1 := args[1]
		str, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.startsWith: first argument must be string", arg0.Line, arg0.Column)
		}
		prefix, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.startsWith: second argument must be string", arg1.Line, arg1.Column)
		}
		return strings.HasPrefix(str, prefix), nil

	case "endsWith":
		if len(args) != 2 {
			return nil, errors.NewParameterError("string.endsWith requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arg1 := args[1]
		str, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.endsWith: first argument must be string", arg0.Line, arg0.Column)
		}
		suffix, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.endsWith: second argument must be string", arg1.Line, arg1.Column)
		}
		return strings.HasSuffix(str, suffix), nil

	case "contains":
		if len(args) != 2 {
			return nil, errors.NewParameterError("string.contains requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arg1 := args[1]
		str, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.contains: first argument must be string", arg0.Line, arg0.Column)
		}
		substr, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.contains: second argument must be string", arg1.Line, arg1.Column)
		}
		return strings.Contains(str, substr), nil

	case "split":
		if len(args) != 2 {
			return nil, errors.NewParameterError("string.split requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arg1 := args[1]
		str, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.split: first argument must be string", arg0.Line, arg0.Column)
		}
		delim, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.split: second argument must be string", arg1.Line, arg1.Column)
		}
		return strings.Split(str, delim), nil

	case "join":
		if len(args) != 2 {
			return nil, errors.NewParameterError("string.join requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arg1 := args[1]
		arr, ok := arg0.Value.([]interface{})
		if !ok {
			return nil, errors.NewTypeError("string.join: first argument must be an array", arg0.Line, arg0.Column)
		}
		sep, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.join: second argument must be string", arg1.Line, arg1.Column)
		}
		var parts []string
		for _, item := range arr {
			s, ok := item.(string)
			if !ok {
				return nil, errors.NewTypeError("string.join: all array elements must be strings", arg0.Line, arg0.Column)
			}
			parts = append(parts, s)
		}
		return strings.Join(parts, sep), nil

	case "substring":
		if len(args) != 3 {
			return nil, errors.NewParameterError("string.substring requires 3 arguments", line, col)
		}
		arg0 := args[0]
		str, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.substring: first argument must be a string", arg0.Line, arg0.Column)
		}
		arg1 := args[1]
		start, ok := types.ToInt(arg1.Value)
		if !ok {
			return nil, errors.NewTypeError("string.substring: second argument must be an integer", arg1.Line, arg1.Column)
		}
		arg2 := args[2]
		length, ok := types.ToInt(arg2.Value)
		if !ok {
			return nil, errors.NewTypeError("string.substring: third argument must be an integer", arg2.Line, arg2.Column)
		}
		runes := []rune(str)
		if int(start) < 0 || int(start) >= len(runes) {
			return nil, errors.NewFunctionCallError("string.substring: start index out of bounds", arg1.Line, arg1.Column)
		}
		end := int(start) + int(length)
		if end > len(runes) {
			end = len(runes)
		}
		return string(runes[int(start):end]), nil

	case "replace":
		if len(args) < 3 || len(args) > 4 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("string.replace requires 3 or 4 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("string.replace requires 3 or 4 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arg1 := args[1]
		arg2 := args[2]
		sArg, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.replace: first argument must be a string", arg0.Line, arg0.Column)
		}
		oldArg, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.replace: second argument must be a string", arg1.Line, arg1.Column)
		}
		newArg, ok := arg2.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.replace: third argument must be a string", arg2.Line, arg2.Column)
		}
		limit := -1
		if len(args) == 4 {
			arg3 := args[3]
			lArg, ok := types.ToInt(arg3.Value)
			if !ok {
				return nil, errors.NewTypeError("string.replace: fourth argument must be numeric", arg3.Line, arg3.Column)
			}
			limit = int(lArg)
		}
		if limit < 0 {
			return strings.ReplaceAll(sArg, oldArg, newArg), nil
		}
		return strings.Replace(sArg, oldArg, newArg, limit), nil

	case "indexOf":
		if len(args) < 2 || len(args) > 3 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("string.indexOf requires 2 or 3 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("string.indexOf requires 2 or 3 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arg1 := args[1]
		str, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.indexOf: first argument must be a string", arg0.Line, arg0.Column)
		}
		substr, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("string.indexOf: second argument must be a string", arg1.Line, arg1.Column)
		}
		fromIndex := 0
		if len(args) == 3 {
			arg2 := args[2]
			idx, ok := types.ToInt(arg2.Value)
			if !ok {
				return nil, errors.NewTypeError("string.indexOf: third argument must be numeric", arg2.Line, arg2.Column)
			}
			fromIndex = int(idx)
		}
		if fromIndex < 0 || fromIndex >= len(str) {
			return -1, nil
		}
		idx := strings.Index(str[fromIndex:], substr)
		if idx < 0 {
			return -1, nil
		}
		return fromIndex + idx, nil

	default:
		return nil, errors.NewFunctionCallError(fmt.Sprintf("unknown string function '%s'", functionName), 0, 0)
	}
}
