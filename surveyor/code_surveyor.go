package surveyor

import (
	"context"
	"log/slog"

	"gerardus/collector"
	"gerardus/parser"
	"gerardus/persister"
	"gerardus/scanner"
	"golang.org/x/sync/errgroup"
)

var _ persister.SurveyAttrs = (*CodeSurveyor)(nil)

type CodeSurveyor struct {
	Codebase  *parser.Codebase
	Files     scanner.Files
	localDir  string
	facetChan chan collector.CodeFacet
}

func (cs *CodeSurveyor) RepoURL() string {
	return cs.Codebase.RepoURL
}

func (cs *CodeSurveyor) LocalDir() string {
	return cs.localDir
}

func NewCodeSurveyor(cb *parser.Codebase, dir string) *CodeSurveyor {
	return &CodeSurveyor{
		Codebase: cb,
		localDir: dir,
	}
}

func (cs *CodeSurveyor) SurveyChan(ctx context.Context, filesChan chan scanner.File, facetChan chan collector.CodeFacet) (err error) {
	var group *errgroup.Group
	cs.facetChan = facetChan
	defer close(filesChan)

	group, ctx = errgroup.WithContext(ctx)
	for f := range filesChan {
		err = cs.SurveyFile(ctx, f, group)
		if err != nil {
			goto end
		}
	}
	err = group.Wait()
end:
	return err
}

func (cs *CodeSurveyor) Survey(ctx context.Context, files scanner.Files) (outFiles scanner.Files, err error) {
	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)
	for _, f := range files {
		slog.Info("Surveying file", "filepath", f.RelPath())
		err = cs.SurveyFile(ctx, f, group)
		if err != nil {
			goto end
		}
	}
	err = group.Wait()
end:
	return cs.Files, err
}

//goland:noinspection GoUnusedParameter
func (cs *CodeSurveyor) SurveyFile(ctx context.Context, f scanner.File, group *errgroup.Group) (err error) {
	switch tf := f.(type) {
	case *parser.ModFile:
		err = cs.SurveyModFile(ctx, tf)

	case *parser.GoFile:
		err = cs.SurveyGoFile(ctx, tf, group)
	}
	return err
}

//goland:noinspection GoUnusedParameter
func (cs *CodeSurveyor) SurveyModFile(ctx context.Context, mf *parser.ModFile) (err error) {
	return err
}

func (cs *CodeSurveyor) SurveyGoFile(ctx context.Context, gf *parser.GoFile, group *errgroup.Group) (err error) {
	// TODO Make this work with Survey() in addition to SurveyChan().
	c := collector.New(gf, cs.facetChan)
	group.Go(func() (err error) {
		return c.CollectFiles(ctx)
	})
	return err
}
