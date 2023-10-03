package collector

type SymbolType int

func (st SymbolType) ID() int {
	return int(st)
}
func (st SymbolType) Name() string {
	return SymbolName(st)
}

const (
	UnclassifiedSymbol SymbolType = 0
	StructSymbol       SymbolType = 1
	InterfaceSymbol    SymbolType = 2
	IdentifierSymbol   SymbolType = 3
	FuncSymbol         SymbolType = 4
	ExprSymbol         SymbolType = 5
	LiteralSymbol      SymbolType = 6
)

//var _ persister.symbolType = (*SymbolType)(nil)

var SymbolTypes = []SymbolType{
	UnclassifiedSymbol,
	IdentifierSymbol,
	InterfaceSymbol,
	FuncSymbol,
	StructSymbol,
	ExprSymbol,
	LiteralSymbol,
}

var SymbolTypeMap = map[string]SymbolType{}

func SymbolName(symType SymbolType) string {
	switch symType {
	case StructSymbol:
		return "struct"
	case IdentifierSymbol:
		return "identifier"
	case InterfaceSymbol:
		return "interface"
	case FuncSymbol:
		return "func"
	case ExprSymbol:
		return "expr"
	}
	return "unclassified"
}

func init() {
	for _, st := range SymbolTypes {
		SymbolTypeMap[SymbolName(st)] = st
	}
}

func SymbolTypeFromExpr(expr Expr) (st SymbolType) {
	switch expr.(type) {
	case StructType:
		st = StructSymbol
	case Ident:
		st = IdentifierSymbol
	case InterfaceType:
		st = InterfaceSymbol
	case FuncType:
		st = FuncSymbol
	case CompositeLit,
		BasicLit,
		FuncLit:
		st = LiteralSymbol
	case SelectorExpr,
		ArrayType,
		MapType,
		StarExpr,
		KeyValueExpr,
		CallExpr,
		ParenExpr,
		BinaryExpr,
		UnaryExpr,
		Ellipsis,
		ChanType,
		TypeAssertExpr,
		IndexExpr,
		IndexListExpr,
		SliceExpr:
		st = ExprSymbol
	case nil:
		panic("Unexpected: AST expr is nil.")
	default:
		st = UnclassifiedSymbol
	}
	return st
}
