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
	args, err := params.Args()
	if err != nil {
		goto end
	}

	cmd, _, err = InvokedCommand(RootCmd, args)
	if err != nil {
		goto end
	}
	invoker.Command = cmd

	err = cmd.setArgValues(args)
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

	flags = cmd.InvokedFlags().Initialize(ctx)

	err = flag.CommandLine.Parse(params.Options().StringSlice())
	if err != nil {
		goto end
	}

	flags = flags.callSetValueFuncs()
	cmd.SetFlags(flags)

end:
	return invoker, err
}
