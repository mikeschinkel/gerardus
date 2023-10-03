package parser

type Variables []*Variable
type Variable struct {
	Name string
	Type string
}

func NewVariable(t string, name string) *Variable {
	return &Variable{
		Name: name,
		Type: t,
	}
}
