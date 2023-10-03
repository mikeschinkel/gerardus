package collector

import (
	"context"
	"fmt"
	"go/ast"
)

type FuncDecl struct {
	File     File
	Name     string
	Receiver *Field
	FuncType
}

func (fd FuncDecl) CodeFacet() {}

func (fd FuncDecl) String() (s string) {
	var r string
	if fd.Receiver != nil {
		r = fmt.Sprintf("(%s)", fd.Receiver.String())
	}
	return fmt.Sprintf("func %s%s%s {...}", r, fd.Name, fd.FuncType)
}

func (c *Collector) CollectFuncDecl(ctx context.Context, fd *ast.FuncDecl) (err error) {
	var r *Field

	f := FuncDecl{
		File: c.File,
		Name: fd.Name.Name,
		FuncType: FuncType{
			Parameters: make(FieldList, 0),
			Results:    make(FieldList, 0),
		},
	}
	r, _, err = c.CollectReceiver(ctx, fd.Recv)
	if err != nil {
		goto end
	}
	f.Receiver = r
	f.FuncType, err = c.CollectFuncType(ctx, fd.Type)
	if err != nil {
		goto end
	}
	c.FacetChan <- f
end:
	return err
}

func (c *Collector) CollectReceiver(ctx context.Context, list *ast.FieldList) (fld *Field, cnt int, err error) {
	var fl FieldList
	if !c.astHasFields(list) {
		goto end
	}
	cnt = len(list.List)
	if cnt > 1 {
		err = fmt.Errorf("unexpected: func has more than one receiver: %#v", list)
		goto end
	}
	fl, err = c.CollectField(ctx, list.List[0])
	if err != nil {
		goto end
	}
	if len(fl) > 0 {
		fld = &Field{}
		*fld = fl[0]
	}
end:
	return fld, cnt, err
}
