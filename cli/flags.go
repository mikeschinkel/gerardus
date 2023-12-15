package cli

import (
	"github.com/mikeschinkel/go-serr"
)

var _ items = (Flags)(nil)

type Flags []*Flag

func (flags Flags) Len() int {
	return len(flags)
}

func (flags Flags) Helpers() (helpers []helper) {
	helpers = make([]helper, len(flags))
	for i, flag := range flags {
		helpers[i] = flag
	}
	return helpers
}

func (flags Flags) DisplayWidth(minWidth int) (width int) {
	width = minWidth
	for _, flag := range flags {
		width = max(width, len(flag.Name))
	}
	return
}

// RequiresSatisfied ensures that values of .Requires are satisfied
func (flags Flags) RequiresSatisfied() (err error) {
	for _, flag := range flags {
		err = flag.RequiresSatisfied()
		if err != nil {
			goto end
		}
	}
end:
	return serr.Cast(err)
}

func (flags Flags) Validate(ctx Context) (err error) {
	for _, f := range flags {
		err = f.Validate(ctx)
		if err != nil {
			goto end
		}
	}
end:
	return err
}

func (flags Flags) SignatureHelp() (s string) {
	for _, flag := range flags {
		s += flag.SignatureHelp()
	}
	return s
}

// callSetValueFuncs calls the user-supplied f.Arg.SetValueFunc() for each flag
// that has been invoked for this command either from the CLI or that has a
// default, and passes that func the value f.Initialize() stored in the flag so
// that whatever value the user wanted to be initialized got initialized.
func (flags Flags) callSetValueFuncs() {
	for _, f := range flags {
		f.callSetValueFunc()
	}
}

// Initialize initializes the flag package flags by calling the flag package's
// flag.<Type>Var() function on a pointer to f.Arg.Value.<type> so that this flag
// 'f' will get the values passed on the command line, or the defaults if not
// passed.
func (flags Flags) Initialize() (err error) {
	for _, f := range flags {
		f.Initialize()
	}
	return err
}
