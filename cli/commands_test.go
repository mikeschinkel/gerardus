package cli_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/mikeschinkel/gerardus/cli"
)

type test struct {
	name      string
	tokens    cli.Tokens
	want      string
	wantCnt   int
	wantDepth int
	wantErr   bool
	panic     bool
	panicMsg  string
	errMsg    string
	wantCmd   *cli.Command
	flags     cli.Flags
}

func (tt test) ShouldPanic() bool {
	return tt.panic
}
func (tt test) ErrorMsg() string {
	return tt.panicMsg
}
func (tt test) Name() string {
	return tt.name
}

func setFlagRequires(flags cli.Flags, name string, r cli.ArgRequires) cli.Flags {
	n := flags.Index(cli.ArgName(name))
	if n != -1 {
		flags[n].Requires = r
		flags[n].Value = cli.NewValue(flags[n].Type, nil)
	}
	return flags
}

func TestCommands(t *testing.T) {
	rootCmd := rootCommand()
	tests := []test{
		{
			name:      "No Args",
			tokens:    cli.Tokens{},
			want:      "",
			wantCnt:   0,
			wantCmd:   nil,
			wantDepth: 0,
		},
		{
			name:      "One (1) Arg",
			tokens:    cli.Tokens{"one"},
			want:      "one",
			wantCnt:   1,
			wantCmd:   oneCommand(rootCmd),
			wantDepth: 1,
		},
		{
			name:      "Two (2) Args",
			tokens:    cli.Tokens{"one", "two"},
			want:      "one two",
			wantCnt:   2,
			wantCmd:   twoCommand(rootCmd),
			wantDepth: 2,
		},
		{
			name:      "Three (3) Args, but only two (2) in command",
			tokens:    cli.Tokens{"one", "two", "three"},
			want:      "one two",
			wantCnt:   2,
			wantCmd:   twoCommand(rootCmd),
			wantDepth: 2,
		},
		{
			name:     "PANIC - Two (2) Args w/opt=test flag",
			tokens:   cli.Tokens{"-opt=test", "one", "two"},
			want:     "one two",
			panic:    true,
			panicMsg: "CommandString(rootRmd,tokens) expects tokens will have '-flags' filtered out, yet flag '-opt=test' found.",
		},
		{
			name:     "-foo on command 'one' but with no value",
			tokens:   cli.Tokens{"-foo=", "one", "two"},
			want:     "one two",
			panic:    true,
			panicMsg: "CommandString(rootRmd,tokens) expects tokens will have '-flags' filtered out, yet flag '-foo=' found.",
			flags:    setFlagRequires(oneCommand(rootCmd).Flags, "foo", cli.NotEmpty),
			errMsg:   "value cannot be empty [option=-foo]",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdName := ""
			ctx := context.Background()
			t.Run("CommandString()", func(t *testing.T) {
				defer OnTestCasePanic(t, tt)
				got, gotCnt, err := cli.CommandString(rootCmd, tt.tokens)
				if (err != nil) != tt.wantErr {
					t.Errorf("CommandString() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("CommandString() got = %v, want = %v", got, tt.want)
				}
				if gotCnt != tt.wantCnt {
					t.Errorf("CommandString() gotCnt = %v, wantCnt = %v", gotCnt, tt.wantCnt)
				}
				cmdName = got
			})
			t.Run("CommandByName()", func(t *testing.T) {
				defer OnTestCasePanic(t, tt)
				gotCmd, gotDepth := cli.CommandByName(rootCmd, cmdName)
				if gotCmd != tt.wantCmd {
					t.Errorf("CommandString() gotCmd = %v, wantCmd = %v", gotCmd, tt.wantCmd)
				}
				if gotDepth != tt.wantDepth {
					t.Errorf("CommandString() gotDepth = %v, wantDepth = %v", gotDepth, tt.wantDepth)
				}
			})
			if tt.flags != nil {
				t.Run("MeetsRequirements()", func(t *testing.T) {
					defer OnTestCasePanic(t, tt)
					err := cli.MeetsRequirements(ctx, cli.FlagType, tt.flags)
					switch {
					case err == nil && tt.wantErr:
						t.Errorf("MeetsRequirements() error = %v, wantErr = %v", err, tt.wantErr)
					case err != nil && !tt.wantErr:
						t.Errorf("MeetsRequirements() error = %v, wantErr = %v", err, tt.wantErr)
					case tt.wantErr && err.Error() != tt.errMsg:
						t.Errorf("MeetsRequirements() error = %v, wantErr = %v", err, tt.errMsg)
					}
				})
			}
		})
	}
}

func stubFunc(context cli.Context, invoker *cli.CommandInvoker) error {
	return nil
}
func rootCommand() *cli.Command {
	rootCmd := cli.NewCommand("", nil).AddFlag(cli.Flag{
		Switch: "opt",
		Arg: cli.Arg{
			Name:  "option",
			Usage: "Option for Testing",
			Type:  reflect.String,
		},
	})
	oneCmd := rootCmd.AddSubCommand("one", stubFunc).AddFlag(cli.Flag{
		Switch: "foo",
		Arg: cli.Arg{
			Name:  "foo",
			Usage: "Option Foo for Testing",
			Type:  reflect.String,
		},
	})
	oneCmd.AddSubCommand("two", stubFunc)
	return rootCmd
}
func oneCommand(rootCmd *cli.Command) *cli.Command {
	return rootCmd.SubCommands["one"]
}
func twoCommand(rootCmd *cli.Command) *cli.Command {
	return oneCommand(rootCmd).SubCommands["two"]
}

// TestCasePanicker allows OnTestCasePanic() to be used by any table test case
// that implements the interface.
//
// TODO: Create a github.com/mikeschinkel/go-testlib with reusable functions for
//
//	testing and put this function there.
type TestCasePanicker interface {
	ShouldPanic() bool
	ErrorMsg() string
	Name() string
}

// OnTestCasePanic can be used on deter statements within a test to capture a
// panic message for checking to see if it matches values for the test case.
// NOTE: Panics should be used RARELY, and NOT as exception handling like in PHP
// and Java. NEVER use panics in long-running server code – unless you have
// infrastructure to restart crashed apps — but they are useful in CLI apps for
// code paths that should never happen but if they do are indicative of a code
// bug that needs to be fixed. In this case it is good to be able to confirm via
// tests that the panic operates as expected when they are triggered.
//
// TODO: Create a github.com/mikeschinkel/go-testlib with reusable functions for
//
//	testing and put this function there.
func OnTestCasePanic(t *testing.T, tt TestCasePanicker) {
	var errMsg string

	t.Helper()
	r := recover()
	if r == nil {
		goto end
	}
	if !tt.ShouldPanic() {
		t.Errorf("CommandString() unexpectedly panicked for test case '%s': %v",
			tt.Name(),
			r,
		)
	}
	errMsg = fmt.Sprintf("%v", r)
	if errMsg != tt.ErrorMsg() {
		t.Errorf("CommandString() incorrect panic message:\nGot:  %s\nWant: %s\n",
			errMsg,
			tt.ErrorMsg(),
		)
	}
end:
}
