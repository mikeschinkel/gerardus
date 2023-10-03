package collector

import (
	"context"
	"fmt"
	"go/ast"
	"strings"
)

type Field struct {
	Name string
	Type Expr
}

func (f Field) String() (s string) {
	if f.Name == "" {
		return f.Type.String()
	}
	s = fmt.Sprintf("%s %s", f.Name, f.Type)
	return s
}

type FieldList []Field

func (fl FieldList) String() (s string) {
	var sb *strings.Builder
	if len(fl) == 0 {
		goto end
	}
	sb = &strings.Builder{}
	for _, f := range fl {
		sb.WriteString(f.String())
		sb.WriteByte(',')
	}
	s = sb.String()
	s = s[:len(s)-1]
end:
	return s
}

func (c *Collector) CollectFieldList(ctx context.Context, list *ast.FieldList) (FieldList, error) {
	return c.CollectFieldSlice(ctx, list.List)
}

func (c *Collector) CollectFieldSlice(ctx context.Context, list []*ast.Field) (fl FieldList, err error) {
	var flds FieldList

	// Ultimate field list will be at least len(list) in size, so start there.
	fl = make(FieldList, 0, len(list))
	for _, fldss := range list {
		flds, err = c.CollectField(ctx, fldss)
		if err != nil {
			goto end
		}
		for _, fld := range flds {
			fl = append(fl, fld)
		}
	}
end:
	return fl, err
}

func (c *Collector) CollectField(ctx context.Context, aFld *ast.Field) (list FieldList, err error) {
	var typ Expr
	typ, err = c.CollectExpr(ctx, aFld.Type)
	if err != nil {
		goto end
	}
	if aFld.Type != nil {
		typ, err = c.CollectExpr(ctx, aFld.Type)
	} else {
		print("Wow. aFld.Type is nil! What now?!?")
	}
	if err != nil {
		goto end
	}
	if len(aFld.Names) == 0 {
		list = make(FieldList, 1)
		list[0] = Field{Type: typ}
		goto end
	}
	list = make(FieldList, len(aFld.Names))
	for i, ident := range aFld.Names {
		list[i] = Field{
			Name: ident.Name,
			Type: typ,
		}
	}
end:
	return list, err
}

func (c *Collector) astHasFields(list *ast.FieldList) (has bool) {
	if list == nil {
		goto end
	}
	if list.List == nil {
		goto end
	}
	if len(list.List) == 0 {
		goto end
	}
	has = true
end:
	return has
}
