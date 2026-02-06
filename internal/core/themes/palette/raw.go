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



type RawPaletteVars struct {
	P *Palette
}

func (r RawPaletteVars) Get(k string) (string, bool) {
	switch k {
		case "foreground": return r.P.Foreground, true;
		case "background": return r.P.Background, true;
		case "cursor": return r.P.Cursor, true;
	}

	if len(k) > 5 && k[:5] == "color" {
		i := 0
		for j := 5; j < len(k); j++ {
			c := k[j]
			if c < '0' || c > '9' {
				return "", false
			}
			i = i*10 + int(c-'0')
		}

		if i >= 0 && i < len(r.P.Colors) {
			return r.P.Colors[i], true
		}
	}

	return "", false
}
