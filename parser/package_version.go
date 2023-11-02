package parser

import (
	"fmt"
	"strings"
)

type PackageVersion struct {
	Package *Package
	Name    string
}

type VersionArgs struct {
	Package *Package
	Name    string
}

func NewPackageVersion(args *VersionArgs) *PackageVersion {
	return &PackageVersion{
		Package: args.Package,
		Name:    args.Name,
	}
}

var goPackageURLFormat = "https://github.com/golang/go/tree/go%s/src/%s"

// StdLibSourceURL returns the source URL for a package for the Go version.
// Adds a fragment #invalid-version if URL resolution failed.
func (v PackageVersion) StdLibSourceURL() (url string) {
	if v.Package.Type != StdLibPackage {
		panicf("Unexpected; package type not StdLib: %s", v.Package.Type.Name())
	}
	ver := v.Name
	if ver == "" {
		ver = "$VERSION"
	}
	return fmt.Sprintf(goPackageURLFormat, ver, v.Package.ImportPath)
}

// ExternalURL returns a URL for an external package at the given version.
// Panics if the package is not external.
func (v PackageVersion) ExternalURL() (url string) {
	if v.Package.Type != ExternalPackage {
		panicf("Unexpected; package type not External: %s", v.Package.Type.Name())
	}
	genericURL := func() string {
		// Just do our best
		return fmt.Sprintf("https://%s@%s#uri-only", v.Package.ImportPath, v.Name)
	}
	host, after, found := strings.Cut(v.Package.ImportPath, "/")
	if !found {
		url = genericURL()
		goto end
	}
	switch host {
	case "github.com":
		url = fmt.Sprintf("https://%s/tree/%s", v.Package.ImportPath, v.Name)
	case "golang.org":
		url = fmt.Sprintf("https://cs.opensource.google/go/%s/+/refs/tags/%s:",
			after,
			v.Name,
		)
	default:
		url = genericURL()
	}
end:
	return url
}

// Source returns a string of its source, complete with a version where applicable
func (v PackageVersion) Source() (src string) {
	switch v.Package.Type {
	case StdLibPackage:
		src = v.StdLibSourceURL()
	case GoModPackage, LocalPackage:
		src = v.Package.Source()
	case ExternalPackage:
		src = v.ExternalURL()
	case InvalidPackage:
		fallthrough
	default:
		panicf("Unexpected invalid package type '%s", v.Package.Type.Name())
	}
	return src
}
