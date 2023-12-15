package cli

import (
	"fmt"
	"log/slog"
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

//func (i *CommandInvoker) ArgCount() (cnt int, err error) {
//	var depth int
//	_, depth, err = InvokedCommand(i.Tokens)
//	if err != nil {
//		goto end
//	}
//	cnt -= depth
//end:
//	return cnt, err
//}

func (i *CommandInvoker) InvokeCommand() (err error) {
	cmd := i.Command
	slog.Info("Invoking command", "command", cmd.String())
	err = cmd.ExecFunc(i)
	return err
}

func InvokedCommand(tokens []string) (_ *Command, _ int, err error) {
	var cmdName string
	var cmd *Command
	var depth int
	var cnt int

	cmdName, cnt, err = CommandString(tokens)
	if err != nil {
		goto end
	}

	if cnt == 0 {
		cmd = RootCmd
		goto end
	}
	cmd, depth = CommandByName(cmdName)
	//if cmd != nil {
	// This assigns "add project" that what which was previously "project"
	// Not sure we actually want this.
	//	cmd.Name = cmdName
	//}

end:
	return cmd, depth, err
}

// CommandString returns the full list of commands minus the flags
func CommandString(tokens []string) (cs string, _ int, err error) {
	var sb strings.Builder
	var cnt int
	var cmds CommandMap

	cmds = RootCmd.SubCommands
	sb = strings.Builder{}
	for _, token := range tokens {
		if _, ok := cmds[token]; !ok {
			err = fmt.Errorf("command '%s' is not valid", token)
			goto end
		}
		sb.WriteString(token)
		sb.WriteByte(' ')
		cnt++
		cmds = cmds[token].SubCommands
		if len(cmds) == 0 {
			break
		}
	}
	cs = sb.String()
	if len(cs) == 0 {
		goto end
	}
	cs = cs[:len(cs)-1]
end:
	return cs, cnt, err
}

func AddCommand(name string) (cmd *Command) {
	return RootCmd.AddSubCommand(name, nil)
}

func AddCommandWithFunc(name string, ef ExecFunc) (cmd *Command) {
	return RootCmd.AddSubCommand(name, ef)
}
