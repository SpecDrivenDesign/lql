package libraries

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/param"
	"github.com/RyanCopley/expression-parser/pkg/types"
	"regexp"

	"github.com/RyanCopley/expression-parser/pkg/errors"
)

// RegexLib implements regex functions.
type RegexLib struct{}

func NewRegexLib() *RegexLib {
	return &RegexLib{}
}

func (r *RegexLib) Call(functionName string, args []param.Arg, line, col, parenLine, parenCol int) (interface{}, error) {
	switch functionName {
	case "match":
		if len(args) != 2 {
			return nil, errors.NewParameterError("regex.match requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arg1 := args[1]
		pattern, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("regex.match: first argument must be a string", arg0.Line, arg0.Column)
		}
		s, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("regex.match: second argument must be a string", arg1.Line, arg1.Column)
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, errors.NewTypeError("regex.match: invalid pattern", arg0.Line, arg0.Column)
		}
		return re.MatchString(s), nil

	case "replace":
		if len(args) < 3 || len(args) > 4 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("regex.replace requires 3 or 4 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("regex.replace requires 3 or 4 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arg1 := args[1]
		arg2 := args[2]
		s, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("regex.replace: first argument must be a string", arg0.Line, arg0.Column)
		}
		pattern, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("regex.replace: second argument must be a string", arg1.Line, arg1.Column)
		}
		replacement, ok := arg2.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("regex.replace: third argument must be a string", arg2.Line, arg2.Column)
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, errors.NewTypeError("regex.replace: invalid pattern", arg1.Line, arg1.Column)
		}
		if len(args) == 3 {
			return re.ReplaceAllString(s, replacement), nil
		}
		arg3 := args[3]
		lArg, ok := types.ToInt(arg3.Value)
		if !ok {
			return nil, errors.NewTypeError("regex.replace: fourth argument must be numeric", arg3.Line, arg3.Column)
		}
		limit := int(lArg)
		result := s
		for i := 0; i < limit; i++ {
			loc := re.FindStringIndex(result)
			if loc == nil {
				break
			}
			replaced := re.ReplaceAllString(result[loc[0]:loc[1]], replacement)
			result = result[:loc[0]] + replaced + result[loc[1]:]
		}
		return result, nil

	case "find":
		if len(args) != 2 {
			return nil, errors.NewParameterError("regex.find requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arg1 := args[1]
		pattern, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("regex.find: first argument must be a string", arg0.Line, arg0.Column)
		}
		s, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("regex.find: second argument must be a string", arg1.Line, arg1.Column)
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, errors.NewTypeError("regex.find: invalid pattern", arg0.Line, arg0.Column)
		}
		match := re.FindString(s)
		if match == "" {
			return "", nil
		}
		return match, nil

	default:
		return nil, errors.NewFunctionCallError(fmt.Sprintf("unknown regex function '%s'", functionName), 0, 0)
	}
}
