package collector

import (
	"context"
	"go/ast"
	"go/token"
)

func (c *Collector) CollectGenDecl(ctx context.Context, d *ast.GenDecl) (err error) {
	switch d.Tok {
	case token.IMPORT:
		for _, spec := range d.Specs {
			err = c.CollectSpec(ctx, spec, nil)
			if err != nil {
				goto end
			}
		}

	case token.TYPE:
		for _, spec := range d.Specs {
			err = c.CollectSpec(ctx, spec, nil)
			if err != nil {
				goto end
			}
		}

	case token.CONST:
		err = c.CollectConst(ctx, d)

	case token.VAR:
		err = c.CollectVar(ctx, d)

	default:
		panicf("Unhandled token type '%s'", d.Tok)
	}
end:
	return err
}

func (c *Collector) CollectVar(ctx context.Context, d *ast.GenDecl) (err error) {
	var typ Expr
	for _, spec := range d.Specs {
		// `typ` MUST be passed as a pointer so that — in some cases — it can be updated
		// first time through and reused for subsequent iterations.
		err = c.CollectSpec(ctx, spec, &typ)
		if err != nil {
			goto end
		}
	}
end:
	return err
}

func (c *Collector) CollectConst(ctx context.Context, d *ast.GenDecl) (err error) {
	return c.CollectVar(ctx, d)
}

func (c *Collector) CollectSpec(ctx context.Context, spec ast.Spec, typ *Expr) (err error) {
	switch t := spec.(type) {
	case *ast.ImportSpec:
		c.FacetChan <- c.CollectImportSpec(ctx, t)

	case *ast.TypeSpec:
		ts, err := c.CollectTypeSpec(ctx, t)
		if err != nil {
			goto end
		}
		if ts != nil {
			c.FacetChan <- *ts
		}

	case *ast.ValueSpec:
		vss, err := c.CollectValueSpec(ctx, t, typ)
		if err != nil {
			goto end
		}
		for _, vs := range vss {
			c.FacetChan <- vs
		}

	default:
		panicf("Unhandled AST type %T", spec)
	}

end:
	return err
}
