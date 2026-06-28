package palette

import "theme-engine/internal/resolver"

func ResolvePalette(raw *Palette) *ResolvedPalette {
	return &ResolvedPalette{
		Foreground: raw.Foreground,
		Background: raw.Background,
		Cursor: raw.Cursor,

		Colors: raw.Colors,

		AccentPrimary: resolver.ResolveVar(raw.Extra.Accent.Primary, RawPaletteVars{raw}),
		AccentSecondary: resolver.ResolveVar(raw.Extra.Accent.Secondary, RawPaletteVars{raw}),
		AccentOn: resolver.ResolveVar(raw.Extra.Accent.OnAccent, RawPaletteVars{raw}),

		TextPrimary: resolver.ResolveVar(raw.Extra.Text.Primary, RawPaletteVars{raw}),
		TextSecondary: resolver.ResolveVar(raw.Extra.Text.Secondary, RawPaletteVars{raw}),
		TextMuted: resolver.ResolveVar(raw.Extra.Text.Muted, RawPaletteVars{raw}),
		TextLink: resolver.ResolveVar(raw.Extra.Text.Link, RawPaletteVars{raw}),

		LayerBase: resolver.ResolveVar(raw.Extra.Layer.Base, RawPaletteVars{raw}),
		LayerMantle: resolver.ResolveVar(raw.Extra.Layer.Mantle, RawPaletteVars{raw}),
		LayerCrust: resolver.ResolveVar(raw.Extra.Layer.Crust, RawPaletteVars{raw}),
		LayerSurface: resolver.ResolveVar(raw.Extra.Layer.Surface, RawPaletteVars{raw}),
		LayerSurfaceRaised: resolver.ResolveVar(raw.Extra.Layer.SurfaceRaised, RawPaletteVars{raw}),
		LayerSurfaceOverlay: resolver.ResolveVar(raw.Extra.Layer.SurfaceOverlay, RawPaletteVars{raw}),

		BorderDefault: resolver.ResolveVar(raw.Extra.Border.Default, RawPaletteVars{raw}),
		BorderActive: resolver.ResolveVar(raw.Extra.Border.Active, RawPaletteVars{raw}),

		StatusSuccess: resolver.ResolveVar(raw.Extra.Status.Success, RawPaletteVars{raw}),
		StatusWarning: resolver.ResolveVar(raw.Extra.Status.Warning, RawPaletteVars{raw}),
		StatusError: resolver.ResolveVar(raw.Extra.Status.Error, RawPaletteVars{raw}),
		StatusCritical: resolver.ResolveVar(raw.Extra.Status.Critical, RawPaletteVars{raw}),
		StatusInfo: resolver.ResolveVar(raw.Extra.Status.Info, RawPaletteVars{raw}),
	}
}
