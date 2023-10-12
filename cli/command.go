package cli

import (
	"flag"
	"os"
	"slices"
	"strings"
)

type Command struct {
	Name          string
	Parent        *Command
	ExecFunc      ExecFunc
	Flags         Flags
	Args          Args
	SubCommands   CommandMap
	invokedFlags  Flags
	argValues     StringMap
	argMap        ArgsMap
	flagValuesMap FlagValuesMap
}

func NewCommand(name string, ef ExecFunc) *Command {
	return &Command{
		Name:          name,
		ExecFunc:      ef,
		Flags:         make(Flags, 0),
		Args:          make(Args, 0),
		SubCommands:   make(CommandMap),
		flagValuesMap: make(FlagValuesMap, 8),
	}
}

func (c *Command) ExecuteFunc(args StringMap) error {
	return c.ExecFunc(args)
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
		switch {
		case flg.SetStringValFunc != nil:
			flg.Default = ""
		case flg.SetIntValFunc != nil:
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
	sb := strings.Builder{}
	sb.WriteString(c.Name)
	//if len(c.InvokedFlags()) > 0 && c.Name != "help" {
	//	for _, flg := range c.InvokedFlags() {
	//		sb.WriteString(" [-")
	//		sb.WriteString(flg.Name)
	//		sb.WriteString("=<")
	//		sb.WriteString(flg.Name)
	//		sb.WriteString(">] ")
	//	}
	//}
	//if len(c.Args) > 0 {
	//	sb.WriteString(" <")
	//	sb.WriteString(strings.Join(c.Args, "> <"))
	//	sb.WriteByte('>')
	//}
	//if len(c.OptArgs) > 0 {
	//	var i int
	//	var arg string
	//	for i, arg = range c.OptArgs {
	//		sb.WriteString(" [<")
	//		sb.WriteString(arg)
	//		sb.WriteByte('>')
	//	}
	//	sb.WriteString(strings.Repeat("]", i+1))
	//}
	return sb.String()
}

// depth returns how deep the command is.
// e.g. `cliapp -a 10 -b hello foo bar baz` would be depth 3
func (c *Command) depth() (n int) {
	for c.Parent != nil {
		c = c.Parent
		n++
	}
	return n
}

func (c *Command) ArgValuesMap() (StringMap, ArgsMap) {
	var depth, index int
	var args Args

	if c.argValues != nil {
		goto end
	}
	depth = c.depth()
	args = c.Args

	c.argValues = make(StringMap)
	c.argMap = make(ArgsMap)
	depth--
	index = len(os.Args) - 1
	for depth >= 0 {
		if index <= 1 {
			goto end
		}
		value := os.Args[index]
		index--
		if value[0] == '-' {
			continue
		}
		arg := args[depth]
		name := arg.Name
		c.argValues[name] = value
		c.argMap[name] = arg
		depth--
	}
end:
	return c.argValues, c.argMap
}

// SetFlagValues the flag value specified by
func (c *Command) SetFlagValues() {
	for _, f := range c.InvokedFlags() {
		fv := flagValuesMap[f.Unique()]
		switch {
		case f.SetStringValFunc != nil:
			f.SetStringValFunc(fv.String)
		case f.SetIntValFunc != nil:
			f.SetIntValFunc(fv.Int)
		default:
			f.noSetFuncAssigned()
		}
	}

}

// AddFlags initializes the flag package flags
func (c *Command) AddFlags() (err error) {
	for _, f := range c.InvokedFlags() {
		fu := FlagUnion{}
		flagValuesMap[f.Unique()] = &fu
		switch {
		case f.SetStringValFunc != nil:
			flag.StringVar(&fu.String, f.Switch, f.Default.(string), f.Usage)
		case f.SetIntValFunc != nil:
			flag.IntVar(&fu.Int, f.Switch, f.Default.(int), f.Usage)
		default:
			f.noSetFuncAssigned()
		}
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
	return c.ArgsCount() - c.RequiredArgsCount()
}

// ArgsCount returns the number of total args; required and optional
func (c *Command) ArgsCount() (cnt int) {
	return len(c.Args)
}

func (c *Command) AddArg(arg *Arg) (cmd *Command) {
	arg.Parent = c

	if arg.SetStringValFunc == nil && arg.SetIntValFunc == nil {
		arg.SetStringValFunc = func(s string) {}
	}
	if arg.Default == nil {
		switch {
		case arg.SetStringValFunc != nil:
			arg.Default = ""
		case arg.SetIntValFunc != nil:
			arg.Default = 0
			//default:
			//	arg.noSetFuncAssigned()
		}
	}
	c.Args = append(c.Args, arg)
	return c
}
