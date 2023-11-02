package parser

import (
	"go/ast"
	"go/types"
)

type MethodMap map[string]*Method
type Method struct {
	Name    string
	GoFile  *GoFile
	field   *ast.Field
	params  Variables
	results Variables
}

func NewMethod(name string, field *ast.Field, file *GoFile) *Method {
	m := &Method{
		Name:    name,
		GoFile:  file,
		field:   field,
		params:  make(Variables, 0),
		results: make(Variables, 0),
	}
	m.collectResults()
	m.collectParams()
	return m
}

type fieldListFunc func(*ast.FuncType) *ast.FieldList

func (m *Method) collectParams() {
	m.params = m.getVars(m.params, func(ft *ast.FuncType) *ast.FieldList {
		return ft.Params
	})
}
func (m *Method) collectResults() {
	m.results = m.getVars(m.results, func(ft *ast.FuncType) *ast.FieldList {
		return ft.Results
	})
}

func (m *Method) IsEmbed() bool {
	return IsSelectorExpr(m.field.Type)
}

func (m *Method) IsPrivate() (private bool) {
	switch t := m.field.Type.(type) {
	case *ast.FuncType:
		private = !isPublicName(m.field.Names[0].Name)
	case *ast.SelectorExpr:
		private = !isPublicName(t.Sel.Name)
	default:
		panicf("Is this method private or not? %v", m)
	}
	return private
}

func (m *Method) getVars(varMap Variables, listFunc fieldListFunc) (vm Variables) {
	var funcType *ast.FuncType
	var list []*ast.Field
	var ok bool
	var lister *ast.FieldList

	vm = make(Variables, 0)
	if varMap == nil {
		goto end
	}
	funcType, ok = m.field.Type.(*ast.FuncType)
	if !ok {
		goto end
	}
	lister = listFunc(funcType)
	if lister == nil {
		goto end
	}
	list = lister.List
	if list == nil {
		goto end
	}
	vm = make(Variables, len(list))
	for i, field := range list {
		vm[i] = &Variable{
			Type: types.ExprString(field.Type),
		}
	}
end:
	return vm
}
