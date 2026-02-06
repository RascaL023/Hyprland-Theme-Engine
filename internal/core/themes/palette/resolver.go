package palette

import (
	"theme-engine/internal/core/context"
	"theme-engine/internal/resolver"
)

func ResolvePalette(raw *Palette, ctx *context.Context) *ResolvedPalette {
	
	return &ResolvedPalette{
		Foreground: resolver.ResolveVar(raw.Foreground, ctx),
	}
}
