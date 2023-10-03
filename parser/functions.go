package parser

import (
	"go/ast"
	"go/types"
)

type Functions []*Function

type Function struct {
	Name       string
	Receiver   *Variable
	Parameters Variables
	Results    Variables
}

func NewFunction(name string, receiver *Variable, fd *ast.FuncDecl) *Function {
	f := &Function{
		Name:       name,
		Receiver:   receiver,
		Parameters: make(Variables, 0),
		Results:    make(Variables, 0),
	}
	f.collectVars(fd.Type)
	return f
}

func (f *Function) collectFields(list *ast.FieldList) (vars Variables) {
	var typeName string

	if list == nil {
		goto end
	}
	if list.List == nil {
		goto end
	}
	vars = make(Variables, len(list.List))
	if len(list.List) == 0 {
		goto end
	}
	for i, fld := range list.List {
		if fld.Type != nil {
			typeName = types.ExprString(fld.Type)
		}
		if fld.Names == nil {
			vars[i] = NewVariable(typeName, "")
			continue
		}
		for _, ident := range fld.Names {
			vars[i] = NewVariable(typeName, ident.Name)
		}
	}
end:
	if vars == nil {
		vars = make(Variables, 0)
	}
	return vars
}

func (f *Function) collectVars(ft *ast.FuncType) {
	f.Parameters = f.collectFields(ft.Params)
	f.Results = f.collectFields(ft.Results)
}
