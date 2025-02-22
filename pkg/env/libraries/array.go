package libraries

import (
	"fmt"
	"github.com/SpecDrivenDesign/lql/pkg/errors"
	"github.com/SpecDrivenDesign/lql/pkg/param"
	"github.com/SpecDrivenDesign/lql/pkg/types"
	"sort"
)

// ArrayLib implements the array library functions.
type ArrayLib struct{}

func NewArrayLib() *ArrayLib {
	return &ArrayLib{}
}

func (a *ArrayLib) Call(functionName string, args []param.Arg, line, col, parenLine, parenCol int) (interface{}, error) {
	switch functionName {
	case "contains":
		if len(args) != 2 {
			return nil, errors.NewParameterError("array.contains requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("array.contains: first argument must be an array", arg0.Line, arg0.Column)
		}
		target := args[1].Value
		for _, item := range arr {
			if types.Equals(item, target) {
				return true, nil
			}
		}
		return false, nil

	case "find":
		if len(args) < 3 || len(args) > 4 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("array.find requires 3 or 4 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("array.find requires 3 or 4 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("array.find: first argument must be an array", arg0.Line, arg0.Column)
		}
		arg1 := args[1]
		subfield, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("array.find: second argument must be string", arg1.Line, arg1.Column)
		}
		matchVal := args[2].Value
		var defaultObj interface{}
		if len(args) == 4 {
			defaultObj = args[3].Value
		}
		for _, elem := range arr {
			obj, ok := types.ConvertToStringMap(elem)
			if !ok {
				continue
			}
			if v, exists := obj[subfield]; exists {
				if types.Equals(v, matchVal) {
					return obj, nil
				}
			}
		}
		if defaultObj != nil {
			return defaultObj, nil
		}
		return nil, errors.NewFunctionCallError("array.find: no match found", arg0.Line, arg0.Column)

	case "first":
		if len(args) < 1 || len(args) > 2 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("array.first requires 1 or 2 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("array.first requires 1 or 2 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("array.first: argument must be an array", arg0.Line, arg0.Column)
		}
		if len(arr) == 0 {
			if len(args) == 2 {
				return args[1].Value, nil
			}
			return nil, errors.NewFunctionCallError("array.first: array is empty", arg0.Line, arg0.Column)
		}
		return arr[0], nil

	case "last":
		if len(args) < 1 || len(args) > 2 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("array.last requires 1 or 2 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("array.last requires 1 or 2 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("array.last: argument must be an array", arg0.Line, arg0.Column)
		}
		if len(arr) == 0 {
			if len(args) == 2 {
				return args[1].Value, nil
			}
			return nil, errors.NewFunctionCallError("array.last: array is empty", arg0.Line, arg0.Column)
		}
		return arr[len(arr)-1], nil

	case "extract":
		if len(args) < 2 || len(args) > 3 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("array.extract requires 2 or 3 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("array.extract requires 2 or 3 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("array.extract: argument must be an array", arg0.Line, arg0.Column)
		}
		arg1 := args[1]
		subfield, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("array.extract: second argument must be string", arg1.Line, arg1.Column)
		}
		var defaultVal interface{}
		if len(args) == 3 {
			defaultVal = args[2].Value
		}
		var result []interface{}
		for _, elem := range arr {
			obj, ok := types.ConvertToStringMap(elem)
			if !ok {
				result = append(result, defaultVal)
			} else {
				if v, exists := obj[subfield]; exists {
					result = append(result, v)
				} else {
					result = append(result, defaultVal)
				}
			}
		}
		return result, nil

	case "sort":
		if len(args) < 1 || len(args) > 2 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("array.sort requires 1 or 2 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("array.sort requires 1 or 2 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("array.sort: first argument must be an array", arg0.Line, arg0.Column)
		}
		ascending := true
		if len(args) == 2 {
			arg1 := args[1]
			asc, ok := arg1.Value.(bool)
			if !ok {
				return nil, errors.NewTypeError("array.sort: second argument must be boolean", arg1.Line, arg1.Column)
			}
			ascending = asc
		}
		if len(arr) == 0 {
			return arr, nil
		}
		first := arr[0]
		isNumeric := false
		isString := false
		if _, ok := types.ToFloat(first); ok {
			isNumeric = true
		} else if _, ok := first.(string); ok {
			isString = true
		} else {
			return nil, errors.NewTypeError("array.sort: elements are not comparable", arg0.Line, arg0.Column)
		}
		sorted := make([]interface{}, len(arr))
		copy(sorted, arr)
		sort.SliceStable(sorted, func(i, j int) bool {
			a := sorted[i]
			b := sorted[j]
			if isNumeric {
				af, _ := types.ToFloat(a)
				bf, _ := types.ToFloat(b)
				if ascending {
					return af < bf
				}
				return af > bf
			}
			if isString {
				as := a.(string)
				bs := b.(string)
				if ascending {
					return as < bs
				}
				return as > bs
			}
			return false
		})
		return sorted, nil

	case "flatten":
		if len(args) != 1 {
			return nil, errors.NewParameterError("array.flatten requires 1 argument", line, col)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("array.flatten: argument must be an array", arg0.Line, arg0.Column)
		}
		var result []interface{}
		for _, elem := range arr {
			if subArr, ok := types.ConvertToInterfaceSlice(elem); ok {
				result = append(result, subArr...)
			} else {
				result = append(result, elem)
			}
		}
		return result, nil

	case "filter":
		if len(args) < 1 || len(args) > 3 {
			return nil, errors.NewParameterError("array.filter requires between 1 and 3 arguments", line, col)
		}
		arg0 := args[0]
		arr, ok := types.ConvertToInterfaceSlice(arg0.Value)
		if !ok {
			return nil, errors.NewTypeError("array.filter: first argument must be an array", arg0.Line, arg0.Column)
		}
		if len(args) == 1 {
			var filtered []interface{}
			for _, elem := range arr {
				if elem != nil {
					filtered = append(filtered, elem)
				}
			}
			return filtered, nil
		}
		arg1 := args[1]
		subfield, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("array.filter: subfield argument must be string", arg1.Line, arg1.Column)
		}
		if len(args) == 2 {
			var filtered []interface{}
			for _, elem := range arr {
				obj, ok := types.ConvertToStringMap(elem)
				if !ok {
					continue
				}
				val, exists := obj[subfield]
				if exists && val != nil {
					filtered = append(filtered, elem)
				}
			}
			return filtered, nil
		}
		matchVal := args[2].Value
		var filtered []interface{}
		for _, elem := range arr {
			obj, ok := types.ConvertToStringMap(elem)
			if !ok {
				continue
			}
			val, exists := obj[subfield]
			if exists && types.Equals(val, matchVal) {
				filtered = append(filtered, elem)
			}
		}
		return filtered, nil

	default:
		return nil, errors.NewFunctionCallError(fmt.Sprintf("unknown array function '%s'", functionName), 0, 0)
	}
}
