package cli

import (
	"slices"
	"strings"
)

type Command struct {
	Name        string
	Parent      *Command
	ExecFunc    ExecFunc
	Flags       Flags
	Args        []string
	OptArgs     []string
	SubCommands CommandMap
}

func NewCommand(name string, ef ExecFunc) *Command {
	return &Command{
		Name:        name,
		ExecFunc:    ef,
		Flags:       make(Flags, 0),
		Args:        make([]string, 0),
		OptArgs:     make([]string, 0),
		SubCommands: make(CommandMap),
	}
}

func (c *Command) ExecuteFunc(args ...string) error {
	return c.ExecFunc(args...)
}

func (c *Command) AddArgs(args ...string) (cmd *Command) {
	c.Args = args
	return c
}
func (c *Command) AddOptArgs(args ...string) (cmd *Command) {
	c.OptArgs = args
	return c
}

func (c *Command) AddSubCommand(name string, ef ExecFunc) (cmd *Command) {
	cmd = NewCommand(name, ef)
	cmd.Parent = c
	c.SubCommands[name] = cmd
	return cmd
}

func (c *Command) AddFlag(flg *Flag) (cmd *Command) {
	flg.Parent = c
	c.Flags = append(c.Flags, flg)
	return c
}

// Unique returns the unique name for a command which includes its ancestor
// commands, e.g.:
//
//   - `help`
//   - `add codebase`
//   - `add project`
//   - `map`
//   - `foo bar baz`
func (c *Command) Unique() (s string) {
	sb := strings.Builder{}
	cmd := c
	cmdNames := make([]string, 1)
	if len(cmd.Name) == 0 {
		cmdNames[0] = "root"
	} else {
		cmdNames[0] = cmd.Name
	}
	if len(cmd.SubCommands) == 0 {
		goto end
	}
	for c.Parent != nil {
		cmd = c.Parent
		cmdNames = append(cmdNames, cmd.Name)
	}
	slices.Reverse(cmdNames)
	for _, name := range cmdNames {
		sb.WriteString(name)
		sb.WriteByte('-')
	}
	s = sb.String()
	s = s[:len(s)-1]
end:
	return s
}

func (c *Command) String() string {
	sb := strings.Builder{}
	sb.WriteString(c.Name)
	if len(c.Args) > 0 {
		sb.WriteString(" <")
		sb.WriteString(strings.Join(c.Args, "> <"))
		sb.WriteByte('>')
	}
	if len(c.OptArgs) > 0 {
		var i int
		var arg string
		for i, arg = range c.OptArgs {
			sb.WriteString(" [<")
			sb.WriteString(arg)
			sb.WriteByte('>')
		}
		sb.WriteString(strings.Repeat("]", i+1))
	}
	if len(c.AllFlags()) > 0 && c.Name != "help" {
		for _, flg := range c.AllFlags() {
			sb.WriteString(" [-")
			sb.WriteString(flg.Name)
			sb.WriteString("=<")
			sb.WriteString(flg.VarName)
			sb.WriteString(">]")
		}
	}
	return sb.String()
}
