package parser

type Project struct {
	Id      int64
	Name    string
	RepoURL string
}

func NewProject(name, repoURL string) *Project {
	return &Project{
		Name:    name,
		RepoURL: repoURL,
	}
}
