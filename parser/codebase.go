package parser

type Codebase struct {
	DBId    int64
	RepoURL string
}

func NewCodebase(repoURL string) *Codebase {
	return &Codebase{
		RepoURL: repoURL,
	}
}
