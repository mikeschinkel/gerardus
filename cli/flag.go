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
	*Arg
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
	//flagsMap[f.Unique()] = fu
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

func (f *Flag) Validate() (err error) {
	if f.CheckFunc == nil {
		goto end
	}
	switch f.Type {
	case reflect.String:
		err = f.CheckFunc(f.Requires, f.Value.string)
	case reflect.Int:
		err = f.CheckFunc(f.Requires, f.Value.int)
	default:
		f.noSetFuncAssigned()
	}
end:
	return serr.Cast(err)
}
