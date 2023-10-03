package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type StructType struct {
	Fields FieldList
}

func (e StructType) String() (s string) {
	s = fmt.Sprintf("struct{%s}", e.Fields)
	return s
}

func (c *Collector) CollectStructType(ctx context.Context, t *ast.StructType) (st StructType, err error) {
	st.Fields, err = c.CollectFieldSlice(ctx, t.Fields.List)
	return st, err
}
