package collector

import (
	"context"
	"go/ast"
)

type BasicLit struct {
	Value Expr
}

func (e BasicLit) String() string {
	return e.Value.String()
}

//goland:noinspection GoUnusedParameter
func (c *Collector) CollectBasicLit(ctx context.Context, e *ast.BasicLit) (BasicLit, error) {
	return BasicLit{Value: String(e.Value)}, nil
}
