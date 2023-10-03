package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type MapType struct {
	KeyType   Expr
	ValueType Expr
}

func (e MapType) String() string {
	return fmt.Sprintf("map[%s]%s", e.KeyType, e.ValueType)
}

func (c *Collector) CollectMapType(ctx context.Context, amt *ast.MapType) (mt MapType, err error) {
	mt.KeyType, err = c.CollectExpr(ctx, amt.Key)
	if err != nil {
		goto end
	}
	switch t := amt.Value.(type) {
	case *ast.FuncType:
		mt.ValueType, err = c.CollectFuncType(ctx, t)
	default:
		mt.ValueType, err = c.CollectExpr(ctx, amt.Value)
	}
end:
	return mt, err
}
