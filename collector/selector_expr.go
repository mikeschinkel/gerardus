package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type SelectorExpr struct {
	Package Expr
	Name    string
}

func (se SelectorExpr) String() string {
	return fmt.Sprintf("%s.%s", se.Package, se.Name)
}

func (c *Collector) CollectSelectorExpr(ctx context.Context, e *ast.SelectorExpr) (se SelectorExpr, err error) {
	var pkg Expr
	pkg, err = c.CollectExpr(ctx, e.X)
	if err != nil {
		goto end
	}
	se = SelectorExpr{
		Package: pkg,
		Name:    e.Sel.Name,
	}
end:
	return se, err
}
