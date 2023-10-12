package cli

import (
	"fmt"
	"strconv"
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
}

func (arg *Arg) noSetFuncAssigned() {
	panicf("No func(<type>) assigned to property `Set*ValFunc` for arg '%s'", arg.Unique())
}

func (arg *Arg) Unique() string {
	return fmt.Sprintf("%s:%s", arg.Parent.Unique(), arg.Name)
}

func (m ArgsMap) validate(sm StringMap) (err error) {
	for name, arg := range m {
		if arg.CheckFunc == nil {
			continue
		}
		switch {
		case arg.SetStringValFunc != nil:
			err = arg.CheckFunc(sm[name])
		case arg.SetIntValFunc != nil:
			var n int
			n, err = strconv.Atoi(sm[name])
			if err != nil {
				goto end
			}
			err = arg.CheckFunc(n)
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
