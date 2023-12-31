package persister

import (
	"context"
	"log/slog"

	"github.com/mikeschinkel/gerardus/channels"
	"github.com/mikeschinkel/gerardus/collector"
	"github.com/mikeschinkel/gerardus/parser"
	"github.com/mikeschinkel/gerardus/scanner"
	"github.com/mikeschinkel/gerardus/surveyor"
	"golang.org/x/sync/errgroup"
)

var _ survey = (*surveyor.CodeSurveyor)(nil)

type survey interface {
	ProjectName() string
	VersionTag() string
	LocalDir() string
	Source() string
	ModuleGraph() *parser.ModuleGraph
}

type SurveyPersister struct {
	survey    survey
	surveyId  int64
	fileId    int64
	filepath  string
	dataStore DataStore
}

func NewSurveyPersister(survey survey, ds DataStore) *SurveyPersister {
	return &SurveyPersister{
		survey:    survey,
		dataStore: ds,
	}
}

//goland:noinspection GoUnusedParameter
func (sp *SurveyPersister) Persist(ctx context.Context, fs scanner.Files) (err error) {
	return nil // TODO Make this work with a slice of files
}

func (sp *SurveyPersister) PersistChan(ctx context.Context, facetChan chan collector.CodeFacet) (err error) {
	var group *errgroup.Group
	var codebaseID int64
	var survey Survey
	ds := sp.dataStore

	codebaseID, err = ds.Queries().LoadCodebaseIDByProjectAndVersion(ctx, LoadCodebaseIDByProjectAndVersionParams{
		Name:       sp.survey.ProjectName(),
		VersionTag: sp.survey.VersionTag(),
	})
	if err != nil {
		goto end
	}
	survey, err = ds.Queries().InsertSurvey(ctx, InsertSurveyParams{
		CodebaseID: codebaseID,
		LocalDir:   sp.survey.LocalDir(),
	})
	if err != nil {
		goto end
	}
	sp.surveyId = survey.ID
	group, ctx = errgroup.WithContext(ctx)
	group.Go(func() (err error) {
		err = sp.persistFacetChan(ctx, facetChan)
		if err != nil {
			err = ErrFailedWhilePersisting.Err(err)
		}
		return
	})
	err = group.Wait()
end:
	return err
}

func (sp *SurveyPersister) persistFacetChan(ctx context.Context, facetChan chan collector.CodeFacet) (err error) {
	var group *errgroup.Group
	var cancel context.CancelFunc
	var insert = func(typ string, f collector.CodeFacet, insert func(context.Context, collector.CodeFacet) error) (err error) {
		args := []any{"spec_type", typ, "spec", f.String()}
		slog.Info("Inserting", args...)
		err = insert(ctx, f)
		if err != nil {
			err = ErrFailedToInsertSpec.Err(err, args...)
		}
		return err
	}
	group, ctx = errgroup.WithContext(ctx)
	ctx, cancel = context.WithCancel(ctx)
	group.Go(func() (err error) {
		return channels.ReadAllFrom(ctx, facetChan, func(facet collector.CodeFacet) error {
			switch ft := facet.(type) {
			case collector.ImportSpec:
				err = insert("import", ft, func(ctx context.Context, facet collector.CodeFacet) error {
					return sp.insertImportSpec(ctx, ft)
				})
			case collector.TypeSpec:
				err = insert("type", ft, func(ctx context.Context, facet collector.CodeFacet) error {
					return sp.insertTypeSpec(ctx, ft)
				})
			case *parser.ModFile:
				err = insert("mod_file", ft, func(ctx context.Context, facet collector.CodeFacet) error {
					return sp.insertModFile(ctx, ft)
				})
			case collector.ValueSpec:
				print()
			case collector.FuncDecl:
				print()
			default:
				panicf("Unhandled CodeFacet type '%T'", ft)
			}
			if err != nil {
				cancel()
			}
			return err
		})
	})
	err = group.Wait()
	if err != nil {
		debugBreakpointHere()
	}
	return err
}

func (sp *SurveyPersister) insertModFile(ctx context.Context, mf *parser.ModFile) (err error) {
	var m Module
	var mv ModuleVersion
	var fileId int64
	var pkg Package

	fileId, err = sp.getFileId(ctx, mf)
	if err != nil {
		goto end
	}
	for _, module := range mf.Modules() {
		//if i == 0 {
		//	// Get the source for the go.mod file
		//	path = mf.Fullpath()
		//} else {
		//	// Get the source for the go.mod file's dependencies
		//}
		pkg, err = sp.upsertPackage(ctx, module.Package)
		if err != nil {
			goto end
		}
		m, mv, err = sp.upsertModule(ctx, module)
		if err != nil {
			goto end
		}
		_, err = sp.dataStore.Queries().UpsertSurveyModule(ctx, UpsertSurveyModuleParams{
			SurveyID:        sp.surveyId,
			ModuleID:        m.ID,
			ModuleVersionID: mv.ID,
			FileID:          fileId,
			PackageID:       pkg.ID,
		})
		if err != nil {
			goto end
		}
	}
end:
	return err
}

func (sp *SurveyPersister) insertTypeSpec(ctx context.Context, ts collector.TypeSpec) (err error) {
	var fileId int64

	fileId, err = sp.getFileId(ctx, ts.File)
	if err != nil {
		goto end
	}
	_, err = sp.dataStore.Queries().InsertType(ctx, InsertTypeParams{
		FileID:       fileId,
		SymbolTypeID: int64(ts.SymbolType),
		Name:         ts.Name,
		Definition:   ts.Definition.String(),
	})
	if err != nil {
		goto end
	}
end:
	return err
}

func (sp *SurveyPersister) upsertPackage(ctx context.Context, pp *parser.Package) (pkg Package, err error) {
	pkg, err = sp.dataStore.Queries().UpsertPackage(ctx, UpsertPackageParams{
		ImportPath: pp.ImportPath,
		Source:     pp.Source(),
		TypeID:     int64(pp.Type),
	})
	if err != nil {
		goto end
	}
	if pp.PackageVersion.Name == "." {
		goto end
	}
	_, err = sp.dataStore.Queries().UpsertPackageVersion(ctx, UpsertPackageVersionParams{
		PackageID: pkg.ID,
		Version:   pp.PackageVersion.Name,
		SourceUrl: pp.PackageVersion.Source(),
	})
	if err != nil {
		goto end
	}
end:
	return pkg, err
}

func (sp *SurveyPersister) upsertModule(ctx context.Context, pm *parser.Module) (m Module, mv ModuleVersion, err error) {
	m, err = sp.dataStore.Queries().UpsertModule(ctx, pm.Name())
	if err != nil {
		goto end
	}
	mv, err = sp.dataStore.Queries().UpsertModuleVersion(ctx, UpsertModuleVersionParams{
		ModuleID: m.ID,
		Version:  pm.VersionName(),
	})
	if err != nil {
		goto end
	}
end:
	return m, mv, err
}

func (sp *SurveyPersister) insertImportSpec(ctx context.Context, is collector.ImportSpec) (err error) {
	var fileId int64
	var pkg Package
	var pp *parser.Package

	fileId, err = sp.getFileId(ctx, is.File)
	if err != nil {
		goto end
	}

	pp = sp.survey.ModuleGraph().DispensePackage(is.Package, is.File.Fullpath())
	if pp == nil {
		goto end
	}

	pkg, err = sp.upsertPackage(ctx, pp)
	if err != nil {
		goto end
	}
	_, err = sp.dataStore.Queries().UpsertImport(ctx, UpsertImportParams{
		FileID:    fileId,
		SurveyID:  sp.surveyId,
		PackageID: pkg.ID,
		Alias:     is.Alias,
	})
	if err != nil {
		goto end
	}
end:
	return err
}

type relPathGetter interface {
	RelPath() string
}

func (sp *SurveyPersister) getFileId(ctx context.Context, f relPathGetter) (fid int64, err error) {
	var mf File

	if f.RelPath() == sp.filepath {
		fid = sp.fileId
		goto end
	}
	mf, err = sp.dataStore.Queries().UpsertFile(ctx, UpsertFileParams{
		SurveyID: sp.surveyId,
		Filepath: f.RelPath(),
	})
	if err != nil {
		goto end
	}
	fid = mf.ID
end:
	return fid, err
}
