package cli_test

import (
	"strings"
	"testing"

	"github.com/mikeschinkel/gerardus/cli"
)

func TestSetExecFunc(t *testing.T) {
	ef := func(cli.ArgsMap) error { return nil }
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
			if !hasCommand(tt.name) {
				t.Errorf("setExecFunc() failed for %s", tt.name)
			}
		})
	}
}

func hasCommand(name string) (has bool) {
	cmd, _ := cli.CommandByName(name)
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
		cmd, ok = a.commands[a.name]
		if !ok {
			cmd = cli.NewCommand(a.name, a.execFunc)
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
		names := cli.MatchSpaces.Split(a.name, -1)
		switch len(names) {
		case 0:
		case 1:
			ensureCommand(args{a.commands, a.name, a.parent, a.execFunc})
		default:
			cmd := ensureCommand(args{cli.Commands(), names[0], a.parent, nil})
			name = strings.Join(names[1:], " ")
			traverseCommands(args{cmd.SubCommands, name, cmd, a.execFunc})
		}
	}
	traverseCommands(args{cli.Commands(), name, nil, ef})
}
