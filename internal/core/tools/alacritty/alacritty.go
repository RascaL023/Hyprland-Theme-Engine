package alacritty

import "theme-engine/internal/core/themes/palette"

type Alacritty struct {
	Palette *palette.ResolvedPalette

	Opacity  float64
	PaddingX float64
	PaddingY float64

	CursorShape string

	FontSize   float64
	FontFamily string
}
