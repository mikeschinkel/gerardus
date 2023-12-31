package parser_test

import (
	"testing"

	"github.com/mikeschinkel/gerardus/parser"
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
		Name:        "github.com/mikeschinkel/gerardus/cmd",
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
			equals(t, "Name", got.Name(), tt.want.ModuleName)
			equals(t, "Version", got.Version().Name, tt.want.Version)
			equals(t, "Version", got.GoModPath(), tt.want.Path)
			equals(t, "Package Type", got.Package.Type, tt.want.PackageType)
			strPtrEquals(t, "Package Dir", got.Package.Directory, tt.want.PackageDir)
			equals(t, "GoVersion", got.GoMod.Version, tt.want.GoVersion)
			equals(t, "GoVersion", got.GoVersion(), tt.want.GoVersion)
		})
	}
}
