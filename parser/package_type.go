package parser

import (
	"fmt"
)

type PackageType int

const (
	InvalidPackage  PackageType = 0
	StdLibPackage   PackageType = 1
	GoModPackage    PackageType = 2
	LocalPackage    PackageType = 3
	ExternalPackage PackageType = 4
)

func (pt PackageType) ID() int {
	return int(pt)
}

func (pt PackageType) Name() string {
	return PackageName(pt)
}

var PackageTypes = []PackageType{
	InvalidPackage,
	StdLibPackage,
	GoModPackage,
	LocalPackage,
	ExternalPackage,
}

var PackageTypeMap = map[string]PackageType{}

func PackageName(pkgType PackageType) string {
	switch pkgType {

	case InvalidPackage:
		return "Invalid"
	case StdLibPackage:
		return "StdLib"
	case GoModPackage:
		return "GoMod"
	case LocalPackage:
		return "Local"
	case ExternalPackage:
		return "External"
	}
	return fmt.Sprintf("Unclassified[%d]", pkgType)
}

func init() {
	for _, pt := range PackageTypes {
		PackageTypeMap[PackageName(pt)] = pt
	}
}
