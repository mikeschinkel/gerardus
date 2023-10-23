package parser

import (
	"gerardus/collector"
	"gerardus/scanner"
	"golang.org/x/mod/modfile"
)

var Modules = make(map[string]struct{})

var _ collector.CodeFacet = (*ModFile)(nil)

type ModFile struct {
	scanner.File
	Content []byte
	ModFile *modfile.File
}

// Name returns the module name as defined in go.mod by the `module` statement
func (mf ModFile) Name() string {
	return mf.ModFile.Module.Mod.Path
}

// GoVersion returns the Go version set in this module
func (mf ModFile) GoVersion() string {
	return mf.ModFile.Go.Version
}

// Require returns the version of the module as defined in go.mod by the `module` statement
func (mf ModFile) Require() []*modfile.Require {
	return mf.ModFile.Require
}

// Version returns the version of the module as defined in go.mod by the `module` statement
func (mf ModFile) Version() (v string) {
	if len(mf.ModFile.Module.Mod.Version) == 0 {
		v = "."
		goto end
	}
	v = mf.ModFile.Module.Mod.Version
end:
	return v
}

func NewModFile(file scanner.File, content []byte) *ModFile {
	return &ModFile{
		File:    file,
		Content: content,
	}
}

// Modules returns a slice of Mod where the first element is the module and the
// rest are its required dependencies. m.ModFile, hence the name.
func (mf *ModFile) Modules() []*Module {
	modFile := mf.ModFile
	modules := make([]*Module, 0, len(modFile.Require)+1)

	m := &Module{
		Name:      mf.Name(),
		Version:   mf.Version(),
		GoVersion: mf.GoVersion(),
	}
	modules = append(modules, m)
	for _, r := range modFile.Require {
		modules = append(modules, &Module{
			Name:      r.Mod.Path,
			Version:   r.Mod.Version,
			GoVersion: modFile.Go.Version,
			Parent:    m,
		})
	}
	return modules
}

// CodeFacet simply marks ModFile as implementing collector.CodeFacet
func (ModFile) CodeFacet() {}

func (mf ModFile) String() string {
	return mf.ModFile.Module.Mod.String()
}
