package surveyor

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"

	"github.com/mikeschinkel/gerardus/channels"
	"github.com/mikeschinkel/gerardus/collector"
	"github.com/mikeschinkel/gerardus/parser"
	"github.com/mikeschinkel/gerardus/scanner"
	"github.com/mikeschinkel/go-typegen"
	"golang.org/x/mod/modfile"
	"golang.org/x/sync/errgroup"
)

type CodeSurveyor struct {
	Codebase    *parser.Codebase
	Project     *parser.Project
	Files       scanner.Files
	localDir    string
	source      string
	facetChan   chan collector.CodeFacet
	moduleGraph *parser.ModuleGraph
}

func (cs *CodeSurveyor) ModuleGraph() *parser.ModuleGraph {
	return cs.moduleGraph
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
		Codebase:    cb,
		Project:     p,
		localDir:    dir,
		moduleGraph: parser.NewModuleGraph(),
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

//goland:noinspection GoUnusedParameter
func (cs *CodeSurveyor) SurveyModFile(ctx context.Context, pmf *parser.ModFile) (err error) {
	var pm *parser.Module
	var modFile *modfile.File

	modFile, err = modfile.Parse("go.mod", pmf.Content, nil)
	if err != nil {
		err = errFailedToParseFile.Err(err, "filename", pmf.Fullpath())
		goto end
	}
	pmf.SetModFile(modFile)
	// Make Modules available w/o having to look up via database to speed insert of imports
	pm = cs.moduleGraph.AddProjectModule(&parser.ModuleArgs{
		PackageType: parser.GoModPackage,
		Name:        pmf.Name(),
		Version:     pmf.Version(),
		GoVersion:   pmf.GoVersion(),
		Path:        filepath.Dir(pmf.Fullpath()),
	})
	for _, r := range pmf.Require() {
		mod := r.Mod
		cs.moduleGraph.AddDependentModule(pm, &parser.ModuleArgs{
			Name:      mod.Path,
			Version:   mod.Version,
			Path:      mod.Path,
			GoVersion: pmf.GoVersion(),
		})
	}
	err = channels.WriteTo(ctx, cs.facetChan, collector.CodeFacet(pmf))
	if err != nil {
		goto end
	}
	showData(pm)
end:
	return err
}
func showData(m *parser.Module) {
	subs := typegen.Substitutions{
		reflect.TypeOf(reflect.Value{}): func(rv *reflect.Value) string {
			return fmt.Sprintf("reflect.ValueOf(%v)", (*rv).Interface())
		},
	}
	nm := typegen.NewNodeMarshaler(subs)
	nodes := nm.Marshal(m)
	fmt.Println(typegen.NewCodeBuilder("getData", "", nodes))
}
func (cs *CodeSurveyor) SurveyGoFile(ctx context.Context, gf *parser.GoFile, group *errgroup.Group) (err error) {
	// TODO Make this work with Survey() in addition to SurveyChan().
	c := collector.New(gf, cs.facetChan)
	group.Go(func() (err error) {
		return c.CollectFacets(ctx)
	})
	return err
}
