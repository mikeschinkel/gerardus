package parser

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"

	"gerardus/channels"
	"gerardus/scanner"
)

type GoFileParser struct {
	Files       scanner.Files
	ModuleGraph *ModuleGraph
	rootDir     string
}

func NewGoFileParser(mg *ModuleGraph, rootDir string) *GoFileParser {
	return &GoFileParser{
		ModuleGraph: mg,
		rootDir:     rootDir,
	}
}

func (p *GoFileParser) ParseChan(ctx context.Context, inFilesChan, outFilesChan chan scanner.File) (err error) {
	defer close(outFilesChan)
	return channels.ReadAllFrom(ctx, inFilesChan, func(inFile scanner.File) (err error) {
		var outFile scanner.File

		slog.Info("Parsing file", "file", inFile.RelPath())
		outFile, err = p.parseFile(ctx, inFile)
		if err != nil {
			goto end
		}
		err = channels.WriteTo(ctx, outFilesChan, outFile)
		if err != nil {
			goto end
		}
	end:
		return err
	})
}

func (p *GoFileParser) Parse(ctx context.Context, files scanner.Files) (outFiles scanner.Files, err error) {
	var outFile scanner.File
	outFiles = make(scanner.Files, len(files))
	for i, inFile := range files {
		outFile, err = p.parseFile(ctx, inFile)
		if err != nil {
			goto end
		}
		outFiles[i] = outFile // Release memory
	}
end:
	return outFiles, err
}

func (p *GoFileParser) parseFile(ctx context.Context, f scanner.File) (_ scanner.File, err error) {
	switch f.Filename() {
	case "go.mod":
		f, err = p.parseModFile(ctx, f)
	default:
		f, err = p.parseGoFile(ctx, f)
	}
	return f, err
}

//goland:noinspection GoUnusedParameter
func (p *GoFileParser) parseModFile(ctx context.Context, file scanner.File) (mf *ModFile, err error) {
	var content []byte
	content, err = os.ReadFile(file.Fullpath())
	if err != nil {
		goto end
	}
	mf = NewModFile(file, content, p.ModuleGraph, p.rootDir)
end:
	return mf, err
}

//goland:noinspection GoUnusedParameter
func (p *GoFileParser) parseGoFile(ctx context.Context, file scanner.File) (gf *GoFile, err error) {
	var pkgName string
	var node *ast.File

	pkgName, err = loadPackageName(file)
	if err != nil {
		goto end
	}

	if pkgName == "." {
		goto end
	}

	gf = NewGoFile(file, pkgName)

	node, err = parser.ParseFile(
		token.NewFileSet(),
		file.Fullpath(), nil,
		parser.ParseComments,
	)
	if err != nil {
		goto end
	}
	gf.ast = node
end:
	return gf, err
}
