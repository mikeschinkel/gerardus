package parser

import (
	"strings"
)

type PackagesMap map[string]*Package

func (pp PackagesMap) HasPackage(name string) bool {
	_, ok := pp[name]
	return ok
}

type Package struct {
	Name string
}

func NewPackage(name string) *Package {
	return &Package{
		Name: name,
	}
}
func (p *Package) LocalName() string {
	pos := strings.LastIndexByte(p.Name, '/')
	if pos == -1 {
		return p.Name
	}
	return p.Name[pos+1:]
}
