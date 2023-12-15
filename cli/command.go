package cli

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type Command struct {
	Name         string
	Parent       *Command
	ExecFunc     ExecFunc
	Flags        Flags
	Args         Args
	SubCommands  CommandMap
	invokedFlags Flags
	argsMap      ArgsMap
}

func NewCommand(name string, ef ExecFunc) *Command {
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
		name = c.Name
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

func (c *Command) AddSubCommand(name string, ef ExecFunc) (cmd *Command) {
	cmd = NewCommand(name, ef)
	cmd.Parent = c
	c.SubCommands[name] = cmd
	return cmd
}

func (c *Command) AddFlag(flg *Flag) (cmd *Command) {
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
		cmdNames[0] = cmd.Name
	}
	if len(cmd.SubCommands) == 0 {
		s = cmdNames[0]
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
	return c.Name
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

//func (c *Command) ArgsMap(args []string) (_ ArgsMap, err error) {
//	var index, depth int
//
//	if len(c.argsMap) > 0 {
//		goto end
//	}
//	c.argsMap = make(ArgsMap)
//
//	depth = c.commandDepth()
//	if depth >= len(args) {
//		goto end
//	}
//	args = args[1+depth:]
//
//	for _, arg := range c.Args {
//		if index < len(args) {
//			value := args[index]
//			if value[0] == '-' {
//				continue
//			}
//			index++
//			arg.Value.string = value
//		}
//		arg.Value.Type = arg.Type
//		c.argsMap[arg.Name] = arg
//	}
//end:
//	return c.argsMap, err
//}

// callSetArgValues sets the Value values
func (c *Command) callSetArgValues(args []string) (err error) {
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
		index++
	}
end:
	return err
}

// callSetArgValueFuncs calls the SetValueFunc for each arg
func (c *Command) callSetArgValueFuncs() (err error) {
	for _, arg := range c.Args {
		arg.callSetValueFunc()
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

func (c *Command) AddArg(arg *Arg) (cmd *Command) {
	arg.Parent = c
	c.Args = append(c.Args, NewArg(arg))
	return c
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
