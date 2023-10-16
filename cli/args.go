package cli

import (
	"fmt"
	"strconv"
	"strings"
)

type Args []*Arg
type ArgsMap map[string]*Arg

type Arg struct {
	Name             string
	Parent           *Command
	Usage            string
	Default          interface{}
	Optional         bool
	CheckFunc        func(any) error
	SetStringValFunc func(string)
	SetIntValFunc    func(int)
	Value            ValueUnion
}

func (arg Arg) IsZero() bool {
	switch {
	case arg.SetStringValFunc != nil:
		return len(arg.Value.String) == 0
	case arg.SetIntValFunc != nil:
		return arg.Value.Int == 0
	default:
		panicf("Unhandled type for arg '%s'", arg.Unique())
	}
	return false
}

func (arg Arg) String() string {
	switch {
	case arg.SetStringValFunc != nil:
		return arg.Value.String
	case arg.SetIntValFunc != nil:
		return strconv.Itoa(arg.Value.Int)
	default:
		panicf("Unhandled type for arg '%s'", arg.Unique())
	}
	return ""
}

func (arg Arg) noSetFuncAssigned() {
	panicf("No func(<type>) assigned to property `Set*ValFunc` for arg '%s'", arg.Unique())
}

func (arg Arg) Unique() string {
	return fmt.Sprintf("%s:%s", arg.Parent.Unique(), arg.Name)
}

func (m ArgsMap) String() (s string) {
	sb := strings.Builder{}
	if len(m) == 0 {
		goto end
	}
	for _, arg := range m {
		sb.WriteString(arg.Name)
		sb.WriteByte(' ')
	}
	s = sb.String()
	s = s[:len(s)-1]
end:
	return s
}

func (m ArgsMap) Validate() (err error) {
	for _, arg := range m {
		if arg.CheckFunc == nil {
			continue
		}
		switch {
		case arg.SetStringValFunc != nil:
			err = arg.CheckFunc(arg.Value.String)
		case arg.SetIntValFunc != nil:
			err = arg.CheckFunc(arg.Value.Int)
			//default:
			//	arg.noSetFuncAssigned()
		}
		if err != nil {
			goto end
		}
	}
end:
	return err
}
