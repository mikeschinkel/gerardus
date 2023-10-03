package collector

import (
	"context"
	"fmt"
	"go/ast"
	"strings"
)

type FuncType struct {
	Parameters FieldList
	Results    FieldList
}

func (ft FuncType) String() (s string) {
	var ps, rs string
	if len(ft.Parameters) > 0 {
		ps = ft.Parameters.String()
	}
	switch len(ft.Results) {
	case 0:
	case 1:
		rs = ft.Results.String()
		if strings.Contains(rs, " ") {
			rs = fmt.Sprintf("(%s)", rs)
		}
	default:
		rs = ft.Results.String()
		rs = fmt.Sprintf("(%s)", rs)
	}
	s = fmt.Sprintf("(%s)%s", ps, rs)
	return s
}

func (c *Collector) CollectFuncType(ctx context.Context, aft *ast.FuncType) (ft FuncType, err error) {
	var p, r FieldList

	if c.astHasFields(aft.Params) {
		p, err = c.CollectFieldList(ctx, aft.Params)
	}
	if err != nil {
		goto end
	}
	if c.astHasFields(aft.Results) {
		r, err = c.CollectFieldList(ctx, aft.Results)
	}
	if err != nil {
		goto end
	}
	ft.Parameters = p
	ft.Results = r
end:
	return ft, err
}
