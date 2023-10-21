package persister

import (
	"context"
	"log/slog"

	"gerardus/collector"
	"gerardus/scanner"
	"golang.org/x/sync/errgroup"
)

type SurveyAttrs interface {
	RepoURL() string
	LocalDir() string
}

type SurveyPersister struct {
	survey    SurveyAttrs
	surveyId  int64
	fileId    int64
	filepath  string
	dataStore *DataStore
}

func NewSurveyPersister(survey SurveyAttrs, ds *DataStore) *SurveyPersister {
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

	codebaseID, err = ds.LoadCodebaseIdByRepoURL(ctx, sp.survey.RepoURL())
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

func (sp *SurveyPersister) persistFacet(ctx context.Context, facetChan chan collector.CodeFacet) (err error) {
	var group *errgroup.Group

	group, ctx = errgroup.WithContext(ctx)
	group.Go(func() (err error) {
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err() // Return the error to terminate this goroutine
				goto end
			case facet, ok := <-facetChan:
				if !ok {
					// Channel is closed, so we're done
					goto end
				}
				switch ft := facet.(type) {
				case collector.FuncDecl:
					print()
				case collector.ImportSpec:
					err = sp.importSpecInsertFunc(ctx, ft)
				case collector.TypeSpec:
					err = sp.typeSpecInsertFunc(ctx, ft)
				case collector.ValueSpec:
					print()
				default:
					panicf("Unhandled CodeFacet type '%T'", ft)
				}
				if err != nil {
					goto end
				}
			}
		}
	end:
		return err
	})
	err = group.Wait()
	if err != nil {
		debugBreakpointHere()
	}
	return err
}

func (sp *SurveyPersister) typeSpecInsertFunc(ctx context.Context, ts collector.TypeSpec) (err error) {
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

func (sp *SurveyPersister) importSpecInsertFunc(ctx context.Context, is collector.ImportSpec) (err error) {
	var fileId int64
	fileId, err = sp.getFileId(ctx, is.File)
	if err != nil {
		goto end
	}
	_, err = sp.dataStore.UpsertImport(ctx, UpsertImportParams{
		FileID:    fileId,
		SurveyID:  sp.surveyId,
		PackageID: "",
		Alias:     is.Alias,
	})
	if err != nil {
		goto end
	}
end:
	return err
}

func (sp *SurveyPersister) getFileId(ctx context.Context, f collector.File) (fid int64, err error) {
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
