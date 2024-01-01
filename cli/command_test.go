package cli_test

import (
	"strings"
	"testing"

	"github.com/mikeschinkel/gerardus/cli"
)

func TestSetExecFunc(t *testing.T) {
	ef := func(*cli.CommandInvoker) error { return nil }
	tests := []struct {
		name string
		ef   cli.ExecFunc
	}{
		{"add subcmd", ef},
		{"add", ef},
		{"add        subcmd", ef},
		{"add subcmd subcmd2", ef},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setExecFunc(tt.name, tt.ef)
			if !hasCommandByName(cli.RootCmd, tt.name) {
				t.Errorf("setExecFunc() failed for %s", tt.name)
			}
		})
	}
}

func hasCommandByName(rootCmd *cli.Command, name string) (has bool) {
	cmd, _ := cli.CommandByName(rootCmd, name)
	return cmd != nil && cmd.ExecFunc != nil
}

func setExecFunc(name string, ef cli.ExecFunc) {
	type args struct {
		commands cli.CommandMap
		name     string
		parent   *cli.Command
		execFunc cli.ExecFunc
	}
	ensureCommand := func(a args) (cmd *cli.Command) {
		var ok bool

		token := cli.Token(a.name)
		cmd, ok = a.commands[token]
		if !ok {
			cmd = cli.NewCommand(token, a.execFunc)
			a.commands[token] = cmd
		} else if a.execFunc != nil {
			cmd.ExecFunc = a.execFunc
		}
		cmd.Name = token
		cmd.Parent = a.parent
		return cmd
	}
	var traverseCommands func(args)
	traverseCommands = func(a args) {
		names := cli.MatchSpaces.Split(string(a.name), -1)
		switch len(names) {
		case 0:
		case 1:
			ensureCommand(args{a.commands, a.name, a.parent, a.execFunc})
		default:
			cmd := ensureCommand(args{cli.Commands(), names[0], a.parent, nil})
			traverseCommands(args{
				cmd.SubCommands,
				strings.Join(names[1:], " "),
				cmd,
				a.execFunc,
			})
		}
	}
	traverseCommands(args{cli.Commands(), name, nil, ef})
}
