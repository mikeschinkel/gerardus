package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"regexp"

	"github.com/mikeschinkel/gerardus/channels"
	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/collector"
	"github.com/mikeschinkel/gerardus/options"
	"github.com/mikeschinkel/gerardus/parser"
	"github.com/mikeschinkel/gerardus/persister"
	"github.com/mikeschinkel/gerardus/scanner"
	"github.com/mikeschinkel/gerardus/surveyor"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdMap = cli.AddCommandWithFunc("map", Root.ExecMap).
	AddArg(projectArg.MustExist()).
	AddArg(versionTagArg.MustExist()).
	AddFlag(&cli.Flag{
		Switch: "src",
		Arg: &cli.Arg{
			Name:         "source_dir",
			Usage:        "Source directory",
			Type:         reflect.String,
			Default:      defaultSourceDir(opts),
			CheckFunc:    checkDir,
			SetValueFunc: options.SetSourceDir,
		},
	})

//goland:noinspection GoUnusedParameter
func (a *App) ExecMap(ctx context.Context, i *cli.CommandInvoker) (err error) {
	var ma mapArgs
	var cs *surveyor.CodeSurveyor
	var cb *parser.Codebase
	var p *parser.Project
	var dir string

	fmt.Printf("Scanning Go source at %s...\n", options.SourceDir())

	project := i.ArgString(ProjectArg)
	versionTag := i.ArgString(VersionTagArg)

	cb = parser.NewCodebase(project, versionTag)
	p = parser.NewProject(project, Check.project.RepoUrl)
	cs = surveyor.NewCodeSurveyor(cb, p, options.SourceDir())
	dir = options.SourceDir()
	ma = mapArgs{
		scanner:   scanner.NewScanner(dir),
		parser:    parser.NewGoFileParser(cs.ModuleGraph(), dir),
		surveyor:  cs,
		persister: persister.NewSurveyPersister(cs, a.dataStore),
	}

	//err = mapWithSlices(ctx,ma)
	err = mapWithChans(ctx, ma)
	if err != nil {
		err = ErrMapCommandFailed.Err(err, "source_dir", options.SourceDir())
		goto end
	}
	fmt.Println("Scanning completed successfully.")
end:
	return err
}

type mapArgs struct {
	scanner   *scanner.Scanner
	parser    *parser.GoFileParser
	surveyor  *surveyor.CodeSurveyor
	persister *persister.SurveyPersister
}

var (
	modFiles = regexp.MustCompile(`^.*/?go\.mod`)
	goFiles  = regexp.MustCompile(`\.go$`)
)

func mapWithChans(ctx context.Context, args mapArgs) (err error) {

	slog.Info("Mapping project files")

	for _, fileType := range []*regexp.Regexp{modFiles, goFiles} {

		slog.Info("Mapping files", "file_type", fileType.String())

		scanFilesChan := make(chan scanner.File, 10)
		parseFilesChan := make(chan scanner.File, 10)
		facetChan := make(chan collector.CodeFacet, 10)

		// Process all the files of fileType
		pipeline := channels.NewPipeline(ctx)
		pipeline.AddStage(func() error { return args.scanner.ScanChan(ctx, fileType, scanFilesChan) })
		pipeline.AddStage(func() error { return args.parser.ParseChan(ctx, scanFilesChan, parseFilesChan) })
		pipeline.AddStage(func() error { return args.surveyor.SurveyChan(ctx, parseFilesChan, facetChan) })
		pipeline.AddStage(func() error { return args.persister.PersistChan(ctx, facetChan) })
		err = pipeline.Go()
		if err != nil {
			goto end
		}
	}

end:
	return err
}

// checkDir validates source directory
func checkDir(ctx Context, requires cli.ArgRequires, dir any) (err error) {
	var info os.FileInfo
	var absDir string

	sDir := dir.(string)
	if len(sDir) == 0 {
		err = ErrDirIsEmpty
		goto end
	}

	absDir, err = makeAbs(sDir)
	if err != nil {
		goto end
	}

	//goland:noinspection GoSwitchMissingCasesForIotaConsts
	switch cli.Existence(requires) {
	case cli.MustExist:
		info, err = os.Stat(absDir)
		if err != nil {
			err = ErrReadingSourceDir.Err(err, "source_dir", absDir)
			goto end
		}
		if !info.IsDir() {
			err = ErrPathNotADir.Err(err, "source_dir", absDir)
			goto end
		}
		dir = absDir // TODO Verify this actually sets the passed parameter
	case cli.OkToExist:
	case cli.MustNotExist:
		panic("Need to implement")
	}

end:
	return err
}

//func mapWithSlices(ctx context.Context, args mapArgs) (err error) {
//	files, err := args.scanner.Scan(ctx)
//	files, err = args.parser.Parse(ctx, files)
//	files, err = args.surveyor.Survey(ctx, files)
//	err = args.persister.Persist(ctx, files)
//	return err
//}
