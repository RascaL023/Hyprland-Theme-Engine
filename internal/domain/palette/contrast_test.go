package palette_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"

	"theme-engine/internal/domain/palette"
)

// contrastPilotThemes is the seed list for the accessibility guard.
//
// Only the pilot theme is asserted for now (VISUAL_FIX_PLAN.md Fase 0-2);
// the remaining themes are migrated and added in Fase 5. New themes should be
// appended here once their palette passes.
var contrastPilotThemes = []string{"kanagawa-dragon"}

var contrastVariants = []string{"dark", "light"}

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

func TestContrastPilotThemes(t *testing.T) {
	for _, theme := range contrastPilotThemes {
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
