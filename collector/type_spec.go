package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type TypeSpec struct {
	File       File
	Name       string
	Definition Expr
	SymbolType SymbolType
}

func (spec TypeSpec) CodeFacet() {}

func (spec TypeSpec) String() (s string) {
	return fmt.Sprintf("type %s %s", spec.Name, spec.Definition)
}

func (c *Collector) CollectTypeSpec(ctx context.Context, ats *ast.TypeSpec) (ts *TypeSpec, err error) {
	var def Expr
	var st SymbolType
	if ats.Name != nil {
		def, err = c.CollectExpr(ctx, ats.Type)
		if err != nil {
			goto end
		}
		st = SymbolTypeFromExpr(def)
		//if st == UnclassifiedSymbol {
		//	debugBreakpointHere()
		//}
		ts = &TypeSpec{
			File:       c.File,
			Name:       ats.Name.Name,
			Definition: def,
			SymbolType: st,
		}
	}
end:
	return ts, err
}
