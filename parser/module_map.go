package parser

type ModuleMap map[string]*ModuleVersions

func (mm ModuleMap) ImportPaths() (ips []string) {
	ips = make([]string, 0, len(mm))
	for ip := range mm {
		ips = append(ips, ip)
	}
	return ips
}

func (mm ModuleMap) GetPackageByImportPath(importPath string) (pkg *Package) {
	var mv *ModuleVersions
	var ok bool

	if mm == nil {
		goto end
	}
	mv, ok = mm[importPath]
	if !ok {
		goto end
	}
	// TODO: Can we share rather than clone the objects?
	// Clone, except maintain pointers for Module and ModuleGraph.
	pkg = mv.Module.Package.PartialClone()
end:
	return pkg
}
