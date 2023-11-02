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

func (is ImportSpec) CodeFacet() {}

func (is ImportSpec) String() (s string) {
	if is.Alias == "" {
		return fmt.Sprintf(`import "%s"`, is.Package)
	}
	return fmt.Sprintf(`import %s "%s"`, is.Alias, is.Package)
}

func (is ImportSpec) Name() string {
	pos := strings.LastIndexByte(is.Package, '/')
	if pos == -1 {
		return is.Package
	}
	return is.Package[pos+1:]
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
