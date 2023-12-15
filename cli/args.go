package cli

import (
	"reflect"
	"strings"

	"github.com/mikeschinkel/go-serr"
)

type Args []*Arg

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

// RequiresSatisfied ensures that values of .Requires are satisfied
func (args Args) RequiresSatisfied() (err error) {
	for _, arg := range args {
		err = arg.RequiresSatisfied()
		if err != nil {
			goto end
		}
	}
end:
	return serr.Cast(err)
}

func (args Args) Validate(ctx Context) (err error) {
	for _, arg := range args {
		if arg.CheckFunc == nil {
			continue
		}
		//goland:noinspection GoSwitchMissingCasesForIotaConsts
		switch arg.Type {
		case reflect.String:
			err = arg.CheckFunc(ctx, arg.Requires, arg.Value.string)
		case reflect.Int:
			err = arg.CheckFunc(ctx, arg.Requires, arg.Value.int)
		default:
			arg.noSetFuncAssigned()
		}
		if err != nil {
			goto end
		}
	}
end:
	return serr.Cast(err)
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
