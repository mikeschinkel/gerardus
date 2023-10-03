package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type InterfaceType struct {
	Methods FieldList
}

func (e InterfaceType) String() (s string) {
	return fmt.Sprintf("interface{%s}", e.Methods)
}

func (c *Collector) CollectInterfaceType(ctx context.Context, iface *ast.InterfaceType) (it InterfaceType, err error) {
	it.Methods, err = c.CollectFieldList(ctx, iface.Methods)
	return it, err
}
