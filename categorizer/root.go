package categorizer

type Root struct {
	Project    Project      `json:"project"`
	Codebase   Codebase     `json:"codebase"`
	Legend     Legend       `json:"legend"`
	Categories []Categories `json:"categories"`
}
