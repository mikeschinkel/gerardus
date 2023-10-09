package categorizer

type Codebase struct {
	Slug        string `json:"slug"`
	RepoURL     string `json:"repo_url"`
	LocalDirVar string `json:"local_dir_var"`
}
