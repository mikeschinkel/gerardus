package parser

import (
	"go/ast"
	"go/types"
)

type FieldMap map[string]*Field
type Field struct {
	Name   string
	Type   string
	GoFile *GoFile
	field  *ast.Field
}

func NewField(name string, field *ast.Field, file *GoFile) *Field {
	return &Field{
		Name:   name,
		Type:   types.ExprString(field.Type),
		GoFile: file,
		field:  field,
	}
}
