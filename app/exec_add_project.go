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
	}).
	AddArg(&cli.Arg{
		Name:     WebsiteArg,
		Optional: true,
	})

func ExecAddProject(i *cli.CommandInvoker) (err error) {
	var ctx context.Context

	name := i.ArgString(ProjectArg)
	repoURL := i.ArgString(RepoURLArg)
	about := i.ArgString(AboutArg)
	website := i.ArgString(WebsiteArg)

	if len(about) == 0 {
		var info *persister.GitHubRepoInfo
		info, err = persister.RequestGitHubRepoInfo(repoURL)
		if err != nil {
			goto end
		}
		about = info.Description
		website = info.Homepage
	}

	ctx = context.Background()
	_, err = persister.GetDataStore().UpsertProject(ctx, persister.UpsertProjectParams{
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
