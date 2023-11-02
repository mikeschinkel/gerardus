package parser

import (
	"path/filepath"
	"sync"
)

type ModulePackageMap map[string]*ModulePackage

type ModulePackage struct {
	Module     *Module
	PackageMap PackageMap
}

type ModuleGraph struct {
	ModuleMap
	pathModuleMap    StringModuleMap
	pathGoModMap     StringModuleMap
	pkgGoModMap      StringModuleMap
	pkgModuleMap     StringModuleMap
	modulePackageMap ModulePackageMap
}

func NewModuleGraph() *ModuleGraph {
	return &ModuleGraph{
		ModuleMap:        make(ModuleMap),
		pathModuleMap:    make(StringModuleMap),
		pathGoModMap:     make(StringModuleMap),
		pkgGoModMap:      make(StringModuleMap),
		pkgModuleMap:     make(StringModuleMap),
		modulePackageMap: make(ModulePackageMap),
	}
}

func (mg *ModuleGraph) AddProjectModule(args *ModuleArgs) *Module {
	m := mg.addModule(args)
	mg.ModuleMap[m.Name()].GoMod = true
	return m
}

func (mg *ModuleGraph) AddDependentModule(pm *Module, args *ModuleArgs) {
	args.Parent = pm
	mg.addModule(args)
}

var mutex sync.Mutex

func (mg *ModuleGraph) addModule(args *ModuleArgs) (m *Module) {
	if args.ModuleGraph == nil {
		args.ModuleGraph = mg
	}
	m = newModule(args)
	mutex.Lock()
	mv, ok := mg.ModuleMap[m.Name()]
	if !ok {
		mv = NewModuleVersions()
		mv.Module = m
		mv.ModPath = m.GoModPath()
		mg.ModuleMap[m.Name()] = mv
	}
	modDir := m.GoModDir()
	_, ok = mv.VersionMap[modDir]
	if !ok {
		mv.VersionMap[modDir] = m
	}
	mutex.Unlock()
	return m
}

func (mg *ModuleGraph) DispenseGoModByImportPath(importPath string) (m *Module) {
	var ok bool
	var mv *ModuleVersions
	var tryIP string

	m, ok = mg.pkgGoModMap[importPath]
	if ok {
		goto end
	}
	tryIP = importPath
	for {
		mv, ok = mg.ModuleMap[tryIP]
		if ok {
			m = mv.Module
			if m.Package.Type != GoModPackage {
				m = nil
			}
			goto end
		}
		if tryIP == "." {
			goto end
		}
		tryIP = filepath.Dir(tryIP)
	}
end:
	if _, ok = mg.pkgGoModMap[importPath]; !ok {
		mg.pkgGoModMap[importPath] = m
	}
	return m
}

func (mg *ModuleGraph) DispenseModuleByImportPath(importPath string) (m *Module) {
	var ok bool
	var tryIP string

	m, ok = mg.pkgModuleMap[importPath]
	if ok {
		goto end
	}
	tryIP = importPath
	for {
		mv, ok := mg.ModuleMap[tryIP]
		if ok {
			m = mv.Module
			goto end
		}
		if tryIP == "." {
			goto end
		}
		tryIP = filepath.Dir(tryIP)
	}
end:
	if _, ok = mg.pkgModuleMap[importPath]; !ok {
		mg.pkgModuleMap[importPath] = m
	}
	return m
}

// DispenseGoModByPath returns a module for a go.mod based upon a path that has
// the go.mod's path as a string prefix, e.g.:
//
//			For /foo/go.mod then /foo/bar/baz will return module for /foo/go.mod,
//	   Except if a /foo/bar/go.mod then it will return module for /foo/bar/go.mod,
//	   OTOH will return `nil` for path="github.com/example/project"
func (mg *ModuleGraph) DispenseGoModByPath(path string) (m *Module) {
	var ok bool

	m, ok = mg.pathGoModMap[path]
	if ok {
		goto end
	}

	for {
		var mv *ModuleVersions
		mv, ok = mg.ModuleMap[path]
		if ok {
			m = mv.Module
			goto end
		}
		if !mv.GoMod {
			continue
		}
		if path == "." {
			goto end
		}
		if path == string(filepath.Separator) {
			goto end
		}
		path = filepath.Dir(path)
	}
end:
	if _, ok = mg.pathGoModMap[path]; !ok {
		mg.pathGoModMap[path] = m
	}
	return m
}

// DispenseModule returns a *Module given a module name and path to source file where imported
func (mg *ModuleGraph) DispenseModule(name, path string) (m *Module) {
	var pmv *ModuleVersions
	var gm *Module
	var ok bool

	pmv, ok = mg.ModuleMap[name]
	if !ok {
		goto end
	}

	gm = mg.DispenseGoModByPath(path)
	if gm == nil {
		// ???
		goto end
	}

	// Get the version of the module from the correct go.map file
	m, ok = pmv.VersionMap[gm.GoModPath()]
	if ok {
		goto end
	}

end:
	return m
}

func (mg *ModuleGraph) DispenseLocalPackage(importPath, source string) (pkg *Package) {
	var gm *Module
	var mp *ModulePackage
	var pm PackageMap
	var pt PackageType
	var ok bool

	// Look for a module (go.mod) with same importPath,
	// or with a subdirectory of the import path.
	gm = mg.DispenseGoModByImportPath(importPath)
	if gm == nil {
		goto end
	}

	mp, ok = mg.modulePackageMap[gm.ImportPath]
	if !ok {
		mg.modulePackageMap[gm.ImportPath] = &ModulePackage{
			Module:     gm,
			PackageMap: make(PackageMap),
		}
	}
	pm = mp.PackageMap
	pkg, ok = pm[importPath]
	if ok {
		goto end
	}

	pt = GoModPackage
	if gm.ImportPath != importPath {
		// If they are different, it must be a subdirectory
		pt = LocalPackage
	}

	// Now create a package based on the above.
	pkg = newPackage(&PackageArgs{
		ImportPath: importPath,
		Type:       pt,
		Version:    ".",
		Module:     gm,
	})
	// pkg.ModuleGraph = ???
	pm[importPath] = pkg

end:
	return pkg
}

// DispensePackage returns a *Module given a package name (w/o alias) and source file where imported
func (mg *ModuleGraph) DispensePackage(importPath, source string) (pkg *Package) {
	var m *Module

	m = mg.DispenseModule(importPath, source)
	if m != nil {
		pkg = m.Package
		goto end
	}

	// Look for a module (go.mod) with same importPath,
	// or with a subdirectory of the import path.
	pkg = mg.DispenseLocalPackage(importPath, source)
	if pkg != nil {
		goto end
	}

	// Look for any module with same importPath, or
	// with a subdirectory of the import path.
	m = mg.DispenseModuleByImportPath(importPath)
	if m == nil {
		// If m == nil then no module so must be a Go standard library package. The
		// return value for `pkg` will be `nil`.
		goto end
	}
	switch m.Package.Type {
	case LocalPackage:
		print()
	case ExternalPackage:
		print()
	default:
		panicf("Unexpected package type: %s", m.Package.Type.Name())
	}
	print()
end:
	return pkg
}

// DispenseModuleByPath will return the applicable module for the given path.
func (mg *ModuleGraph) DispenseModuleByPath(path string) (m *Module) {
	var ok bool

	m, ok = mg.pathModuleMap[path]
	if ok {
		goto end
	}
	if path == "." {
		goto end
	}
	if path == string(filepath.Separator) {
		goto end
	}
	for _, mv := range mg.ModuleMap {
		m, ok = mv.VersionMap[path]
		if ok {
			goto end
		}
	}
	m = mg.DispenseModuleByPath(filepath.Dir(path))
end:
	if _, ok = mg.pathModuleMap[path]; !ok {
		mg.pathModuleMap[path] = m
	}
	return m
}
