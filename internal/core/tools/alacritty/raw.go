package alacritty

type Raw struct {
	Opacity float64 `json:"opacity"`

	Padding struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"padding"`

	Cursor struct {
		Style struct {
			Shape string `json:"shape"`
		} `json:"style"`
	} `json:"cursor"`

	FontSize float64 `json:"fontSize"`
}
