package renderctx

import (
	"theme-engine/internal/domain/palette"
	"theme-engine/internal/domain/theme"
)

type Context struct {
	Palette   *palette.ResolvedPalette
	Theme     *theme.Theme
	ThemeType string
}
