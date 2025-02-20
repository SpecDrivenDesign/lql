package errors

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	stdErrors "errors"
)

// PositionalError interface for errors that include positional information.
type PositionalError interface {
	error
	GetLine() int
	GetColumn() int
	Kind() string
}

// TypeError
type TypeError struct {
	Msg    string
	Line   int
	Column int
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("TypeError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *TypeError) GetLine() int   { return e.Line }
func (e *TypeError) GetColumn() int { return e.Column }
func (e *TypeError) Kind() string   { return "TypeError" }

func NewTypeError(msg string, line, column int) error {
	return &TypeError{Msg: msg, Line: line, Column: column}
}

// DivideByZeroError
type DivideByZeroError struct {
	Msg    string
	Line   int
	Column int
}

func (e *DivideByZeroError) Error() string {
	return fmt.Sprintf("DivideByZeroError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *DivideByZeroError) GetLine() int   { return e.Line }
func (e *DivideByZeroError) GetColumn() int { return e.Column }
func (e *DivideByZeroError) Kind() string   { return "DivideByZeroError" }

func NewDivideByZeroError(msg string, line, column int) error {
	return &DivideByZeroError{Msg: msg, Line: line, Column: column}
}

// ReferenceError
type ReferenceError struct {
	Msg    string
	Line   int
	Column int
}

func (e *ReferenceError) Error() string {
	return fmt.Sprintf("ReferenceError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *ReferenceError) GetLine() int   { return e.Line }
func (e *ReferenceError) GetColumn() int { return e.Column }
func (e *ReferenceError) Kind() string   { return "ReferenceError" }

func NewReferenceError(msg string, line, column int) error {
	return &ReferenceError{Msg: msg, Line: line, Column: column}
}

// UnknownIdentifierError
type UnknownIdentifierError struct {
	Msg    string
	Line   int
	Column int
}

func (e *UnknownIdentifierError) Error() string {
	return fmt.Sprintf("UnknownIdentifierError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *UnknownIdentifierError) GetLine() int   { return e.Line }
func (e *UnknownIdentifierError) GetColumn() int { return e.Column }
func (e *UnknownIdentifierError) Kind() string   { return "UnknownIdentifierError" }

func NewUnknownIdentifierError(msg string, line, column int) error {
	return &UnknownIdentifierError{Msg: msg, Line: line, Column: column}
}

// UnknownOperatorError
type UnknownOperatorError struct {
	Msg    string
	Line   int
	Column int
}

func (e *UnknownOperatorError) Error() string {
	return fmt.Sprintf("UnknownOperatorError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *UnknownOperatorError) GetLine() int   { return e.Line }
func (e *UnknownOperatorError) GetColumn() int { return e.Column }
func (e *UnknownOperatorError) Kind() string   { return "UnknownOperatorError" }

func NewUnknownOperatorError(msg string, line, column int) error {
	return &UnknownOperatorError{Msg: msg, Line: line, Column: column}
}

// FunctionCallError
type FunctionCallError struct {
	Msg    string
	Line   int
	Column int
}

func (e *FunctionCallError) Error() string {
	return fmt.Sprintf("FunctionCallError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *FunctionCallError) GetLine() int   { return e.Line }
func (e *FunctionCallError) GetColumn() int { return e.Column }
func (e *FunctionCallError) Kind() string   { return "FunctionCallError" }

func NewFunctionCallError(msg string, line, column int) error {
	return &FunctionCallError{Msg: msg, Line: line, Column: column}
}

// ParameterError
type ParameterError struct {
	Msg    string
	Line   int
	Column int
}

func (e *ParameterError) Error() string {
	return fmt.Sprintf("ParameterError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *ParameterError) GetLine() int   { return e.Line }
func (e *ParameterError) GetColumn() int { return e.Column }
func (e *ParameterError) Kind() string   { return "ParameterError" }

func NewParameterError(msg string, line, column int) error {
	return &ParameterError{Msg: msg, Line: line, Column: column}
}

// LexicalError
type LexicalError struct {
	Msg    string
	Line   int
	Column int
}

func (e *LexicalError) Error() string {
	return fmt.Sprintf("LexicalError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *LexicalError) GetLine() int   { return e.Line }
func (e *LexicalError) GetColumn() int { return e.Column }
func (e *LexicalError) Kind() string   { return "LexicalError" }

func NewLexicalError(msg string, line, column int) error {
	return &LexicalError{Msg: msg, Line: line, Column: column}
}

// SyntaxError
type SyntaxError struct {
	Msg    string
	Line   int
	Column int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *SyntaxError) GetLine() int   { return e.Line }
func (e *SyntaxError) GetColumn() int { return e.Column }
func (e *SyntaxError) Kind() string   { return "SyntaxError" }

func NewSyntaxError(msg string, line, column int) error {
	return &SyntaxError{Msg: msg, Line: line, Column: column}
}

// SemanticError
type SemanticError struct {
	Msg    string
	Line   int
	Column int
}

func (e *SemanticError) Error() string {
	return fmt.Sprintf("SemanticError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *SemanticError) GetLine() int   { return e.Line }
func (e *SemanticError) GetColumn() int { return e.Column }
func (e *SemanticError) Kind() string   { return "SemanticError" }

func NewSemanticError(msg string, line, column int) error {
	return &SemanticError{Msg: msg, Line: line, Column: column}
}

// ArrayOutOfBoundsError
type ArrayOutOfBoundsError struct {
	Msg    string
	Line   int
	Column int
}

func (e *ArrayOutOfBoundsError) Error() string {
	return fmt.Sprintf("ArrayOutOfBoundsError: %s at line %d, column %d", e.Msg, e.Line, e.Column)
}

func (e *ArrayOutOfBoundsError) GetLine() int   { return e.Line }
func (e *ArrayOutOfBoundsError) GetColumn() int { return e.Column }
func (e *ArrayOutOfBoundsError) Kind() string   { return "ArrayOutOfBoundsError" }

func NewArrayOutOfBoundsError(msg string, line, column int) error {
	return &ArrayOutOfBoundsError{Msg: msg, Line: line, Column: column}
}

// GetErrorContext returns a formatted error context string showing the line and a pointer to the error column.
func GetErrorContext(expr string, errLine, errColumn int, colored bool) string {
	lines := strings.Split(expr, "\n")
	if errLine-1 < 0 || errLine-1 >= len(lines) {
		return ""
	}
	lineText := lines[errLine-1]
	if errColumn > len(lineText) {
		errColumn = len(lineText)
	}
	pointer := ""
	for i := 0; i < errColumn-1 && i < len(lineText); i++ {
		if lineText[i] == '\t' {
			pointer += "\t"
		} else {
			pointer += "-"
		}
	}
	pointer += "^"
	if colored {
		pointer = "\033[31m" + pointer + "\033[0m"
	}
	return fmt.Sprintf("    %s\n    %s", lineText, pointer)
}

// GetErrorPosition attempts to extract the line and column from an error.
func GetErrorPosition(err error) (int, int) {
	type positioner interface {
		Position() (int, int)
	}
	if pe, ok := err.(positioner); ok {
		return pe.Position()
	}
	var ep PositionalError
	if stdErrors.As(err, &ep) {
		return ep.GetLine(), ep.GetColumn()
	}
	v := reflect.ValueOf(err)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	fieldLine := v.FieldByName("Line")
	fieldColumn := v.FieldByName("Column")
	if fieldLine.IsValid() && fieldColumn.IsValid() && fieldLine.CanInt() && fieldColumn.CanInt() {
		return int(fieldLine.Int()), int(fieldColumn.Int())
	}
	r := regexp.MustCompile(`at line (\d+), column (\d+)`)
	matches := r.FindStringSubmatch(err.Error())
	if len(matches) == 3 {
		line, err1 := strconv.Atoi(matches[1])
		col, err2 := strconv.Atoi(matches[2])
		if err1 == nil && err2 == nil {
			return line, col
		}
	}
	return 0, 0
}
