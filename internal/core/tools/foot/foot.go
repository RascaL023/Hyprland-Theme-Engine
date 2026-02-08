package foot

import "theme-engine/internal/core/themes/palette"

type Foot struct {
	Font string
	FontSize float64

	PaddingX float64
	PaddingY float64

	Palette *palette.ResolvedPalette
}
