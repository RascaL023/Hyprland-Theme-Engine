package context

import (
	"theme-engine/internal/core/themes/palette"
	"theme-engine/internal/core/themes/theme"
)

type Context struct {
	Palette   *palette.ResolvedPalette
	Theme     *theme.Theme
	ThemeType string
}
