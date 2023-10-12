package cli

import (
	"fmt"
	"slices"
	"strings"
)

//goland:noinspection ALL
var CmdHelp = AddCommandWithFunc("help", ExecHelp).
	AddArg(&Arg{
		Name:     "command",
		Usage:    "Specifies the command to show help for",
		Optional: true,
	})

func ExecHelp(args StringMap) error {
	var output func(*Command, string) []string

	output = func(cmd *Command, prefix string) (out []string) {
		if len(cmd.SubCommands) == 0 {
			out = []string{fmt.Sprintf("\t   %s%s\n", prefix, cmd.Name)}
			goto end
		}
		for _, sc := range cmd.SubCommands {
			s := fmt.Sprintf("%s %s ", prefix, cmd.Name)
			s = strings.TrimLeft(s, " ")
			out = append(out, output(sc, s)...)
		}
	end:
		return out
	}
	StdErr("\tUsage: %s [<options>] <command> [<args>]\n", AppName)
	StdErr("\tCommands:\n")
	cmdHelp := make([]string, 0)
	for _, cmd := range Commands() {
		cmdHelp = append(cmdHelp, output(cmd, "")...)
	}
	slices.Sort(cmdHelp)
	StdErr("%s\n", strings.Join(cmdHelp, ""))
	return nil
}
