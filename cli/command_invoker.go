package cli

import (
	"errors"
	"log/slog"

	"github.com/mikeschinkel/go-serr"
)

type CommandInvoker struct {
	Tokens    Tokens
	AppName   string
	EnvPrefix string
	Command   *Command
	Args      Args
}

func NewCommandInvoker(params Params) *CommandInvoker {
	return &CommandInvoker{
		Tokens:    params.Tokens(),
		AppName:   params.AppName,
		EnvPrefix: params.EnvPrefix,
	}
}

func (i *CommandInvoker) ArgValue(name ArgName) (value *Value) {
	for _, arg := range i.Args {
		if arg.Name != name {
			continue
		}
		value = arg.Value
		goto end
	}
end:
	return value
}
func (i *CommandInvoker) ArgString(name ArgName) (s string) {
	value := i.ArgValue(name)
	if value == nil {
		goto end
	}
	s = value.String()
end:
	return s
}
func (i *CommandInvoker) ArgInt(name ArgName) (n int) {
	value := i.ArgValue(name)
	if value == nil {
		goto end
	}
	n = value.Int()
end:
	return n
}

func (i *CommandInvoker) SubCommands() CommandMap {
	return i.Command.SubCommands
}

// Validate validates args and options passed on the CLI
func (i *CommandInvoker) Validate(ctx Context) (err error) {
	slog.Info("Validating CLI Args and Flags")
	cmd := i.Command

	if cmd.ExecFunc == nil {
		if cmd == RootCmd {
			err = ErrNoCommandSpecified
			goto end
		}
		// For when using a partial command like 'add' when the command is 'add project'.
		err = ErrNoExecFuncFound
		goto end
	}

	err = i.EmptyStateSatisfied()
	if err != nil {
		goto end
	}
	err = i.ValidateArgs(ctx)
	if err != nil {
		goto end
	}
	err = i.ValidateFlags(ctx)
	if err != nil {
		goto end
	}

end:
	return err
}

// EmptyStateSatisfied ensures that values of .Requires are satisfied for both args and options.
func (i *CommandInvoker) EmptyStateSatisfied() (err error) {
	cmd := i.Command
	err = cmd.Args.EmptyStateSatisfied()
	if err != nil {
		goto end
	}
	err = cmd.InvokedFlags().EmptyStateSatisfied()
	if err != nil {
		goto end
	}
end:
	return err
}

// ValidateFlags validates options passed on the CLI
func (i *CommandInvoker) ValidateFlags(ctx Context) (err error) {
	return i.Command.InvokedFlags().Validate(ctx)
}

// ValidateArgs validates args passed on the CLI
func (i *CommandInvoker) ValidateArgs(ctx Context) (err error) {
	var expected, got int

	cmd := i.Command
	err = cmd.Args.Validate(ctx)
	if err != nil {
		if errors.Is(err, ErrRequiresCheckFailed) {
			err = err.(serr.SError).Unwrap()
		}
		goto end
	}
	expected = cmd.RequiredArgsCount()
	got = i.ReceivedArgsCount()
	if got < expected {
		// TODO: Add 'missing'
		err = ErrTooFewArgsPassed.Args("expected", expected, "passed", got)
		goto end
	}
	expected = cmd.DeclaredArgsCount()
	if got > expected {
		// TODO: Add 'extra'
		err = ErrTooManyArgsPassed.Args("expected", expected, "passed", got)
		goto end
	}
end:
	return err
}

// ReceivedArgsCount returns number of args received on command line. Example: If
// the command is "make widget" and the os.Args has:
//
//	"/path/to/maker make widget -v -n foo bar baz"
//
// then ReceivedArgsCount returns 3 for "foo bar baz."
func (i *CommandInvoker) ReceivedArgsCount() int {
	return len(i.Tokens.Args())
}
