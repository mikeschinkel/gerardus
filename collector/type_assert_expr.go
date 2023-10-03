package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type TypeAssertExpr struct {
	Expr Expr
	Type Expr
}

func (e TypeAssertExpr) String() (s string) {
	return fmt.Sprintf("%s.(%s)", e.Expr, e.Type)
}

func (c *Collector) CollectTypeAssertExpr(ctx context.Context, atae *ast.TypeAssertExpr) (tae TypeAssertExpr, err error) {
	tae.Expr, err = c.CollectExpr(ctx, atae.X)
	tae.Type, err = c.CollectExpr(ctx, atae.Type)
	return tae, err
}
