package persister

import (
	"context"
	"log/slog"

	"gerardus/channels"
	"gerardus/collector"
	"gerardus/scanner"
	"golang.org/x/sync/errgroup"
)

var _ survey = (*surveyor.CodeSurveyor)(nil)

type survey interface {
	ProjectName() string
	VersionTag() string
	LocalDir() string
	Source() string
}

type SurveyPersister struct {
	survey    survey
	surveyId  int64
	fileId    int64
	filepath  string
	dataStore *DataStore
}

func NewSurveyPersister(survey survey, ds *DataStore) *SurveyPersister {
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

	codebaseID, err = ds.LoadCodebaseByProjectNameAndVersionTag(ctx, LoadCodebaseByProjectNameAndVersionTagParams{
		Name:       sp.survey.ProjectName(),
		VersionTag: sp.survey.VersionTag(),
	})
	if err != nil {
		goto end
	}
	survey, err = ds.InsertSurvey(ctx, InsertSurveyParams{
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
			err = errFailedWhilePersisting.Err(err)
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
			err = errFailedToInsertSpec.Err(err, args...)
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

func (sp *SurveyPersister) insertTypeSpec(ctx context.Context, ts collector.TypeSpec) (err error) {
	var fileId int64

	fileId, err = sp.getFileId(ctx, ts.File)
	if err != nil {
		goto end
	}
	_, err = sp.dataStore.InsertType(ctx, InsertTypeParams{
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

func (sp *SurveyPersister) insertImportSpec(ctx context.Context, is collector.ImportSpec) (err error) {
	var fileId int64
	var p Package

	fileId, err = sp.getFileId(ctx, is.File)
	if err != nil {
		goto end
	}
	p, err = sp.dataStore.UpsertPackage(ctx, UpsertPackageParams{
		Path:   is.Package,
		Source: sp.survey.Source(),
	})
	if err != nil {
		goto end
	}
	_, err = sp.dataStore.UpsertImport(ctx, UpsertImportParams{
		FileID:    fileId,
		SurveyID:  sp.surveyId,
		PackageID: p.ID,
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
	mf, err = sp.dataStore.UpsertFile(ctx, UpsertFileParams{
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
