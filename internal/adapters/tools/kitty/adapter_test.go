package kitty

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractColourLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ui.conf")

	content := `##########
# Colors #
##########

foreground                #c9d3dc
background                #161c22
cursor                    #8fb8d0

selection_background      #232b34
selection_foreground      #c9d3dc

color0                    #10151a
color15                   #e2e9ef

active_tab_foreground     #c9d3dc
inactive_tab_background   #161c22

active_border_color       #8fb8d0
inactive_border_color     #2d3742


##########
# Layout #
##########

cursor_shape          beam
background_opacity    0.95

tab_bar_style         powerline

window_border_width   2

font_size        11

font_family      MapleMono NF
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := extractColourLines(path)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(got), "\n")
	want := []string{
		"foreground                #c9d3dc",
		"background                #161c22",
		"cursor                    #8fb8d0",
		"selection_background      #232b34",
		"selection_foreground      #c9d3dc",
		"color0                    #10151a",
		"color15                   #e2e9ef",
		"active_tab_foreground     #c9d3dc",
		"inactive_tab_background   #161c22",
		"active_border_color       #8fb8d0",
		"inactive_border_color     #2d3742",
	}

	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d:\n%s", len(lines), len(want), got)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}

	for _, excluded := range []string{"cursor_shape", "background_opacity", "tab_bar_style", "window_border_width", "font_size", "font_family"} {
		if strings.Contains(got, excluded) {
			t.Errorf("non-colour key %q leaked into extraction", excluded)
		}
	}
}

func TestIsColourKey(t *testing.T) {
	colours := []string{
		"foreground", "background", "cursor", "cursor_text_color",
		"selection_foreground", "selection_background",
		"active_tab_foreground", "active_tab_background",
		"inactive_tab_foreground", "inactive_tab_background",
		"active_border_color", "inactive_border_color",
		"color0", "color9", "color15",
	}
	for _, k := range colours {
		if !isColourKey(k) {
			t.Errorf("isColourKey(%q) = false, want true", k)
		}
	}

	nonColours := []string{
		"cursor_shape", "background_opacity", "tab_bar_style",
		"tab_powerline_style", "window_border_width",
		"window_margin_width", "window_padding_width",
		"font_size", "font_family", "bold_font",
		"colors", "color", "colorx", "color1x",
	}
	for _, k := range nonColours {
		if isColourKey(k) {
			t.Errorf("isColourKey(%q) = true, want false", k)
		}
	}
}
