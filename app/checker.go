package app

import (
	"github.com/mikeschinkel/gerardus/persister"
)

type checker struct {
	project *persister.Project
	App     *App
}

var Check = checker{}
