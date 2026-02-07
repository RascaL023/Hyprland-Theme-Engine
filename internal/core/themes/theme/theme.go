package theme

import "encoding/json"

type Theme struct {
	Name string `json:"theme"`

	Tools map[string]json.RawMessage `json:"tools"`
	// Tools struct {
	// 	Cava struct {
	// 		Gradients []string `json:"gradients"`
	// 	} `json:"cava"`
	// } `json:"tools"`
}
