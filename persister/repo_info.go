package persister

// RepoInfo contains relevant information about a GitHub repository
type RepoInfo struct {
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}
