package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mikeschinkel/go-lib"
)

type Help struct {
	invoker             *CommandInvoker
	SetStderrWriterFunc func(io.Writer)
}

func NewHelp(invoker *CommandInvoker) Help {
	return Help{
		invoker: invoker,
	}
}

func (h Help) Usage(err error, w io.Writer) {
	// Set the stdErr writer
	h.SetStderrWriterFunc(w)
	StdErr(h.GetUsage(err))
}

func (h Help) GetUsage(err error) (help string) {

	switch h.commandType(err) {
	case RootCommand:
		help = h.rootCmdHelp(err)
	case BranchCommand:
		help = h.branchCmdHelp(err)
	case LeafCommand:
		help = h.leafCmdHelp(err)
	case HelpCommand:
		help = h.helpCmdHelp()
	default:
		panicf("Undefined command type %d for '%s'",
			h.commandType(err),
			h.command().Name,
		)
	}
	return help
}

func (h Help) command() *Command {
	return h.invoker.Command
}

func (h Help) withCommand(cmd *Command) Help {
	h.invoker.Command = cmd
	return h
}
func (h Help) withTokens(tokens Tokens) Help {
	h.invoker.Tokens = tokens
	return h
}

func (h Help) cmdForHelp() (cmd *Command, tokens Tokens) {
	var first, name string

	tokens = h.invoker.Tokens
	if len(tokens) >= 2 {
		first = string(tokens[1])
	}
	switch {
	case ArgName(first) != HelpArg:
		panicf("Help.cmdForHelp() called when 'help' is not the top level command.")

	case len(tokens) <= 2:
		// We have no command, or just a 'help' command (though the former should never happen.)
		cmd = RootCmd
		// Make sure tokens is just the executable file name.
		tokens = tokens[:1]

	default:
		// Save the executable filepath
		exe := tokens[0]
		// Omit executable file name (tokens[0]) and the "help" command (tokens[1]).
		tokens = tokens[2:]
		// Convert to a string for command lookup via CommandByName
		name = tokens.Join(" ")
		// Get command for help
		cmd, _ = CommandByName(RootCmd, name)
		// Rebuild tokens for the command, so we can get help for that command.
		tokens = append(Tokens{exe}, tokens...)
	}
	return cmd, tokens
}

func (h Help) helpCmdHelp() string {
	var err error

	cmd, tokens := h.cmdForHelp()
	if cmd == nil {
		// If we do not have a command, or it is invalid, then note the error.
		err = ErrCommandNotValid.Args("command", tokens[1:].Join(" "))
		// Use the RootCmd for spoofing GetUsage().
		cmd = RootCmd
	}
	// Set command and tokens to spoof GetUsage() into providing help for that
	return h.withCommand(cmd).
		withTokens(tokens).
		GetUsage(err)
}

func (h Help) rootCmdHelp(err error) string {
	var sb = strings.Builder{}
	sb.WriteString(h.usageHeader(err))
	sb.WriteString(h.command().Help())
	sb.WriteString(h.globalOptionsHelp(HelpOpts{
		indent: strings.Repeat(Indent, 2),
	}))
	return sb.String()
}

func (h Help) branchCmdHelp(err error) string {
	var sb = strings.Builder{}
	cmd := h.command()
	if errors.Is(err, ErrNoExecFuncFound) {
		err = fmt.Errorf("there is no '%s' command, but there are these commands", cmd.Name)
	}
	sb.WriteString(h.usageHeader(err))
	sb.WriteString(cmd.Help())
	sb.WriteString(h.globalOptionsHelp(HelpOpts{
		indent: strings.Repeat(Indent, 2),
	}))
	return sb.String()
}

func (h Help) leafCmdHelp(err error) string {
	var sb = strings.Builder{}
	indent := strings.Repeat(Indent, 4)
	sb.WriteString(h.usageHeader(err))
	cmd := h.command()
	sb.WriteString(cmd.SignatureHelp())
	sb.WriteString(leafItemsHelp(cmd.Flags, HelpOpts{
		indent:         indent,
		label:          "Options",
		includeDefault: true,
	}))
	sb.WriteString(leafItemsHelp(cmd.Args, HelpOpts{
		indent: indent,
		label:  "Args",
		width:  len(Indent) + len("Default"),
	}))
	sb.WriteString(h.globalOptionsHelp(HelpOpts{
		indent: indent,
	}))
	return sb.String()
}

func (h Help) globalOptionsHelp(opts HelpOpts) string {
	opts.label = "Global Options"
	return leafItemsHelp(RootCmd.Flags, opts)
}

func (h Help) numSubCommands() int {
	return len(h.subCommands())
}

func (h Help) subCommands() CommandMap {
	return h.invoker.SubCommands()
}

func (h Help) isRootCommand() bool {
	return h.command() == RootCmd
}

func (h Help) AppName() string {
	return h.invoker.AppName
}

func (h Help) usageHeader(err error) string {
	var sb = strings.Builder{}
	if err == nil {
		sb.WriteString("\n")
	} else {
		msg := lib.UpperFirst(err.Error())
		sb.WriteString(fmt.Sprintf("\nERROR: %s:\n\n", msg))
	}
	sb.WriteString(fmt.Sprintf("%sUsage: %s [<options>] <command> [<args>]\n\n", Indent, h.AppName()))
	if h.numSubCommands() > 0 {
		sb.WriteString(fmt.Sprintf("%sCommands:\n\n", Indent))
	} else {
		sb.WriteString(fmt.Sprintf("%sCommand:\n\n", Indent))
	}
	return sb.String()
}

func (h Help) commandType(err error) (ct CommandType) {
	if errors.Is(err, ErrHelpSentinel) {
		ct = HelpCommand
		goto end
	}
	if h.isRootCommand() {
		ct = RootCommand
		goto end
	}
	if h.numSubCommands() == 0 {
		ct = LeafCommand
		goto end
	}
	ct = BranchCommand
end:
	return ct
}

type items interface {
	DisplayWidth(int) int
	Len() int
	Helpers() []helper
}

type HelpOpts struct {
	label          string
	indent         string
	width          int
	signature      string
	includeDefault bool
}
type helper interface {
	Help(HelpOpts) string
}

func leafItemsHelp(items items, opts HelpOpts) (help string) {
	opts.width = items.DisplayWidth(opts.width)
	if items.Len() > 0 {
		help += fmt.Sprintf("\n%s%s:\n\n", opts.indent, opts.label)
		for _, helper := range items.Helpers() {
			help += helper.Help(opts)
		}
	}
	return help
}
