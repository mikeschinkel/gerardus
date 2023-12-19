package app

import (
	"context"
	"fmt"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/persister"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdAddProject = CmdAdd.
	AddSubCommand("project", Root.ExecAddProject).
	AddArg(projectArg.NotEmpty().MustFailCheck()).
	AddArg(repoURLArg.NotEmpty().MustPassCheck()).
	AddArg(cli.Arg{
		Name:     AboutArg,
		Optional: true,
		Usage:    "Repo description. Defaults to 'about' from the GitHub API",
		Requires: cli.AndRequires(cli.EmptyOk, cli.IgnoreCheck),
	}).
	AddArg(cli.Arg{
		Name:     WebsiteArg,
		Optional: true,
		Usage:    "Project website URL. Defaults to 'website' from the GitHub API",
		Requires: cli.AndRequires(cli.EmptyOk, cli.IgnoreCheck),
	})

func (a *App) ExecAddProject(ctx context.Context, i *cli.CommandInvoker) (err error) {
	var injector FI

	name := i.ArgString(ProjectArg)
	repoURL := i.ArgString(RepoURLArg)
	about := i.ArgString(AboutArg)
	website := i.ArgString(WebsiteArg)

	if len(about) == 0 {
		var info *persister.RepoInfo
		injector = AssignFI(ctx, FI{Persister: PersisterFI{RepoInfoRequesterFunc: persister.RequestGitHubRepoInfo}})
		info, err = injector.Persister.RepoInfoRequesterFunc(repoURL)
		if err != nil {
			goto end
		}
		about = info.Description
		website = info.Homepage
	}

	ctx = context.Background()
	injector = AssignFI(ctx, FI{Persister: PersisterFI{UpsertProjectFunc: a.Queries().UpsertProject}})
	_, err = injector.Persister.UpsertProject(ctx, persister.UpsertProjectParams{
		Name:    name,
		About:   about,
		RepoUrl: repoURL,
		Website: website,
	})
	if err != nil {
		err = ErrFailedToAddProject.Err(err, "project", name, "repo_url", repoURL)
		goto end
	}
	fmt.Printf("\nSuccessfully added project '%s' with repo URL %s.\n",
		name,
		repoURL,
	)
end:
	return err
}
