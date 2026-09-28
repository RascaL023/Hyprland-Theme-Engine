package palette

import (
	"fmt"

	"theme-engine/internal/infra/pathenv"
)

func fallback(val, fallback string) string {
	if val == "" {
		return fallback
	}
	return val
}

// maxColorIndex is the highest raw.Colors index used by syntax fallbacks
// (Colors[6]). Checked once up front: a single int compare, no per-field cost.
const maxColorIndex = 6

func ResolvePalette(raw *Palette) (*ResolvedPalette, error) {
	if len(raw.Colors) <= maxColorIndex {
		return nil, fmt.Errorf("palette needs at least %d colors for syntax fallbacks, got %d",
			maxColorIndex+1, len(raw.Colors))
	}
	rp := &ResolvedPalette{
		Foreground: raw.Foreground,
		Background: raw.Background,
		Cursor:     raw.Cursor,

		Colors: raw.Colors,

		AccentPrimary:   pathenv.ResolveVar(raw.Extra.Accent.Primary, RawPaletteVars{raw}),
		AccentSecondary: pathenv.ResolveVar(raw.Extra.Accent.Secondary, RawPaletteVars{raw}),
		AccentOn:        pathenv.ResolveVar(raw.Extra.Accent.OnAccent, RawPaletteVars{raw}),

		TextPrimary:   pathenv.ResolveVar(raw.Extra.Text.Primary, RawPaletteVars{raw}),
		TextSecondary: pathenv.ResolveVar(raw.Extra.Text.Secondary, RawPaletteVars{raw}),
		TextMuted:     pathenv.ResolveVar(raw.Extra.Text.Muted, RawPaletteVars{raw}),
		TextLink:      pathenv.ResolveVar(raw.Extra.Text.Link, RawPaletteVars{raw}),

		LayerBase:           pathenv.ResolveVar(raw.Extra.Layer.Base, RawPaletteVars{raw}),
		LayerMantle:         pathenv.ResolveVar(raw.Extra.Layer.Mantle, RawPaletteVars{raw}),
		LayerCrust:          pathenv.ResolveVar(raw.Extra.Layer.Crust, RawPaletteVars{raw}),
		LayerSurface:        pathenv.ResolveVar(raw.Extra.Layer.Surface, RawPaletteVars{raw}),
		LayerSurfaceRaised:  pathenv.ResolveVar(raw.Extra.Layer.SurfaceRaised, RawPaletteVars{raw}),
		LayerSurfaceOverlay: pathenv.ResolveVar(raw.Extra.Layer.SurfaceOverlay, RawPaletteVars{raw}),

		BorderDefault: pathenv.ResolveVar(raw.Extra.Border.Default, RawPaletteVars{raw}),
		BorderActive:  pathenv.ResolveVar(raw.Extra.Border.Active, RawPaletteVars{raw}),

		StatusSuccess:  pathenv.ResolveVar(raw.Extra.Status.Success, RawPaletteVars{raw}),
		StatusWarning:  pathenv.ResolveVar(raw.Extra.Status.Warning, RawPaletteVars{raw}),
		StatusError:    pathenv.ResolveVar(raw.Extra.Status.Error, RawPaletteVars{raw}),
		StatusCritical: pathenv.ResolveVar(raw.Extra.Status.Critical, RawPaletteVars{raw}),
		StatusInfo:     pathenv.ResolveVar(raw.Extra.Status.Info, RawPaletteVars{raw}),

		FgDim:      pathenv.ResolveVar(raw.Extra.Fg.Dim, RawPaletteVars{raw}),
		FgDimMuted: pathenv.ResolveVar(raw.Extra.Fg.DimMuted, RawPaletteVars{raw}),
		FgDisabled: pathenv.ResolveVar(raw.Extra.Fg.Disabled, RawPaletteVars{raw}),

		BgConflict:  pathenv.ResolveVar(raw.Extra.Bg.Conflict, RawPaletteVars{raw}),
		BgDiskUsage: pathenv.ResolveVar(raw.Extra.Bg.DiskUsage, RawPaletteVars{raw}),
	}

	s := raw.Extra.Syntax
	rp.SyntaxPurple = fallback(pathenv.ResolveVar(s.Purple, RawPaletteVars{raw}), rp.AccentPrimary)
	rp.SyntaxMagenta2 = fallback(pathenv.ResolveVar(s.Magenta2, RawPaletteVars{raw}), rp.AccentSecondary)
	rp.SyntaxBlue0 = fallback(pathenv.ResolveVar(s.Blue0, RawPaletteVars{raw}), rp.LayerSurface)
	rp.SyntaxBlue1 = fallback(pathenv.ResolveVar(s.Blue1, RawPaletteVars{raw}), rp.TextLink)
	rp.SyntaxBlue5 = fallback(pathenv.ResolveVar(s.Blue5, RawPaletteVars{raw}), raw.Colors[4])
	rp.SyntaxBlue6 = fallback(pathenv.ResolveVar(s.Blue6, RawPaletteVars{raw}), raw.Colors[6])
	rp.SyntaxBlue7 = fallback(pathenv.ResolveVar(s.Blue7, RawPaletteVars{raw}), rp.BorderDefault)
	rp.SyntaxGreen1 = fallback(pathenv.ResolveVar(s.Green1, RawPaletteVars{raw}), rp.StatusSuccess)
	rp.SyntaxGreen2 = fallback(pathenv.ResolveVar(s.Green2, RawPaletteVars{raw}), raw.Colors[2])
	rp.SyntaxOrange = fallback(pathenv.ResolveVar(s.Orange, RawPaletteVars{raw}), rp.StatusWarning)
	rp.SyntaxRed1 = fallback(pathenv.ResolveVar(s.Red1, RawPaletteVars{raw}), rp.StatusCritical)
	rp.SyntaxTeal = fallback(pathenv.ResolveVar(s.Teal, RawPaletteVars{raw}), raw.Colors[6])

	// Tokyonight reuses `terminal_black` as the inline-code background
	// (@markup.raw.markdown_inline), so it needs a subtle dark surface, not the
	// ANSI bright-black colour. Fallback keeps older themes unchanged.
	rp.SyntaxTerminalBlack = fallback(pathenv.ResolveVar(s.TerminalBlack, RawPaletteVars{raw}), rp.LayerSurfaceOverlay)

	rp.BgStatusline = fallback(pathenv.ResolveVar(raw.Extra.UI.BgStatusline, RawPaletteVars{raw}), rp.LayerSurfaceRaised)

	// Tokyonight's fg_gutter backs lualine section B and the gutter UI, so it is
	// a subtle surface of its own rather than the brightest layer.
	rp.BgGutter = fallback(pathenv.ResolveVar(raw.Extra.UI.Gutter, RawPaletteVars{raw}), rp.LayerSurfaceOverlay)

	rp.LinkVisited = fallback(pathenv.ResolveVar(raw.Extra.Text.Visited, RawPaletteVars{raw}), rp.AccentSecondary)
	rp.BorderMedium = fallback(pathenv.ResolveVar(raw.Extra.Border.Medium, RawPaletteVars{raw}), rp.BorderDefault)

	rp.FgDim = fallback(rp.FgDim, rp.TextMuted)
	rp.FgDimMuted = fallback(rp.FgDimMuted, rp.TextMuted)
	rp.FgDisabled = fallback(rp.FgDisabled, rp.TextMuted)
	rp.BgConflict = fallback(rp.BgConflict, rp.StatusWarning)
	rp.BgDiskUsage = fallback(rp.BgDiskUsage, rp.LayerSurfaceOverlay)

	return rp, nil
}
