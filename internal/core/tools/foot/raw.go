package foot

type Raw struct {
	FontSize float64 `json:"fontSize"`
	Opacity float64 `json:"opacity"`
	Padding struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"padding"`
}
