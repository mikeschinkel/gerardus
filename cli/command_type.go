package cli

type CommandType int

//goland:noinspection GoUnusedConst
const (
	Undefined CommandType = iota
	RootCommand
	BranchCommand
	LeafCommand
	HelpCommand
)
