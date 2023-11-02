package parser

import (
	"bufio"
	"go/ast"
	"os"
	"regexp"

	"gerardus/paths"
	"gerardus/scanner"
)

type GoFile struct {
	scanner.File
	Package   *Package
	ImportMap ImportMap
	Functions Functions
	//Types      _archive.Types
	References  map[string]struct{}
	ast         *ast.File
	debugString string
}

func NewGoFile(file scanner.File, pkgPath string) *GoFile {
	return &GoFile{
		File:        file,
		debugString: file.RelPath(),
		Package:     NewFlyweightPackage(pkgPath),
		References:  make(map[string]struct{}),
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
//		gf.ImportMap = make(gerardus.ImportMap)
//		return
//	}
//	for i, imp := range gf.ImportMap {
//		_, ok := gf.References[imp.Package.Name()]
//		if ok {
//			continue
//		}
//		delete(gf.ImportMap, i)
//	}
//}

func (gf *GoFile) PackageImportPath() string {
	return gf.Package.ImportPath
}

func (gf *GoFile) PackageName() string {
	return gf.Package.Name()
}

//func (gf *GoFile) HasTypes() bool {
//	return len(gf.Types) > 0
//}

func (gf *GoFile) HasFunctions() bool {
	return len(gf.Functions) > 0
}

var matchPackage = regexp.MustCompile(`^\s*package\s+(.+)\s*$`)

func loadPackageName(file scanner.File) (name string, err error) {
	var exists bool
	var scnr *bufio.Scanner
	var fh *os.File
	var match []string

	exists, err = paths.FileExists(file.Fullpath())
	if err != nil {
		panicf("Checking file existence caused an error: '%s'; %s",
			file.Fullpath(),
			err.Error(),
		)
	}
	if !exists {
		goto end
	}
	fh, err = os.Open(file.Fullpath())
	if err != nil {
		panicf("Cannot read file: '%s'; %s", file.Fullpath(), err.Error())
	}
	defer Close(fh, WarnOnError)

	scnr = bufio.NewScanner(fh)
	for scnr.Scan() {
		line := scnr.Text()
		// Use regex to match the line
		match = matchPackage.FindStringSubmatch(line)
		if match != nil {
			name = match[1]
			break
		}
	}
	err = scnr.Err()
end:
	return name, err
}
