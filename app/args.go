package app

import (
	"reflect"
	"strings"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/options"
	"github.com/mikeschinkel/gerardus/persister"
	"github.com/mikeschinkel/go-serr"
)

var projectArg = &cli.Arg{
	Name:         ProjectArg,
	Usage:        "Project name, e.g. 'golang'",
	Type:         reflect.String,
	CheckFunc:    Check.projectName,
	SetValueFunc: options.SetProjectName,
}

func (c *checker) projectName(ctx Context, requires cli.ArgRequires, project any) (err error) {
	var p persister.Project
	var existence = cli.Existence(requires)

	projName := project.(string)
	if projName == "" && existence == cli.MustExist {
		err = ErrProjectNotFound
		goto end
	}
	//goland:noinspection GoSwitchMissingCasesForIotaConsts
	switch existence {
	case cli.MustExist:
		injector := AssignFI(ctx, FI{Persister: PersisterFI{
			LoadProjectByNameFunc: c.App.Queries().LoadProjectByName,
		}})
		p, err = injector.Persister.LoadProjectByName(ctx, projName)
		if err != nil {
			err = ErrProjectNotFound.Err(err, "project", project)
			goto end
		}
		c.project = &p
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
	CheckFunc:    Check.versionTag,
	SetValueFunc: options.SetVersionTag,
}

func (c *checker) versionTag(ctx Context, requires cli.ArgRequires, tag any) (err error) {
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

	//goland:noinspection GoSwitchMissingCasesForIotaConsts
	switch cli.Existence(requires) {
	case cli.MustExist:
		injector := AssignFI(ctx, FI{Persister: PersisterFI{
			LoadCodebaseIDByProjectAndVersionFunc: c.App.Queries().LoadCodebaseIDByProjectAndVersion,
		}})
		_, err = injector.Persister.LoadCodebaseIDByProjectAndVersion(ctx, persister.LoadCodebaseIDByProjectAndVersionParams{
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
	CheckFunc:    Check.repoURL,
	SetValueFunc: options.SetRepoURL,
}

func (c *checker) repoURL(ctx Context, requires cli.ArgRequires, url any) (err error) {
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
		injector := AssignFI(ctx, FI{CheckURLFunc: cli.CheckURL})
		err = injector.CheckURL(repoURL)
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
