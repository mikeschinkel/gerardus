package cli

import (
	"fmt"
	"reflect"
	"strings"
)

type ArgName string

type ArgRequires int

const (
	MustPassCheck ArgRequires = 1 << iota
	IgnoreCheck
	MustFailCheck
	EmptyOk
	MustBeEmpty
	NotEmpty
)

var _ helper = (*Arg)(nil)

type Arg struct {
	Name         ArgName
	Parent       *Command
	Usage        string
	Default      interface{}
	Optional     bool
	CheckFunc    func(Context, any, *Arg) error
	Type         reflect.Kind
	SetValueFunc func(*Value)
	Value        *Value
	Requires     ArgRequires
	Message      string
}

func NewArg(arg Arg) Arg {
	if arg.SetValueFunc == nil {
		arg.SetValueFunc = func(*Value) {}
	}
	if arg.Type == reflect.Invalid {
		arg.Type = reflect.String
	}
	if arg.Value == nil {
		arg.Value = &Value{Type: arg.Type}
	}
	return arg
}

func (arg Arg) Check(requires ArgRequires) bool {
	return arg.Requires&requires != 0
}

func (arg Arg) EmptyStateSatisfied() (err error) {
	e := ArgEmptiness(arg.Requires)
	isZero := arg.Value.IsZero()
	name := fmt.Sprintf("<%s>", arg.Name)
	switch {
	case e == NotEmpty && isZero:
		err = ErrArgCannotBeEmpty.Args("arg_name", name)
		goto end
	case e == MustBeEmpty && !isZero:
		err = ErrArgMustBeEmpty.Args("arg_name", name)
		goto end
	}
end:
	return err
}

func (arg Arg) IsZero() bool {
	switch arg.Type {
	case reflect.String:
		return len(arg.Value.string) == 0
	case reflect.Int:
		return arg.Value.int == 0
	default:
		panicf("Unhandled type for arg '%s'", arg.Unique())
	}
	return false
}

func (arg Arg) EmptyOk() Arg {
	arg.Requires &= ^NotEmpty
	arg.Requires &= ^MustBeEmpty
	arg.Requires |= EmptyOk
	return arg
}
func (arg Arg) NotEmpty() Arg {
	arg.Requires &= ^EmptyOk
	arg.Requires &= ^MustBeEmpty
	arg.Requires |= NotEmpty
	return arg
}
func (arg Arg) MustBeEmpty() Arg {
	arg.Requires &= ^EmptyOk
	arg.Requires &= ^NotEmpty
	arg.Requires |= MustBeEmpty
	return arg
}
func (arg Arg) MustPassCheck() Arg {
	arg.Requires &= ^IgnoreCheck
	arg.Requires &= ^MustFailCheck
	arg.Requires |= MustPassCheck
	return arg
}
func (arg Arg) IgnoreCheck() Arg {
	arg.Requires &= ^MustPassCheck
	arg.Requires &= ^MustFailCheck
	arg.Requires |= IgnoreCheck
	return arg
}
func (arg Arg) MustFailCheck() Arg {
	arg.Requires &= ^MustPassCheck
	arg.Requires &= ^IgnoreCheck
	arg.Requires |= MustFailCheck
	return arg
}

func (arg Arg) ClearCheckFunc() Arg {
	arg.CheckFunc = nil
	return arg
}

func (arg Arg) DefaultHelp(opts HelpOpts) (help string) {
	space := strings.Repeat(" ", opts.width-len("Default")-len(Indent))
	return fmt.Sprintf("%s%s%sDefault: %s%s\n",
		opts.indent,
		Indent,
		Indent,
		space,
		arg.Default,
	)
}

func (arg Arg) Help(opts HelpOpts) (help string) {
	opts.signature = string(arg.Name)
	return arg.help(opts)
}

func (arg Arg) help(opts HelpOpts) string {
	sb := strings.Builder{}
	sb.WriteString(opts.indent)
	sb.WriteString(Indent)
	sb.WriteString(opts.signature)
	sb.WriteByte(':')
	if arg.Usage != "" {
		sb.WriteByte(' ')
		sb.WriteString(strings.Repeat(" ", opts.width-len(arg.Name)))
		sb.WriteString(arg.Usage)
	}
	sb.WriteByte('\n')
	if arg.Default == nil {
		goto end
	}
	if !opts.includeDefault {
		goto end
	}
	if opts.signature[0] == '-' {
		// This is a hack to better align "Default:" with flag signature
		opts.indent = opts.indent[1:]
	}
	sb.WriteString(arg.DefaultHelp(opts))
end:
	return sb.String()
}

func (arg Arg) String() string {
	return arg.Value.String()
}

func (arg Arg) Unique() string {
	return fmt.Sprintf("%s:%s", arg.Parent.Unique(), arg.Name)
}

func (arg Arg) SignatureHelp() (s string) {
	if arg.Optional {
		s = fmt.Sprintf(" [<%s>", arg.Name)
	} else {
		s = fmt.Sprintf(" <%s>", arg.Name)
	}
	return s
}

func (arg Arg) noSetFuncAssigned() {
	panicf("No func(<type>) assigned to property `Set<type>ValFunc` for arg '%s'", arg.Unique())
}

// callSetArgValueFunc sets the Value for one arg
func callSetArgValueFunc(arg Arg) Arg {
	arg.SetValueFunc(arg.Value)
	return arg
}
