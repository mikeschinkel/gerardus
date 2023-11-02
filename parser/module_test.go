package parser_test

import (
	"testing"

	"gerardus/parser"
)

type newModuleWant struct {
	ModuleName  string
	PackageDir  string
	Version     string
	Path        string
	GoVersion   string
	PackageType parser.PackageType
}

var moduleGraph = parser.NewModuleGraph()

func init() {

	pm := moduleGraph.AddProjectModule(&parser.ModuleArgs{
		ModuleGraph: moduleGraph,
		Name:        "gerardus",
		PackageType: parser.GoModPackage,
		GoVersion:   "1.21",
		Version:     ".",
		PackageDir:  rootPath(""),
		Path:        rootPath("go.mod"),
	})

	moduleGraph.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "github.com/mattn/go-sqlite3",
		Version: "v1.14.17",
	})

	moduleGraph.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "golang.org/x/mod",
		Version: "v0.13.0",
	})

	moduleGraph.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "golang.org/x/sync",
		Version: "v0.3.0",
	})

	pm = moduleGraph.AddProjectModule(&parser.ModuleArgs{
		ModuleGraph: moduleGraph,
		Name:        "gerardus/cmd",
		PackageType: parser.GoModPackage,
		GoVersion:   "1.20",
		Version:     ".",
		PackageDir:  rootPath("cmd"),
		Path:        rootPath("cmd/go.mod"),
	})

	moduleGraph.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "github.com/mattn/go-sqlite3",
		Version: "v1.14.17",
	})

	moduleGraph.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "golang.org/x/mod",
		Version: "v0.13.0",
	})

	moduleGraph.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "golang.org/x/sync",
		Version: "v0.3.0",
	})
}
func TestNewModule(t *testing.T) {
	tests := []struct {
		args *parser.ModuleArgs
		want newModuleWant
	}{
		{
			args: &parser.ModuleArgs{
				ModuleGraph: moduleGraph,
				Name:        "gerardus",
				PackageDir:  rootPath(""),
				Version:     ".",
				Path:        rootPath("go.mod"),
				GoVersion:   "1.21",
				PackageType: parser.GoModPackage,
			},
			want: newModuleWant{
				ModuleName:  "geradus",
				PackageDir:  rootPath(""),
				Version:     ".",
				Path:        rootPath("go.mod"),
				GoVersion:   "1.21",
				PackageType: parser.GoModPackage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.Name, func(t *testing.T) {
			got := parser.NewModule(tt.args)
			equals(t, "ModuleName", got.Name, tt.want.ModuleName)
		})
	}
}
