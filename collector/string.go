package collector

type String string

func (s String) String() string {
	return string(s)
}
