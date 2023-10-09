package categorizer

import (
	"encoding/json"
	"os"
)

type Categorizer struct {
	Root
	filepath string
}

func NewCategorizer(filepath string) *Categorizer {
	return &Categorizer{
		filepath: filepath,
	}
}

func (c *Categorizer) LoadJSON() error {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		goto end
	}
	err = json.Unmarshal(bytes, &c.Root)
	if err != nil {
		goto end
	}

end:
	return err
}

func (c *Categorizer) WriteJSON(filename string) error {
	bytes, err := json.MarshalIndent(c.Root, "", "   ")
	if err != nil {
		goto end
	}

	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		goto end
	}

end:
	return err
}

//func main() {
//	// Read example
//	c := NewCategorizer("input.json")
//	err := c.LoadJSON()
//	if err != nil {
//		panic(err)
//	}
//
//	// Here you can operate with the data stored in the `root variable`
//	// For an example, let's just write the same data into another file.
//
//	// Writing example
//	err = c.WriteJSON("output.json")
//	if err != nil {
//		panic(err)
//	}
//}
