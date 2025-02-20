package env

import "github.com/RyanCopley/expression-parser/pkg/param"

// ILibrary is the interface for DSL libraries.
type ILibrary interface {
	Call(functionName string, args []param.Arg, line, column, parenLine, parenColumn int) (interface{}, error)
}
