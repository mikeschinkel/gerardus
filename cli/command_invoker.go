package cli

import (
	"log/slog"
)

type CommandInvoker struct {
	Tokens    Tokens
	AppName   string
	EnvPrefix string
	Command   *Command
}

func NewCommandInvoker(params Params) *CommandInvoker {
	return &CommandInvoker{
		Tokens:    params.Tokens(),
		AppName:   params.AppName,
		EnvPrefix: params.EnvPrefix,
	}
}

// MeetsRequirements validates args and options passed on the CLI
func (i *CommandInvoker) MeetsRequirements(ctx Context) (err error) {
	slog.Info("Validating CLI Args and Flags")
	return i.Command.MeetsRequirements(ctx, i.Tokens.Count())
}

func (i *CommandInvoker) ArgValue(name ArgName) (value *Value) {
	for _, arg := range i.Args() {
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

func (i *CommandInvoker) Args() Args {
	return i.Command.Args
}
