package parser

import (
	"fmt"
	"go/ast"
	"io"
	"net/http"
	"os"
)

var StderrWriter io.Writer = os.Stderr

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
		_, _ = fmt.Fprintln(StderrWriter, err.Error())
	}
}

func CheckURL(url string) (status int, err error) {
	var resp *http.Response

	resp, err = http.Get(url)
	if err != nil {
		goto end
	}
	defer Close(resp.Body, WarnOnError) // Make sure to close the response body.

	status = resp.StatusCode
	if status < 200 || status > 299 {
		err = fmt.Errorf("received HTTP status code %d from %s", status, url)
		goto end
	}
end:
	return status, err
}
