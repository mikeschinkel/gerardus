package app

import (
	"context"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/persister"
)

//goland:noinspection GoUnusedGlobalVariable
var CmdAddCodebase = CmdAdd.
	AddSubCommand("codebase", Root.ExecAddCodebase).
	AddArg(projectArg.NotEmpty().MustExist()).
	AddArg(versionTagArg.NotEmpty().MustValidate().NotExist())

func (a *App) ExecAddCodebase(ctx context.Context, i *cli.CommandInvoker) (err error) {
	var p persister.Project

	project := i.ArgString(ProjectArg)
	versionTag := i.ArgString(VersionTagArg)
	sourceURL := i.ArgString(SourceURLArg)

	p, err = a.Queries().LoadProjectByName(ctx, project)
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
			cli.HelpArg, "Potentially bad project name, version tag, or repo URL.",
		)
		goto end
	}
	_, err = a.Queries().UpsertCodebase(ctx, persister.UpsertCodebaseParams{
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
	cli.StdOut("\nSuccessfully added codebase for '%s' version '%s'.\n",
		project,
		versionTag,
	)
end:
	return err
}
