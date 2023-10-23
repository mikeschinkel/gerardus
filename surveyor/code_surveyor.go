package surveyor

import (
	"context"
	"log/slog"
	"sync"

	"gerardus/channels"
	"gerardus/collector"
	"gerardus/parser"
	"gerardus/scanner"
	"golang.org/x/mod/modfile"

	"golang.org/x/sync/errgroup"
)

type CodeSurveyor struct {
	Codebase  *parser.Codebase
	Project   *parser.Project
	Files     scanner.Files
	localDir  string
	source    string
	facetChan chan collector.CodeFacet
}

func (cs *CodeSurveyor) ProjectName() string {
	return cs.Codebase.Project
}

func (cs *CodeSurveyor) VersionTag() string {
	return cs.Codebase.VersionTag
}

func (cs *CodeSurveyor) LocalDir() string {
	return cs.localDir
}

func (cs *CodeSurveyor) Source() string {
	if len(cs.source) > 0 {
		goto end
	}
	if cs.VersionTag() == "." {
		cs.source = cs.localDir
		goto end
	}
	cs.source = cs.Project.RepoURL
end:
	return cs.source
}

type Project interface {
	RepoURL() string
}

func NewCodeSurveyor(cb *parser.Codebase, p *parser.Project, dir string) *CodeSurveyor {
	return &CodeSurveyor{
		Codebase: cb,
		Project:  p,
		localDir: dir,
	}
}

func (cs *CodeSurveyor) Survey(ctx context.Context, files scanner.Files) (outFiles scanner.Files, err error) {
	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)
	for _, f := range files {
		err = cs.SurveyFile(ctx, f, group)
		if err != nil {
			goto end
		}
	}
	err = group.Wait()
end:
	return cs.Files, err
}

func (cs *CodeSurveyor) SurveyChan(ctx context.Context, filesChan chan scanner.File, facetChan chan collector.CodeFacet) (err error) {
	var group *errgroup.Group
	var cancel context.CancelFunc

	cs.facetChan = facetChan
	defer close(facetChan)
	group, ctx = errgroup.WithContext(ctx)
	ctx, cancel = context.WithCancel(ctx)
	err = channels.ReadAllFrom(ctx, filesChan, func(f scanner.File) (err error) {
		slog.Info("Surveying file", "filepath", f.RelPath())
		err = cs.SurveyFile(ctx, f, group)
		if err != nil {
			cancel()
		}
		return err
	})
	if err != nil {
		goto end
	}
	err = group.Wait()
end:
	return err
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

var mutex sync.Mutex

//goland:noinspection GoUnusedParameter
func (cs *CodeSurveyor) SurveyModFile(ctx context.Context, mf *parser.ModFile) (err error) {
	null := struct{}{}

	mf.ModFile, err = modfile.Parse("go.mod", mf.Content, nil)
	if err != nil {
		err = errFailedToParseFile.Err(err, "filename", mf.Fullpath())
		goto end
	}

	// Make Modules available w/o having to look up via database to speed insert of imports
	mutex.Lock()
	parser.Modules[mf.Name()] = null
	for _, r := range mf.Require() {
		parser.Modules[r.Mod.Path] = null
	}
	mutex.Unlock()

	err = channels.WriteTo(ctx, cs.facetChan, collector.CodeFacet(mf))
	if err != nil {
		goto end
	}

end:
	return err
}

func (cs *CodeSurveyor) SurveyGoFile(ctx context.Context, gf *parser.GoFile, group *errgroup.Group) (err error) {
	// TODO Make this work with Survey() in addition to SurveyChan().
	c := collector.New(gf, cs.facetChan)
	group.Go(func() (err error) {
		return c.CollectFacets(ctx)
	})
	return err
}
