package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type StarExpr struct {
	Expr Expr
}

func (e StarExpr) String() string {
	return fmt.Sprintf("*%s", e.Expr)
}

func (c *Collector) CollectStarExpr(ctx context.Context, ase *ast.StarExpr) (se StarExpr, err error) {
	se.Expr, err = c.CollectExpr(ctx, ase.X)
	return se, err
}
