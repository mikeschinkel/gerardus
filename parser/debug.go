//go:build debug

package parser

import (
	"gerardus/scanner"
)

func init() {
	var s string
	f := &GoFile{File: scanner.NewFile("", &s)}
	f.DebugString()
}
func (gf *GoFile) DebugString() string {
	return gf.Path
}
