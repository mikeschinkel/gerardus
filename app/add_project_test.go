package app_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mikeschinkel/gerardus/app"
	"github.com/mikeschinkel/gerardus/persister"
)

func TestAddProject(t *testing.T) {}

func addProjectTests() []test {
	return []test{
		{
			name:   "FAIL — NO PROJECT ARG",
			fail:   true,
			args:   []string{"add", "project"},
			output: "\nERROR: Argument cannot be empty [arg_name='<project>']:\n" + addProjectUsage(),
			errStr: "argument cannot be empty [arg_name='<project>']",
		},
		{
			name:   "FAIL — NO REPO URL ARG",
			fail:   true,
			args:   []string{"add", "project", "golang"},
			output: "\nERROR: Argument cannot be empty [arg_name='<repo_url>']:\n" + addProjectUsage(),
			errStr: "argument cannot be empty [arg_name='<repo_url>']",
		},
		{
			name:   "FAIL — NO GITHUB URL",
			fail:   true,
			args:   []string{"add", "project", "golang", "https://not.there"},
			output: "\nERROR: Not a valid GitHub repo URL [repo_url='https://not.there']:\n" + addProjectUsage(),
			errStr: "not a valid GitHub repo URL [repo_url='https://not.there']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadMissingProjectByNameStub,
			},
		},
		{
			name:   "FAIL — NO HTTPS",
			fail:   true,
			args:   []string{"add", "project", "golang", "http://github.com/not/there"},
			output: "\nERROR: Repo URL does not begin with https://github.com [repo_url='http://github.com/not/there']:\n" + addProjectUsage(),
			errStr: "repo URL does not begin with https://github.com [repo_url='http://github.com/not/there']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadMissingProjectByNameStub,
			},
		},
		{
			name:   "FAIL — PROJECT EXISTS",
			fail:   true,
			args:   []string{"add", "project", "golang", "https://github.com/golang/go"},
			output: "\nERROR: Project exists [project='golang']:\n" + addProjectUsage(),
			errStr: "project exists [project='golang']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadFoundProjectByNameStub,
			},
		},
		{
			name:   "FAIL — CANNOT DEREFERENCE PROJECT URL",
			fail:   true,
			args:   []string{"add", "project", "golang", "https://github.com/not/there"},
			output: "\nERROR: URL could not be dereferenced [repo_url='https://github.com/not/there']:\n" + addProjectUsage(),
			errStr: "URL could not be dereferenced [repo_url='https://github.com/not/there']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadMissingProjectByNameStub,
			},
			fi: func(fi app.FI) app.FI {
				fi.Persister.RequestGitHubRepoInfoFunc = RequestGitHubRepoInfoStub // TODO: Verify this is called by this test
				return fi
			},
		},
		{
			name:   "SUCCESS — ADD PROJECT",
			fail:   false,
			args:   []string{"add", "project", "golang", "https://github.com/golang/go"},
			output: "\nSuccessfully added project 'golang' with repo URL https://github.com/golang/go.\n",
			errStr: "<n/a>",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadMissingProjectByNameStub,
				UpsertProjectFunc:     SuccessfulUpsertProjectStub,
			},
		},
	}
}

func addProjectUsage() string {
	return `
  Usage: gerardus [<options>] <command> [<args>]

  Command:

    - add project <project> <repo_url> [<about> [<website>]]

        Args:

          project:   Project name, e.g. 'golang'
          repo_url:  The full GitHub repository URL for the project, e.g. https://github.com/golang/go
          about:     Repo description. Defaults to 'about' from the GitHub API
          website:   Project website URL. Defaults to 'website' from the GitHub API

        Global Options:

          -data=<data_file>: Data file (sqlite3)
`
}

func SuccessfulUpsertProjectStub(ctx context.Context, arg persister.UpsertProjectParams) (persister.Project, error) {
	return projectID1NameGolangStub(), nil
}

func LoadFoundProjectByNameStub(ctx context.Context, name string) (persister.Project, error) {
	return projectID1NameGolangStub(), nil
}

func LoadMissingProjectByNameStub(ctx context.Context, name string) (persister.Project, error) {
	return persister.Project{}, sql.ErrNoRows
}

func projectID1NameGolangStub() persister.Project {
	return persister.Project{
		ID:      1,
		Name:    "golang",
		RepoUrl: "https://github.com/golang/go",
	}
}
