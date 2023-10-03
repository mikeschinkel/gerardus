package persister

import (
	"gerardus/collector"
)

type SQLGenerator interface {
	SQL() []string
}

var _ SQLGenerator = (*FuncDecl)(nil)
var _ SQLGenerator = (*ImportSpec)(nil)
var _ SQLGenerator = (*TypeSpec)(nil)
var _ SQLGenerator = (*ValueSpec)(nil)

type FuncDecl struct {
	FuncDecl collector.FuncDecl
}

func (f FuncDecl) SQL() []string {
	//TODO implement me
	panic("implement me")
}

type ImportSpec struct {
	ImportSpec collector.ImportSpec
}

func (i ImportSpec) SQL() []string {
	//TODO implement me
	panic("implement me")
}

type TypeSpec struct {
	TypeSpec collector.TypeSpec
}

func (t TypeSpec) SQL() []string {
	//TODO implement me
	panic("implement me")
}

type ValueSpec struct {
	ValueSpec collector.ValueSpec
}

func (v ValueSpec) SQL() []string {
	//TODO implement me
	panic("implement me")
}
