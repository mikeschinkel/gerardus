package cli

import (
	"flag"
	"regexp"
	"slices"
	"strings"
)

var MatchSpaces = regexp.MustCompile(`\s+`)
var AppName string

func Initialize(appName string) (err error) {

	AppName = appName

	cmd, _, err := InvokedCommand()
	if err != nil {
		goto end
	}

	err = cmd.AddFlags()
	if err != nil {
		goto end
	}
	flag.Parse()
	cmd.SetFlagValues()

end:
	return err
}

// ValidateInput validates flags and args passed on the CLI
func ValidateInput() (err error) {
	var cmd *Command
	var sm StringMap
	var am ArgsMap

	cmd, _, err = InvokedCommand()
	if err != nil {
		goto end
	}
	sm, am = cmd.ArgValuesMap()
	err = am.validate(sm)
	if err != nil {
		goto end
	}
	err = cmd.InvokedFlags().validate()
	if err != nil {
		goto end
	}
end:
	return err
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

func (f *Flag) noSetFuncAssigned() {
	panicf("No func(<type>) assigned to property `Set*ValFunc` for flag '%s'", f.Unique())
}

func SetExecFunc(name string, ef ExecFunc) {
	type args struct {
		commands CommandMap
		name     string
		parent   *Command
		execFunc ExecFunc
	}
	ensureCommand := func(a args) (cmd *Command) {
		var ok bool
		cmd, ok = a.commands[a.name]
		if !ok {
			cmd = NewCommand(a.name, a.execFunc)
			a.commands[a.name] = cmd
		} else if a.execFunc != nil {
			cmd.ExecFunc = a.execFunc
		}
		cmd.Name = a.name
		cmd.Parent = a.parent
		return cmd
	}
	var traverseCommands func(args)
	traverseCommands = func(a args) {
		names := MatchSpaces.Split(a.name, -1)
		switch len(names) {
		case 0:
		case 1:
			ensureCommand(args{a.commands, a.name, a.parent, a.execFunc})
		default:
			cmd := ensureCommand(args{Commands(), names[0], a.parent, nil})
			name = strings.Join(names[1:], " ")
			traverseCommands(args{cmd.SubCommands, name, cmd, a.execFunc})
		}
	}
	traverseCommands(args{Commands(), name, nil, ef})
}
