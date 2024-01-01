package cli

var _ items = (Flags)(nil)

type Flags []Flag

func (flags Flags) Len() int {
	return len(flags)
}

func (flags Flags) Remove(n int) Flags {
	var end, i int
	newFlags := flags
	if n < 0 {
		goto end
	}
	end = len(newFlags) - 1
	for i = n; i < end; i++ {
		newFlags[i] = newFlags[i+1]
	}
	newFlags = newFlags[:end]
end:
	return newFlags
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
func (flags Flags) Initialize(ctx Context) Flags {
	for i, f := range flags {
		flags[i] = f.Initialize(ctx)
	}
	return flags
}
