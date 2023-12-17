package app

import (
	"context"

	"github.com/mikeschinkel/gerardus/persister"
)

// DataStoreQueries ...
type DataStoreQueries interface {
	DeleteCategory(ctx context.Context, id int64) error
	DeleteCodebase(ctx context.Context, id int64) error
	DeleteCodebaseByProjectIdAndVersionTag(ctx context.Context, arg persister.DeleteCodebaseByProjectIdAndVersionTagParams) error
	DeleteCodebaseSurveys(ctx context.Context, codebaseID int64) error
	DeleteFile(ctx context.Context, id int64) error
	DeleteImport(ctx context.Context, id int64) error
	DeleteMethod(ctx context.Context, id int64) error
	DeleteModule(ctx context.Context, id int64) error
	DeleteModuleVersion(ctx context.Context, id int64) error
	DeletePackage(ctx context.Context, id int64) error
	DeletePackageType(ctx context.Context, id int64) error
	DeletePackageVersion(ctx context.Context, id int64) error
	DeleteProject(ctx context.Context, id int64) error
	DeleteProjectByName(ctx context.Context, name string) error
	DeleteSurvey(ctx context.Context, id int64) error
	DeleteSurveyModule(ctx context.Context, id int64) error
	DeleteSymbolType(ctx context.Context, id int64) error
	DeleteType(ctx context.Context, id int64) error
	DeleteVariable(ctx context.Context, id int64) error
	InsertCategory(ctx context.Context, arg persister.InsertCategoryParams) (persister.Category, error)
	InsertCodebase(ctx context.Context, arg persister.InsertCodebaseParams) (persister.Codebase, error)
	InsertFile(ctx context.Context, arg persister.InsertFileParams) (persister.File, error)
	InsertImport(ctx context.Context, arg persister.InsertImportParams) (persister.Import, error)
	InsertMethod(ctx context.Context, arg persister.InsertMethodParams) (persister.Method, error)
	InsertModule(ctx context.Context, name string) (persister.Module, error)
	InsertModuleVersion(ctx context.Context, arg persister.InsertModuleVersionParams) (persister.ModuleVersion, error)
	InsertPackage(ctx context.Context, arg persister.InsertPackageParams) (persister.Package, error)
	InsertPackageType(ctx context.Context, arg persister.InsertPackageTypeParams) (persister.PackageType, error)
	InsertPackageVersion(ctx context.Context, arg persister.InsertPackageVersionParams) (persister.PackageVersion, error)
	InsertProject(ctx context.Context, arg persister.InsertProjectParams) (persister.Project, error)
	InsertSurvey(ctx context.Context, arg persister.InsertSurveyParams) (persister.Survey, error)
	InsertSurveyModule(ctx context.Context, arg persister.InsertSurveyModuleParams) (persister.SurveyModule, error)
	InsertSymbolType(ctx context.Context, arg persister.InsertSymbolTypeParams) (persister.SymbolType, error)
	InsertType(ctx context.Context, arg persister.InsertTypeParams) (persister.Type, error)
	InsertVariable(ctx context.Context, arg persister.InsertVariableParams) (persister.Variable, error)
	ListCategories(ctx context.Context) ([]persister.Category, error)
	ListCodebaseSurveys(ctx context.Context, codebaseID int64) ([]persister.Survey, error)
	ListCodebases(ctx context.Context) ([]persister.Codebase, error)
	ListFiles(ctx context.Context) ([]persister.File, error)
	ListFilesBySurvey(ctx context.Context, surveyID int64) ([]persister.File, error)
	ListImports(ctx context.Context) ([]persister.Import, error)
	ListMethods(ctx context.Context) ([]persister.Method, error)
	ListModuleVersions(ctx context.Context) ([]persister.ModuleVersion, error)
	ListModules(ctx context.Context) ([]persister.Module, error)
	ListPackageTypes(ctx context.Context) ([]persister.PackageType, error)
	ListPackageTypesByName(ctx context.Context) ([]persister.PackageType, error)
	ListPackageVersions(ctx context.Context) ([]persister.PackageVersion, error)
	ListPackages(ctx context.Context) ([]persister.Package, error)
	ListProjects(ctx context.Context) ([]persister.Project, error)
	ListSurveyModules(ctx context.Context) ([]persister.SurveyModule, error)
	ListSurveys(ctx context.Context) ([]persister.ListSurveysRow, error)
	ListSymbolTypes(ctx context.Context) ([]persister.SymbolType, error)
	ListSymbolTypesByName(ctx context.Context) ([]persister.SymbolType, error)
	ListTypes(ctx context.Context) ([]persister.TypeView, error)
	ListTypesByFile(ctx context.Context, fileID int64) ([]persister.TypeView, error)
	ListTypesBySurvey(ctx context.Context, surveyID int64) ([]persister.TypeView, error)
	ListVariables(ctx context.Context) ([]persister.Variable, error)
	LoadCategory(ctx context.Context, id int64) (persister.Category, error)
	LoadCodebase(ctx context.Context, id int64) (persister.Codebase, error)
	LoadCodebaseIDByProjectAndVersion(ctx context.Context, arg persister.LoadCodebaseByProjectNameAndVersionTagParams) (int64, error)
	LoadCodebaseIdByRepoURL(ctx context.Context, repoUrl string) (int64, error)
	LoadFile(ctx context.Context, id int64) (persister.File, error)
	LoadImport(ctx context.Context, id int64) (persister.Import, error)
	LoadMethod(ctx context.Context, id int64) (persister.Method, error)
	LoadModule(ctx context.Context, id int64) (persister.Module, error)
	LoadModuleVersion(ctx context.Context, id int64) (persister.ModuleVersion, error)
	LoadPackage(ctx context.Context, id int64) (persister.Package, error)
	LoadPackageType(ctx context.Context, id int64) (persister.PackageType, error)
	LoadPackageVersion(ctx context.Context, id int64) (persister.PackageVersion, error)
	LoadProject(ctx context.Context, id int64) (persister.Project, error)
	LoadProjectByName(ctx context.Context, name string) (persister.Project, error)
	LoadProjectByRepoURL(ctx context.Context, repoUrl string) (persister.Project, error)
	LoadProjectRepoURL(ctx context.Context, id int64) (string, error)
	LoadSurvey(ctx context.Context, id int64) (persister.Survey, error)
	LoadSurveyByRepoURL(ctx context.Context, repoUrl string) (persister.LoadSurveyByRepoURLRow, error)
	LoadSurveyModule(ctx context.Context, id int64) (persister.SurveyModule, error)
	LoadSymbolType(ctx context.Context, id int64) (persister.SymbolType, error)
	LoadType(ctx context.Context, id int64) (persister.Type, error)
	LoadVariable(ctx context.Context, id int64) (persister.Variable, error)
	UpdateCategory(ctx context.Context, arg persister.UpdateCategoryParams) error
	UpdateCodebase(ctx context.Context, arg persister.UpdateCodebaseParams) error
	UpdateCodebaseByProjectIdAndVersionTag(ctx context.Context, arg persister.UpdateCodebaseByProjectIdAndVersionTagParams) error
	UpdateFile(ctx context.Context, arg persister.UpdateFileParams) error
	UpdateImport(ctx context.Context, arg persister.UpdateImportParams) error
	UpdateMethod(ctx context.Context, arg persister.UpdateMethodParams) error
	UpdateModule(ctx context.Context, arg persister.UpdateModuleParams) error
	UpdateModuleVersion(ctx context.Context, arg persister.UpdateModuleVersionParams) error
	UpdatePackage(ctx context.Context, arg persister.UpdatePackageParams) error
	UpdatePackageType(ctx context.Context, arg persister.UpdatePackageTypeParams) error
	UpdatePackageVersion(ctx context.Context, arg persister.UpdatePackageVersionParams) error
	UpdateProject(ctx context.Context, arg persister.UpdateProjectParams) error
	UpdateProjectByName(ctx context.Context, arg persister.UpdateProjectByNameParams) error
	UpdateSurveyModule(ctx context.Context, arg persister.UpdateSurveyModuleParams) error
	UpdateSymbolType(ctx context.Context, arg persister.UpdateSymbolTypeParams) error
	UpdateType(ctx context.Context, arg persister.UpdateTypeParams) error
	UpdateVariable(ctx context.Context, arg persister.UpdateVariableParams) error
	UpsertCategory(ctx context.Context, arg persister.UpsertCategoryParams) (persister.Category, error)
	UpsertCodebase(ctx context.Context, arg persister.UpsertCodebaseParams) (persister.Codebase, error)
	UpsertFile(ctx context.Context, arg persister.UpsertFileParams) (persister.File, error)
	UpsertImport(ctx context.Context, arg persister.UpsertImportParams) (persister.Import, error)
	UpsertModule(ctx context.Context, name string) (persister.Module, error)
	UpsertModuleVersion(ctx context.Context, arg persister.UpsertModuleVersionParams) (persister.ModuleVersion, error)
	UpsertPackage(ctx context.Context, arg persister.UpsertPackageParams) (persister.Package, error)
	UpsertPackageType(ctx context.Context, arg persister.UpsertPackageTypeParams) (persister.PackageType, error)
	UpsertPackageVersion(ctx context.Context, arg persister.UpsertPackageVersionParams) (persister.PackageVersion, error)
	UpsertProject(ctx context.Context, arg persister.UpsertProjectParams) (persister.Project, error)
	UpsertSurveyModule(ctx context.Context, arg persister.UpsertSurveyModuleParams) (persister.SurveyModule, error)
	UpsertSymbolType(ctx context.Context, arg persister.UpsertSymbolTypeParams) (persister.SymbolType, error)
}
