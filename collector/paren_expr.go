package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type ParenExpr struct {
	Expr Expr
}

func (e ParenExpr) String() string {
	return fmt.Sprintf("(%s)", e.Expr)
}

func (c *Collector) CollectParenExpr(ctx context.Context, ape *ast.ParenExpr) (pe ParenExpr, err error) {
	pe.Expr, err = c.CollectExpr(ctx, ape.X)
	return pe, err
}
