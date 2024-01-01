package cli

import (
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strings"
)

type Command struct {
	Name         Token
	Parent       *Command
	ExecFunc     ExecFunc
	Flags        Flags
	Args         Args
	SubCommands  CommandMap
	invokedFlags Flags
}

func NewCommand(name Token, ef ExecFunc) *Command {
	return &Command{
		Name:        name,
		ExecFunc:    ef,
		Flags:       make(Flags, 0),
		Args:        make(Args, 0),
		SubCommands: make(CommandMap),
	}
}

func (c *Command) FullName() (name string) {
	if c.Parent == nil {
		name = string(c.Name)
		goto end
	}
	name = fmt.Sprintf("%s %s", c.Parent.FullName(), c.Name)
	if name[0] == ' ' {
		name = name[1:]
	}
end:

	return name
}

func (c *Command) IsLeaf() bool {
	return len(c.SubCommands) == 0
}

// MeetsRequirements validates args and options passed on the CLI for full command path.
func (c *Command) MeetsRequirements(ctx Context, tokenCnt int) (err error) {
	var expected, got int
	var cmds []*Command

	slog.Info("Validating CLI Args and Flags")

	if c.ExecFunc == nil {
		if c == RootCmd {
			err = ErrNoCommandSpecified
			goto end
		}
		// For when using a partial command like 'add' when the command is 'add project'.
		err = ErrNoExecFuncFound
		goto end
	}

	cmds = c.commandPath()
	for _, cmd := range cmds {
		err = cmd.meetsRequirements(ctx, tokenCnt)
		if err != nil {
			goto end
		}
	}

	expected = c.RequiredArgsCount()
	got = c.ReceivedArgsCount(tokenCnt)
	if got < expected {
		// TODO: Add 'missing'
		err = ErrTooFewArgsPassed.Args("expected", expected, "passed", got)
		goto end
	}
	expected = c.DeclaredArgsCount()
	if got > expected {
		// TODO: Add 'extra'
		err = ErrTooManyArgsPassed.Args("expected", expected, "passed", got)
		goto end
	}

end:
	return err
}

type TokenType string

const (
	ArgType  TokenType = "arg"
	FlagType TokenType = "option"
)

// meetsRequirements validates all args and options meet requirements for one command.
func (c *Command) meetsRequirements(ctx Context, tokenCnt int) (err error) {

	err = MeetsRequirements(ctx, ArgType, c.Args)
	if err != nil {
		goto end
	}

	err = MeetsRequirements(ctx, FlagType, c.Flags)
	if err != nil {
		goto end
	}

end:
	return err
}

func (c *Command) Help() string {
	var helpCmds []string
	helpCmds = make([]string, 0)
	for _, subCmd := range c.SubCommands {
		//if !c.IsLeaf() {
		//	continue
		//}
		helpCmds = append(helpCmds, subCmd.subHelp()...)
	}
	slices.Sort(helpCmds)
	return strings.Join(helpCmds, "")
}

func (c *Command) subHelp() (help []string) {
	if len(c.SubCommands) == 0 {
		help = []string{c.SignatureHelp()}
		goto end
	}
	for _, subCmd := range c.SubCommands {
		//if !c.IsLeaf() {
		//	continue
		//}
		help = append(help, subCmd.subHelp()...)
	}
end:
	return help
}

func (c *Command) SignatureHelp() string {
	var sb = strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s%s- %s", Indent, Indent, c.FullName()))
	sb.WriteString(c.Flags.SignatureHelp())
	sb.WriteString(c.Args.SignatureHelp())
	sb.WriteByte('\n')
	return sb.String()
}

func (c *Command) AddSubCommand(name Token, ef ExecFunc) (cmd *Command) {
	cmd = NewCommand(name, ef)
	cmd.Parent = c
	c.SubCommands[name] = cmd
	return cmd
}

func (c *Command) AddFlag(flg Flag) (cmd *Command) {
	flg.Parent = c
	if flg.Default == nil {
		switch flg.Type {
		case reflect.String:
			flg.Default = ""
		case reflect.Int:
			flg.Default = 0
		default:
			flg.noSetFuncAssigned()
		}
	}
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
		cmdNames[0] = string(cmd.Name)
	}
	if len(cmd.SubCommands) == 0 {
		s = cmdNames[0]
		goto end
	}
	for c.Parent != nil {
		cmd = c.Parent
		cmdNames = append(cmdNames, string(cmd.Name))
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
	return string(c.Name)
}

// commandDepth returns how deep the command is.
// e.g. `myapp -a 10 -b hello foo bar baz` would be commandDepth 3 for `foo bar baz`
func (c *Command) commandDepth() (n int) {
	p := c
	for p.Parent != nil {
		p = p.Parent
		n++
	}
	return n
}

// setArgValues sets the Value values
func (c *Command) setArgValues(args Tokens) (err error) {
	var index, depth int
	depth = c.commandDepth()
	if depth <= len(args) {
		args = args[depth:]
	}
	// Loop through all args defined for this command
	for _, arg := range c.Args {
		if index >= len(args) {
			goto end
		}
		// If we received the arg on the CLI then assign it
		arg.Value = NewValue(arg.Type, args[index])
		c.Args[index] = arg
		index++
	}
end:
	return err
}

// callSetArgValueFuncs calls the SetValueFunc for each arg
func (c *Command) callSetArgValueFuncs() (err error) {
	for i, arg := range c.Args {
		c.Args[i] = callSetArgValueFunc(arg)
	}
	return err
}

// RequiredArgsCount returns the number of required args
func (c *Command) RequiredArgsCount() (cnt int) {
	for _, arg := range c.Args {
		if arg.Optional {
			continue
		}
		cnt++
	}
	return cnt
}

// OptionalArgsCount returns the number of optional args
func (c *Command) OptionalArgsCount() (cnt int) {
	return c.DeclaredArgsCount() - c.RequiredArgsCount()
}

// DeclaredArgsCount returns the number of total args; required and optional
func (c *Command) DeclaredArgsCount() (cnt int) {
	return len(c.Args)
}

func (c *Command) AddArg(arg Arg) (cmd *Command) {
	arg.Parent = c
	c.Args = append(c.Args, NewArg(arg))
	return c
}

func (c *Command) SetFlags(flags Flags) Flags {
	unsetFlags := flags
	for i, flag := range flags {
		n := c.Flags.Index(flag.Name)
		if n != -1 {
			c.Flags[n] = flag
			unsetFlags = unsetFlags.Remove(i)
			continue
		}
		if c.Parent == nil {
			continue
		}
		unsetFlags = c.Parent.SetFlags(unsetFlags)
	}
	return unsetFlags
}

// InvokedFlags returns all the flags for th invoked command, including all
// parent flags including the root flags.
func (c *Command) InvokedFlags() (flags Flags) {
	var cmds []*Command
	var cmd *Command

	if c.invokedFlags != nil {
		goto end
	}
	cmds = []*Command{c}
	cmd = c
	for cmd.Parent != nil {
		cmds = append(cmds, cmd.Parent)
		cmd = cmd.Parent
	}
	slices.Reverse(cmds)
	for _, cmd = range cmds {
		if len(cmd.Flags) == 0 {
			continue
		}
		flags = append(flags, cmd.Flags...)
	}
end:
	return flags
}
