package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type IndexExpr struct {
	Expr  Expr
	Index Expr
}

func (c *Collector) CollectIndexExpr(ctx context.Context, aie *ast.IndexExpr) (ie IndexExpr, err error) {
	ie.Expr, err = c.CollectExpr(ctx, aie.X)
	if err != nil {
		goto end
	}
	ie.Index, err = c.CollectExpr(ctx, aie.Index)
end:
	return ie, err
}
func (e IndexExpr) String() (s string) {
	return fmt.Sprintf("%s[%s]", e.Expr, e.Index)
}
