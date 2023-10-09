package main

import (
	"context"
	"fmt"

	"gerardus/persister"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdAddCodebase = CmdAdd.
	AddSubCommand("codebase", ExecAddCodebase).
	AddArgs("project", "version_tag").
	AddOptArgs("source_url")

func ExecAddCodebase(args ...string) (err error) {
	var sourceURL string
	var p persister.Project

	ds := persister.GetDataStore()

	project := args[0]
	versionTag := args[1]
	if len(args) > 2 {
		// If passed as 3rd argument of `add codebase` then set sourceURL.
		sourceURL = args[2]
	}

	ctx := context.Background()
	p, err = ds.LoadProjectByName(ctx, project)
	if err != nil {
		err = fmt.Errorf("project '%s' does not exist; %w", project, err)
		goto end
	}
	if len(sourceURL) == 0 {
		// If not yet set, compose the URL for GitHub
		sourceURL, err = persister.CodebaseSourceURL(p.RepoUrl, versionTag)
	}
	if err != nil {
		err = fmt.Errorf("invalid URL; %w. Potnentially bad repo url (%s) or bad version tag (%s)",
			err, p.RepoUrl, versionTag)
		goto end
	}
	_, err = ds.UpsertCodebase(ctx, persister.UpsertCodebaseParams{
		ProjectID:  p.ID,
		VersionTag: versionTag,
		SourceUrl:  sourceURL,
	})
	if err != nil {
		err = fmt.Errorf("fail to add new project for %#v; %w", args, err)
		goto end
	}
	fmt.Printf("\nSuccessfully added codebase for '%s' version '%s' with source URL %s.\n",
		project,
		versionTag,
		sourceURL,
	)
end:
	return err
}
