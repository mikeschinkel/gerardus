package cli

//goland:noinspection ALL
var CmdHelp = AddCommandWithFunc(Token(HelpArg), ExecHelp).
	AddArg(Arg{
		Name:     CommandArg,
		Usage:    "Specifies the command to show help for",
		Optional: true,
		Variadic: true,
	}.EmptyOk())

func ExecHelp(ctx Context, i *CommandInvoker) (err error) {
	// TODO Implement help
	return ErrHelpSentinel
}
