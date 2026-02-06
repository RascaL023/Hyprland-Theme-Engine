package palette

type ResolvedPalette struct {
 	Foreground string
	Background string
	Cursor     string

	Colors     []string

	Extra struct {
		Accent struct {
			Primary string
			Secondary string
			Teritary string
		}

		Text struct {
			Primary string
			Secondary string
			Teritary string
		}

		Overlay []string
		Surface []string

		Base string
		Mantle string
		Crust string

		Warning string
		Critical string
		Charging string 
	}	

	Flat map[string]string
}

