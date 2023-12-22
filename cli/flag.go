package cli

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/mikeschinkel/go-serr"
)

var _ helper = (*Flag)(nil)

type Flag struct {
	Switch string
	Arg
}

func (f *Flag) noSetFuncAssigned() {
	panicf("No func(<type>) assigned to property `Set<type>ValFunc` for flag '%s'", f.Unique())
}

func (f *Flag) Initialize() {
	fu := &Value{}
	switch f.Type {
	case reflect.String:
		flag.StringVar(&fu.string, f.Switch, f.Default.(string), f.Usage)
	case reflect.Int:
		flag.IntVar(&fu.int, f.Switch, f.Default.(int), f.Usage)
	default:
		f.noSetFuncAssigned()
	}
	f.Value = fu
}

func (f *Flag) String() string {
	return fmt.Sprintf(" [-%s=<%s>]", f.Switch, f.Name)
}

// Unique returns a string that uniquely identifies a flag for its command
func (f *Flag) Unique() string {
	return fmt.Sprintf("%s:%s", f.Parent.Unique(), f.Name)
}

func (f *Flag) Help(opts HelpOpts) (help string) {
	opts.signature = f.signature()
	return f.help(opts)
}

func (f *Flag) SignatureHelp() string {
	return fmt.Sprintf(" [%s]", f.signature())
}

func (f *Flag) signature() string {
	return fmt.Sprintf("-%s=<%s>", f.Switch, f.Name)
}

func (f *Flag) CheckExists(ctx Context) (err error) {
	if f.ExistsFunc == nil {
		goto end
	}
	switch f.Type {
	case reflect.String:
		err = f.ExistsFunc(ctx, f.Value.string, &f.Arg)
	case reflect.Int:
		err = f.ExistsFunc(ctx, f.Value.int, &f.Arg)
	default:
		f.noSetFuncAssigned()
	}
end:
	if err != nil && f.SuccessMsg != "" {
		err = serr.New(f.SuccessMsg).Err(err)
	}
	return serr.Cast(err)
}
