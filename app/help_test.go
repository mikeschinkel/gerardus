package app_test

import (
	"testing"
)

func TestHelpCommand(t *testing.T) {}

//goland:noinspection GoUnusedFunction
func helpTests() []test {
	return []test{
		{
			name:   "FAIL - INVALID COMMAND",
			fail:   true,
			args:   []string{"help", "foo", "bar"},
			output: "\nERROR: Command is not valid [command='foo bar']:\n" + helpUsage(),
			errStr: "help",
		},
		{
			name:   "SUCCESS - HELP",
			fail:   false,
			args:   []string{"help"},
			output: helpUsage(),
			errStr: "<n/a>",
		},
		{
			name:   "SUCCESS - ADD HELP",
			fail:   false,
			args:   []string{"help", "add"},
			output: helpAddUsage(),
			errStr: "<n/a>",
		},
		{
			name:   "SUCCESS - ADD PROJECT HELP",
			fail:   false,
			args:   []string{"help", "add", "project"},
			output: helpAddProjectUsage(),
			errStr: "<n/a>",
		},
		{
			name:   "SUCCESS - ADD CODEBASE HELP",
			fail:   false,
			args:   []string{"help", "add", "codebase"},
			output: helpAddCodebaseUsage(),
			errStr: "<n/a>",
		},
		{
			name:   "SUCCESS - MAP HELP",
			fail:   false,
			args:   []string{"help", "map"},
			output: helpMapUsage(),
			errStr: "<n/a>",
		},
	}
}

func helpUsage() string {
	return `
  Usage: gerardus [<options>] <command> [<args>]

  Commands:

    - add codebase <project> <version_tag>
    - add project <project> <repo_url> [<about> [<website>]]
    - help [<command>]
    - map [-src=<source_dir>] <project> <version_tag>

    Global Options:

      -data=<data_file>: Data file (sqlite3)
`
}
func helpAddUsage() string {
	return `
  Usage: gerardus [<options>] <command> [<args>]

  Commands:

    - add codebase <project> <version_tag>
    - add project <project> <repo_url> [<about> [<website>]]

    Global Options:

      -data=<data_file>: Data file (sqlite3)
`
}
func helpAddProjectUsage() string {
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
func helpAddCodebaseUsage() string {
	return `
  Usage: gerardus [<options>] <command> [<args>]

  Command:

    - add codebase <project> <version_tag>

        Args:

          project:     Project name, e.g. 'golang'
          version_tag: Git version tag

        Global Options:

          -data=<data_file>: Data file (sqlite3)
`
}
func helpMapUsage() string {
	return `
  Usage: gerardus [<options>] <command> [<args>]

  Command:

    - map [-src=<source_dir>] <project> <version_tag>

        Options:

          -src=<source_dir>: Source directory
           Default:  /tmp/test/dir

        Args:

          project:     Project name, e.g. 'golang'
          version_tag: Git version tag

        Global Options:

          -data=<data_file>: Data file (sqlite3)
`
}
