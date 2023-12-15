package app

import (
	"context"
	"fmt"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/persister"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdAddProject = CmdAdd.
	AddSubCommand("project", ExecAddProject).
	AddArg(projectArg.MustExist()).
	AddArg(repoURLArg.MustExist()).
	AddArg(&cli.Arg{
		Name:     AboutArg,
		Optional: true,
		Usage:    "Repo description. Defaults to 'about' from the GitHub API",
	}).
	AddArg(&cli.Arg{
		Name:     WebsiteArg,
		Optional: true,
		Usage:    "Project website URL. Defaults to 'website' from the GitHub API",
	})

func ExecAddProject(ctx context.Context, i *cli.CommandInvoker) (err error) {

	di := ctx.Value(DIKey).(*DI)

	name := i.ArgString(ProjectArg)
	repoURL := i.ArgString(RepoURLArg)
	about := i.ArgString(AboutArg)
	website := i.ArgString(WebsiteArg)

	di.Assign(DI{RepoInfoRequesterFunc: persister.RequestGitHubRepoInfo})

	if len(about) == 0 {
		var info *persister.RepoInfo
		info, err = di.RepoInfoRequesterFunc(repoURL)
		if err != nil {
			goto end
		}
		about = info.Description
		website = info.Homepage
	}

	di.Assign(DI{UpsertProjectFunc: persister.GetDataStore().UpsertProject})
	ctx = context.Background()
	_, err = di.UpsertProjectFunc(ctx, persister.UpsertProjectParams{
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
