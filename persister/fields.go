package persister

import (
	"fmt"
	"strings"
)

type Fields []string

func NewFieldsFromString(flds string) Fields {
	return strings.Split(flds, ",")
}

func (fs Fields) Names() (s string) {
	sb := strings.Builder{}
	for _, f := range fs {
		sb.WriteString(f)
		sb.WriteByte(',')
	}
	s = sb.String()
	s = s[:len(s)-1]
	return s
}
func (fs Fields) PlaceHolders() (s string) {
	return strings.Repeat("?", len(fs))
}
func (fs Fields) DoUpdateSet() (s string) {
	sb := strings.Builder{}
	for _, f := range fs {
		f = fmt.Sprintf("%s=excluded.%s", f, f)
		sb.WriteString(f)
		sb.WriteByte(',')
	}
	s = sb.String()
	s = s[:len(s)-1]
	return s
}
