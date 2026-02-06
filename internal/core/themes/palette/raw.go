package palette

type Raw struct {
	Palettes struct {
		Dark Palette `json:"dark"`
		Light Palette `json:"light"`
	} `json:"palettes"`
}

type Palette struct {
	Foreground string   `json:"foreground"`
	Background string   `json:"background"`
	Cursor     string   `json:"cursor"`

	Colors     []string `json:"colors"`

	Extra struct {
		Accent struct {
			Primary string `json:"primary"`
			Secondary string `json:"secondary"`
			Teritary string `json:"teritary"`
		} `json:"accent"`

		Text struct {
			Primary string `json:"primary"`
			Secondary string `json:"secondary"`
			Teritary string `json:"teritary"`
		} `json:"text"`

		Overlay []string `json:"overlay"`
		Surface []string `json:"surface"`

		Base string `json:"base"`
		Mantle string `json:"mantle"`
		Crust string `json:"crust"`

		Warning string `json:"warning"`
		Critical string `json:"critical"`
		Charging string `json:"charging"`
	}
}

func (r *Raw) ResolveSelected(selected string) *ResolvedPalette {
	switch selected {
		case "light": return ResolvePalette(&r.Palettes.Light);
		default: return ResolvePalette(&r.Palettes.Dark);
	}
}
