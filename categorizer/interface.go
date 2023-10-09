package categorizer

type Interface struct {
	Package    string `json:"package"`
	Interface  string `json:"interface,omitempty"`
	Interfaces string `json:"interfaces,omitempty"`
}
