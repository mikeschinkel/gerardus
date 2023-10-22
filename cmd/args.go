package main

import (
	"context"
	"strings"

	"gerardus/cli"
	"gerardus/options"
	"gerardus/persister"
	"gerardus/serr"
)

type checker struct{}

var check = checker{}

var projectArg = &cli.Arg{
	Name:             ProjectArg,
	Usage:            "Project name, e.g. 'golang'",
	CheckFunc:        check.projectName,
	SetStringValFunc: options.SetProjectName,
}

func (checker) projectName(mode cli.ArgCheckMode, project any) (err error) {
	projName := project.(string)
	ds := persister.GetDataStore()
	ctx := context.Background()
	switch mode {
	case cli.MustExist:
		_, err = ds.LoadProjectByName(ctx, projName)
		if err != nil {
			err = errProjectNotFound.Err(err, "project", project)
			goto end
		}
	case cli.OkToExist:
	case cli.MustNotExist:
		panic("Need to implement")
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

func (checker) versionTag(mode cli.ArgCheckMode, tag any) (err error) {
	var ds *persister.DataStore
	var ctx context.Context
	var verTag string

	projName := options.ProjectName()
	if len(projName) == 0 {
		err = errNoProjectSpecified
		goto end
	}
	verTag = tag.(string)
	if len(verTag) == 0 {
		err = errNoVersionTagSpecified.Args("project", projName)
		goto end
	}
	ds = persister.GetDataStore()
	ctx = context.Background()

	switch mode {
	case cli.MustExist:
		_, err = ds.LoadCodebaseByProjectNameAndVersionTag(ctx, persister.LoadCodebaseByProjectNameAndVersionTagParams{
			Name:       projName,
			VersionTag: verTag,
		})
		if err != nil {
			err = errFailedToAddCodebase.Err(err, "project", projName, "version_tag", verTag)
			goto end
		}
	case cli.OkToExist:
	case cli.MustNotExist:
		panic("Need to implement")
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

func (checker) repoURL(mode cli.ArgCheckMode, url any) (err error) {
	var parts []string
	var numParts int
	repoURL := url.(string)
	if len(repoURL) == 0 {
		err = errNoRepoURLSpecified
		goto end
	}
	parts = strings.Split(strings.TrimRight(repoURL, "/"), "/")
	numParts = len(parts)
	if numParts != 5 {
		err = errInvalidGitHubRepoURL
		goto end
	}
	if strings.Join(parts[:3], "/") != "https://github.com" {
		err = errInvalidGitHubRepoRootURL
		goto end
	}
	switch mode {
	case cli.MustExist:
		err = cli.CheckURL(repoURL)
		if err != nil {
			err = errURLCouldNotBeDereferenced
		}
	case cli.OkToExist:
	case cli.MustNotExist:
		panic("Need to implement")
	}
end:
	if err != nil && len(repoURL) > 0 {
		err = err.(serr.SError).Args("repo_url", repoURL)
	}
	return err
}
