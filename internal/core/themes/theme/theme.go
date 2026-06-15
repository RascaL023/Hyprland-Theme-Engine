package theme

import "encoding/json"

type Theme struct {
	Theme struct {
		Name string `json:"name"`

		Fonts struct {
			Primary   string  `json:"primary"`
			Secondary string  `json:"secondary"`
			Size      float64 `json:"size"`
		} `json:"fonts"`
	} `json:"theme"`

	Tools map[string]json.RawMessage `json:"tools"`
	// Tools struct {
	// 	Cava struct {
	// 		Gradients []string `json:"gradients"`
	// 	} `json:"cava"`
	// } `json:"tools"`
}
