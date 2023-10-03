package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type KeyValueExpr struct {
	Key   Expr
	Value Expr
}

func (e KeyValueExpr) String() string {
	return fmt.Sprintf("%s:%s", e.Key, e.Value)
}

func (c *Collector) CollectKeyValueExpr(ctx context.Context, e *ast.KeyValueExpr) (kve KeyValueExpr, err error) {
	var key, value Expr

	key, err = c.CollectExpr(ctx, e.Key)
	if err != nil {
		goto end
	}
	value, err = c.CollectExpr(ctx, e.Value)
	if err != nil {
		goto end
	}
	kve = KeyValueExpr{
		Key:   key,
		Value: value,
	}
end:
	return kve, err
}
