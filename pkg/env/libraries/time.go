package libraries

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/param"
	"github.com/RyanCopley/expression-parser/pkg/types"
	"strconv"
	"time"

	"github.com/RyanCopley/expression-parser/pkg/errors"
)

// TimeValue represents a time value.
type TimeValue struct {
	EpochMillis int64
	Zone        string
}

func newTimeValue(t time.Time) TimeValue {
	return TimeValue{
		EpochMillis: t.UnixNano() / int64(time.Millisecond),
		Zone:        t.Location().String(),
	}
}

// TimeLib implements time-related functions.
type TimeLib struct{}

func NewTimeLib() *TimeLib {
	return &TimeLib{}
}

func (t *TimeLib) Call(functionName string, args []param.Arg, line, col, parenLine, parenCol int) (interface{}, error) {
	switch functionName {
	case "now":
		if len(args) != 0 {
			return nil, errors.NewParameterError("time.now() takes no arguments", line, col)
		}
		now := time.Now().UTC()
		return newTimeValue(now), nil

	case "parse":
		if len(args) < 2 {
			if len(args) == 0 {
				return nil, errors.NewParameterError("time.parse requires at least 2 arguments", parenLine, parenCol)
			}
			lastArg := args[len(args)-1]
			return nil, errors.NewParameterError("time.parse requires at least 2 arguments", lastArg.Line, lastArg.Column)
		}
		arg0 := args[0]
		inputStr, ok := arg0.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("time.parse: first argument must be a string", arg0.Line, arg0.Column)
		}
		arg1 := args[1]
		format, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("time.parse: second argument must be a string", arg1.Line, arg1.Column)
		}
		var tTime time.Time
		var err error
		switch format {
		case "iso8601":
			tTime, err = time.Parse(time.RFC3339Nano, inputStr)
		case "dateOnly":
			tTime, err = time.Parse("2006-01-02", inputStr)
		case "epochMillis":
			ms, err2 := strconv.ParseInt(inputStr, 10, 64)
			if err2 != nil {
				return nil, errors.NewTypeError("time.parse: invalid epochMillis", arg0.Line, arg0.Column)
			}
			return TimeValue{EpochMillis: ms, Zone: "UTC"}, nil
		case "rfc2822":
			tTime, err = time.Parse(time.RFC1123Z, inputStr)
		case "custom":
			if len(args) != 3 {
				return nil, errors.NewParameterError("time.parse with 'custom' requires a formatDetails argument", line, col)
			}
			arg2 := args[2]
			formatDetails, ok := arg2.Value.(string)
			if !ok {
				return nil, errors.NewTypeError("time.parse: formatDetails must be a string", arg2.Line, arg2.Column)
			}
			tTime, err = time.Parse(formatDetails, inputStr)
		default:
			return nil, errors.NewTypeError("time.parse: unknown format", arg1.Line, arg1.Column)
		}
		if err != nil {
			return nil, errors.NewTypeError("time.parse error: "+err.Error(), arg0.Line, arg0.Column)
		}
		return newTimeValue(tTime.UTC()), nil

	case "add":
		if len(args) != 2 {
			return nil, errors.NewParameterError("time.add requires 2 arguments", line, col)
		}
		arg0 := args[0]
		tv, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.add: first argument must be Time", arg0.Line, arg0.Column)
		}
		arg1 := args[1]
		dur, ok := types.ToInt(arg1.Value)
		if !ok {
			return nil, errors.NewTypeError("time.add: second argument must be numeric", arg1.Line, arg1.Column)
		}
		return TimeValue{EpochMillis: tv.EpochMillis + dur, Zone: tv.Zone}, nil

	case "subtract":
		if len(args) != 2 {
			return nil, errors.NewParameterError("time.subtract requires 2 arguments", line, col)
		}
		arg0 := args[0]
		tv, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.subtract: first argument must be Time", arg0.Line, arg0.Column)
		}
		arg1 := args[1]
		dur, ok := types.ToInt(arg1.Value)
		if !ok {
			return nil, errors.NewTypeError("time.subtract: second argument must be numeric", arg1.Line, arg1.Column)
		}
		return TimeValue{EpochMillis: tv.EpochMillis - dur, Zone: tv.Zone}, nil

	case "diff":
		if len(args) != 2 {
			return nil, errors.NewParameterError("time.diff requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arg1 := args[1]
		tv1, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.diff: first argument must be Time", arg0.Line, arg0.Column)
		}
		tv2, ok := arg1.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.diff: second argument must be Time", arg1.Line, arg1.Column)
		}
		return tv1.EpochMillis - tv2.EpochMillis, nil

	case "isBefore":
		if len(args) != 2 {
			return nil, errors.NewParameterError("time.isBefore requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arg1 := args[1]
		tv1, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.isBefore: first argument must be Time", arg0.Line, arg0.Column)
		}
		tv2, ok := arg1.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.isBefore: second argument must be Time", arg1.Line, arg1.Column)
		}
		return tv1.EpochMillis < tv2.EpochMillis, nil

	case "isAfter":
		if len(args) != 2 {
			return nil, errors.NewParameterError("time.isAfter requires 2 arguments", line, col)
		}
		arg0 := args[0]
		arg1 := args[1]
		tv1, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.isAfter: first argument must be Time", arg0.Line, arg0.Column)
		}
		tv2, ok := arg1.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.isAfter: second argument must be Time", arg1.Line, arg1.Column)
		}
		return tv1.EpochMillis > tv2.EpochMillis, nil

	case "isEqual":
		if len(args) != 2 {
			return nil, errors.NewParameterError("time.isEqual requires 2 arguments", line, col)
		}
		tv1, ok1 := args[0].Value.(TimeValue)
		if !ok1 {
			return nil, errors.NewTypeError("time.isEqual: first argument must be Time", args[0].Line, args[0].Column)
		}
		tv2, ok2 := args[1].Value.(TimeValue)
		if !ok2 {
			return nil, errors.NewTypeError("time.isEqual: second argument must be Time", args[1].Line, args[1].Column)
		}
		return tv1.EpochMillis == tv2.EpochMillis, nil

	case "toEpochMillis":
		if len(args) != 1 {
			return nil, errors.NewParameterError("time.toEpochMillis requires 1 argument", line, col)
		}
		arg0 := args[0]
		tv, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.toEpochMillis: argument must be Time", arg0.Line, arg0.Column)
		}
		return tv.EpochMillis, nil

	case "format":
		if len(args) != 2 {
			return nil, errors.NewParameterError("time.format requires 2 arguments", line, col)
		}
		arg0 := args[0]
		tv, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.format: first argument must be Time", arg0.Line, arg0.Column)
		}
		arg1 := args[1]
		formatStr, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("time.format: second argument must be string", arg1.Line, arg1.Column)
		}
		loc, err := time.LoadLocation(tv.Zone)
		if err != nil {
			loc = time.UTC
		}
		tTime := time.Unix(0, tv.EpochMillis*int64(time.Millisecond)).In(loc)
		return tTime.Format(formatStr), nil

	case "getYear":
		if len(args) != 1 {
			return nil, errors.NewParameterError("time.getYear requires 1 argument", line, col)
		}
		arg0 := args[0]
		tv, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.getYear: argument must be Time", arg0.Line, arg0.Column)
		}
		loc, err := time.LoadLocation(tv.Zone)
		if err != nil {
			loc = time.UTC
		}
		tTime := time.Unix(0, tv.EpochMillis*int64(time.Millisecond)).In(loc)
		return int64(tTime.Year()), nil

	case "getMonth":
		if len(args) != 1 {
			return nil, errors.NewParameterError("time.getMonth requires 1 argument", line, col)
		}
		arg0 := args[0]
		tv, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.getMonth: argument must be Time", arg0.Line, arg0.Column)
		}
		loc, err := time.LoadLocation(tv.Zone)
		if err != nil {
			loc = time.UTC
		}
		tTime := time.Unix(0, tv.EpochMillis*int64(time.Millisecond)).In(loc)
		return int64(tTime.Month()), nil

	case "getDay":
		if len(args) != 1 {
			return nil, errors.NewParameterError("time.getDay requires 1 argument", line, col)
		}
		arg0 := args[0]
		tv, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.getDay: argument must be Time", arg0.Line, arg0.Column)
		}
		loc, err := time.LoadLocation(tv.Zone)
		if err != nil {
			loc = time.UTC
		}
		tTime := time.Unix(0, tv.EpochMillis*int64(time.Millisecond)).In(loc)
		return int64(tTime.Day()), nil

	case "startOfDay":
		if len(args) != 1 {
			return nil, errors.NewParameterError("time.startOfDay requires 1 argument", line, col)
		}
		arg0 := args[0]
		tv, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.startOfDay: argument must be Time", arg0.Line, arg0.Column)
		}
		loc, err := time.LoadLocation(tv.Zone)
		if err != nil {
			loc = time.UTC
		}
		tTime := time.Unix(0, tv.EpochMillis*int64(time.Millisecond)).In(loc)
		start := time.Date(tTime.Year(), tTime.Month(), tTime.Day(), 0, 0, 0, 0, loc)
		return newTimeValue(start), nil

	case "endOfDay":
		if len(args) != 1 {
			return nil, errors.NewParameterError("time.endOfDay requires 1 argument", line, col)
		}
		arg0 := args[0]
		tv, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.endOfDay: argument must be Time", arg0.Line, arg0.Column)
		}
		loc, err := time.LoadLocation(tv.Zone)
		if err != nil {
			loc = time.UTC
		}
		tTime := time.Unix(0, tv.EpochMillis*int64(time.Millisecond)).In(loc)
		end := time.Date(tTime.Year(), tTime.Month(), tTime.Day(), 23, 59, 59, int(time.Millisecond*999), loc)
		return newTimeValue(end), nil

	case "withZone":
		if len(args) != 2 {
			return nil, errors.NewParameterError("time.withZone requires 2 arguments", line, col)
		}
		arg0 := args[0]
		tv, ok := arg0.Value.(TimeValue)
		if !ok {
			return nil, errors.NewTypeError("time.withZone: first argument must be Time", arg0.Line, arg0.Column)
		}
		arg1 := args[1]
		zoneName, ok := arg1.Value.(string)
		if !ok {
			return nil, errors.NewTypeError("time.withZone: second argument must be string", arg1.Line, arg1.Column)
		}
		loc, err := time.LoadLocation(zoneName)
		if err != nil {
			return nil, errors.NewTypeError("time.withZone: invalid zone name", arg1.Line, arg1.Column)
		}
		return TimeValue{EpochMillis: tv.EpochMillis, Zone: loc.String()}, nil

	default:
		return nil, errors.NewFunctionCallError(fmt.Sprintf("unknown time function '%s'", functionName), 0, 0)
	}
}
