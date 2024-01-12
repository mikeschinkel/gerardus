package cli

type CommandType int

const (
	Undefined CommandType = iota
	RootCommand
	BranchCommand
	LeafCommand
	HelpCommand
)
