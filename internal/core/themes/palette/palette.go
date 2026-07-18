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

	SyntaxPurple   string
	SyntaxMagenta2 string
	SyntaxBlue0    string
	SyntaxBlue1    string
	SyntaxBlue5    string
	SyntaxBlue6    string
	SyntaxBlue7    string
	SyntaxGreen1   string
	SyntaxGreen2   string
	SyntaxOrange   string
	SyntaxRed1     string
	SyntaxTeal     string

	BgStatusline string

	Flat map[string]string
}


type ResolvedPaletteVars struct {
  P *ResolvedPalette
}

func (r ResolvedPaletteVars) Get(k string) (string, bool) {
	v, ok := r.P.Flat[k]
	return v, ok
}
