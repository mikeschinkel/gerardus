package categorizer

type Categories struct {
	Slug       string      `json:"slug"`
	Name       string      `json:"name"`
	SeeAlso    []string    `json:"see_also,omitempty"`
	Interfaces []Interface `json:"interfaces"`
}
