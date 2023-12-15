package cli

type CommandType int

const (
	Undefined CommandType = iota
	Root
	Branch
	Leaf
)
