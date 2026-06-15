package theme

import "theme-engine/internal/core/themes/palette"

func BuildFlattenPalette(t *palette.ResolvedPalette) {
	t.Flat = map[string]string{
		"extra.foreground":       t.Foreground,
		"extra.background":       t.Background,
		"extra.cursor":           t.Cursor,
		"extra.primaryaccent":    t.PrimaryAccent,
		"extra.secondaryaccent":  t.SecondaryAccent,
		"extra.primarytext":      t.PrimaryText,
		"extra.secondarytext":    t.SecondaryText,
		"extra.teritarytext":     t.TeritaryText,
		"extra.primaryoverlay":   t.PrimaryOverlay,
		"extra.secondaryoverlay": t.SecondaryOverlay,
		"extra.teritaryoverlay":  t.TeritaryOverlay,
		"extra.primarysurface":   t.PrimarySurface,
		"extra.secondarysurface": t.SecondarySurface,
		"extra.teritarysurface":  t.TeritarySurface,
		"extra.base":             t.Base,
		"extra.mantle":           t.Mantle,
		"extra.crust":            t.Crust,
		"extra.warning":          t.Warning,
		"extra.critical":         t.Critical,
		"extra.charging":         t.Charging,
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
