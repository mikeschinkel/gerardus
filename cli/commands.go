package cli

import (
	"fmt"
	"os"
	"strings"
)

type CommandMap map[string]*Command

// RootCmd is parent of all top level command and is where global flags will
// be attached.
var RootCmd = NewCommand("", nil)

func Commands() CommandMap {
	return RootCmd.SubCommands
}

func CommandByName(name string) (cmd *Command, depth int) {
	var ok bool

	cmds := RootCmd.SubCommands
	cmdNames := MatchSpaces.Split(name, -1)
	for _, cmdName := range cmdNames {
		depth++
		cmd, ok = cmds[cmdName]
		if !ok {
			cmd = nil
			goto end
		}
		if len(cmd.SubCommands) == 0 {
			goto end
		}
		cmds = cmd.SubCommands
	}
end:
	return cmd, depth
}

var argsCount *int

func ArgsCount() (_ int, err error) {
	var n, depth int
	if argsCount != nil {
		goto end
	}
	argsCount = &n
	for _, arg := range os.Args[1:] {
		if arg[0] == '-' {
			continue
		}
		*argsCount++
	}
	_, depth, err = InvokedCommand()
	if err != nil {
		goto end
	}
	*argsCount -= depth
end:
	return *argsCount, err
}

func ExecInvokedCommand() (err error) {
	var am ArgsMap
	var expected, got int

	cmd, _, err := InvokedCommand()
	if err != nil {
		goto end
	}
	expected = cmd.RequiredArgsCount()
	got, err = ArgsCount()
	if err != nil {
		goto end
	}
	if got < expected {
		err = fmt.Errorf("not enough CLI args passed; expected at least %d, got %d", expected, got)
		goto end
	}
	expected = cmd.ArgsCount()
	if got > expected {
		err = fmt.Errorf("too many CLI args passed; expected no more than %d, got %d", expected, got)
		goto end
	}
	am, err = cmd.ArgValuesMap()
	err = cmd.ExecuteFunc(am)
end:
	return err
}

var invokedCommand *Command
var commandDepth int

func InvokedCommand() (_ *Command, _ int, err error) {
	var arg string
	var cnt int

	if invokedCommand != nil {
		goto end
	}

	arg, cnt, err = CommandString()
	if err != nil {
		goto end
	}

	if cnt == 0 {
		invokedCommand = RootCmd.SubCommands["help"]
		goto end
	}
	invokedCommand, commandDepth = CommandByName(arg)
	if invokedCommand != nil {
		invokedCommand.Name = arg
	}
end:
	return invokedCommand, commandDepth, err
}

var commandCount *int

// CommandCount returns the number of commands minus the flags
func CommandCount() (cnt int, err error) {
	if commandCount != nil {
		goto end
	}
	_, cnt, err = CommandString()
	if err != nil {
		goto end
	}
	commandCount = &cnt
end:
	return *commandCount, err
}

var commandString *string

// CommandString returns the full list of commands minus the flags
func CommandString() (cs string, _ int, err error) {
	var sb strings.Builder
	var n int
	var cmds CommandMap

	if commandString != nil {
		goto end
	}
	commandString = &cs
	cmds = RootCmd.SubCommands
	commandCount = &n
	sb = strings.Builder{}
	for _, arg := range os.Args[1:] {
		if arg[0] == '-' {
			continue
		}
		if _, ok := cmds[arg]; !ok {
			err = fmt.Errorf("command '%s' is not valid", arg)
			goto end
		}
		sb.WriteString(arg)
		sb.WriteByte(' ')
		*commandCount++
		cmds = cmds[arg].SubCommands
	}
	cs = sb.String()
	if len(cs) == 0 {
		goto end
	}
	cs = cs[:len(cs)-1]
	commandString = &cs
end:
	return *commandString, *commandCount, err
}

func AddCommand(name string) (cmd *Command) {
	return RootCmd.AddSubCommand(name, nil)
}

func AddCommandWithFunc(name string, ef ExecFunc) (cmd *Command) {
	return RootCmd.AddSubCommand(name, ef)
}
