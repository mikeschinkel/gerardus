package cli

import (
	"fmt"
)

type FlagUnion struct {
	String string
	Int    int
}

type FlagValuesMap map[string]*FlagUnion

var flagValuesMap = make(FlagValuesMap, 32)

type Flags []*Flag
type Flag struct {
	Switch string
	Arg
}

func (f *Flag) String() string {
	return fmt.Sprintf(" [-%s=<%s>]", f.Switch, f.Name)
}

// Unique returns a string that uniquely identifies a flag for its command
func (f *Flag) Unique() string {
	return fmt.Sprintf("%s:%s", f.Parent.Unique(), f.Name)
}

func (flgs Flags) validate() (err error) {
	cmd, _, err := InvokedCommand()
	if err != nil {
		goto end
	}
	for _, f := range cmd.InvokedFlags() {
		if f.CheckFunc == nil {
			continue
		}
		switch {
		case f.SetStringValFunc != nil:
			err = f.CheckFunc(flagValuesMap[f.Unique()].String)
		case f.SetIntValFunc != nil:
			err = f.CheckFunc(flagValuesMap[f.Unique()].Int)
		default:
			f.noSetFuncAssigned()
		}
		if err != nil {
			goto end
		}
	}
end:
	return err
}

// InvokedFlagValuesMap returns a map of the invoked flags
func (f *Flag) InvokedFlagValuesMap() (m FlagValuesMap, err error) {
	var flags Flags

	cmd, _, err := InvokedCommand()
	if err != nil {
		goto end
	}
	flags = cmd.InvokedFlags()
	m = make(FlagValuesMap)
	for _, flg := range flags {
		m[flg.Switch] = flagValuesMap[flg.Unique()]
	}
end:
	return m, err
}
