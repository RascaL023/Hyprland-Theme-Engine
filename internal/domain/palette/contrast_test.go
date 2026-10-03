package palette_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"theme-engine/internal/domain/palette"
)

// contrastVariants are the variants every active theme must ship.
var contrastVariants = []string{"dark", "light"}

// activeThemes discovers the themes the guard covers.
//
// Discovery (instead of a hand-kept list) means a new theme cannot silently
// skip the guard: dropping themes/<name>/palette.json in is enough to be
// asserted. "*.bak" directories are backups, not shipped themes.
func activeThemes(t *testing.T) []string {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join("..", "..", "..", "themes"))
	if err != nil {
		t.Fatalf("read themes dir: %v", err)
	}

	var themes []string
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasSuffix(entry.Name(), ".bak") {
			continue
		}
		if _, err := os.Stat(filepath.Join("..", "..", "..", "themes", entry.Name(), "palette.json")); err != nil {
			continue
		}
		themes = append(themes, entry.Name())
	}

	if len(themes) < 2 {
		t.Fatalf("expected at least 2 active themes, found %d: %v", len(themes), themes)
	}
	return themes
}

const (
	// minUIRatio is the WCAG threshold for UI accents and large text.
	minUIRatio = 3.0
	// minTextRatio is the WCAG threshold for normal body text.
	minTextRatio = 4.5
)

type contrastCheck struct {
	name   string
	fg     string
	bg     string
	min    float64
	reason string
}

// resolveTheme loads themes/<theme>/palette.json and runs it through the real
// resolve + flatten path so the test exercises the same values the renderer
// sees (including fallbacks).
func resolveTheme(t *testing.T, theme, variant string) *palette.ResolvedPalette {
	t.Helper()

	path := filepath.Join("..", "..", "..", "themes", theme, "palette.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	var p palette.Raw
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	rp, err := p.ResolveSelected(variant)
	if err != nil {
		t.Fatalf("resolve %s/%s: %v", theme, variant, err)
	}
	palette.BuildFlattenPalette(rp)
	return rp
}

func TestContrastThemes(t *testing.T) {
	for _, theme := range activeThemes(t) {
		for _, variant := range contrastVariants {
			t.Run(theme+"/"+variant, func(t *testing.T) {
				rp := resolveTheme(t, theme, variant)

				checks := []contrastCheck{
					{
						name:   "nvim inline code (backtick)",
						fg:     rp.Colors[4],
						bg:     rp.SyntaxTerminalBlack,
						min:    minUIRatio,
						reason: "tokyonight @markup.raw.markdown_inline = { bg = terminal_black, fg = blue }",
					},
					{
						name:   "editor main text",
						fg:     rp.TextPrimary,
						bg:     rp.LayerBase,
						min:    minTextRatio,
						reason: "RULE.md: text.primary must contrast with layer.base",
					},
					{
						name:   "lazygit selected line",
						fg:     rp.TextPrimary,
						bg:     rp.LayerSurfaceOverlay,
						min:    minUIRatio,
						reason: "lazygit draws defaultFgColor over selectedLineBgColor",
					},
				}

				// lualine section B draws a mode colour over the gutter surface.
				for _, m := range []struct{ name, hex string }{
					{"red", rp.Colors[1]},
					{"green", rp.Colors[2]},
					{"yellow", rp.Colors[3]},
					{"blue", rp.Colors[4]},
					{"magenta", rp.Colors[5]},
					{"green1", rp.SyntaxGreen1},
				} {
					checks = append(checks, contrastCheck{
						name:   "lualine section B " + m.name,
						fg:     m.hex,
						bg:     rp.BgGutter,
						min:    minUIRatio,
						reason: "tokyonight lualine b = { bg = fg_gutter, fg = mode colour }",
					})
				}

				for _, c := range checks {
					if c.fg == "" || c.bg == "" {
						t.Errorf("%s: empty colour (fg=%q bg=%q)", c.name, c.fg, c.bg)
						continue
					}
					got := palette.ContrastRatio(c.fg, c.bg)
					t.Logf("%-28s %.2f:1 (min %.1f:1) fg=%s bg=%s", c.name, got, c.min, c.fg, c.bg)
					if got < c.min {
						t.Errorf("%s: contrast %.2f:1 < %.1f:1 (fg=%s bg=%s) — %s",
							c.name, got, c.min, c.fg, c.bg, c.reason)
					}
				}
			})
		}
	}
}

func TestContrastRatio(t *testing.T) {
	cases := []struct {
		name string
		a    string
		b    string
		want float64
		tol  float64
	}{
		{"black on white is maximal", "#000000", "#ffffff", 21, 0.01},
		{"identical colours are minimal", "#8ba4b0", "#8ba4b0", 1, 0.001},
		{"shorthand equals longhand", "#abc", "#aabbcc", 1, 0.001},
		{"uppercase is case-insensitive", "#ABC", "#abc", 1, 0.001},
		{"invalid colour yields zero", "not-a-colour", "#ffffff", 0, 0.001},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := palette.ContrastRatio(c.a, c.b)
			if math.Abs(got-c.want) > c.tol {
				t.Errorf("ContrastRatio(%q, %q) = %.4f, want %.4f", c.a, c.b, got, c.want)
			}
		})
	}
}
