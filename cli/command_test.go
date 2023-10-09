package cli_test

import (
	"testing"

	"gerardus/cli"
)

func TestSetExecFunc(t *testing.T) {
	ef := func(...string) error { return nil }
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
			cli.SetExecFunc(tt.name, tt.ef)
			if !hasCommand(tt.name) {
				t.Errorf("SetExecFunc() failed for %s", tt.name)
			}
		})
	}
}

func hasCommand(name string) (has bool) {
	cmd, _ := cli.CommandByName(name)
	return cmd != nil && cmd.ExecFunc != nil
}
