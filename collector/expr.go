package collector

import (
	"context"
	"go/ast"
	"strings"
)

type ExprList []Expr

func (el ExprList) String() (s string) {
	var sb strings.Builder
	if len(el) == 0 {
		goto end
	}
	sb = strings.Builder{}
	for _, e := range el {
		sb.WriteString(e.String())
		sb.WriteByte(',')
	}
	s = sb.String()
	s = s[:len(s)-1]
end:
	return s
}

func (c *Collector) CollectExprSlice(ctx context.Context, astExprs []ast.Expr) (el ExprList, err error) {
	if len(astExprs) == 0 {
		goto end
	}
	el = make(ExprList, len(astExprs))
	for i, astExpr := range astExprs {
		el[i], err = c.CollectExpr(ctx, astExpr)
		//switch t := astExpr.(type) {
		//case *ast.KeyValueExpr:
		//	c.CollectKeyValueExpr(ctx, t)
		//case *ast.Ident:
		//	el[i], err = c.CollectIdent(ctx, t)
		//case *ast.BasicLit:
		//	el[i], err = c.CollectBasicLit(ctx, t)
		//case *ast.CallExpr:
		//	el[i], err = c.CollectCallExpr(ctx, t)
		//default:
		//	panicf("Unexpected AST property type for *ast.Expr slice '%T'.", astExpr)
		//}
		if err != nil {
			goto end
		}
	}
end:
	return el, err
}

func (c *Collector) CollectExprString(ctx context.Context, ae ast.Expr) (s string, err error) {
	var expr Expr
	expr, err = c.CollectExpr(ctx, ae)
	if err != nil {
		goto end
	}
	s = expr.String()
end:
	return s, err
}

func (c *Collector) CollectExpr(ctx context.Context, expr ast.Expr) (e Expr, err error) {
	switch t := expr.(type) {
	case *ast.StructType:
		e, err = c.CollectStructType(ctx, t)
	case *ast.Ident:
		e, err = c.CollectIdent(ctx, t)
	case *ast.SelectorExpr:
		e, err = c.CollectSelectorExpr(ctx, t)
	case *ast.ArrayType:
		e, err = c.CollectArrayType(ctx, t)
	case *ast.MapType:
		e, err = c.CollectMapType(ctx, t)
	case *ast.StarExpr:
		e, err = c.CollectStarExpr(ctx, t)
	case *ast.FuncType:
		e, err = c.CollectFuncType(ctx, t)
	case *ast.CompositeLit:
		e, err = c.CollectCompositeLit(ctx, t)
	case *ast.BasicLit:
		e, err = c.CollectBasicLit(ctx, t)
	case *ast.KeyValueExpr:
		e, err = c.CollectKeyValueExpr(ctx, t)
	case *ast.InterfaceType:
		e, err = c.CollectInterfaceType(ctx, t)
	case *ast.CallExpr:
		e, err = c.CollectCallExpr(ctx, t)
	case *ast.ParenExpr:
		e, err = c.CollectParenExpr(ctx, t)
	case *ast.BinaryExpr:
		e, err = c.CollectBinaryExpr(ctx, t)
	case *ast.UnaryExpr:
		e, err = c.CollectUnaryExpr(ctx, t)
	case *ast.Ellipsis:
		e, err = c.CollectEllipsis(ctx, t)
	case *ast.ChanType:
		e, err = c.CollectChanType(ctx, t)
	case *ast.FuncLit:
		e, err = c.CollectFuncLit(ctx, t)
	case *ast.TypeAssertExpr:
		e, err = c.CollectTypeAssertExpr(ctx, t)
	case *ast.IndexExpr:
		e, err = c.CollectIndexExpr(ctx, t)
	case *ast.IndexListExpr:
		e, err = c.CollectIndexListExpr(ctx, t)
	case *ast.SliceExpr:
		e, err = c.CollectSliceExpr(ctx, t)
	case nil:
		panic("Unexpected: AST expr is nil.")
	default:
		panicf("Unhandled AST expr type %T", expr)
	}
	return e, err
}

//func ExprString(expr Expr) (s string) {
//	switch t := expr.(type) {
//	case StructType:
//		s = t.String()
//	case Ident:
//		s = t.String()
//	case SelectorExpr:
//		s = t.String()
//	case ArrayType:
//		s = t.String()
//	case CallExpr:
//		s = t.String()
//	case MapType:
//		s = t.String()
//	case StarExpr:
//		s = t.String()
//	case CompositeLit:
//		s = t.String()
//	case BasicLit:
//		s = t.String()
//	case KeyValueExpr:
//		s = t.String()
//	case InterfaceType:
//		s = t.String()
//	case ParenExpr:
//		s = t.String()
//	case BinaryExpr:
//		s = t.String()
//	case UnaryExpr:
//		s = t.String()
//	case Ellipsis:
//		s = t.String()
//	case ChanType:
//		s = t.String()
//	case FuncLit:
//		s = t.String()
//	case TypeAssertExpr:
//		s = t.String()
//	case IndexExpr:
//		s = t.String()
//	case IndexListExpr:
//		s = t.String()
//	case SliceExpr:
//		s = t.String()
//	case ExprList:
//		s = t.String()
//	case nil:
//		panic("Unexpected: Expr is nil.")
//	default:
//		panicf("Unhandled Expr type %T", expr)
//	}
//	return s
//}
