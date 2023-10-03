package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type BinaryExpr struct {
	LeftExpr  Expr
	Operator  Operator
	RightExpr Expr
}

func (e BinaryExpr) String() string {
	return fmt.Sprintf("%s%s%s",
		e.LeftExpr,
		e.Operator,
		e.RightExpr,
	)
}

func (c *Collector) CollectBinaryExpr(ctx context.Context, abe *ast.BinaryExpr) (be BinaryExpr, err error) {
	be.LeftExpr, err = c.CollectExpr(ctx, abe.X)
	be.Operator = Operator(abe.Op)
	be.RightExpr, err = c.CollectExpr(ctx, abe.Y)
	return be, err
}
