package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type UnaryExpr struct {
	Operator Operator
	Expr     Expr
}

func (e UnaryExpr) String() string {
	return fmt.Sprintf("%s", "")
}

func (c *Collector) CollectUnaryExpr(ctx context.Context, aue *ast.UnaryExpr) (ue UnaryExpr, err error) {
	ue.Operator = Operator(aue.Op)
	ue.Expr, err = c.CollectExpr(ctx, aue.X)
	return ue, err
}
