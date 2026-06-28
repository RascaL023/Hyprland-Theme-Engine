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
			Primary   string `json:"primary"`
			Secondary string `json:"secondary"`
			OnAccent  string `json:"on_accent"`
		} `json:"accent"`

		Text struct {
			Primary   string `json:"primary"`
			Secondary string `json:"secondary"`
			Muted     string `json:"muted"`
			Link      string `json:"link"`
		} `json:"text"`

		Layer struct {
			Base           string `json:"base"`
			Mantle         string `json:"mantle"`
			Crust          string `json:"crust"`
			Surface        string `json:"surface"`
			SurfaceRaised  string `json:"surface_raised"`
			SurfaceOverlay string `json:"surface_overlay"`
		} `json:"layer"`

		Border struct {
			Default string `json:"default"`
			Active  string `json:"active"`
		} `json:"border"`

		Status struct {
			Success  string `json:"success"`
			Warning  string `json:"warning"`
			Error    string `json:"error"`
			Critical string `json:"critical"`
			Info     string `json:"info"`
		} `json:"status"`
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
