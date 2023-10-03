package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type SliceExpr struct {
	Expr Expr
	Low  Expr
	High Expr
}

func (c *Collector) CollectSliceExpr(ctx context.Context, ase *ast.SliceExpr) (se SliceExpr, err error) {
	se.Expr, err = c.CollectExpr(ctx, ase.X)
	if err != nil {
		goto end
	}
	switch {
	case ase.High != nil && ase.Low != nil:
		se.Low, err = c.CollectExpr(ctx, ase.Low)
		if err != nil {
			goto end
		}
		se.High, err = c.CollectExpr(ctx, ase.High)
	case ase.Low != nil:
		se.Low, err = c.CollectExpr(ctx, ase.Low)
	case ase.High != nil:
		se.High, err = c.CollectExpr(ctx, ase.High)
	}
end:
	return se, err
}
func (se SliceExpr) String() (s string) {
	switch {
	case se.High != nil && se.Low != nil:
		s = fmt.Sprintf("%s[%s:%s]", se.Expr, se.Low, se.High)
	case se.Low != nil:
		s = fmt.Sprintf("%s[%s:]", se.Expr, se.Low)
	case se.High != nil:
		s = fmt.Sprintf("%s[:%s]", se.Expr, se.High)
	}
	return s
}
