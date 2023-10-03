package collector

import (
	"context"
	"go/ast"
)

type Idents []Ident

type Ident struct {
	Name Expr
}

func (e Ident) String() string {
	return e.Name.String()
}

//goland:noinspection GoUnusedParameter
func (c *Collector) CollectIdent(ctx context.Context, e *ast.Ident) (Ident, error) {
	return Ident{Name: String(e.Name)}, nil
}

func (c *Collector) CollectIdentSlice(ctx context.Context, astExprs []ast.Expr) (is Idents, err error) {
	is = make(Idents, len(astExprs))
	for i, expr := range astExprs {
		is[i], err = c.CollectIdent(ctx, expr.(*ast.Ident))
		if err != nil {
			goto end
		}
	}
end:
	return is, err
}
