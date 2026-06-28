package palette

type ResolvedPalette struct {
	Foreground string
	Background string
	Cursor     string

	Colors     []string

	AccentPrimary   string
	AccentSecondary string
	AccentOn        string

	TextPrimary   string
	TextSecondary string
	TextMuted     string
	TextLink      string

	LayerBase           string
	LayerMantle         string
	LayerCrust          string
	LayerSurface        string
	LayerSurfaceRaised  string
	LayerSurfaceOverlay string

	BorderDefault string
	BorderActive  string

	StatusSuccess  string
	StatusWarning  string
	StatusError    string
	StatusCritical string
	StatusInfo     string

	Flat map[string]string
}


type ResolvedPaletteVars struct {
  P *ResolvedPalette
}

func (r ResolvedPaletteVars) Get(k string) (string, bool) {
	v, ok := r.P.Flat[k]
	return v, ok
}
