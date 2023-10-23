package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"gerardus/channels"
	"gerardus/cli"
	"gerardus/collector"
	"gerardus/options"
	"gerardus/parser"
	"gerardus/persister"
	"gerardus/scanner"
	"gerardus/surveyor"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdMap = cli.AddCommandWithFunc("map", ExecMap).
	AddArg(projectArg.MustExist()).
	AddArg(versionTagArg.MustExist()).
	AddFlag(cli.Flag{
		Switch: "src",
		Arg: cli.Arg{
			Name:             "source_dir",
			Usage:            "Source directory",
			Default:          defaultSourceDir(opts),
			CheckFunc:        checkDir,
			SetStringValFunc: options.SetSourceDir,
		},
	})

//goland:noinspection GoUnusedParameter
func ExecMap(args cli.ArgsMap) (err error) {
	var ma mapArgs
	var cs *surveyor.CodeSurveyor
	var cb *parser.Codebase
	var p *parser.Project
	var ctx context.Context

	fmt.Printf("Scanning Go source at %s...\n", options.SourceDir())

	project := options.ProjectName()
	versionTag := options.VersionTag()

	cb = parser.NewCodebase(project, versionTag)
	p = parser.NewProject(project, check.project.RepoUrl)
	cs = surveyor.NewCodeSurveyor(cb, p, options.SourceDir())
	ma = mapArgs{
		scanner:   scanner.NewScanner(options.SourceDir()),
		parser:    parser.NewGoFileParser(),
		surveyor:  cs,
		persister: persister.NewSurveyPersister(cs, persister.GetDataStore()),
	}

	ctx = context.Background()
	//err = mapWithSlices(ctx,ma)
	err = mapWithChans(ctx, ma)
	if err != nil {
		err = errMapCommandFailed.Err(err, "source_dir", options.SourceDir())
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
	modFiles = regexp.MustCompile(`^go\.mod$`)
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
func checkDir(mode cli.ArgCheckMode, dir any) (err error) {
	var info os.FileInfo
	var absDir string

	sDir := dir.(string)
	if len(sDir) == 0 {
		err = errDirIsEmpty
		goto end
	}

	absDir, err = makeAbs(sDir)
	if err != nil {
		goto end
	}

	switch mode {
	case cli.MustExist:
		info, err = os.Stat(absDir)
		if err != nil {
			err = errReadingSourceDir.Err(err, "source_dir", absDir)
			goto end
		}
		if !info.IsDir() {
			err = errPathNotADir.Err(err, "source_dir", absDir)
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
