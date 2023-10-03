package _archive

//import (
//	"fmt"
//	"go/ast"
//	"sort"
//
//	"gerardus"
//	"gerardus/collector"
//	"gerardus/parser"
//)
//
//type ASTNodeType int
//
//const (
//	OtherTypeNode ASTNodeType = iota
//	InterfaceTypeNode
//)
//
//type TypesMap map[string]*Type
//type Types []*Type
//
//func (tt Types) Sorted() Types {
//	sort.SliceStable(tt, func(i, j int) bool {
//		return tt[i].name < tt[j].name
//	})
//	return tt
//
//}
//
//type Type struct {
//	file        *parser.GoFile
//	node        ast.Node
//	name        string
//	Methods     gerardus.MethodMap
//	Fields      gerardus.FieldMap
//	symType     collector.SymbolType
//	Key         *Type
//	subType     *Type
//	debugString string
//}
//
//func NewType(st collector.SymbolType, name string, file *parser.GoFile, node ast.Node) *Type {
//	t := &Type{
//		symType:     st,
//		name:        name,
//		file:        file,
//		node:        node,
//		Methods:     make(gerardus.MethodMap),
//		Fields:      make(gerardus.FieldMap),
//		debugString: name,
//	}
//	if t.node != nil {
//		t.collectChildren()
//	}
//	return t
//}
//
//func (t *Type) HasMethods() (has bool) {
//	if len(t.Methods) == 0 {
//		goto end
//	}
//	// Report it has no methods if any are private
//	has = true
//	for _, m := range t.Methods {
//		if m.IsPrivate() {
//			has = false
//			break
//		}
//	}
//end:
//	return has
//}
//
//func (t *Type) UniqueKey() string {
//	return fmt.Sprintf("%s.%s", t.PackageName(), t.Name())
//}
//func (t *Type) Name() string {
//	return t.name
//}
//func (t *Type) File() parser.File {
//	return t.file
//}
//
//func (t *Type) PackageName() string {
//	return t.file.FullPackageName()
//}
//
//func (t *Type) Type() ASTNodeType {
//	switch t.node.(type) {
//	case *ast.InterfaceType:
//		return InterfaceTypeNode
//	default:
//		return OtherTypeNode
//	}
//}
//
//func (t *Type) maybeCollectReference(expr ast.Expr) {
//	switch et := expr.(type) {
//	case *ast.SelectorExpr:
//		ref := et.X.(*ast.Ident).Name
//		t.file.AddReference(ref)
//	case *ast.FuncType:
//		if et.Params != nil {
//			for _, p := range et.Params.List {
//				t.maybeCollectReference(p.Type)
//			}
//		}
//		if et.Results != nil {
//			for _, r := range et.Results.List {
//				t.maybeCollectReference(r.Type)
//			}
//		}
//	case *ast.StarExpr: // Pointers
//		t.maybeCollectReference(et.X)
//	case *ast.ArrayType:
//		t.maybeCollectReference(et.Elt)
//	case *ast.Ident:
//		print()
//		// Nothing to see here, move along, nothing to see
//	default:
//		print()
//	}
//}
//
//func (t *Type) collectChildren() {
//	switch mt := t.node.(type) {
//	case *ast.StructType:
//		for _, fld := range mt.Fields.List {
//			t.collectFieldOrFields(fld)
//		}
//	case *ast.InterfaceType:
//		for _, m := range mt.Methods.List {
//			t.collectFieldOrFields(m)
//		}
//	case *ast.FuncType:
//		print()
//	case *ast.Ident:
//		t.subType = NewType(collector.IdentifierSymbol, mt.Name, t.file, nil)
//		// This is because complex expressions won't display with DebugString()
//		t.debugString += " " + t.subType.name
//		if mt.Obj != nil {
//			print()
//		}
//	default:
//		gerardus.panicf("Type not handling child '%T'; %v", t.node)
//	}
//}
//
//// collectFieldOrFields collects one or more fields from an *ast.Field value.
//// There will be more than one field if multiple fields are declared on the same
//// line with a single type, e.g.
////
////	struct{
////		x, y int
////	}
//func (t *Type) collectFieldOrFields(fld *ast.Field) {
//	for _, name := range fld.Names {
//		// These are for the struct fields
//		t.Fields[name.Name] = gerardus.NewField(name.Name, fld, t.file)
//	}
//}
//
//// collectMethods collects a method from an *ast.Field value. Expect only one.
//func (t *Type) collectMethods(m *ast.Field) {
//	if len(m.Names) > 1 {
//		gerardus.panicf("More than one method: %v!", m)
//	}
//	// These are for the interface method signatures
//	for _, name := range m.Names {
//		t.Methods[name.Name] = gerardus.NewMethod(name.Name, m, t.file)
//		t.maybeCollectReference(m.Type)
//	}
//	// These are for embedded interfaces
//	if se, ok := m.Type.(*ast.SelectorExpr); ok {
//		ident, ok := se.X.(*ast.Ident)
//		if !ok {
//			goto end
//		}
//		name := fmt.Sprintf("%s.%s", ident.Name, se.Sel.Name)
//		t.Methods[name] = gerardus.NewMethod(name, m, t.file)
//		t.file.AddReference(ident.Name)
//	}
//end:
//}
//
//type collectFunc func(*Type, string, *ast.Field)
//
//func (t *Type) collectListChildren(list *ast.FieldList, f collectFunc) {
//	var field *ast.Field
//
//	for _, field = range list.List {
//		for _, name := range field.Names {
//			f(t, name.Name, field)
//		}
//		if se, ok := field.Type.(*ast.SelectorExpr); ok {
//			ident, ok := se.X.(*ast.Ident)
//			if !ok {
//				continue
//			}
//			name := fmt.Sprintf("%s.%s", ident.Name, se.Sel.Name)
//			f(t, name, field)
//		}
//	}
//}
