package collector

type ValueType int

const (
	UnsetValue ValueType = 0
	VarValue   ValueType = 1
	ConstValue ValueType = 2
)

func (vt ValueType) String() (s string) {
	switch vt {
	case VarValue:
		s = "var"
	case ConstValue:
		s = "const"
	default:
		panicf("Invalid ValueType '%d'", vt)
	}
	return s
}
