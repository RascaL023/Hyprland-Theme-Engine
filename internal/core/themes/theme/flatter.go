package theme

import "theme-engine/internal/core/themes/palette"

func BuildFlattenPalette(t *palette.ResolvedPalette) {
	t.Flat = map[string]string{
		"extra.foreground":              t.Foreground,
		"extra.background":              t.Background,
		"extra.cursor":                  t.Cursor,
		"extra.accent.primary":          t.AccentPrimary,
		"extra.accent.secondary":        t.AccentSecondary,
		"extra.accent.on_accent":        t.AccentOn,
		"extra.text.primary":            t.TextPrimary,
		"extra.text.secondary":          t.TextSecondary,
		"extra.text.muted":              t.TextMuted,
		"extra.text.link":               t.TextLink,
		"extra.layer.base":              t.LayerBase,
		"extra.layer.mantle":            t.LayerMantle,
		"extra.layer.crust":             t.LayerCrust,
		"extra.layer.surface":           t.LayerSurface,
		"extra.layer.surface_raised":    t.LayerSurfaceRaised,
		"extra.layer.surface_overlay":   t.LayerSurfaceOverlay,
		"extra.border.default":          t.BorderDefault,
		"extra.border.active":           t.BorderActive,
		"extra.status.success":          t.StatusSuccess,
		"extra.status.warning":          t.StatusWarning,
		"extra.status.error":            t.StatusError,
		"extra.status.critical":         t.StatusCritical,
		"extra.status.info":             t.StatusInfo,
	}

	for i, color := range t.Colors {
		t.Flat["color"+itoa(i)] = color
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}

	return string(buf[i:])
}
