package palette

import "theme-engine/internal/resolver"

func fallback(val, fallback string) string {
	if val == "" {
		return fallback
	}
	return val
}

func ResolvePalette(raw *Palette) *ResolvedPalette {
	rp := &ResolvedPalette{
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

	s := raw.Extra.Syntax
	rp.SyntaxPurple   = fallback(resolver.ResolveVar(s.Purple, RawPaletteVars{raw}), rp.AccentPrimary)
	rp.SyntaxMagenta2 = fallback(resolver.ResolveVar(s.Magenta2, RawPaletteVars{raw}), rp.AccentSecondary)
	rp.SyntaxBlue0    = fallback(resolver.ResolveVar(s.Blue0, RawPaletteVars{raw}), rp.LayerSurface)
	rp.SyntaxBlue1    = fallback(resolver.ResolveVar(s.Blue1, RawPaletteVars{raw}), rp.TextLink)
	rp.SyntaxBlue5    = fallback(resolver.ResolveVar(s.Blue5, RawPaletteVars{raw}), raw.Colors[4])
	rp.SyntaxBlue6    = fallback(resolver.ResolveVar(s.Blue6, RawPaletteVars{raw}), raw.Colors[6])
	rp.SyntaxBlue7    = fallback(resolver.ResolveVar(s.Blue7, RawPaletteVars{raw}), rp.BorderDefault)
	rp.SyntaxGreen1   = fallback(resolver.ResolveVar(s.Green1, RawPaletteVars{raw}), rp.StatusSuccess)
	rp.SyntaxGreen2   = fallback(resolver.ResolveVar(s.Green2, RawPaletteVars{raw}), raw.Colors[2])
	rp.SyntaxOrange   = fallback(resolver.ResolveVar(s.Orange, RawPaletteVars{raw}), rp.StatusWarning)
	rp.SyntaxRed1     = fallback(resolver.ResolveVar(s.Red1, RawPaletteVars{raw}), rp.StatusCritical)
	rp.SyntaxTeal     = fallback(resolver.ResolveVar(s.Teal, RawPaletteVars{raw}), raw.Colors[6])

	rp.BgStatusline = fallback(resolver.ResolveVar(raw.Extra.UI.BgStatusline, RawPaletteVars{raw}), rp.LayerSurfaceRaised)

	return rp
}
