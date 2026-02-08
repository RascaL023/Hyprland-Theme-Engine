package kitty

import "theme-engine/internal/core/themes/palette"

type Kitty struct {
	Palette *palette.ResolvedPalette
	
	SelectionBackground string
	SelectionForeground string

	ActiveTabBg string
	ActiveTabFg string
	InActiveTabBg string
	InActiveTabFg string

	ActiveBorder string
	InActiveBorder string

	CursorShape string
	Opacity float64

	TabBar string
	TabPowerline string

	WindowBorder int
	WindowMargin int
	WindowPadding int

	FontSize float64
	FontFamily string
}
