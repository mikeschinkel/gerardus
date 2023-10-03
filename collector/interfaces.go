package collector

import (
	"fmt"
	"go/ast"
)

type Expr interface {
	fmt.Stringer
}

type ASTGetter interface {
	AST() *ast.File
}

type File interface {
	ASTGetter
	RelPath() string
}

var _ CodeFacet = (*FuncDecl)(nil)
var _ CodeFacet = (*ImportSpec)(nil)
var _ CodeFacet = (*TypeSpec)(nil)
var _ CodeFacet = (*ValueSpec)(nil)

type CodeFacet interface {
	CodeFacet()
	fmt.Stringer
}
