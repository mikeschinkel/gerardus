package cli

//goland:noinspection ALL
var CmdHelp = AddCommandWithFunc("help", ExecHelp).
	AddArg(&Arg{
		Name:     CommandArg,
		Usage:    "Specifies the command to show help for",
		Optional: true,
	})

func ExecHelp(i *CommandInvoker) (err error) {
	return err
}
