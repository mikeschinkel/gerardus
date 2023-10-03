package parser

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

const timeoutSec = 10

type Parser struct {
	Files       Files
	GoFilesChan chan *GoFile
}

func NewParser() *Parser {
	return &Parser{
		GoFilesChan: make(chan *GoFile, 10),
	}
}

func (p *Parser) Parse(ctx context.Context, files Files) (outFiles Files, err error) {
	var outFile File
	//var gf *GoFile
	//var ok bool

	defer close(p.GoFilesChan)
	outFiles = make(Files, len(files))
	for i, inFile := range files {
		outFile, err = p.parseFile(ctx, inFile)
		if err != nil {
			goto end
		}
		//gf, ok = outFile.(*GoFile)
		//if !ok {
		//	continue
		//}
		//select {
		//case p.GoFilesChan <- gf:
		//case <-ctx.Done():
		//	err = ctx.Err()
		//	goto end
		//case <-time.After(timeoutSec * time.Second): // Timeout after 2 seconds
		//	err = fmt.Errorf("timeout while sending")
		//	goto end
		//}
		outFiles[i] = outFile // Release memory
	}
end:
	return outFiles, err
}

func (p *Parser) parseFile(ctx context.Context, f File) (_ File, err error) {
	switch f.Filename() {
	case "go.mod":
		f, err = p.parseModFile(ctx, f)
	default:
		f, err = p.parseGoFile(ctx, f)
	}
	return f, err
}

//goland:noinspection GoUnusedParameter
func (p *Parser) parseModFile(ctx context.Context, file File) (mf *ModFile, err error) {
	var content []byte
	content, err = os.ReadFile(file.Fullpath())
	if err != nil {
		goto end
	}
	mf = NewModFile(file, content)
end:
	return mf, err
}

//goland:noinspection GoUnusedParameter
func (p *Parser) parseGoFile(ctx context.Context, file File) (gf *GoFile, err error) {
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
