package parser

type Codebase struct {
	DBId    int64
	Project string
	RepoURL string
}

func NewCodebase(project, repoURL string) *Codebase {
	return &Codebase{
		Project: project,
		RepoURL: repoURL,
	}
}
