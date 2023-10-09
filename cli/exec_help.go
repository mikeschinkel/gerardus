package cli

import (
	"fmt"
	"slices"
	"strings"
)

//goland:noinspection ALL
var CmdHelp = AddCommandWithFunc("help", ExecHelp).AddOptArgs("command")

func ExecHelp(...string) error {
	var output func(*Command, string) []string
	output = func(cmd *Command, prefix string) (out []string) {
		if len(cmd.SubCommands) == 0 {
			out = []string{fmt.Sprintf("\t   %s%s\n", prefix, cmd)}
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
	StdErr("\tUsage: gerardus <options> <command> [<args>]\n")
	StdErr("\tCommands:\n")
	cmdHelp := make([]string, 0)
	for _, cmd := range Commands() {
		cmdHelp = append(cmdHelp, output(cmd, "")...)
	}
	slices.Sort(cmdHelp)
	StdErr("%s\n", strings.Join(cmdHelp, ""))
	return nil
}
