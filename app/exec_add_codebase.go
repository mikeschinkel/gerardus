package app

import (
	"context"
	"fmt"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/persister"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdAddCodebase = CmdAdd.
	AddSubCommand("codebase", ExecAddCodebase).
	AddArg(projectArg.MustExist()).
	AddArg(versionTagArg.OkToExist())

//AddArg(cli.Arg{
//	Name:             "source_url",
//	Usage:            "URL for versioned source of a Project repo",
//	Optional:         true,
//	CheckFunc:        checkSourceURL,
//	SetStringValueFunc: options.SetSourceURL,
//})

func ExecAddCodebase(ctx context.Context, i *cli.CommandInvoker) (err error) {
	var p persister.Project

	ds := persister.GetDataStore()

	project := i.ArgString(ProjectArg)
	versionTag := i.ArgString(VersionTagArg)
	sourceURL := i.ArgString(SourceURLArg)

	p, err = ds.LoadProjectByName(ctx, project)
	if err != nil {
		err = ErrProjectNotFound.Err(err, "project", project)
		goto end
	}
	if versionTag != "." && len(sourceURL) == 0 {
		// If not yet set, compose the URL for GitHub
		sourceURL, err = persister.ComposeCodebaseSourceURL(p.RepoUrl, versionTag)
	}
	if err != nil {
		err = ErrInvalidCodebaseSourceURL.Err(err,
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
		err = ErrAddingCodebase.Err(err,
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

//func checkSourceURL(requires cli.ArgRequires, url any) (err error) {
//	sourceURL := url.(string)
//	switch requires {
//	case cli.MustExist:
//		err = cli.CheckURL(sourceURL)
//		if err != nil {
//			err = ErrSourceURLAppearsInvalid.Err(err, "source_url", sourceURL)
//			goto end
//		}
//	case cli.OkToExist:
//	case cli.MustNotExist:
//		// TODO: Implement
//	}
//end:
//	return err
//}
