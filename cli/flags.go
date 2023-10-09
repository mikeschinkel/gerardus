package cli

import (
	"flag"
	"fmt"
)

type FlagUnion struct {
	String string
	Int    int
}

var flagValues = make(map[string]*FlagUnion, 32)

type Flags []*Flag

type Flag struct {
	Name             string
	Parent           *Command
	VarName          string
	Default          any
	Usage            string
	CheckFunc        func(any) error
	SetStringValFunc func(string)
	SetIntValFunc    func(int)
}

func (f *Flag) String() string {
	return fmt.Sprintf(" [-%s=<%s>]", f.Name, f.VarName)
}

// Unique returns a string that uniquely identifies a flag for its command
func (f *Flag) Unique() string {
	return fmt.Sprintf("%s:%s", f.Parent.Unique(), f.VarName)
}

func addFlags() (err error) {
	cmd, _ := InvokedCommand()
	if cmd == nil {
		err = fmt.Errorf("invalid command '%s'", CommandString())
		goto end
	}
	for _, f := range cmd.AllFlags() {
		fu := FlagUnion{}
		flagValues[f.Unique()] = &fu
		switch {
		case f.SetStringValFunc != nil:
			flag.StringVar(&fu.String, f.Name, f.Default.(string), f.Usage)
		case f.SetIntValFunc != nil:
			flag.IntVar(&fu.Int, f.Name, f.Default.(int), f.Usage)
		default:
			noSetFuncAssigned(f)
		}
	}
end:
	return err
}

func checkFlags() (err error) {
	cmd, _ := InvokedCommand()
	for _, f := range cmd.AllFlags() {
		if f.CheckFunc == nil {
			continue
		}
		switch {
		case f.SetStringValFunc != nil:
			err = f.CheckFunc(flagValues[f.Unique()].String)
		case f.SetIntValFunc != nil:
			err = f.CheckFunc(flagValues[f.Unique()].Int)
		default:
			noSetFuncAssigned(f)
		}
		if err != nil {
			goto end
		}
	}
end:
	return err
}
