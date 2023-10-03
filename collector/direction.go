package collector

import (
	"go/ast"
)

type Direction ast.ChanDir

func (cd Direction) String() (s string) {
	switch ast.ChanDir(cd) {
	case ast.SEND: // ast.SEND==1
		s = "->"
	case ast.RECV: // ast.RECV==2
		s = "<-"
	}
	return s
}
