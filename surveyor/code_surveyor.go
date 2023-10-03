package surveyor

import (
	"context"

	"gerardus/collector"
	"gerardus/parser"
	"gerardus/persister"
	"golang.org/x/sync/errgroup"
)

var _ persister.SurveyAttrs = (*CodeSurveyor)(nil)

type CodeSurveyor struct {
	Codebase  *parser.Codebase
	Files     parser.Files
	localDir  string
	facetChan chan collector.CodeFacet
}

func (cs *CodeSurveyor) RepoURL() string {
	return cs.Codebase.RepoURL
}

func (cs *CodeSurveyor) LocalDir() string {
	return cs.localDir
}

func NewCodeSurveyor(cb *parser.Codebase, files parser.Files, dir string) *CodeSurveyor {
	return &CodeSurveyor{
		Codebase: cb,
		Files:    files,
		localDir: dir,
	}
}

func (cs *CodeSurveyor) Survey(ctx context.Context, facetChan chan collector.CodeFacet) (err error) {
	var group *errgroup.Group
	cs.facetChan = facetChan
	defer close(cs.facetChan)

	group, ctx = errgroup.WithContext(ctx)
	for _, f := range cs.Files {
		if err != nil {
			goto end
		}
		switch tf := f.(type) {
		case *parser.ModFile:
			err = cs.SurveyModFile(ctx, tf)

		case *parser.GoFile:
			err = cs.SurveyGoFile(ctx, tf, group)
		}
	}
	err = group.Wait()
end:
	return err
}

//goland:noinspection GoUnusedParameter
func (cs *CodeSurveyor) SurveyModFile(ctx context.Context, mf *parser.ModFile) (err error) {
	return err
}

func (cs *CodeSurveyor) SurveyGoFile(ctx context.Context, gf *parser.GoFile, group *errgroup.Group) (err error) {
	c := collector.New(gf, cs.facetChan)
	group.Go(func() (err error) {
		return c.CollectFiles(ctx)
	})
	return err
}
