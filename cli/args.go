package cli

import (
	"strings"
)

type Args []Arg

var _ items = (Args)(nil)

func (args Args) DisplayWidth(minWidth int) (width int) {
	width = minWidth
	for _, arg := range args {
		width = max(width, len(arg.Name))
	}
	return width
}

func (args Args) Len() int {
	return len(args)
}

func (args Args) Helpers() (helpers []helper) {
	helpers = make([]helper, len(args))
	for i, flag := range args {
		helpers[i] = flag
	}
	return helpers
}

func (args Args) SignatureHelp() (help string) {
	optCnt := 0
	for _, arg := range args {
		help += arg.SignatureHelp()
		if arg.Optional {
			optCnt++
		}
	}
	if optCnt > 0 {
		help += strings.Repeat("]", optCnt)
	}
	return help
}

func (args Args) String() (s string) {
	sb := strings.Builder{}
	if len(args) == 0 {
		goto end
	}
	for _, arg := range args {
		sb.WriteString(string(arg.Name))
		sb.WriteByte(' ')
	}
	s = sb.String()
	s = s[:len(s)-1]
end:
	return s
}
