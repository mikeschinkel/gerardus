package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mikeschinkel/go-lib"
)

type Help struct {
	invoker *CommandInvoker
}

func NewHelp(invoker *CommandInvoker) Help {
	return Help{
		invoker: invoker,
	}
}

func (h Help) Usage(err error, w io.Writer) {
	var help string

	StderrWriter = w

	switch h.commandType() {
	case Root:
		help = h.rootCmdHelp(err)
	case Branch:
		help = h.branchCmdHelp(err)
	case Leaf:
		help = h.leafCmdHelp(err)
	default:
		panicf("Undefined command type %d for '%s'",
			h.commandType(),
			h.command().Name,
		)
	}
	StdErr(help)
}

func (h Help) command() *Command {
	return h.invoker.Command
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
	msg := lib.UpperFirst(err.Error())
	sb.WriteString(fmt.Sprintf("\nERROR: %s:\n\n", msg))
	sb.WriteString(fmt.Sprintf("%sUsage: %s [<options>] <command> [<args>]\n\n", Indent, h.AppName()))
	if h.numSubCommands() > 0 {
		sb.WriteString(fmt.Sprintf("%sCommands:\n\n", Indent))
	} else {
		sb.WriteString(fmt.Sprintf("%sCommand:\n\n", Indent))
	}
	return sb.String()
}

func (h Help) commandType() (ct CommandType) {
	if h.isRootCommand() {
		ct = Root
		goto end
	}
	if h.numSubCommands() == 0 {
		ct = Leaf
		goto end
	}
	ct = Branch
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
