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
			Visited   string `json:"visited"`
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
			Medium  string `json:"medium"`
		} `json:"border"`

		Status struct {
			Success  string `json:"success"`
			Warning  string `json:"warning"`
			Error    string `json:"error"`
			Critical string `json:"critical"`
			Info     string `json:"info"`
		} `json:"status"`

		Fg struct {
			Dim      string `json:"dim"`
			DimMuted string `json:"dim_muted"`
			Disabled string `json:"disabled"`
		} `json:"fg"`

		Bg struct {
			Conflict  string `json:"conflict"`
			DiskUsage string `json:"disk_usage"`
		} `json:"bg"`

		Syntax struct {
			Purple   string `json:"purple"`
			Magenta2 string `json:"magenta2"`
			Blue0    string `json:"blue0"`
			Blue1    string `json:"blue1"`
			Blue5    string `json:"blue5"`
			Blue6    string `json:"blue6"`
			Blue7    string `json:"blue7"`
			Green1   string `json:"green1"`
			Green2   string `json:"green2"`
			Orange   string `json:"orange"`
			Red1     string `json:"red1"`
			Teal     string `json:"teal"`
		} `json:"syntax"`

		UI struct {
			BgStatusline string `json:"bg_statusline"`
		} `json:"ui"`
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
