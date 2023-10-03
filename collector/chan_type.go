package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type ChanType struct {
	Direction Direction
	Type      Expr
}

func (e ChanType) String() (s string) {
	return fmt.Sprintf("chan%s %s", e.Direction, e.Type)
}

func (c *Collector) CollectChanType(ctx context.Context, act *ast.ChanType) (ct ChanType, err error) {
	ct.Direction = Direction(act.Dir)
	ct.Type, err = c.CollectExpr(ctx, act.Value)
	return ct, err
}
