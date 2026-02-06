package palette

import "theme-engine/internal/resolver"

func ResolvePalette(raw *Palette) *ResolvedPalette {
	return &ResolvedPalette{
		Foreground: raw.Foreground,
		Background: raw.Background,
		Cursor: raw.Cursor,

		Colors: raw.Colors,

		PrimaryAccent: resolver.ResolveVar(raw.Extra.Accent.Primary, RawPaletteVars{raw}),
		SecondaryAccent: resolver.ResolveVar(raw.Extra.Accent.Secondary, RawPaletteVars{raw}),

		PrimaryText: resolver.ResolveVar(raw.Extra.Text.Primary, RawPaletteVars{raw}),
		SecondaryText: resolver.ResolveVar(raw.Extra.Text.Secondary, RawPaletteVars{raw}),
		TeritaryText: resolver.ResolveVar(raw.Extra.Text.Teritary, RawPaletteVars{raw}),

		PrimarySurface: resolver.ResolveVar(raw.Extra.Surface[2], RawPaletteVars{raw}),
		SecondarySurface: resolver.ResolveVar(raw.Extra.Surface[1], RawPaletteVars{raw}),
		TeritarySurface: resolver.ResolveVar(raw.Extra.Surface[0], RawPaletteVars{raw}),

		PrimaryOverlay: resolver.ResolveVar(raw.Extra.Overlay[2], RawPaletteVars{raw}),
		SecondaryOverlay: resolver.ResolveVar(raw.Extra.Overlay[1], RawPaletteVars{raw}),
		TeritaryOverlay: resolver.ResolveVar(raw.Extra.Overlay[0], RawPaletteVars{raw}),

		Base: resolver.ResolveVar(raw.Extra.Base, RawPaletteVars{raw}),
		Mantle: resolver.ResolveVar(raw.Extra.Mantle, RawPaletteVars{raw}),
		Crust: resolver.ResolveVar(raw.Extra.Crust, RawPaletteVars{raw}),

		Warning: resolver.ResolveVar(raw.Extra.Warning, RawPaletteVars{raw}),
		Critical: resolver.ResolveVar(raw.Extra.Critical, RawPaletteVars{raw}),
		Charging: resolver.ResolveVar(raw.Extra.Charging, RawPaletteVars{raw}),
	}
}
