package param

// Arg represents an argument passed to a library function.
type Arg struct {
	Value  interface{}
	Line   int
	Column int
}
