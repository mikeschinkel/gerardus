package app

import (
	"context"
	"reflect"
	"strings"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/options"
	"github.com/mikeschinkel/gerardus/persister"
	"github.com/mikeschinkel/go-serr"
)

type checker struct {
	project *persister.Project
}

var check = checker{}

var projectArg = &cli.Arg{
	Name:         ProjectArg,
	Usage:        "Project name, e.g. 'golang'",
	Type:         reflect.String,
	CheckFunc:    check.projectName,
	SetValueFunc: options.SetProjectName,
}

func (checker) projectName(requires cli.ArgRequires, project any) (err error) {
	var p persister.Project
	var ds *persister.DataStore
	var existence = cli.Existence(requires)
	var ctx context.Context

	projName := project.(string)
	if projName == "" && existence == cli.MustExist {
		err = ErrProjectNotFound
		goto end
	}
	ds = persister.GetDataStore()
	ctx = context.Background()
	//goland:noinspection GoSwitchMissingCasesForIotaConsts
	switch existence {
	case cli.MustExist:
		p, err = ds.LoadProjectByName(ctx, projName)
		if err != nil {
			err = ErrProjectNotFound.Err(err, "project", project)
			goto end
		}
		check.project = &p
	case cli.OkToExist:
	case cli.MustNotExist:
		panic("Need to implement")
	}
end:
	return err
}

var versionTagArg = &cli.Arg{
	Name:         VersionTagArg,
	Usage:        "Git version tag",
	Type:         reflect.String,
	CheckFunc:    check.versionTag,
	SetValueFunc: options.SetVersionTag,
}

func (checker) versionTag(requires cli.ArgRequires, tag any) (err error) {
	var ds *persister.DataStore
	var ctx context.Context
	var verTag string

	projName := options.ProjectName()
	if len(projName) == 0 {
		err = ErrNoProjectSpecified
		goto end
	}
	verTag = tag.(string)
	if len(verTag) == 0 {
		err = ErrNoVersionTagSpecified.Args("project", projName)
		goto end
	}
	ds = persister.GetDataStore()
	ctx = context.Background()

	//goland:noinspection GoSwitchMissingCasesForIotaConsts
	switch cli.Existence(requires) {
	case cli.MustExist:
		_, err = ds.LoadCodebaseByProjectNameAndVersionTag(ctx, persister.LoadCodebaseByProjectNameAndVersionTagParams{
			Name:       projName,
			VersionTag: verTag,
		})
		if err != nil {
			err = ErrFailedToAddCodebase.Err(err, "project", projName, "version_tag", verTag)
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
	Name:         RepoURLArg,
	Usage:        "The full GitHub repository URL for the project, e.g. https://github.com/golang/go",
	Type:         reflect.String,
	CheckFunc:    check.repoURL,
	SetValueFunc: options.SetRepoURL,
}

func (checker) repoURL(requires cli.ArgRequires, url any) (err error) {
	var parts []string
	var numParts int
	repoURL := url.(string)
	if len(repoURL) == 0 {
		err = ErrNoRepoURLSpecified
		goto end
	}
	if len(repoURL) == 1 && repoURL == "." {
		// A dot/period means ignore the repo
		goto end
	}
	parts = strings.Split(strings.TrimRight(repoURL, "/"), "/")
	numParts = len(parts)
	if numParts != 5 {
		err = ErrInvalidGitHubRepoURL
		goto end
	}
	if strings.Join(parts[:3], "/") != "https://github.com" {
		err = ErrInvalidGitHubRepoRootURL
		goto end
	}
	//goland:noinspection GoSwitchMissingCasesForIotaConsts
	switch cli.Existence(requires) {
	case cli.MustExist:
		err = cli.CheckURL(repoURL)
		if err != nil {
			err = ErrURLCouldNotBeDereferenced
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
