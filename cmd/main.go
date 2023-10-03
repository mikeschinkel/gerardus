package main

import (
	"context"
	"fmt"

	"gerardus/collector"
	"gerardus/options"
	"gerardus/parser"
	"gerardus/persister"
	"gerardus/scanner"
	"gerardus/surveyor"
	"golang.org/x/sync/errgroup"
)

const GoStdLibURL = "https://github.com/golang/go/tree/master/src"

func main() {

	err := options.InitOptions()
	if err != nil {
		usage("Failed to initialize; %s", err.Error())
	}

	err = persister.Initialize(context.Background(),
		options.GetDataFile(),
		collector.SymbolTypes,
	)
	if err != nil {
		usage("Failed to initialize Symbol types; %s", err.Error())
	}

	//printSymbolTypes()
	fmt.Printf("Scanning Go source at %s...\n", options.SourceDir)
	s := scanner.NewScanner(options.SourceDir)
	files, err := s.Scan()
	if err != nil {
		usage(err.Error())
	}

	group, ctx := errgroup.WithContext(context.Background())

	p := parser.NewParser()
	files, err = p.Parse(ctx, files)
	if err != nil {
		usage("Failed to parse source code at %s; %s", options.SourceDir, err.Error())
	}

	cb := parser.NewCodebase(GoStdLibURL)
	cs := surveyor.NewCodeSurveyor(cb, files, options.GetSourceDir())
	sp := persister.NewSurveyPersister(cs, persister.GetDataStore())
	facetChan := make(chan collector.CodeFacet, 10)
	group.Go(func() error {
		return sp.Persist(ctx, facetChan)
	})
	group.Go(func() error {
		return cs.Survey(ctx, facetChan)
	})
	err = group.Wait()
	if err != nil {
		usage("Failed to survey and/or persist survey results %s; %s", options.SourceDir, err.Error())
		usage(err.Error())
	} else {
		println("Done!")
	}

	//dumper := gerardus.NewDumper(files)
	//dumper.Dump()
	//return
	//
	//
	//cg := gerardus.NewCodeGenerator(gerardus.Options.OutputDir, cs)
	//err = cg.Generate()
	//if err != nil {
	//	ErrOut(err)
	//}
	//fmt.Printf("\nSUCCESS! Output written to %s", gerardus.Options.OutputDir)
}
