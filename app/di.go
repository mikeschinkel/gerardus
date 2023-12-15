package app

import (
	"context"

	"github.com/mikeschinkel/gerardus/persister"
)

const DIKey = "di"

type DI struct {
	RepoInfoRequesterFunc func(string) (*persister.RepoInfo, error)
	UpsertProjectFunc     func(context.Context, persister.UpsertProjectParams) (persister.Project, error)
	CheckURLFunc          func(string) error
}

func (di *DI) Assign(new DI) *DI {
	if di.RepoInfoRequesterFunc == nil {
		di.RepoInfoRequesterFunc = new.RepoInfoRequesterFunc
	}
	if di.UpsertProjectFunc == nil {
		di.UpsertProjectFunc = new.UpsertProjectFunc
	}
	if di.CheckURLFunc == nil {
		di.CheckURLFunc = new.CheckURLFunc
	}
	return di
}
