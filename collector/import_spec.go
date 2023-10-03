package collector

import (
	"context"
	"fmt"
	"go/ast"
	"strings"
)

type ImportSpec struct {
	File    File
	Package string
	Alias   string
}

func (spec ImportSpec) CodeFacet() {}

func (spec ImportSpec) String() (s string) {
	if spec.Alias == "" {
		return fmt.Sprintf(`import "%s"`, spec.Package)
	}
	return fmt.Sprintf(`import %s "%s"`, spec.Alias, spec.Package)
}

func (spec ImportSpec) Name() string {
	pos := strings.LastIndexByte(spec.Package, '/')
	if pos == -1 {
		return spec.Package
	}
	return spec.Package[pos+1:]
}

//goland:noinspection GoUnusedParameter
func (c *Collector) CollectImportSpec(ctx context.Context, ais *ast.ImportSpec) (is ImportSpec) {
	var name string
	if ais.Name != nil {
		name = ais.Name.Name
	}
	return ImportSpec{
		File:    c.File,
		Package: strings.Trim(ais.Path.Value, `"`),
		Alias:   name,
	}
}
