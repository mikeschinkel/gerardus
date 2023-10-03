package collector

import (
	"go/token"
)

type Operator token.Token

func (op Operator) String() (s string) {
	switch token.Token(op) {
	case token.ADD:
		s = "+"
	case token.ADD_ASSIGN:
		s = "+="
	case token.AND:
		s = "&"
	case token.AND_ASSIGN:
		s = "&="
	case token.AND_NOT:
		s = "&^"
	case token.AND_NOT_ASSIGN:
		s = "&^="
	case token.ARROW:
		s = "<-"
	case token.ASSIGN:
		s = "="
	case token.COLON:
		s = ":"
	case token.COMMA:
		s = ","
	case token.DEC:
		s = "--"
	case token.DEFINE:
		s = ":="
	case token.ELLIPSIS:
		s = "..."
	case token.EQL:
		s = "=="
	case token.GEQ:
		s = ">="
	case token.GTR:
		s = ">"
	case token.INC:
		s = "++"
	case token.LAND:
		s = "&&"
	case token.LBRACE:
		s = "{"
	case token.LBRACK:
		s = "["
	case token.LEQ:
		s = "<="
	case token.LOR:
		s = "||"
	case token.LPAREN:
		s = "("
	case token.LSS:
		s = "<"
	case token.MUL:
		s = "*"
	case token.MUL_ASSIGN:
		s = "*="
	case token.NEQ:
		s = "!="
	case token.NOT:
		s = "!"
	case token.OR:
		s = "|"
	case token.OR_ASSIGN:
		s = "|="
	case token.PERIOD:
		s = "."
	case token.QUO:
		s = "/"
	case token.QUO_ASSIGN:
		s = "/="
	case token.RBRACE:
		s = "}"
	case token.RBRACK:
		s = "]"
	case token.REM:
		s = "%"
	case token.REM_ASSIGN:
		s = "%="
	case token.RPAREN:
		s = ")"
	case token.SEMICOLON:
		s = ";"
	case token.SHL:
		s = "<<"
	case token.SHL_ASSIGN:
		s = "<<="
	case token.SHR:
		s = ">>"
	case token.SHR_ASSIGN:
		s = ">>="
	case token.SUB:
		s = "-"
	case token.SUB_ASSIGN:
		s = "-="
	case token.XOR:
		s = "^"
	case token.XOR_ASSIGN:
		s = "^="
	case token.TILDE:
		s = "~"

	default:
		panicf("Unhandled OP token %T. See: ./go/token/token.go.",
			token.Token(op),
		)
	}
	return s
}
