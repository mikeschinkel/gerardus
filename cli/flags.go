package cli

import (
	"github.com/mikeschinkel/go-serr"
)

var _ items = (Flags)(nil)

type Flags []*Flag

func (flags Flags) Len() int {
	return len(flags)
}

func (flags Flags) Index(name ArgName) (n int) {
	n = -1
	for i, flag := range flags {
		if flag.Name != name {
			continue
		}
		n = i
		goto end
	}
end:
	return n
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

// EmptyStateSatisfied ensures that values of .Requires are satisfied
func (flags Flags) EmptyStateSatisfied() (err error) {
	for _, flag := range flags {
		err = flag.EmptyStateSatisfied()
		if err != nil {
			goto end
		}
	}
end:
	return serr.Cast(err)
}

func (flags Flags) Validate(ctx Context) (err error) {
	for _, f := range flags {
		err = f.CheckExists(ctx)
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
func (flags Flags) callSetValueFuncs() Flags {
	for i, f := range flags {
		flags[i].Arg = callSetArgValueFunc(f.Arg)
	}
	return flags
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
