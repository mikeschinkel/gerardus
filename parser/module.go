package parser

import (
	"fmt"
)

type Module struct {
	Name      string
	Version   string
	GoVersion string
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

	if len(source) == 0 {
		// Dependency
		path = m.Name
		goto end
	}

	// Internal dependency, e.g. name of ./go.mod referenced in ./cmd/go.mod
	path = source
end:
	return path
}
