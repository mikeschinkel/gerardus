package parser

type PackageMap map[string]*Package

func (pm PackageMap) HasPackage(name string) bool {
	_, ok := pm[name]
	return ok
}
