package app

import (
	"context"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/persister"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdAddProject = CmdAdd.
	AddSubCommand("project", Root.ExecAddProject).
	AddArg(projectArg.NotEmpty().NotExist()).
	AddArg(repoURLArg.NotEmpty().MustExist()).
	AddArg(cli.Arg{
		Name:     AboutArg,
		Optional: true,
		Usage:    "Repo description. Defaults to 'about' from the GitHub API",
		Requires: cli.AndRequires(cli.EmptyOk, cli.IgnoreExists),
	}).
	AddArg(cli.Arg{
		Name:     WebsiteArg,
		Optional: true,
		Usage:    "Project website URL. Defaults to 'website' from the GitHub API",
		Requires: cli.AndRequires(cli.EmptyOk, cli.IgnoreExists),
	})

func (a *App) ExecAddProject(ctx context.Context, i *cli.CommandInvoker) (err error) {

	name := i.ArgString(ProjectArg)
	repoURL := i.ArgString(RepoURLArg)
	about := i.ArgString(AboutArg)
	website := i.ArgString(WebsiteArg)

	if len(about) == 0 {
		info := a.repoInfo
		about = info.Description
		website = info.Homepage
	}

	_, err = a.Queries().UpsertProject(ctx, persister.UpsertProjectParams{
		Name:    name,
		About:   about,
		RepoUrl: repoURL,
		Website: website,
	})
	if err != nil {
		err = ErrFailedToAddProject.Err(err, "project", name, "repo_url", repoURL)
		goto end
	}
	cli.StdOut("\nSuccessfully added project '%s' with repo URL %s.\n",
		name,
		repoURL,
	)
end:
	return err
}
