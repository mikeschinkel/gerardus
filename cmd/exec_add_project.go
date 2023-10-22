package main

import (
	"context"
	"fmt"

	"gerardus/cli"
	"gerardus/persister"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdAddProject = CmdAdd.
	AddSubCommand("project", ExecAddProject).
	AddArg(projectArg.OkToExist()).
	AddArg(repoURLArg.MustExist()).
	AddArg(cli.Arg{
		Name:     AboutArg,
		Optional: true,
	}).
	AddArg(cli.Arg{
		Name:     WebsiteArg,
		Optional: true,
	})

func ExecAddProject(args cli.ArgsMap) (err error) {
	var about, website string
	var ctx context.Context

	name := args[ProjectArg].String()

	repoURL := args[RepoURLArg].String()

	if !args[AboutArg].IsZero() {
		about = args[AboutArg].String()
	}
	if !args[WebsiteArg].IsZero() {
		website = args[WebsiteArg].String()
	}
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
		err = errFailedToAddProject.Err(err, "project", name, "repo_url", repoURL)
		goto end
	}
	fmt.Printf("\nSuccessfully added project '%s' with repo URL %s.\n",
		name,
		repoURL,
	)
end:
	return err
}
