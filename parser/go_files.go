package parser

import (
	"bufio"
	"go/ast"
	"os"
	"regexp"

	"gerardus/paths"
)

type GoFile struct {
	File
	Package   *Package
	Imports   ImportsMap
	Functions Functions
	//Types      _archive.Types
	References map[string]struct{}
	ast        *ast.File
	Path       string // Only needed for DebugString
}

func NewGoFile(file File, pkgName string) *GoFile {
	return &GoFile{
		File:       file,
		Path:       file.RelPath(),
		Package:    NewPackage(pkgName),
		References: make(map[string]struct{}),
	}
}

func (gf *GoFile) AST() *ast.File {
	return gf.ast
}

func (gf *GoFile) String() string {
	return gf.File.RelPath()
}

func (gf *GoFile) AddReference(ref string) {
	gf.References[ref] = struct{}{}
}

//func (gf *GoFile) ReduceImports() {
//	if len(gf.Types) == 0 {
//		gf.Imports = make(gerardus.ImportsMap)
//		return
//	}
//	for i, imp := range gf.Imports {
//		_, ok := gf.References[imp.Package.LocalName()]
//		if ok {
//			continue
//		}
//		delete(gf.Imports, i)
//	}
//}

func (gf *GoFile) FullPackageName() string {
	return gf.Package.Name
}

func (gf *GoFile) LocalPackageName() string {
	return gf.Package.LocalName()
}

func (gf *GoFile) HasImports() bool {
	return len(gf.Imports) > 0
}

//func (gf *GoFile) HasTypes() bool {
//	return len(gf.Types) > 0
//}

func (gf *GoFile) HasFunctions() bool {
	return len(gf.Functions) > 0
}

func AsGoFile(file File) *GoFile {
	goFile, ok := file.(*GoFile)
	if !ok {
		panicf("Attempting to generate package on '%T'; %v", file)
	}
	return goFile
}

var matchPackage = regexp.MustCompile(`^\s*package\s+(.+)\s*$`)

func loadPackageName(file File) (name string, err error) {
	var exists bool
	var scanner *bufio.Scanner
	var fh *os.File
	var match []string

	exists, err = paths.FileExists(file.Fullpath())
	if err != nil {
		panicf("Checking file existence caused an error: '%f'; %s",
			file.Fullpath(),
			err.Error(),
		)
	}
	if !exists {
		goto end
	}
	fh, err = os.Open(file.Fullpath())
	if err != nil {
		panicf("Cannot read file: '%f'; %s", file.Fullpath(), err.Error())
	}
	defer Close(fh, WarnOnError)

	scanner = bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		// Use regex to match the line
		match = matchPackage.FindStringSubmatch(line)
		if match != nil {
			name = match[1]
			break
		}
	}
	err = scanner.Err()
end:
	return name, err
}

//func receiverName(fd *ast.FuncDecl) (s string) {
//	var fld *ast.Field
//	var ident *ast.Ident
//	fld = receiver(fd)
//	if fld == nil {
//		goto end
//	}
//	if fld.Names == nil {
//		goto end
//	}
//	if len(fld.Names) == 0 {
//		goto end
//	}
//	ident = fld.Names[0]
//	if ident == nil {
//		goto end
//	}
//	s = ident.Value
//end:
//	if s == "" {
//		s = "<unknown>"
//	}
//	return s
//}
//
//func receiverType(fd *ast.FuncDecl) (s string) {
//	var fld *ast.Field
//	fld = receiver(fd)
//	if fld == nil {
//		goto end
//	}
//	if fld.Type == nil {
//		goto end
//	}
//	s = types.ExprString(fld.Type)
//end:
//	if s == "" {
//		s = "<unknown>"
//	}
//	return s
//}
//func receiver(fd *ast.FuncDecl) (r *ast.Field) {
//	if fd.Recv == nil {
//		goto end
//	}
//	if fd.Recv.List == nil {
//		goto end
//	}
//	if len(fd.Recv.List) == 0 {
//		goto end
//	}
//	r = fd.Recv.List[0]
//	if r == nil {
//		goto end
//	}
//end:
//	return r
//}
//func (gf *GoFile) collectImport(imp *ast.ImportSpec, cs *CodeSurveyor) {
//	// Remove double quotes that surround the value as it is found in AST.
//	path := strings.Trim(imp.Path.Value, "\"")
//	var alias string
//	if imp.Value != nil {
//		alias = imp.Value.Value
//	}
//	pkg := cs.collectPackage(path)
//	gf.Imports[path] = NewImport(pkg, alias)
//}
//
//// collectTypes collects all types from an *ast.File into a map
//// from type keyword (e.g., "struct", "interface") to slices of corresponding ast.Node instances.
//func (gf *GoFile) collectTypes(cs *CodeSurveyor) {
//	// Walk through all declarations in the file.
//	for _, decl := range gf.AST.Decls {
//		gf.collectDecl(decl, cs)
//	}
//}
//
//func (gf *GoFile) collectDecl(decl ast.Decl, cs *CodeSurveyor) {
//	switch dt := decl.(type) {
//	case *ast.FuncDecl:
//		gf.collectFuncDecl(dt, cs)
//	case *ast.GenDecl:
//		gf.collectGenDecl(dt, cs)
//	default:
//		panicf("Unhandled AST type '%T'", dt)
//	}
//}
//
//func (gf *GoFile) collectFuncDecl(fd *ast.FuncDecl, cs *CodeSurveyor) {
//	gf.AddFunction(NewFunction(
//		fd.Value.Value,
//		NewVariable(receiverName(fd), receiverType(fd)),
//		fd,
//	))
//}
//
//func (gf *GoFile) collectGenDecl(gd *ast.GenDecl, cs *CodeSurveyor) {
//	// Walk through all specs in the declaration.
//	for _, spec := range gd.Specs {
//		gf.collectSpec(gd.Tok, spec, cs)
//	}
//}
//
//func (gf *GoFile) collectSpec(tok token.Token, spec ast.Spec, cs *CodeSurveyor) {
//	switch t := spec.(type) {
//	case *ast.TypeSpec:
//		gf.CollectTypeSpec(t.Value.Value, t.Type)
//	case *ast.ImportSpec:
//		gf.collectImport(t, cs)
//	case *ast.ValueSpec:
//		//Value Collect these
//		// We don't care about them right now
//		print()
//	default:
//		panicf("Unhandled AST type '%T'", t)
//	}
//}
//
//func (gf *GoFile) AddType(t *Type) *Type {
//	gf.Types = append(gf.Types, t)
//	return t
//}
//
//func (gf *GoFile) AddFunction(f *Function) *Function {
//	gf.Functions = append(gf.Functions, f)
//	return f
//}
//
//func (gf *GoFile) CollectTypeSpec(name string, expr ast.Expr) (t *Type) {
//	switch et := expr.(type) {
//	case *ast.InterfaceType:
//		t = gf.AddType(NewType(InterfaceType, name, gf, et))
//	case *ast.StructType:
//		t = gf.AddType(NewType(StructType, name, gf, et))
//	case *ast.FuncType:
//		t = gf.AddType(NewType(FuncType, name, gf, et))
//	case *ast.Ident:
//		t = gf.CollectPublicIdent(name, gf, et)
//	case *ast.SelectorExpr:
//		// Handle qualified identifiers like pkg.Type
//		if xIdent, ok := et.X.(*ast.Ident); ok {
//			t = gf.CollectPublicIdent(name, gf, xIdent)
//		}
//	case *ast.StarExpr: // Pointers
//		t = gf.CollectTypeSpec(name, et.X)
//	case *ast.IndexExpr: // Pointers
//		t = gf.CollectTypeSpec(name, et.X)
//	case *ast.IndexListExpr: // Pointers
//		t = gf.CollectTypeSpec(name, et.X)
//	case *ast.ParenExpr: // Pointers
//		t = gf.CollectTypeSpec(name, et.X)
//	case *ast.ArrayType: // Arrays
//		t = gf.CollectTypeSpec(name, et.Elt)
//		if et.Len != nil {
//			t.Key = gf.CollectTypeSpec(name, et.Len)
//		}
//	case *ast.MapType:
//		// Maps
//		t = gf.CollectTypeSpec(name, et.Value)
//		if t != nil {
//			t.Key = gf.CollectTypeSpec(name, et.Key)
//		}
//	case *ast.ChanType:
//		// Channels
//		t = gf.CollectTypeSpec(name, et.Value)
//	case *ast.BasicLit:
//		t = NewType(LiteralType, name, gf, et)
//	default:
//		panicf("Unhandled AST type '%T': %#v", et, et)
//		// Move along, nothing to see here, move along
//	}
//	return t
//}
//
//func (gf *GoFile) CollectPublicIdent(name string, _ *GoFile, ident *ast.Ident) (t *Type) {
//	t = NewType(IdentifierType, name, gf, ident)
//	gf.Types = append(gf.Types, t)
//	return t
//}
