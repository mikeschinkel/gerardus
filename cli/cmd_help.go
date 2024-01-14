package cli

//goland:noinspection ALL
var CmdHelp = AddCommandWithFunc(Token(HelpArg), ExecHelp).
	AddArg(Arg{
		Name:     CommandArg,
		Usage:    "Specifies the command to show help for",
		Optional: true,
		Variadic: true,
	}.EmptyOk())

// ExecHelp "implements" the help command by simply delegating to help.Usage() by
// returning an sentinel error which help.Usage() will recognize and generate the
// appropriate output.
func ExecHelp(Context, *CommandInvoker) (err error) {
	return ErrHelpSentinel
}
