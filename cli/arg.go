package cli

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mikeschinkel/go-serr"
)

type ArgName string

type ArgRequires int

const (
	MustExist ArgRequires = 1 << iota
	IgnoreExists
	NotExist
	EmptyOk
	MustBeEmpty
	NotEmpty
	MustValidate
)

var _ helper = (*Arg)(nil)

type Arg struct {
	Name         ArgName
	Parent       *Command
	Usage        string
	Default      interface{}
	Optional     bool
	Variadic     bool
	ExistsFunc   func(Context, any, *Arg) error
	ValidateFunc func(Context, any, *Arg) error
	Type         reflect.Kind
	SetValueFunc func(*Value)
	Value        *Value
	Requires     ArgRequires
	SuccessMsg   string
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

// EmptyStateSatisfied ensures that values of .Requires for .Args are satisfied
//
//goland:noinspection GoUnusedParameter
func (arg Arg) EmptyStateSatisfied(ctx Context, tt TokenType) (err error) {
	e := ArgEmptiness(arg.Requires)
	isZero := arg.Value.IsZero()
	value := arg.Name
	name := string(tt) + "_name"
	if tt == FlagType {
		value = "-" + value
	}
	switch {
	case e == NotEmpty && isZero:
		err = ErrTokenValueCannotBeEmpty.Args(name, value)
		goto end
	case e == MustBeEmpty && !isZero:
		err = ErrTokenValueMustBeEmpty.Args(name, value)
		goto end
	}
end:
	if err != nil {
		err = ErrEmptyStateNotSatisfied.Err(err, "arg_name", arg.Name)
	}
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
func (arg Arg) MustExist() Arg {
	arg.Requires &= ^IgnoreExists
	arg.Requires &= ^NotExist
	arg.Requires |= MustExist
	return arg
}
func (arg Arg) IgnoreExists() Arg {
	arg.Requires &= ^MustExist
	arg.Requires &= ^NotExist
	arg.Requires |= IgnoreExists
	return arg
}
func (arg Arg) NotExist() Arg {
	arg.Requires &= ^MustExist
	arg.Requires &= ^IgnoreExists
	arg.Requires |= NotExist
	return arg
}
func (arg Arg) MustValidate() Arg {
	arg.Requires |= MustValidate
	return arg
}

func (arg Arg) ClearCheckFunc() Arg {
	arg.ExistsFunc = nil
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

func (arg Arg) CheckExistence(ctx Context) (err error) {
	var value any
	var successMsg string
	var emptyState ArgRequires

	if arg.ExistsFunc == nil {
		goto end
	}
	emptyState = ArgExistence(arg.Requires)
	if emptyState == IgnoreExists {
		goto end
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
	successMsg = arg.SuccessMsg
	switch {
	case emptyState == MustExist && err != nil:
		err = ErrDoesNotExist.Err(err, "arg_name", arg.Name, "value", value)
	case emptyState == NotExist && err == nil:
		err = ErrAlreadyExists.Args("arg_name", arg.Name, "value", value)
	default:
		err = nil
	}
end:
	if err != nil && successMsg != "" {
		err = serr.New(successMsg).Err(err)
	}
	return err
}

func (arg Arg) ValidateByFunc(ctx Context) (err error) {
	var validateState ArgRequires

	if arg.ValidateFunc == nil {
		goto end
	}
	validateState = ArgValidation(arg.Requires)
	if validateState != MustValidate {
		goto end
	}
	//goland:noinspection GoSwitchMissingCasesForIotaConsts
	switch arg.Type {
	case reflect.String:
		err = arg.ValidateFunc(ctx, arg.Value.string, &arg)
	case reflect.Int:
		err = arg.ValidateFunc(ctx, arg.Value.int, &arg)
	default:
		panicf("No func assigned to `ValidateFunc` for '%s'", arg.Name)
	}
end:
	if err != nil {
		err = ErrDoesNotValidate.Err(err, "arg_name", arg.Name)
	}
	return err
}
