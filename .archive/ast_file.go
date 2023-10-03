package _archive

//type ASTFile struct {
//	*ast.File
//	pkg *parser.Package
//}
//
//func NewASTFile(file *ast.File, pkg *parser.Package) *ASTFile {
//	return &ASTFile{
//		File: file,
//		pkg:  pkg,
//	}
//}
//
//type ASTArgs struct {
//	Prefix   string
//	TypeName string
//}
//
//func (f *ASTFile) hasFields(list *ast.FieldList) (has bool) {
//	if list == nil {
//		goto end
//	}
//	if list.List == nil {
//		goto end
//	}
//	if len(list.List) == 0 {
//		goto end
//	}
//	has = true
//end:
//	return has
//}
//
//func (f *ASTFile) String() (s string) {
//	sb := strings.Builder{}
//	sb.WriteString(f.PackageString())
//	for _, decl := range f.Decls {
//		switch dt := decl.(type) {
//		case *ast.FuncDecl:
//			sb.WriteString(f.FuncDeclString(dt))
//		case *ast.GenDecl:
//			sb.WriteString(f.GenDeclString(dt))
//		default:
//			panicf("Unhandled AST type %T", dt)
//		}
//		sb.WriteByte('\n')
//	}
//	return sb.String()
//}
//
//func (f *ASTFile) PackageString() (s string) {
//	return fmt.Sprintf("package %s\n", f.pkg.LocalName())
//}
//
//func (f *ASTFile) FuncDeclString(fd *ast.FuncDecl) (s string) {
//	s = "func "
//	if f.hasFields(fd.Recv) {
//		if len(fd.Recv.List) > 1 {
//			panicf("Unexpected: func has more than one receiver: %#v", fd)
//		}
//		s += fmt.Sprintf("(%s) ", f.FieldString(fd.Recv.List[0]))
//	}
//	s += fmt.Sprintf("%s%s",
//		fd.Name.Name,
//		f.FuncTypeString(fd.Type),
//	)
//	s += f.BlockStmtString(fd.Body)
//	return s
//}
//
////goland:noinspection GoUnusedParameter
//func (f *ASTFile) BlockStmtString(b *ast.BlockStmt) (s string) {
//	// Value Maybe expand on block statements
//	return " {...}"
//}
//
//func (f *ASTFile) GenDeclString(d *ast.GenDecl) (s string) {
//	switch d.Tok {
//	case token.IMPORT:
//		args := &ASTArgs{}
//		s = "import ("
//		for _, spec := range d.Specs {
//			s += fmt.Sprintf("%s; ", f.SpecString(spec, args))
//		}
//		s = fmt.Sprintf("%s)", s[:len(s)-2])
//
//	case token.TYPE:
//		args := &ASTArgs{Prefix: "type"}
//		for _, spec := range d.Specs {
//			s += f.SpecString(spec, args)
//		}
//
//	case token.CONST:
//		s += f.ConstString(d)
//
//	case token.VAR:
//		s += f.VarString(d)
//
//	default:
//		panicf("Unhandled token type '%s'", d.Tok.String())
//	}
//	return s
//}
//
//func (f *ASTFile) VarString(d *ast.GenDecl) (s string) {
//	return f.VarStringWithPrefix(d, "var")
//}
//
//func (f *ASTFile) ConstString(d *ast.GenDecl) (s string) {
//	return f.VarStringWithPrefix(d, "const")
//}
//
//func (f *ASTFile) VarStringWithPrefix(d *ast.GenDecl, prefix string) (s string) {
//	switch len(d.Specs) {
//	case 0:
//	case 1:
//		s = fmt.Sprintf("%s ", prefix)
//	default:
//		s = fmt.Sprintf("%s (", prefix)
//	}
//	args := ASTArgs{
//		Prefix: prefix,
//		// TypeName will get set first time through the loop.
//	}
//	for _, spec := range d.Specs {
//		// args MUST be passed as a pointer.
//		c := f.SpecString(spec, &args)
//		s += fmt.Sprintf("%s; ", c)
//	}
//	if len(s) >= 2 {
//		s = s[:len(s)-2]
//	}
//	if len(d.Specs) > 1 {
//		s += ")"
//	}
//	return s
//}
//
//func (f *ASTFile) SpecString(spec ast.Spec, args *ASTArgs) (s string) {
//	switch t := spec.(type) {
//	case *ast.TypeSpec:
//		s = f.TypeSpecString(t, args)
//	case *ast.ImportSpec:
//		s = f.ImportSpecString(t)
//	case *ast.ValueSpec:
//		s = f.ValueSpecString(t, args)
//	default:
//		panicf("Unhandled AST type %T", spec)
//	}
//	return s
//}
//
//func (f *ASTFile) ImportSpecString(is *ast.ImportSpec) (s string) {
//	if is.Name != nil {
//		s = fmt.Sprintf("%s %s", is.Name.Name, is.Path.Value)
//		goto end
//	}
//	s = is.Path.Value
//end:
//	return s
//}
//
//func (f *ASTFile) ValueSpecString(vs *ast.ValueSpec, args *ASTArgs) (s string) {
//	var value string
//	if vs.Type != nil {
//		// This if for when iota is used to define constants
//		// This sets args the first time through after which it continues using it
//		args.TypeName = f.ExprString(vs.Type)
//	}
//	for i, ident := range vs.Names {
//		if vs.Values != nil && len(vs.Values) > i {
//			value = f.ExprString(vs.Values[i])
//		}
//		if value != "" {
//			s += fmt.Sprintf("%s = %s; ", ident.Name, value)
//		} else {
//			s += fmt.Sprintf("%s %s; ", ident.Name, args.TypeName)
//		}
//	}
//	if len(s) >= 2 {
//		s = s[:len(s)-2]
//	}
//	return s
//}
//
//func (f *ASTFile) TypeSpecString(ts *ast.TypeSpec, args *ASTArgs) (s string) {
//	if ts.Name == nil {
//		goto end
//	}
//	s = fmt.Sprintf("%s %s %s",
//		args.Prefix,
//		ts.Name.Name,
//		f.ExprString(ts.Type),
//	)
//end:
//	return s
//}
//
//func (f *ASTFile) ExprString(expr ast.Expr) (s string) {
//	switch t := expr.(type) {
//	case *ast.StructType:
//		s = "struct{"
//		for _, fld := range t.Fields.List {
//			s += fmt.Sprintf("%s; ", f.FieldString(fld))
//		}
//		if len(s) >= 2 {
//			s = s[:len(s)-2]
//		}
//		s += "}"
//	case *ast.Ident:
//		s += t.Name
//	case *ast.SelectorExpr:
//		s += fmt.Sprintf("%s.%s", f.ExprString(t.X), t.Sel.Name)
//	case *ast.ArrayType:
//		if t.Len != nil {
//			s += fmt.Sprintf("[%s]%s",
//				f.ExprString(t.Len),
//				f.ExprString(t.Elt),
//			)
//		} else {
//			s += fmt.Sprintf("[]%s", f.ExprString(t.Elt))
//		}
//	case *ast.MapType:
//		s += f.MapTypeString(t)
//	case *ast.StarExpr:
//		s += fmt.Sprintf("*%s", f.ExprString(t.X))
//	case *ast.FuncType:
//		s += f.FuncTypeString(t)
//	case *ast.CompositeLit:
//		s += f.CompositeLitString(t)
//	case *ast.BasicLit:
//		s += t.Value
//	case *ast.KeyValueExpr:
//		s += fmt.Sprintf("%s:%s",
//			f.ExprString(t.Key),
//			f.ExprString(t.Value),
//		)
//	case *ast.InterfaceType:
//		s += f.InterfaceTypeString(t)
//	case *ast.CallExpr:
//		s += f.CallExprString(t)
//	case *ast.ParenExpr:
//		s += f.ParenExprString(t)
//	case *ast.BinaryExpr:
//		s += f.BinaryExprString(t)
//	case *ast.UnaryExpr:
//		s += f.UnaryExprString(t)
//	case *ast.Ellipsis:
//		s += f.EllipsisString(t)
//	case *ast.ChanType:
//		s += f.ChanTypeString(t)
//	case *ast.FuncLit:
//		s += f.FuncLitString(t)
//	case *ast.TypeAssertExpr:
//		s += f.TypeAssertExprString(t)
//	case *ast.IndexExpr:
//		s += f.IndexExprString(t)
//	case *ast.IndexListExpr:
//		s += f.IndexListExprString(t)
//	case *ast.SliceExpr:
//		s += f.SliceExprString(t)
//	case nil:
//		panic("Unexpected: AST expr is nil.")
//	default:
//		panicf("Unhandled AST expr type %T", expr)
//	}
//	return s
//}
//
//func (f *ASTFile) SliceExprString(e *ast.SliceExpr) (s string) {
//	name := f.ExprString(e.X)
//	switch {
//	case e.High != nil && e.Low != nil:
//		s = fmt.Sprintf("%s[%s:%s]", name,
//			f.ExprString(e.Low),
//			f.ExprString(e.High),
//		)
//	case e.Low != nil:
//		s = fmt.Sprintf("%s[%s:]", name, f.ExprString(e.Low))
//	case e.High != nil:
//		s = fmt.Sprintf("%s[:%s]", name, f.ExprString(e.High))
//	}
//
//	return s
//}
//func (f *ASTFile) IndexListExprString(e *ast.IndexListExpr) (s string) {
//	return fmt.Sprintf("%s[%s]",
//		f.ExprString(e.X),
//		f.ExprSliceString(e.Indices),
//	)
//}
//
//func (f *ASTFile) IndexExprString(e *ast.IndexExpr) (s string) {
//	return fmt.Sprintf("%s[%s]",
//		f.ExprString(e.X),
//		f.ExprString(e.Index),
//	)
//}
//
//func (f *ASTFile) TypeAssertExprString(e *ast.TypeAssertExpr) (s string) {
//	return fmt.Sprintf("%s.(%s)",
//		f.ExprString(e.X),
//		f.ExprString(e.Type),
//	)
//}
//
//func (f *ASTFile) FuncLitString(fl *ast.FuncLit) (s string) {
//	return fmt.Sprintf("func%s%s",
//		f.FuncTypeString(fl.Type),
//		f.BlockStmtString(fl.Body),
//	)
//}
//
//func (f *ASTFile) CompositeLitString(cl *ast.CompositeLit) (s string) {
//	s = fmt.Sprintf("{%s}", f.ExprSliceString(cl.Elts))
//	if cl.Type == nil {
//		goto end
//	}
//	s = fmt.Sprintf("%s%s", f.ExprString(cl.Type), s)
//end:
//	return s
//}
//
//func (f *ASTFile) ChanTypeString(e *ast.ChanType) (s string) {
//	return fmt.Sprintf("chan%s %s",
//		chanDirString(e.Dir),
//		f.ExprString(e.Value),
//	)
//}
//func (f *ASTFile) EllipsisString(e *ast.Ellipsis) (s string) {
//	if e.Elt == nil {
//		return "..."
//	}
//	return fmt.Sprintf("...%s", f.ExprString(e.Elt))
//}
//func (f *ASTFile) ParenExprString(e *ast.ParenExpr) (s string) {
//	return fmt.Sprintf("(%s)", f.ExprString(e.X))
//}
//
//func (f *ASTFile) UnaryExprString(e *ast.UnaryExpr) (s string) {
//	return fmt.Sprintf("%s%s",
//		opString(e.Op),
//		f.ExprString(e.X),
//	)
//}
//func (f *ASTFile) BinaryExprString(e *ast.BinaryExpr) (s string) {
//	return fmt.Sprintf("%s%s%s",
//		f.ExprString(e.X),
//		opString(e.Op),
//		f.ExprString(e.Y),
//	)
//}
//
//func (f *ASTFile) MapTypeString(mt *ast.MapType) (s string) {
//	var valueType string
//	switch t := mt.Value.(type) {
//	case *ast.FuncType:
//		valueType = fmt.Sprintf("func%s", f.FuncTypeString(t))
//	default:
//		valueType = f.ExprString(mt.Value)
//	}
//	return fmt.Sprintf("map[%s]%s",
//		f.ExprString(mt.Key),
//		valueType,
//	)
//}
//
//func (f *ASTFile) InterfaceTypeString(iface *ast.InterfaceType) (s string) {
//	return fmt.Sprintf("interface{%s}",
//		f.FieldListString(iface.Methods),
//	)
//}
//
//func (f *ASTFile) CallExprString(call *ast.CallExpr) (s string) {
//	return fmt.Sprintf("%s(%s)",
//		f.ExprString(call.Fun),
//		f.ExprSliceString(call.Args),
//	)
//}
//
//func (f *ASTFile) ExprSliceString(exprs []ast.Expr) (s string) {
//	for _, expr := range exprs {
//		s += fmt.Sprintf("%s,", f.ExprString(expr))
//	}
//	if len(s) >= 1 {
//		s = s[:len(s)-1]
//	}
//	return s
//}
//
//func (f *ASTFile) FuncTypeString(ft *ast.FuncType) (s string) {
//
//	// Params
//	switch {
//	case !f.hasFields(ft.Params):
//		s += "()"
//	default:
//		s += "(" + f.FieldListString(ft.Params) + ")"
//	}
//
//	// Results
//	switch {
//	case !f.hasFields(ft.Results):
//		// Do nothing as nothing is needed when no results
//	case len(ft.Results.List) == 1:
//		s += f.FieldListString(ft.Results)
//	default:
//		s += "(" + f.FieldListString(ft.Results) + ")"
//	}
//	return s
//}
//
//func (f *ASTFile) FieldListString(list *ast.FieldList) (s string) {
//	return f.FieldSliceString(list.List)
//}
//
//func (f *ASTFile) FieldSliceString(list []*ast.Field) (s string) {
//	for _, fld := range list {
//		s += fmt.Sprintf("%s,", f.FieldString(fld))
//	}
//	if len(s) >= 1 {
//		s = s[:len(s)-1]
//	}
//	return s
//}
//
//func (f *ASTFile) FieldString(fld *ast.Field) (s string) {
//	typeName := f.ExprString(fld.Type)
//	if len(fld.Names) == 0 {
//		// If there is no var name, often the case for results
//		s = typeName
//		goto end
//	}
//	for _, ident := range fld.Names {
//		if len(typeName) > 0 && typeName[0] == '(' {
//			s += fmt.Sprintf("%s%s,", ident.Name, typeName)
//		} else {
//			s += fmt.Sprintf("%s %s,", ident.Name, typeName)
//		}
//	}
//	if len(s) >= 1 {
//		s = s[:len(s)-1]
//	}
//end:
//	return s
//}
