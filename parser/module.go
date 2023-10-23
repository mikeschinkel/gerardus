package parser

import (
	"fmt"
	"strings"
)

type Module struct {
	Name      string
	Version   string
	GoVersion string
	Parent    *Module
}

var goPackageURLFormat = "https://github.com/golang/go/tree/go%s/src/%s"

// OriginPath returns the composed module path for a package
func (m Module) OriginPath() (path string) {
	_, ok := Modules[m.Name]
	if ok {
		path = m.Name
		goto end
	}
	path = fmt.Sprintf("%s/%s", m.Name, path)
	if m.Parent != nil && strings.HasPrefix(m.Name, m.Parent.Name) {
		// Local package
		goto end
	}
	// Go standard library package
	path = fmt.Sprintf(goPackageURLFormat, m.GoVersion, m.Name)
end:
	return path
}
