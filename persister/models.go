// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package persister

import ()

type Category struct {
	ID       int64
	SurveyID int64
	Name     string
}

type CategoryType struct {
	ID         int64
	Name       string
	CategoryID int64
	SurveyID   int64
	TypeID     int64
}

type Codebase struct {
	ID         int64
	ProjectID  int64
	VersionTag string
	SourceUrl  string
}

type File struct {
	ID       int64
	SurveyID int64
	Filepath string
}

type Import struct {
	ID        int64
	FileID    int64
	SurveyID  int64
	PackageID int64
	Alias     string
}

type Method struct {
	ID       int64
	Name     string
	Params   string
	Results  string
	FileID   int64
	SurveyID int64
	TypeID   int64
}

type Module struct {
	ID   int64
	Name string
}

type ModuleVersion struct {
	ID       int64
	ModuleID int64
	Version  string
}

type ModuleVersionView struct {
	ID       int64
	ModuleID int64
	Name     string
	Version  string
}

type Package struct {
	ID         int64
	ImportPath string
	Source     string
	TypeID     int64
	Name       interface{}
}

type PackageType struct {
	ID   int64
	Name string
}

type PackageVersion struct {
	ID        int64
	PackageID int64
	Version   string
	SourceUrl string
}

type Project struct {
	ID      int64
	Name    string
	RepoUrl string
	About   string
	Website string
}

type Survey struct {
	ID         int64
	CodebaseID int64
	LocalDir   string
	Timestamp  string
}

type SurveyModule struct {
	ID              int64
	SurveyID        int64
	ModuleID        int64
	ModuleVersionID int64
	PackageID       int64
	FileID          int64
}

type SurveyView struct {
	ID         int64
	Project    string
	CodebaseID int64
	ProjectID  int64
	RepoUrl    string
	VersionTag string
	SourceUrl  string
	LocalDir   string
	Timestamp  string
}

type SymbolType struct {
	ID   int64
	Name string
}

type Type struct {
	ID           int64
	FileID       int64
	SurveyID     int64
	SymbolTypeID int64
	Name         string
	Definition   string
}

type TypeView struct {
	ID           int64
	Project      string
	Filepath     string
	Name         string
	SymbolName   string
	SurveyID     int64
	FileID       int64
	SymbolTypeID int64
	CodebaseID   int64
	Definition   string
	Timestamp    string
	SourceUrl    string
	RepoUrl      string
	AboutProject string
	LocalDir     string
}

type Variable struct {
	ID       int64
	Name     string
	SurveyID int64
	TypeID   int64
	Usage    int64
	IsParam  interface{}
	IsResult interface{}
}
