package parser

type Codebase struct {
	Id         int64
	Project    string
	VersionTag string
}

func NewCodebase(project, tag string) *Codebase {
	return &Codebase{
		Project:    project,
		VersionTag: tag,
	}
}
