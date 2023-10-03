package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type ValueSpecs []ValueSpec
type ValueSpec struct {
	File      File
	Name      string
	Type      Expr
	ValueType ValueType
}

func (spec ValueSpec) String() (s string) {
	return fmt.Sprintf("%s %s = %s",
		spec.ValueType,
		spec.Name,
		spec.Type,
	)
}

func (spec ValueSpec) CodeFacet() {}

func (c *Collector) CollectValueSpec(ctx context.Context, avs *ast.ValueSpec, typ *Expr) (vs ValueSpecs, err error) {
	var thisType Expr
	if avs.Type != nil {
		// This is for when iota is used to define constants
		// This sets args the first time through after which it continues using it
		*typ, err = c.CollectExpr(ctx, avs.Type)
	}
	if err != nil {
		goto end
	}
	vs = make(ValueSpecs, len(avs.Names))
	for i, ident := range avs.Names {
		if avs.Values == nil || len(avs.Values) <= i {
			thisType = *typ
		} else {
			thisType, err = c.CollectExpr(ctx, avs.Values[i])
		}
		if err != nil {
			goto end
		}
		vs[i] = ValueSpec{
			File: c.File,
			Name: ident.Name,
			Type: thisType,
		}
	}
end:
	return vs, err
}
