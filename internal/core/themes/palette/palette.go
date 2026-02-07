package palette

type ResolvedPalette struct {
 	Foreground string
	Background string
	Cursor     string

	Colors     []string

	PrimaryAccent string
	SecondaryAccent string

	PrimaryText string
	SecondaryText string
	TeritaryText string

	PrimaryOverlay string
	SecondaryOverlay string
	TeritaryOverlay string

	PrimarySurface string
	SecondarySurface string
	TeritarySurface string

	Base string
	Mantle string
	Crust string

	Warning string
	Critical string
	Charging string 

	Flat map[string]string
}


type ResolvedPaletteVars struct {
  P *ResolvedPalette
}

func (r ResolvedPaletteVars) Get(k string) (string, bool) {
	v, ok := r.P.Flat[k]
	return v, ok
}

