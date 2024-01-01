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

	fs = flag.NewFlagSet(
		ExecutableFilepath(params.AppName),
		flag.ContinueOnError,
	)
	fs.SetOutput(StderrWriter)
	flag.CommandLine = fs

	flags = cmd.InvokedFlags().Initialize()

	err = flag.CommandLine.Parse(params.Options().StringSlice())
	if err != nil {
		goto end
	}

	flags = flags.callSetValueFuncs()
	cmd.SetFlags(flags)

end:
	return invoker, err
}
