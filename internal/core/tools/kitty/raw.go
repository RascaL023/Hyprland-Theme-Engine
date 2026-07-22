package kitty

type Raw struct {
	CursorShape string  `json:"cursorShape"`
	Opacity     float64 `json:"opacity"`

	Tab struct {
		Active struct {
			Foreground string `json:"foreground"`
			Background string `json:"background"`
		} `json:"active"`

		InActive struct {
			Foreground string `json:"foreground"`
			Background string `json:"background"`
		} `json:"inactive"`

		Bar   string `json:"bar"`
		Style string `json:"style"`
	} `json:"tab"`

	Window struct {
		Border struct {
			Active   string `json:"active"`
			InActive string `json:"inactive"`
			Width    int    `json:"width"`
		} `json:"border"`

		Margin  int `json:"margin"`
		Padding int `json:"padding"`
	} `json:"window"`
}
