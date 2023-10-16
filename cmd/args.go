package main

import (
	"context"
	"fmt"
	"strings"

	"gerardus/cli"
	"gerardus/options"
	"gerardus/persister"
)

type checker struct{}

var check = checker{}

var projectArg = &cli.Arg{
	Name:             ProjectArg,
	Usage:            "Project name, e.g. 'golang'",
	CheckFunc:        check.projectName,
	SetStringValFunc: options.SetProjectName,
}

func (checker) projectName(project any) (err error) {
	projName := project.(string)
	ds := persister.GetDataStore()
	ctx := context.Background()
	_, err = ds.LoadProjectByName(ctx, projName)
	if err != nil {
		err = fmt.Errorf("project '%s' has not been added; %w", project, err)
		goto end
	}
end:
	return err
}

var versionTagArg = &cli.Arg{
	Name:             VersionTagArg,
	Usage:            "Git version tag",
	CheckFunc:        check.versionTag,
	SetStringValFunc: options.SetVersionTag,
}

func (checker) versionTag(tag any) (err error) {
	var ds *persister.DataStore
	var ctx context.Context
	var verTag string

	projName := options.ProjectName()
	if len(projName) == 0 {
		err = fmt.Errorf("no project has been specified. Use the `-prj=<project_name>` switch to specify the option")
		goto end
	}
	verTag = tag.(string)
	ds = persister.GetDataStore()
	ctx = context.Background()
	_, err = ds.LoadCodebaseByProjectNameAndVersionTag(ctx, persister.LoadCodebaseByProjectNameAndVersionTagParams{
		Name:       projName,
		VersionTag: verTag,
	})
	if err != nil {
		err = fmt.Errorf("codebase for project '%s' and version tag '%s' has not been added; %w", projName, verTag, err)
		goto end
	}
end:
	return err
}

var repoURLArg = &cli.Arg{
	Name:             RepoURLArg,
	Usage:            "The full GitHub repository URL for the project, e.g. https://github.com/golang/go",
	CheckFunc:        check.repoURL,
	SetStringValFunc: options.SetRepoURL,
}

func (checker) repoURL(url any) (err error) {
	repoURL := url.(string)
	parts := strings.Split(strings.TrimRight(repoURL, "/"), "/")
	numParts := len(parts)
	if numParts != 5 {
		err = fmt.Errorf("repo URL %s not a valid Github repo URL", repoURL)
		goto end
	}
	if strings.Join(parts[:3], "/") != "https://github.com" {
		err = fmt.Errorf("repo URL %s not a https://github.com URL", repoURL)
		goto end
	}
	err = cli.CheckURL(repoURL)
	if err != nil {
		goto end
	}
end:
	return err
}
