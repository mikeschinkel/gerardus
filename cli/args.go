package cli

import (
	"reflect"
	"strings"

	"github.com/mikeschinkel/go-serr"
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

// EmptyStateSatisfied ensures that values of .Requires are satisfied
func (args Args) EmptyStateSatisfied() (err error) {
	for _, arg := range args {
		err = arg.EmptyStateSatisfied()
		if err != nil {
			goto end
		}
	}
end:
	return serr.Cast(err)
}

func (args Args) Validate(ctx Context) (err error) {
	for _, arg := range args {
		if arg.ValidateFunc == nil {
			continue
		}
		validateState := ArgValidation(arg.Requires)
		if validateState != MustValidate {
			continue
		}
		//goland:noinspection GoSwitchMissingCasesForIotaConsts
		switch arg.Type {
		case reflect.String:
			err = arg.ValidateFunc(ctx, arg.Value.string, &arg)
		case reflect.Int:
			err = arg.ValidateFunc(ctx, arg.Value.int, &arg)
		default:
			panicf("No func assigned to `ValidateFunc` for arg '%s'", arg.Unique())
		}
		if err != nil {
			err = ErrDoesNotValidate.Err(err, "arg_name", arg.Name)
			goto end
		}
	}
end:
	return serr.Cast(err)
}

func (args Args) CheckExistence(ctx Context) (err error) {
	var onSuccess string
	var value any

	for _, arg := range args {
		if arg.ExistsFunc == nil {
			continue
		}
		emptyState := ArgExistence(arg.Requires)
		if emptyState == IgnoreExists {
			continue
		}
		//goland:noinspection GoSwitchMissingCasesForIotaConsts
		switch arg.Type {
		case reflect.String:
			err = arg.ExistsFunc(ctx, arg.Value.string, &arg)
			value = arg.Value.string
		case reflect.Int:
			err = arg.ExistsFunc(ctx, arg.Value.int, &arg)
			value = arg.Value.int
		default:
			panicf("No func assigned to `ExistsFunc` for arg '%s'", arg.Unique())
		}
		onSuccess = arg.OnSuccess
		switch {
		case emptyState == MustExist && err != nil:
			err = ErrDoesNotExist.Err(err, "arg_name", arg.Name, "value", value)
		case emptyState == NotExist && err == nil:
			err = ErrAlreadyExists.Args("arg_name", arg.Name, "value", value)
		default:
			err = nil
			continue
		}
		goto end
	}
end:
	if err != nil && onSuccess != "" {
		err = serr.New(onSuccess).Err(err)
	}
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
