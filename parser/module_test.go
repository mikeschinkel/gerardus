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

func getModuleGraph() *parser.ModuleGraph {
	mg := parser.NewModuleGraph()
	pm := mg.AddProjectModule(&parser.ModuleArgs{
		ModuleGraph: mg,
		Name:        "gerardus",
		PackageType: parser.GoModPackage,
		GoVersion:   "1.21",
		Version:     ".",
		PackageDir:  rootPath(""),
		Path:        rootPath("go.mod"),
	})

	mg.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "github.com/mattn/go-sqlite3",
		Version: "v1.14.17",
	})

	mg.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "golang.org/x/mod",
		Version: "v0.13.0",
	})

	mg.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "golang.org/x/sync",
		Version: "v0.3.0",
	})

	pm = mg.AddProjectModule(&parser.ModuleArgs{
		ModuleGraph: mg,
		Name:        "gerardus/cmd",
		PackageType: parser.GoModPackage,
		GoVersion:   "1.20",
		Version:     ".",
		PackageDir:  rootPath("cmd"),
		Path:        rootPath("cmd/go.mod"),
	})

	mg.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "github.com/mattn/go-sqlite3",
		Version: "v1.14.17",
	})

	mg.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "golang.org/x/mod",
		Version: "v0.13.0",
	})

	mg.AddDependentModule(pm, &parser.ModuleArgs{
		Name:    "golang.org/x/sync",
		Version: "v0.3.0",
	})
	return mg
}
func TestNewModule(t *testing.T) {
	var moduleGraph = getModuleGraph()
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
				ModuleName:  "gerardus",
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
			equals(t, "ModuleName", got.Name(), tt.want.ModuleName)
		})
	}
}
