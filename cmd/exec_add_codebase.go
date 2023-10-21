package main

import (
	"context"
	"fmt"

	"gerardus/cli"
	"gerardus/options"
	"gerardus/persister"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdAddCodebase = CmdAdd.
	AddSubCommand("codebase", ExecAddCodebase).
	AddArg(projectArg).
	AddArg(versionTagArg).
	AddArg(&cli.Arg{
		Name:             "source_url",
		Usage:            "URL for versioned source of a Project repo",
		Optional:         true,
		CheckFunc:        checkSourceURL,
		SetStringValFunc: options.SetSourceURL,
	})

func ExecAddCodebase(args cli.ArgsMap) (err error) {
	var p persister.Project

	ds := persister.GetDataStore()

	project := options.ProjectName()
	versionTag := options.VersionTag()
	sourceURL := options.SourceURL()

	ctx := context.Background()
	p, err = ds.LoadProjectByName(ctx, project)
	if err != nil {
		err = errProjectNotFound.Err(err, "project", project)
		goto end
	}
	if len(sourceURL) == 0 {
		// If not yet set, compose the URL for GitHub
		sourceURL, err = persister.CodebaseSourceURL(p.RepoUrl, versionTag)
	}
	if err != nil {
		err = errInvalidCodebaseSourceURL.Err(err,
			"project", project,
			"version_tag", versionTag,
			"repo_url", p.RepoUrl,
			"help", "Potentially bad project name, version tag, or repo URL.",
		)
		goto end
	}
	_, err = ds.UpsertCodebase(ctx, persister.UpsertCodebaseParams{
		ProjectID:  p.ID,
		VersionTag: versionTag,
		SourceUrl:  sourceURL,
	})
	if err != nil {
		err = errAddingCodebase.Err(err,
			"project_id", p.ID,
			"version_tag", versionTag,
			"project", project,
			"repo_url", p.RepoUrl,
			"source_url", sourceURL,
		)
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

func checkSourceURL(url any) (err error) {
	sourceURL := url.(string)
	err = cli.CheckURL(sourceURL)
	if err != nil {
		err = fmt.Errorf("source URL does not appear to be valid; %w", err)
		goto end
	}
end:
	return err
}
