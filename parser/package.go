package parser

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Package holds the name used early in the pipeline; during parsing
type Package struct {
	*PackageDetails
	ImportPath  string
	debugString string
}

// PackageDetails holds values used later in the pipeline
type PackageDetails struct {
	PackageVersion *PackageVersion
	Module         *Module
	ModuleGraph    *ModuleGraph

	// Directory contains the full path to the package's directory
	// TODO: Change this to relative path when all paths are changed to relative
	Directory *string
	Type      PackageType
}

type PackageArgs struct {
	ModuleGraph *ModuleGraph
	Module      *Module
	ImportPath  string
	Directory   string
	Type        PackageType
	Version     string
}

func NewFlyweightPackage(importPath string) *Package {
	return &Package{
		ImportPath: importPath,
	}
}

func dispensePackage(args *PackageArgs) (pkg *Package) {
	pkg = args.ModuleGraph.GetPackageByImportPath(args.ImportPath)
	if pkg != nil {
		goto end
	}
	pkg = newPackage(args)
end:
	return pkg
}

func newPackage(args *PackageArgs) *Package {
	var pkg *Package
	var dir string

	print()

	pkg = &Package{
		ImportPath: args.ImportPath,
		PackageDetails: &PackageDetails{
			Type:        args.Type,
			Module:      args.Module,
			ModuleGraph: args.ModuleGraph,
			PackageVersion: NewPackageVersion(&VersionArgs{
				Name: args.Version,
			}),
		},
	}
	if pkg.Type == InvalidPackage {
		switch {
		case strings.Contains(args.ImportPath, "."):
			pkg.Type = ExternalPackage

		case pkg.VersionName() == ".":
			pkg.Type = LocalPackage

		default:
			pkg.Type = StdLibPackage
		}
	}

	if args.Directory != "" {
		switch pkg.Type {
		case GoModPackage, LocalPackage:
			dir = args.Directory
			pkg.Directory = &dir
		}
	}
	pkg.debugString = pkg.String()
	pkg.PackageVersion.Package = pkg

	return pkg
}

// PartialClone returns a cloned version of Package completely disconnected from
// any pointers, EXCEPT for maintaining Module and ModuleGraph as pointers.
func (p Package) PartialClone() (pkg *Package) {
	clone := p
	pd := *p.PackageDetails
	pv := *p.PackageDetails.PackageVersion
	pkg = &Package{
		ImportPath:     p.ImportPath,
		PackageDetails: &pd,
	}
	p.PackageVersion = &pv
	p.PackageVersion.Package = &clone
	p.debugString = p.String()
	return pkg
}

// StdLibSourceURL returns the source URL for a Go package.
// Adds a fragment #invalid-version if URL resolution failed.
func (p Package) StdLibSourceURL() (url string) {
	svURL := p.PackageVersion.StdLibSourceURL()
	url, _, _ = strings.Cut(svURL, "/src/")
	return url
}

// ExternalURL returns a URL for an external package. Panics if the package
// is not external.
func (p Package) ExternalURL() (url string) {
	genericURL := func() string {
		// Just do our best
		url, _, _ = strings.Cut(p.PackageVersion.ExternalURL(), "@")
		return url + "#uri-only"
	}
	host, _, found := strings.Cut(p.ImportPath, "/")
	if !found {
		url = genericURL()
		goto end
	}
	switch host {
	case "github.com":
		url, _, _ = strings.Cut(p.PackageVersion.ExternalURL(), "/tree/")
	case "golang.org":
		url, _, _ = strings.Cut(p.PackageVersion.ExternalURL(), "/+/")
	default:
		url = genericURL()
	}
end:
	return url
}

// Source returns a string of its source
func (p Package) Source() (src string) {
	if p.PackageDetails == nil {
		src = "unknown"
		goto end
	}
	switch p.Type {
	case StdLibPackage:
		src = p.StdLibSourceURL()
	case GoModPackage, LocalPackage:
		src = p.LocalSourceDir()
	case ExternalPackage:
		src = p.ExternalURL()
	case InvalidPackage:
		fallthrough
	default:
		panicf("Unexpected invalid package type %d", p.Type)
	}
end:
	return src
}

func (p Package) LocalSourceDir() (dir string) {
	dir = fmt.Sprintf("%s/%s", filepath.Dir(p.GoModPath()), p.ImportPath)
	return dir
}

func (p Package) VersionName() (s string) {
	return p.PackageVersion.Name
}

func (p Package) SetVersion(ver string) {
	p.PackageVersion.Name = ver
}

func (p Package) GoModPath() (s string) {
	if p.Module == nil {
		return "~module_not_set~"
	}
	return p.Module.GoModPath()
}

// Name returns the last segemnt of the name, e.g. `baz` for `foo/bar/baz`.
func (p Package) Name() (s string) {
	pos := strings.LastIndexByte(p.ImportPath, '/')
	if pos == -1 {
		return p.ImportPath
	}
	s = p.ImportPath[pos+1:]
	if strings.HasPrefix(s, "go-") {
		s = s[3:]
	}
	return s
}

// String returns a textual representation of Package for error messages and
// database fields.
func (p Package) String() (s string) {
	var name string
	s = p.ImportPath
	if p.PackageDetails == nil {
		goto end
	}
	name = p.VersionName()
	if name == "" {
		goto end
	}
	if name == "." {
		goto end
	}
	s = fmt.Sprintf("%s@%s", p.ImportPath, name)
end:
	return s
}
