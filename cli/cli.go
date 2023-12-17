package cli

import (
	"flag"
	"log/slog"
	"regexp"
)

var MatchSpaces = regexp.MustCompile(`\s+`)

func Initialize(ctx Context, params Params) (invoker *CommandInvoker, err error) {
	var flags Flags

	slog.Info("Initializing commands")

	invoker = NewCommandInvoker(params)
	args := params.Args()

	cmd, _, err := InvokedCommand(args)
	if err != nil {
		goto end
	}
	invoker.Command = cmd

	err = cmd.callSetArgValues(args)
	if err != nil {
		goto end
	}

	err = cmd.callSetArgValueFuncs()
	if err != nil {
		goto end
	}

	flag.CommandLine = flag.NewFlagSet(
		ExecutableFilepath(params.AppName),
		flag.ExitOnError,
	)
	flags = cmd.InvokedFlags()
	err = flags.Initialize()
	if err != nil {
		goto end
	}
	err = flag.CommandLine.Parse(params.Options())
	if err != nil {
		goto end
	}
	flags.callSetValueFuncs()

end:
	return invoker, err
}
