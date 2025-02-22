package env

import "github.com/SpecDrivenDesign/lql/pkg/param"

// ILibrary is the interface for DSL libraries.
type ILibrary interface {
	Call(functionName string, args []param.Arg, line, column, parenLine, parenColumn int) (interface{}, error)
}
