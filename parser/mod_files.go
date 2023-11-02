package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	"gerardus/collector"
	"gerardus/scanner"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

var _ collector.CodeFacet = (*ModFile)(nil)

type ModFile struct {
	scanner.File
	Content     []byte
	modFile     *modfile.File
	ModuleGraph *ModuleGraph
	rootDir     string
	debugString string
}

// Name returns the module name as defined in go.mod by the `module` statement
func (mf ModFile) Name() string {
	return mf.modFile.Module.Mod.Path
}

// UniqueID returns a unique ID string that incorporates a version, if non-zero.
func (mf ModFile) UniqueID() (uId string) {
	return mf.modFile.Module.Mod.String()
}

// PathVersion returns the module.Version stored in mf.modFile.ModuleArgs.Mod,
// however, we replace the .Path value with mf.Fullpath().
func (mf ModFile) PathVersion() module.Version {
	mv := mf.modFile.Module.Mod
	mv.Path = filepath.Dir(mf.Fullpath())
	if len(mv.Version) == 0 {
		mv.Version = "."
	}
	return mv
}

// GoVersion returns the Go version set in this module
func (mf ModFile) GoVersion() string {
	return mf.modFile.Go.Version
}

func (mf *ModFile) SetModFile(f *modfile.File) {
	mf.modFile = f
	mf.debugString = mf.String()
}

func (mf *ModFile) ModFile() *modfile.File {
	return mf.modFile
}

// Require returns the version of the module as defined in go.mod by the `module` statement
func (mf ModFile) Require() []*modfile.Require {
	return mf.modFile.Require
}

// Version returns the version of the module as defined in go.mod by the `module` statement
func (mf ModFile) Version() (v string) {
	if len(mf.modFile.Module.Mod.Version) == 0 {
		v = "."
		goto end
	}
	v = mf.modFile.Module.Mod.Version
end:
	return v
}

func NewModFile(file scanner.File, content []byte, mg *ModuleGraph, rootDir string) *ModFile {
	return &ModFile{
		File:        file,
		Content:     content,
		ModuleGraph: mg,
		rootDir:     rootDir,
	}
}

func (mf *ModFile) ImportPath() (s string) {
	d, found := mf.File.RelDir(filepath.Dir(mf.rootDir))
	if !found {
		panicf("Unexpected pattern: go.mod's dir %s not equal to, or a parent of source dir %s.",
			mf.File.SourceDir(),
			mf.rootDir,
		)
	}
	return d
}

// Modules returns a slice of Mod where the first element is the module and the
// rest are its required dependencies. m.modFile, hence the name.
func (mf *ModFile) Modules() []*Module {
	modFile := mf.modFile
	modules := make([]*Module, 0, len(modFile.Require)+1)
	oops := func(ip string) {
		panicf("Unexpected: module file import path '%s' not found in module graph's module map.", ip)
	}
	mm := mf.ModuleGraph.ModuleMap
	mv, ok := mm[mf.ImportPath()]
	if !ok {
		oops(mf.ImportPath())
	}
	modules = append(modules, mv.Module)
	for _, r := range modFile.Require {
		mv, ok = mm[r.Mod.Path]
		if !ok {
			oops(r.Mod.Path)
		}
		m := mf.ModuleGraph.DispenseModule(r.Mod.Path, mf.SourceDir())
		if m == nil {
			oops(fmt.Sprintf("%s [%s]",
				r.Mod.Path,
				mf.SourceDir(),
			))
		}
		modules = append(modules, m)
	}
	return modules

}
func derivePackageType(mv module.Version) (pt PackageType) {
	switch {
	case strings.HasPrefix(mv.Version, "v0.0.0-"):
		pt = LocalPackage
	case mv.Version == ".":
		pt = LocalPackage
	case strings.Contains(mv.Path, "."):
		pt = ExternalPackage
	default:
		pt = StdLibPackage
	}
	return pt
}

// CodeFacet simply marks ModFile as implementing collector.CodeFacet
func (ModFile) CodeFacet() {}

// String returns a string representation of ModFile for debugging and error messages.
func (mf ModFile) String() string {
	modFile := mf.modFile
	return fmt.Sprintf("[go%s] %s",
		modFile.Go.Version,
		modFile.Module.Mod.String(),
	)
}
