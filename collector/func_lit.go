package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type FuncLit struct {
	FuncType Expr
}

func (e FuncLit) String() string {
	return fmt.Sprintf("func%s", e.FuncType)
}

func (c *Collector) CollectFuncLit(ctx context.Context, fl *ast.FuncLit) (_ FuncLit, err error) {
	ft, err := c.CollectFuncType(ctx, fl.Type)
	return FuncLit{FuncType: ft}, err
}
