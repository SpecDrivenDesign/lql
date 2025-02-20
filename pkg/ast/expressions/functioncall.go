package expressions

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/ast"
	"github.com/RyanCopley/expression-parser/pkg/env"
	"github.com/RyanCopley/expression-parser/pkg/errors"
	"github.com/RyanCopley/expression-parser/pkg/param"
	"strings"
)

// FunctionCallExpr represents a function call.
type FunctionCallExpr struct {
	Namespace   []string
	Args        []ast.Expression
	Line        int
	Column      int
	ParenLine   int
	ParenColumn int
}

func (f *FunctionCallExpr) Eval(ctx map[string]interface{}, env *env.Environment) (interface{}, error) {
	if len(f.Namespace) < 2 {
		return nil, errors.NewParameterError("function call missing namespace", f.Line, f.Column)
	}
	libName := f.Namespace[0]
	funcName := f.Namespace[1]
	lib, ok := env.GetLibrary(libName)
	if !ok {
		return nil, errors.NewReferenceError(fmt.Sprintf("library '%s' not found", libName), f.Line, f.Column)
	}
	var args []param.Arg
	for _, argExpr := range f.Args {
		val, err := argExpr.Eval(ctx, env)
		if err != nil {
			return nil, err
		}
		l, c := argExpr.Pos()
		args = append(args, param.Arg{Value: val, Line: l, Column: c})
	}
	return lib.Call(funcName, args, f.Line, f.Column, f.ParenLine, f.ParenColumn)
}

func (f *FunctionCallExpr) Pos() (int, int) {
	return f.Line, f.Column
}
func (f *FunctionCallExpr) String() string {
	var sb strings.Builder

	if len(f.Namespace) == 0 {
		return "(missing function call)"
	}

	// The first item in the Namespace is the "library" name.
	libraryName := f.Namespace[0]
	if ColorEnabled {
		libraryName = LibraryColor + libraryName + ColorReset
	}

	// If there is more than one item, the rest are the "function" name(s).
	// We'll join them with '.' in a single string and color them all as FunctionColor.
	var functionName string
	if len(f.Namespace) > 1 {
		rest := f.Namespace[1:]
		fnStr := strings.Join(rest, ".")
		if ColorEnabled {
			fnStr = FunctionColor + fnStr + ColorReset
		}

		// Insert a "." (punctuation) between library and function portion
		dot := "."
		if ColorEnabled {
			dot = PunctuationColor + "." + ColorReset
		}
		functionName = dot + fnStr
	}

	// parentheses and commas
	openParen := "("
	closeParen := ")"
	comma := ", "
	if ColorEnabled {
		openParen = PunctuationColor + "(" + ColorReset
		closeParen = PunctuationColor + ")" + ColorReset
		comma = PunctuationColor + "," + ColorReset + " "
	}

	// Write out library + function portion
	sb.WriteString(libraryName)
	sb.WriteString(functionName)

	sb.WriteString(openParen)

	for i, arg := range f.Args {
		if i > 0 {
			sb.WriteString(comma)
		}
		sb.WriteString(arg.String())
	}

	sb.WriteString(closeParen)
	return sb.String()
}
