package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type CompositeLit struct {
	Name       Expr
	Properties ExprList
}

func (e CompositeLit) String() (s string) {
	return fmt.Sprintf("%s{%s}", e.Name, e.Properties)
}

func (c *Collector) CollectCompositeLit(ctx context.Context, acl *ast.CompositeLit) (cl CompositeLit, err error) {
	cl.Properties, err = c.CollectExprSlice(ctx, acl.Elts)
	if err != nil {
		goto end
	}
	if acl.Type == nil {
		goto end
	}
	cl.Name, err = c.CollectExpr(ctx, acl.Type)
end:
	return cl, err
}
