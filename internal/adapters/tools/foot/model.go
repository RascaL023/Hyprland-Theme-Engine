package foot

type Raw struct {
	Opacity float64 `json:"opacity"`
	Padding struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"padding"`
}
