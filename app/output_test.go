package app_test

import (
	"testing"
)

// TestOutput just here to stop Go and Goland from bickering
func TestOutput(t *testing.T) {}

func noCLIArgsOutput() string {
	return `
ERROR: No command specified:

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

func addArgsOutput() string {
	return `
ERROR: There is no 'add' command, but there are these commands:

  Usage: gerardus [<options>] <command> [<args>]

  Commands:

    - add codebase <project> <version_tag>
    - add project <project> <repo_url> [<about> [<website>]]

    Global Options:

      -data=<data_file>: Data file (sqlite3)
`
}

func projectUsage() string {
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
