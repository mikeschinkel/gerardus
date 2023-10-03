package parser

import (
	"fmt"
	"go/ast"
	"io"

	"gerardus/options"
)

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func debugBreakpointHere(...any) {
	// just a function for debugging
}

func IsSelectorExpr(expr ast.Expr) bool {
	_, ok := expr.(*ast.SelectorExpr)
	return ok
}

func isPublicName(name string) (isPublic bool) {
	if len(name) == 0 {
		goto end
	}
	if isLower(name[0]) {
		goto end
	}
	isPublic = true
end:
	return isPublic
}

func isLower(ch byte) bool {
	return 'a' <= ch && ch <= 'z'
}

func Close(c io.Closer, f func(err error)) {
	f(c.Close())
}

func WarnOnError(err error) {
	if err != nil {
		options.StdErr(err.Error())
	}
}
