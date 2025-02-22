package expressions

import (
	"fmt"
	"github.com/SpecDrivenDesign/lql/pkg/env"
	"github.com/SpecDrivenDesign/lql/pkg/errors"
)

// IdentifierExpr represents an identifier.
type IdentifierExpr struct {
	Name   string
	Line   int
	Column int
}

func (i *IdentifierExpr) Eval(ctx map[string]interface{}, env *env.Environment) (interface{}, error) {
	return nil, errors.NewUnknownIdentifierError(fmt.Sprintf("Bare identifier '%s' is not allowed", i.Name), i.Line, i.Column)
}

func (i *IdentifierExpr) Pos() (int, int) {
	return i.Line, i.Column
}
func (i *IdentifierExpr) String() string {
	if ColorEnabled {
		return IdentifierColor + i.Name + ColorReset
	}
	return i.Name
}
