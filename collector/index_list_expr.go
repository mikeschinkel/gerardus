package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type IndexListExpr struct {
	Expr    Expr
	Indices Expr
}

func (c *Collector) CollectIndexListExpr(ctx context.Context, aile *ast.IndexListExpr) (ile IndexListExpr, err error) {
	ile.Expr, err = c.CollectExpr(ctx, aile.X)
	ile.Indices, err = c.CollectExprSlice(ctx, aile.Indices)
	return ile, err
}
func (e IndexListExpr) String() (s string) {
	return fmt.Sprintf("%s[%s]", e.Expr, e.Indices)
}
