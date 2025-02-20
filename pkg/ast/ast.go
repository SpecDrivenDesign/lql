package ast

import (
	"github.com/RyanCopley/expression-parser/pkg/env"
)

// Expression interface represents an AST node that can be evaluated.
type Expression interface {
	Eval(ctx map[string]interface{}, env *env.Environment) (interface{}, error)
	Pos() (int, int)
	String() string
}
