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
	AddArg(projectArg).
	AddArg(repoURLArg).
	AddArg(&cli.Arg{
		Name:     "about",
		Optional: true,
	}).
	AddArg(&cli.Arg{
		Name:     "website",
		Optional: true,
	})

func ExecAddProject(args cli.StringMap) (err error) {
	var about, website string
	var ok bool
	var ctx context.Context

	name := args["name"]

	repoURL := args["repo_url"]

	if len(args["about"]) != 0 {
		about = args["about"]
	}
	if len(args["website"]) != 0 {
		website = args["website"]
	}
	if len(about) == 0 {
		var info *persister.GitHubRepoInfo
		info, err = persister.RequestGitHubRepoInfo(repoURL)
		if !ok {
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
		err = fmt.Errorf("fail to add new project for %#v; %w", args, err)
		goto end
	}
	fmt.Printf("\nSuccessfully added project '%s' with repo URL %s.\n",
		name,
		repoURL,
	)
end:
	return err
}
