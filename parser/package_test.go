package parser_test

import (
	"testing"

	"github.com/mikeschinkel/gerardus/parser"
)

type newPackageWant struct {
	ImportPath  string
	PackageName string
	PackageType parser.PackageType
	Module      *parser.Module
	ModuleGraph *parser.ModuleGraph
	Directory   string
	Version     string
	String      string
	Sources     sourceWant
}
type sourceWant struct {
	Source        string
	SourceVersion string
}

func TestNewPackage(t *testing.T) {
	tests := []struct {
		args *parser.PackageArgs
		want newPackageWant
	}{
		{
			args: &parser.PackageArgs{
				ImportPath: "gerardus",
				Type:       parser.GoModPackage,
				Version:    ".",
				Directory:  rootPath(""),
			},
			want: newPackageWant{
				ImportPath:  "gerardus",
				PackageName: "gerardus",
				PackageType: parser.GoModPackage,
				Directory:   rootPath(""),
				Version:     ".",
				Sources: sourceWant{
					Source:        "./gerardus",
					SourceVersion: "./gerardus",
				},
			},
		},
		{
			args: &parser.PackageArgs{
				ImportPath: "github.com/mikeschinkel/gerardus/cmd",
				Type:       parser.GoModPackage,
				Version:    ".",
				Directory:  rootPath("cmd"),
			},
			want: newPackageWant{
				ImportPath:  "github.com/mikeschinkel/gerardus/cmd",
				PackageName: "cmd",
				PackageType: parser.GoModPackage,
				Directory:   rootPath("cmd"),
				Version:     ".",
				Sources: sourceWant{
					Source:        "https://github.com/mikeschinkel/gerardus/cmd",
					SourceVersion: "https://github.com/mikeschinkel/gerardus/cmd",
				},
			},
		},
		{
			args: &parser.PackageArgs{
				ImportPath: "github.com/mattn/go-sqlite3",
				Version:    "v1.14.17",
			},
			want: newPackageWant{
				ImportPath:  "github.com/mattn/go-sqlite3",
				PackageName: "sqlite3",
				PackageType: parser.ExternalPackage,
				Version:     "v1.14.17",
				Sources: sourceWant{
					Source:        "https://github.com/mattn/go-sqlite3",
					SourceVersion: "https://github.com/mattn/go-sqlite3/tree/v1.14.17",
				},
			},
		},
		{
			args: &parser.PackageArgs{
				ImportPath: "fmt",
			},
			want: newPackageWant{
				ImportPath:  "fmt",
				String:      "fmt",
				PackageName: "fmt",
				PackageType: parser.StdLibPackage,
				Sources: sourceWant{
					Source:        "https://github.com/golang/go/tree/go$VERSION",
					SourceVersion: "https://github.com/golang/go/tree/go$VERSION/src/fmt",
				},
			},
		},
		{
			args: &parser.PackageArgs{
				ImportPath: "encoding/json",
			},
			want: newPackageWant{
				ImportPath:  "encoding/json",
				String:      "encoding/json",
				PackageName: "json",
				PackageType: parser.StdLibPackage,
				Sources: sourceWant{
					Source:        "https://github.com/golang/go/tree/go$VERSION",
					SourceVersion: "https://github.com/golang/go/tree/go$VERSION/src/encoding/json",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.ImportPath+" — New()", func(t *testing.T) {
			got := parser.NewPackage(tt.args)
			equals(t, "Source", got.Source(), tt.want.Sources.Source)
			equals(t, "Source Version", got.PackageVersion.Source(), tt.want.Sources.SourceVersion)
			equals(t, "ImportPath", got.ImportPath, tt.want.ImportPath)
			equals(t, "Version", got.PackageVersion.Name, tt.want.Version)
			equals(t, "Name", got.Name(), tt.want.PackageName)
			equals(t, "Package Type", got.Type.Name(), tt.want.PackageType.Name())
			equals(t, "Module", got.Module, tt.want.Module)
			equals(t, "ModuleGraph", got.ModuleGraph, tt.want.ModuleGraph)
			strPtrEquals(t, "Directory", got.Directory, tt.want.Directory)
		})
		t.Run(tt.args.ImportPath+" — PartialClone()", func(t *testing.T) {
			pkg := parser.NewPackage(tt.args)
			pkg.Module = &parser.Module{}
			got := pkg.PartialClone()
			equals(t, "ImportPath", got.ImportPath, pkg.ImportPath)
			equals(t, "Module", got.Module, pkg.Module)
			notEquals(t, "Details", got.PackageDetails, pkg.PackageDetails)
			notEquals(t, "Version", got.PackageDetails.PackageVersion, pkg.PackageDetails.PackageVersion)
			equals(t, "Version Package", got.PackageVersion.Package, pkg)
		})
	}
}
