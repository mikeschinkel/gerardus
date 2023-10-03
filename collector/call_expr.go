package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type CallExpr struct {
	Call Expr
	Args Expr
}

func (ce CallExpr) String() string {
	return fmt.Sprintf("%s(%s)", ce.Call, ce.Args)
}

func (c *Collector) CollectCallExpr(ctx context.Context, call *ast.CallExpr) (ce CallExpr, err error) {
	ce.Call, err = c.CollectExpr(ctx, call.Fun)
	if err != nil {
		goto end
	}
	ce.Args, err = c.CollectExprSlice(ctx, call.Args)
end:
	return ce, err
}
