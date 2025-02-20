package env

import (
	libraries2 "github.com/RyanCopley/expression-parser/pkg/env/libraries"
)

// Environment holds the available libraries.
type Environment struct {
	Libraries map[string]ILibrary
}

// NewEnvironment creates a new Environment with default libraries.
func NewEnvironment() *Environment {
	env := &Environment{Libraries: make(map[string]ILibrary)}
	env.Libraries["time"] = libraries2.NewTimeLib()
	env.Libraries["math"] = libraries2.NewMathLib()
	env.Libraries["string"] = libraries2.NewStringLib()
	env.Libraries["regex"] = libraries2.NewRegexLib()
	env.Libraries["array"] = libraries2.NewArrayLib()
	env.Libraries["cond"] = libraries2.NewCondLib()
	env.Libraries["type"] = libraries2.NewTypeLib()
	return env
}

// GetLibrary retrieves a library by name.
func (e *Environment) GetLibrary(name string) (ILibrary, bool) {
	lib, ok := e.Libraries[name]
	return lib, ok
}
