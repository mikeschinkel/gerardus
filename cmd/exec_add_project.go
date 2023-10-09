package main

import (
	"context"
	"fmt"

	"gerardus/persister"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdAddProject = CmdAdd.
	AddSubCommand("project", ExecAddProject).
	AddArgs("name", "repo_url").
	AddOptArgs("about", "website")

func ExecAddProject(args ...string) (err error) {
	var about, website string

	name := args[0]
	repoURL := args[1]
	if len(args) > 2 {
		about = args[2]
	}
	if len(args) > 3 {
		website = args[3]
	}
	if len(about) == 0 {
		var info *persister.GitHubRepoInfo
		info, err = persister.RequestGitHubRepoInfo(repoURL)
		about = info.Description
		website = info.Homepage
	}
	ctx := context.Background()
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
