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

func ArgsCount() int {
	var n int
	if argsCount != nil {
		goto end
	}
	// Calling this ensures commandDepth is set
	InvokedCommand()
	argsCount = &n
	for _, arg := range os.Args[1:] {
		if arg[0] == '-' {
			continue
		}
		*argsCount++
	}
	*argsCount -= commandDepth
end:
	return *argsCount
}

func ExecInvokedCommand() (err error) {
	cmd, depth := InvokedCommand()
	expected := len(cmd.Args)
	got := ArgsCount()
	if got < expected {
		err = fmt.Errorf("not enough CLI args passed; expected at least %d, got %d", expected, got)
		goto end
	}
	expected += len(cmd.OptArgs)
	if got > expected {
		err = fmt.Errorf("too many CLI args passed; expected no more than %d, got %d", expected, got)
		goto end
	}
	err = cmd.ExecuteFunc(os.Args[depth:]...)
end:
	return err
}

var invokedCommand *Command
var commandDepth int

func InvokedCommand() (*Command, int) {
	var arg string
	if invokedCommand != nil {
		goto end
	}

	if CommandCount() == 0 {
		invokedCommand = RootCmd.SubCommands["help"]
		goto end
	}
	arg = CommandString()
	invokedCommand, commandDepth = CommandByName(arg)
	if invokedCommand != nil {
		invokedCommand.Name = arg
	}
end:
	return invokedCommand, commandDepth
}

var commandCount *int

// CommandCount returns the number of commands minus the flags
func CommandCount() int {
	if commandCount != nil {
		goto end
	}
	CommandString()
end:
	return *commandCount
}

var commandString *string

// CommandString returns the full list of commands minus the flags
func CommandString() string {
	var sb strings.Builder
	var cs string
	var n int
	var cmds CommandMap
	if commandString != nil {
		goto end
	}
	cmds = RootCmd.SubCommands
	commandCount = &n
	sb = strings.Builder{}
	for _, arg := range os.Args[1:] {
		if arg[0] == '-' {
			continue
		}
		if _, ok := cmds[arg]; !ok {
			break
		}
		sb.WriteString(arg)
		sb.WriteByte(' ')
		*commandCount++
		cmds = cmds[arg].SubCommands
	}
	cs = sb.String()
	cs = cs[:len(cs)-1]
	commandString = &cs
end:
	return *commandString
}

func AddCommand(name string) (cmd *Command) {
	return RootCmd.AddSubCommand(name, nil)
}

func AddCommandWithFunc(name string, ef ExecFunc) (cmd *Command) {
	return RootCmd.AddSubCommand(name, ef)
}
