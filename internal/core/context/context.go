package context

import "theme-engine/internal/core/themes/palette"


type Context struct {
	Palette *palette.ResolvedPalette
}
