package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"gerardus/cli"
	"gerardus/collector"
	"gerardus/options"
	"gerardus/parser"
	"gerardus/persister"
	"gerardus/scanner"
	"gerardus/surveyor"
	"golang.org/x/sync/errgroup"
)

const GoStdLibURL = "https://github.com/golang/go/tree/go1.21.1/src"

//goland:noinspection GoUnusedGlobalVariable
var CmdMap = cli.AddCommandWithFunc("map", ExecMap).
	AddArg(projectArg.MustExist()).
	AddArg(versionTagArg.MustExist()).
	AddFlag(cli.Flag{
		Switch: "src",
		Arg: cli.Arg{
			Name:             "source_dir",
			Usage:            "Source directory",
			Default:          defaultSourceDir(),
			CheckFunc:        checkDir,
			SetStringValFunc: options.SetSourceDir,
		},
	})

//goland:noinspection GoUnusedParameter
func ExecMap(args cli.ArgsMap) (err error) {
	var ma mapArgs
	var cs *surveyor.CodeSurveyor
	var cb *parser.Codebase
	var ctx context.Context

	fmt.Printf("Scanning Go source at %s...\n", options.SourceDir())

	project := options.ProjectName()
	versionTag := options.VersionTag()

	cb = parser.NewCodebase(project, versionTag)
	cs = surveyor.NewCodeSurveyor(cb, options.SourceDir())
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
		err = fmt.Errorf("failed for %s; %w", options.SourceDir(), err)
		goto end
	}
end:
	return err
}

type mapArgs struct {
	scanner   *scanner.Scanner
	parser    *parser.GoFileParser
	surveyor  *surveyor.CodeSurveyor
	persister *persister.SurveyPersister
}

func mapWithSlices(ctx context.Context, args mapArgs) (err error) {
	files, err := args.scanner.Scan()
	files, err = args.parser.Parse(ctx, files)
	files, err = args.surveyor.Survey(ctx, files)
	err = args.persister.Persist(ctx, files)
	return err
}

func mapWithChans(ctx context.Context, args mapArgs) (err error) {
	var group *errgroup.Group
	slog.Info("Mapping project files")

	group, ctx = errgroup.WithContext(ctx)

	scanFilesChan := make(chan scanner.File, 10)
	parseFilesChan := make(chan scanner.File, 10)
	facetChan := make(chan collector.CodeFacet, 10)

	funcs := []func() error{
		func() error {
			return args.scanner.ScanChan(scanFilesChan)
		},
		func() error {
			return args.parser.ParseChan(ctx, scanFilesChan, parseFilesChan)
		},
		func() error {
			return args.surveyor.SurveyChan(ctx, parseFilesChan, facetChan)
		},
		func() error {
			return args.persister.PersistChan(ctx, facetChan)
		},
	}
	for i := len(funcs) - 1; i >= 0; i-- {
		// Call in reverse order do the dowstream function will be ready before the
		// upstream function starts.
		group.Go(funcs[i])
	}
	err = group.Wait()
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
