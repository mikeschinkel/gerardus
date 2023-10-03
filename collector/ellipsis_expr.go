package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type Ellipsis struct {
	Type Expr
}

func (e Ellipsis) String() string {
	if e.Type == nil {
		return "..."
	}
	return fmt.Sprintf("...%s", e.Type.String())
}

func (c *Collector) CollectEllipsis(ctx context.Context, ae *ast.Ellipsis) (e Ellipsis, err error) {
	if ae.Elt != nil {
		e.Type, err = c.CollectExpr(ctx, ae.Elt)
	}
	return e, err
}
