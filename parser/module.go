package parser

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
)

type GoMod struct {
	Version string
	Path    string
	Tag     string
}
type Module struct {
	*Package
	*GoMod
	Parent      *Module
	Module      *Module
	ModuleGraph *ModuleGraph
	pathMap     PathMap
	debugString string
}

type ModuleArgs struct {
	ModuleGraph *ModuleGraph
	Parent      *Module
	Module      *Module
	Name        string
	PackageDir  string
	Version     string
	Path        string
	GoVersion   string
	PackageType PackageType
}

func newModule(args *ModuleArgs) *Module {
	var gm *GoMod
	var err error

	print()
	if args.PackageDir == "" {
		args.PackageDir = filepath.Dir(args.Path)
	}
	m := &Module{}
	pkg := newPackage(&PackageArgs{
		ImportPath:  args.Name,
		Directory:   args.PackageDir,
		Type:        args.PackageType,
		Version:     args.Version,
		ModuleGraph: args.ModuleGraph,
		Module:      m,
	})
	if args.Version == "." {
		// Only create this for Go mod that are being loaded.
		gm = &GoMod{
			Version: args.GoVersion,
			Path:    args.Path,
		}
	} else {
		// Go Mods that are dependent on other go mods will reference the source go mod
		// and will have a long version, e.g. v0.0.0-00010101000000-000000000000, so we
		// do n\ot need to recreate GoMod as we already have it.
		gm = args.Parent.GoMod
	}
	m.Package = pkg
	m.GoMod = gm
	m.ModuleGraph = args.ModuleGraph
	m.Parent = args.Parent
	m.pathMap = make(PathMap)
	m.debugString = fmt.Sprintf("[go%s] %s",
		args.GoVersion,
		args.Name,
	)

	switch pkg.Type {
	case StdLibPackage:
		// Sometimes the version names in go.mod don't match the version tags on GitHub
		// for golang/go repository. So this code tries to fix the version number.
		err = m.MaybeFixGoVersion()
		if err != nil {
			slog.Error("Unable to get package URL for version; likely a programming logic problem",
				"version", pkg.VersionName(),
				"error", err.Error(),
			)
		}
		pkg.SetVersion(m.GoMod.Tag)
		pkg.debugString = pkg.String()

	default:
		goto end
	}

end:
	return m
}

// GoModDir returns the path to the go.mod file for the module, WITHOUT the suffix "/go.mod".
func (m *Module) GoModDir() string {
	return filepath.Dir(m.GoModPath())
}

// GoModPath returns the path to the go.mod file for the module, WITH the suffix "/go.mod".
func (m *Module) GoModPath() string {
	if m.GoMod == nil {
		return "$GOMODPATH"
	}
	return m.GoMod.Path
}

// GoVersion returns Go version as a string
func (m *Module) GoVersion() (s string) {
	return m.GoMod.Version
}

func (m *Module) Name() string {
	return m.Package.ImportPath
}

func (m *Module) VersionName() string {
	return m.Package.PackageVersion.Name
}

func (m *Module) Version() *PackageVersion {
	return m.Package.PackageVersion
}

//// OriginPath returns the composed module path for a package
//func (m *Module) OriginPath() (path string) {
//	var mv *ModuleVersions
//
//	cacheKey := m.String()
//	p, ok := m.pathMap[cacheKey]
//	if ok {
//		// Previously found and cached
//		path = p
//		goto end
//	}
//	mv, ok = m.ModuleGraph.ModuleMap[m.Name()]
//	if !ok {
//		panicf("Unexpected origin path found for %s: %#v", m.Name(), m)
//		//path = m.Name
//		//m.pathMap[cacheKey] = path
//		//goto end
//	}
//	path, ok = mv.PathForVersion(m.VersionName())
//	if ok {
//		// Internal dependency, e.g. name of ./go.mod referenced in ./cmd/go.mod
//		m.pathMap[cacheKey] = path
//		goto end
//	}
//	// Panic since this should only happen is there is a logic bug in the program.
//	panicf("No origin path found for %s: %#v", m.String(), m)
//end:
//	return path
//}

var goVersionMap = NewSafeMap[string, string]()

// MaybeFixGoVersion checks to see if Go has a GitHub tag matching the version in
// the go.mod file, or if we need to embellish it. Embellishing it means either
// adding a ".0" to the end, or removing it. More embellishment may be needed in
// the future, if the Go team gets sloppier with their tagging.
func (m *Module) MaybeFixGoVersion() (err error) {
	var fixedVer string
	var status int
	var ok bool
	var url string

	goVer := m.GoVersion()

	fixedVer, ok = goVersionMap.Load(goVer)
	if ok {
		goto end
	}
	url = fmt.Sprintf(goVersionURLFormat, goVer)
	status, _ = CheckURL(url)
	if status == http.StatusOK {
		goto end
	}
	// Because the Go team has been inconsistent in tagging, we have to deal with it here.
	switch strings.Count(goVer, ".") {
	case 1:
		// Add trailing `.0 `and try again, e.g. 1.21 => 1.21.0
		fixedVer = goVer + ".0"
	case 2:
		// Remove trailing `.0 `and try again, e.g. 1.20.0 => 1.20
		fixedVer, ok = strings.CutSuffix(goVer, ".0")
		if !ok {
			// Accept that the URL is wrong, but use it with an "invalid" tag
			err = fmt.Errorf("version URL is invalid: %s", url)
			goto end
		}
	}
	url = fmt.Sprintf(goVersionURLFormat, fixedVer)
	status, _ = CheckURL(url)
	switch status {
	case http.StatusOK:
		m.GoMod.Tag = fixedVer
	default:
		m.GoMod.Tag = "invalid"
		// Accept that the URL is wrong, but use it with an "invalid" tag
		err = fmt.Errorf("version URL is invalid: %s", url)
		goto end
	}
end:
	if !goVersionMap.Has(goVer) {
		goVersionMap.Save(goVer, fixedVer)
	}
	return err
}
