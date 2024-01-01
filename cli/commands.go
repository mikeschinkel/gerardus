package cli

import (
	"log/slog"
	"strings"
)

type CommandMap map[Token]*Command

// RootCmd is parent of all top level command and is where global flags will
// be attached.
var RootCmd = NewCommand("", nil)

func Commands() CommandMap {
	return RootCmd.SubCommands
}

func CommandByName(rootCmd *Command, name string) (cmd *Command, depth int) {
	var ok bool
	var cmds CommandMap
	var cmdNames []string

	if name == "" {
		goto end
	}
	cmds = rootCmd.SubCommands
	cmdNames = MatchSpaces.Split(name, -1)
	for _, cmdName := range cmdNames {
		depth++
		cmd, ok = cmds[Token(cmdName)]
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

func (i *CommandInvoker) InvokeCommand(ctx Context) (err error) {
	cmd := i.Command
	slog.Info("Invoking command", "command", cmd.String())
	err = cmd.ExecFunc(ctx, i)
	return err
}

func InvokedCommand(rootCmd *Command, tokens Tokens) (_ *Command, _ int, err error) {
	var cmdName string
	var cmd *Command
	var depth int
	var cnt int

	cmdName, cnt, err = CommandString(rootCmd, tokens)
	if err != nil {
		goto end
	}

	if cnt == 0 {
		cmd = RootCmd
		goto end
	}
	cmd, depth = CommandByName(rootCmd, cmdName)
	//if cmd != nil {
	// This assigns "add project" that what which was previously "project"
	// Not sure we actually want this.
	//	cmd.Name = cmdName
	//}

end:
	return cmd, depth, err
}

// CommandString returns the full list of commands minus the flags
func CommandString(rootCmd *Command, tokens Tokens) (cs string, cnt int, err error) {
	var sb strings.Builder
	var cmds CommandMap

	cmds = rootCmd.SubCommands
	sb = strings.Builder{}
	for i, token := range tokens {
		if token[0] == '-' {
			// It's a flag, this function doesn't deal with that. So panic.
			panicf("CommandString(rootRmd,tokens) expects tokens will have '-flags' filtered out, yet flag '%s' found.", token)
		}
		if _, ok := cmds[Token(token)]; !ok {
			err = ErrCommandNotValid.Args("command", tokens[:i+1])
			goto end
		}
		sb.WriteString(string(token))
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

func AddCommand(name Token) (cmd *Command) {
	return RootCmd.AddSubCommand(name, nil)
}

func AddCommandWithFunc(name Token, ef ExecFunc) (cmd *Command) {
	return RootCmd.AddSubCommand(name, ef)
}
