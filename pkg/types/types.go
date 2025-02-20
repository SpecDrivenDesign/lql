package types

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/errors"
	"math"
	"strconv"
	"strings"
)

// ToFloat converts a numeric value to a float64.
func ToFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	}
	return 0, false
}

// ToInt converts a numeric value to an int64.
func ToInt(val interface{}) (int64, bool) {
	switch v := val.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case float64:
		return int64(v), true
	}
	return 0, false
}

// IsInt checks if a value is of an integer type.
func IsInt(val interface{}) bool {
	switch val.(type) {
	case int, int64:
		return true
	}
	return false
}

// Equals compares two values for equality.
func Equals(left, right interface{}) bool {
	lf, lok := ToFloat(left)
	rf, rok := ToFloat(right)
	if lok && rok {
		return math.Abs(lf-rf) < 1e-9
	}
	return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
}

// Compare compares two values using the given operator.
func Compare(left, right interface{}, op string, line, column int) (bool, error) {
	lf, lok := ToFloat(left)
	rf, rok := ToFloat(right)
	if lok && rok {
		switch op {
		case "<":
			return lf < rf, nil
		case ">":
			return lf > rf, nil
		case "<=":
			return lf <= rf, nil
		case ">=":
			return lf >= rf, nil
		}
	}
	ls, lok := left.(string)
	rs, rok := right.(string)
	if lok && rok {
		switch op {
		case "<":
			return ls < rs, nil
		case ">":
			return ls > rs, nil
		case "<=":
			return ls <= rs, nil
		case ">=":
			return ls >= rs, nil
		}
	}
	return false, errors.NewSemanticError(fmt.Sprintf("'%s' operator not allowed on given types", op), line, column)
}

// ParseNumber parses a numeric literal string.
func ParseNumber(lit string) interface{} {
	if strings.ContainsAny(lit, ".eE") {
		f, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			return 0.0
		}
		return f
	} else {
		i, err := strconv.ParseInt(lit, 10, 64)
		if err != nil {
			return int64(0)
		}
		return i
	}
}

// ConvertToInterfaceSlice converts various slice types to []interface{}.
func ConvertToInterfaceSlice(val interface{}) ([]interface{}, bool) {
	switch v := val.(type) {
	case []interface{}:
		return v, true
	case []int:
		s := make([]interface{}, len(v))
		for i, e := range v {
			s[i] = e
		}
		return s, true
	case []int64:
		s := make([]interface{}, len(v))
		for i, e := range v {
			s[i] = e
		}
		return s, true
	case []float64:
		s := make([]interface{}, len(v))
		for i, e := range v {
			s[i] = e
		}
		return s, true
	}
	return nil, false
}

// ConvertToStringMap converts various map types to map[string]interface{}.
func ConvertToStringMap(val interface{}) (map[string]interface{}, bool) {
	switch v := val.(type) {
	case map[string]interface{}:
		return v, true
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for key, value := range v {
			m[fmt.Sprintf("%v", key)] = value
		}
		return m, true
	}
	return nil, false
}
