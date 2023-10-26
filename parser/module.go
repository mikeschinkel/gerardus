package parser

import (
	"fmt"
	"path/filepath"
)

type Module struct {
	Name      string
	Version   string
	GoVersion string
	Filepath  string
}

var goPackageURLFormat = "https://github.com/golang/go/tree/go%s/src/%s"

// OriginPath returns the composed module path for a package
func (m Module) OriginPath() (path string) {
	source, ok := Modules[m.Name]
	if !ok {
		// Go standard library package(?)
		path = fmt.Sprintf(goPackageURLFormat, m.GoVersion, m.Name)
		goto end
	}
	switch source {
	case ModuleFile:
		path = filepath.Dir(m.Filepath)
	case ModuleDependency:
		path = m.Name
	}
end:
	return path
}
