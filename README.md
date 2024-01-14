# gerardus
Golang Source Code Cartographer 

Gerardus reads all the `.go` files in a directory recursively and then generates a Sqlite database containing what it found.

It defaults to reading the directory found in `GOROOT` if you don't specify otherwise.

Still very much a **work-in-progress**, and thus **pre-alpha.**


## Roadmap

In no particular order:

### CI/CD
1. Install commands for `go generate` in GitHub Actions `yaml`
2. Vendor dependencies

### Testing
1. ~~Implement robust integration tests~~ 
    1. ~~Simulate calling from the command line.~~
    2. ~~Stubbing out the database~~
    3. ~~Add tests for `add project` command.~~
    4. ~~Add tests for `add codebase` command.~~
    5. ~~Add tests for `map` command.~~
    6. ~~Add tests for `help` command.~~

### Miscellaneous
1. ~~Ensure chans refactor works~~
2. ~~Break out Package from Import~~
3. ~~Break out Project from Codebase~~
4. ~~Add GitHub URL parser for Name, Version~~
5. ~~Add Web client to get About and Website~~
    1. ~~https://chat.openai.com/share/cb482a49-aad6-4f3b-abf9-201d60055054~~
6. Get slices vs. chans working
7. Leverage https://deps.dev/

### CLI
1. ~~Add other commands besides `map`~~
2. Add switch for local directory to clone repos to
3. Add ability to clone a Git repo
4. Use Cobra
5. Try Kong

### SQL/Store
1. ~~Add survey_id to all relevant tables~~
2. ~~Add Name, About to Project~~
3. ~~Add VersionTag, Website fields to Codebase~~
4. ~~Capture package data~~
5. Add schema migration tool
6. Write SQL for Category Interfaces
7. Capture constraint interfaces
8. Capture interface methods
9. Capture struct fields
10. Break out `survey_file` from `file`
11. Support Postgres as an alternative

### Categories
1. Implement categories
2. Implement category types
3. Implement type category
4. 
### Config
1. Create JSON config files
2. Embed default config files
3. Support YAML config files.
4. Add local clone root directory config

### Env
1. Allow options to be set by env
2. Especially Codebase Repo URL

### Templating
1. Write Markdown template for interfaces for a category
2. `text/template` or `html/template`?

### Web Service
1. Serve Markdown for Obsidian
2. Http server to serve Markdown snippets
3. Command to serve
4. Command to add serve to `at` command
5. Command to add serve to `cron` command
6. Configuration to serve via SystemD service

### O/Ses
1. Mac
2. Linux
3. Windows
