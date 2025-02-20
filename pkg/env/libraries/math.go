package libraries

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/errors"
	"github.com/RyanCopley/expression-parser/pkg/param"
	"github.com/RyanCopley/expression-parser/pkg/types"
	"math"
)

// MathLib implements math library functions.
type MathLib struct{}

func NewMathLib() *MathLib {
	return &MathLib{}
}

func (m *MathLib) Call(functionName string, args []param.Arg, line, col, parenLine, parenCol int) (interface{}, error) {
	switch functionName {
	case "abs":
		if len(args) != 1 {
			return nil, errors.NewParameterError("math.abs requires 1 argument", line, col)
		}
		arg0 := args[0]
		num, ok := types.ToFloat(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("math.abs: argument must be numeric", arg0.Line, arg0.Column)
		}
		if num < 0 {
			if types.IsInt(arg0.Value) {
				return int64(-num), nil
			}
			return -num, nil
		}
		if types.IsInt(arg0.Value) {
			return int64(num), nil
		}
		return num, nil

	case "sqrt":
		if len(args) != 1 {
			return nil, errors.NewParameterError("math.sqrt requires 1 argument", line, col)
		}
		arg0 := args[0]
		num, ok := types.ToFloat(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("math.sqrt: argument must be numeric", arg0.Line, arg0.Column)
		}
		if num < 0 {
			return nil, errors.NewFunctionCallError("math.sqrt: argument must be non‑negative", arg0.Line, arg0.Column)
		}
		return math.Sqrt(num), nil

	case "floor":
		if len(args) != 1 {
			return nil, errors.NewParameterError("math.floor requires 1 argument", line, col)
		}
		arg0 := args[0]
		num, ok := types.ToFloat(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("math.floor: argument must be numeric", arg0.Line, arg0.Column)
		}
		return math.Floor(num), nil

	case "round":
		if len(args) != 1 {
			return nil, errors.NewParameterError("math.round requires 1 argument", line, col)
		}
		arg0 := args[0]
		num, ok := types.ToFloat(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("math.round: argument must be numeric", arg0.Line, arg0.Column)
		}
		return math.Round(num), nil

	case "ceil":
		if len(args) != 1 {
			return nil, errors.NewParameterError("math.ceil requires 1 argument", line, col)
		}
		arg0 := args[0]
		num, ok := types.ToFloat(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("math.ceil: argument must be numeric", arg0.Line, arg0.Column)
		}
		return math.Ceil(num), nil

	case "pow":
		if len(args) != 2 {
			return nil, errors.NewParameterError("math.pow requires 2 arguments", line, col)
		}
		arg0 := args[0]
		base, ok := types.ToFloat(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("math.pow: first argument must be numeric", arg0.Line, arg0.Column)
		}
		arg1 := args[1]
		exp, ok := types.ToFloat(arg1.Value)
		if !ok {
			return nil, errors.NewTypeError("math.pow: second argument must be numeric", arg1.Line, arg1.Column)
		}
		return math.Pow(base, exp), nil

	case "sum":
		if len(args) < 1 || len(args) > 3 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("math.sum requires 1 to 3 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("math.sum requires 1 to 3 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("math.sum: first argument must be an array", arg0.Line, arg0.Column)
		}
		var subfield string
		var defaultVal interface{}
		if len(args) >= 2 {
			arg1 := args[1]
			sf, ok := arg1.Value.(string)
			if !ok {
				return nil, errors.NewTypeError("math.sum: second argument must be string", arg1.Line, arg1.Column)
			}
			subfield = sf
		}
		if len(args) == 3 {
			defaultVal = args[2].Value
		}
		sum := 0.0
		for _, elem := range arr {
			var num interface{}
			if subfield != "" {
				obj, ok := types.ConvertToStringMap(elem)
				if !ok {
					if defaultVal != nil {
						num = defaultVal
					} else {
						return nil, errors.NewFunctionCallError("math.sum: element is not an object and subfield specified", arg0.Line, arg0.Column)
					}
				} else {
					if v, exists := obj[subfield]; exists {
						num = v
					} else {
						if defaultVal != nil {
							num = defaultVal
						} else {
							return nil, errors.NewFunctionCallError(fmt.Sprintf("math.sum: field '%s' missing in element", subfield), arg0.Line, arg0.Column)
						}
					}
				}
			} else {
				num = elem
			}
			nf, ok := types.ToFloat(num)
			if !ok {
				return nil, errors.NewTypeError("math.sum: element is not numeric", arg0.Line, arg0.Column)
			}
			sum += nf
		}
		return sum, nil

	case "min":
		if len(args) < 1 || len(args) > 3 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("math.min requires 1 to 3 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("math.min requires 1 to 3 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("math.min: first argument must be an array", arg0.Line, arg0.Column)
		}
		var subfield string
		var defaultVal interface{}
		if len(args) >= 2 {
			arg1 := args[1]
			sf, ok := arg1.Value.(string)
			if !ok {
				return nil, errors.NewTypeError("math.min: second argument must be string", arg1.Line, arg1.Column)
			}
			subfield = sf
		}
		if len(args) == 3 {
			defaultVal = args[2].Value
		}
		if len(arr) == 0 {
			if defaultVal != nil {
				return defaultVal, nil
			}
			return nil, errors.NewFunctionCallError("math.min: array is empty", arg0.Line, arg0.Column)
		}
		var m float64
		first := true
		for _, elem := range arr {
			var num interface{}
			if subfield != "" {
				obj, ok := types.ConvertToStringMap(elem)
				if !ok {
					if defaultVal != nil {
						num = defaultVal
					} else {
						return nil, errors.NewFunctionCallError("math.min: element is not an object and subfield specified", arg0.Line, arg0.Column)
					}
				} else {
					if v, exists := obj[subfield]; exists {
						num = v
					} else {
						if defaultVal != nil {
							num = defaultVal
						} else {
							return nil, errors.NewFunctionCallError(fmt.Sprintf("math.min: field '%s' missing in element", subfield), arg0.Line, arg0.Column)
						}
					}
				}
			} else {
				num = elem
			}
			nf, ok := types.ToFloat(num)
			if !ok {
				return nil, errors.NewTypeError("math.min: element is not numeric", arg0.Line, arg0.Column)
			}
			if first {
				m = nf
				first = false
			} else {
				if nf < m {
					m = nf
				}
			}
		}
		return m, nil

	case "max":
		if len(args) < 1 || len(args) > 3 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("math.max requires 1 to 3 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("math.max requires 1 to 3 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("math.max: first argument must be an array", arg0.Line, arg0.Column)
		}
		var subfield string
		var defaultVal interface{}
		if len(args) >= 2 {
			arg1 := args[1]
			sf, ok := arg1.Value.(string)
			if !ok {
				return nil, errors.NewTypeError("math.max: second argument must be string", arg1.Line, arg1.Column)
			}
			subfield = sf
		}
		if len(args) == 3 {
			defaultVal = args[2].Value
		}
		if len(arr) == 0 {
			if defaultVal != nil {
				return defaultVal, nil
			}
			return nil, errors.NewFunctionCallError("math.max: array is empty", arg0.Line, arg0.Column)
		}
		var m float64
		first := true
		for _, elem := range arr {
			var num interface{}
			if subfield != "" {
				obj, ok := types.ConvertToStringMap(elem)
				if !ok {
					if defaultVal != nil {
						num = defaultVal
					} else {
						return nil, errors.NewFunctionCallError("math.max: element is not an object and subfield specified", arg0.Line, arg0.Column)
					}
				} else {
					if v, exists := obj[subfield]; exists {
						num = v
					} else {
						if defaultVal != nil {
							num = defaultVal
						} else {
							return nil, errors.NewFunctionCallError(fmt.Sprintf("math.max: field '%s' missing in element", subfield), arg0.Line, arg0.Column)
						}
					}
				}
			} else {
				num = elem
			}
			nf, ok := types.ToFloat(num)
			if !ok {
				return nil, errors.NewTypeError("math.max: element is not numeric", arg0.Line, arg0.Column)
			}
			if first {
				m = nf
				first = false
			} else {
				if nf > m {
					m = nf
				}
			}
		}
		return m, nil

	case "avg":
		if len(args) < 1 || len(args) > 3 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("math.avg requires 1 to 3 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("math.avg requires 1 to 3 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("math.avg: first argument must be an array", arg0.Line, arg0.Column)
		}
		var subfield string
		var defaultVal interface{}
		if len(args) >= 2 {
			arg1 := args[1]
			sf, ok := arg1.Value.(string)
			if !ok {
				return nil, errors.NewTypeError("math.avg: second argument must be string", arg1.Line, arg1.Column)
			}
			subfield = sf
		}
		if len(args) == 3 {
			defaultVal = args[2].Value
		}
		if len(arr) == 0 {
			if defaultVal != nil {
				return defaultVal, nil
			}
			return nil, errors.NewFunctionCallError("math.avg: array is empty", arg0.Line, arg0.Column)
		}
		sum := 0.0
		count := 0
		for _, elem := range arr {
			var num interface{}
			if subfield != "" {
				obj, ok := types.ConvertToStringMap(elem)
				if !ok {
					return nil, errors.NewFunctionCallError("math.avg: element is not an object and subfield specified", arg0.Line, arg0.Column)
				}
				if v, exists := obj[subfield]; exists {
					num = v
				} else {
					return nil, errors.NewFunctionCallError(fmt.Sprintf("math.avg: field '%s' missing in element", subfield), arg0.Line, arg0.Column)
				}
			} else {
				num = elem
			}
			nf, ok := types.ToFloat(num)
			if !ok {
				return nil, errors.NewTypeError("math.avg: element is not numeric", arg0.Line, arg0.Column)
			}
			sum += nf
			count++
		}
		return sum / float64(count), nil

	default:
		return nil, errors.NewFunctionCallError(fmt.Sprintf("unknown math function '%s'", functionName), 0, 0)
	}
}
