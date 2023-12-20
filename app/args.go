package app

import (
	"reflect"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/options"
)

var projectArg = cli.Arg{
	Name:         ProjectArg,
	Usage:        "Project name, e.g. 'golang'",
	Type:         reflect.String,
	CheckFunc:    Root.checkProjectName,
	SetValueFunc: options.SetProjectName,
}

var versionTagArg = &cli.Arg{
	Name:         VersionTagArg,
	Usage:        "Git version tag",
	Type:         reflect.String,
	CheckFunc:    Root.checkVersionTag,
	SetValueFunc: options.SetVersionTag,
}

var repoURLArg = &cli.Arg{
	Name:         RepoURLArg,
	Usage:        "The full GitHub repository URL for the project, e.g. https://github.com/golang/go",
	Type:         reflect.String,
	CheckFunc:    Root.checkRepoURL,
	SetValueFunc: options.SetRepoURL,
}
