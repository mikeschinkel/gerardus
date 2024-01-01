package cli

import (
	"flag"
	"log/slog"
	"regexp"
)

var MatchSpaces = regexp.MustCompile(`\s+`)

func Initialize(ctx Context, params Params) (invoker *CommandInvoker, err error) {
	var flags Flags
	var cmd *Command
	var fs *flag.FlagSet

	slog.Info("Initializing commands")

	invoker = NewCommandInvoker(params)
	args := params.Args()

	cmd, _, err = InvokedCommand(RootCmd, args)
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
	flags = flags.Initialize()
	if err != nil {
		goto end
	}
	err = flag.CommandLine.Parse(params.Options().StringSlice())
	if err != nil {
		goto end
	}

	flags = flags.callSetValueFuncs()
	cmd.SetFlags(flags)

end:
	return invoker, err
}
