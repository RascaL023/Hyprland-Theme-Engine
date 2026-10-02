package firefox

type Raw struct {
	Wallpaper			string `json:"wallpaper"`

	Special struct {
		Background		string	`json:"background"`
		Foreground		string	`json:"foreground"`
		Cursor			string	`json:"cursor"`
	} `json:"special"`

	Colors	struct {
		Color0			string `json:"color0"`
		Color1			string `json:"color1"`
		Color2			string `json:"color2"`
		Color3			string `json:"color3"`
		Color4			string `json:"color4"`
		Color5			string `json:"color5"`
		Color6			string `json:"color6"`
		Color7			string `json:"color7"`
		Color8			string `json:"color8"`
		Color9			string `json:"color9"`
		Color10			string `json:"color10"`
		Color11			string `json:"color11"`
		Color12			string `json:"color12"`
		Color13			string `json:"color13"`
		Color14			string `json:"color14"`
		Color15			string `json:"color15"`
	} `json:"colors"`
}
