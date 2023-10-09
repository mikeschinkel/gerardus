package cli

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

var MatchSpaces = regexp.MustCompile(`\s+`)

func Initialize() (err error) {
	err = addFlags()
	if err != nil {
		goto end
	}
	flag.Parse()
	setValues()
	if len(os.Args) >= 2 {
		cmd := os.Args[1]
		if _, ok := RootCmd.SubCommands[cmd]; ok {
			goto end
		}
		err = fmt.Errorf("command '%s' is not a valid command", cmd)
		goto end
	}
	err = checkFlags()
end:
	return err
}

// AllFlags returns all the flags for this command, including all parent flags
// including the root flags.
func (c *Command) AllFlags() (flags Flags) {
	cmd := c
	flags = c.Flags
	for cmd.Parent != nil {
		cmd = cmd.Parent
		if len(cmd.Flags) == 0 {
			continue
		}
		flags = append(flags, cmd.Flags...)
	}
	slices.Reverse(flags)
	return flags
}

// Sets the value specified by
func setValues() {
	cmd, _ := InvokedCommand()
	for _, f := range cmd.AllFlags() {
		fv := flagValues[f.Unique()]
		switch {
		case f.SetStringValFunc != nil:
			f.SetStringValFunc(fv.String)
		case f.SetIntValFunc != nil:
			f.SetIntValFunc(fv.Int)
		default:
			noSetFuncAssigned(f)
		}
	}

}

func noSetFuncAssigned(f *Flag) {
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
