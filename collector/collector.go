package collector

import (
	"context"
	"go/ast"
)

type Collector struct {
	File      File
	FacetChan chan CodeFacet // replace FuncDecl with the actual type you intend to use
}

func New(file File, facetChan chan CodeFacet) *Collector {
	return &Collector{
		File:      file,
		FacetChan: facetChan,
	}
}

func (c *Collector) CollectFiles(ctx context.Context) (err error) {
	for _, decl := range c.File.AST().Decls {
		switch dt := decl.(type) {
		case *ast.FuncDecl:
			err = c.CollectFuncDecl(ctx, dt)
		case *ast.GenDecl:
			err = c.CollectGenDecl(ctx, dt)
		default:
			panicf("Unhandled AST type %T", dt)
		}
		if err != nil {
			goto end
		}
	}
end:
	return err
}
