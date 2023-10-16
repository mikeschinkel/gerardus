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
		err = fmt.Errorf("project '%s' does not exist; %w", project, err)
		goto end
	}
	if len(sourceURL) == 0 {
		// If not yet set, compose the URL for GitHub
		sourceURL, err = persister.CodebaseSourceURL(p.RepoUrl, versionTag)
	}
	if err != nil {
		err = fmt.Errorf("invalid URL; %w. Potentially bad repo url (%s) or bad version tag (%s)",
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
