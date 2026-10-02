package foot

import (
	"strconv"
	"strings"
	"testing"

	"theme-engine/internal/domain/palette"
)

func testPalette() *palette.ResolvedPalette {
	return &palette.ResolvedPalette{
		Foreground: "#c9d3dc",
		Background: "#161c22",
		Cursor:     "#8fb8d0",
		Colors: []string{
			"#10151a", "#d08f8f", "#93b89b", "#d9c08a",
			"#8fb8d0", "#c0a3c0", "#8cc0bb", "#c9d3dc",
			"#6b7885", "#e0a0a0", "#a8cdae", "#ecd9a8",
			"#a6cbe0", "#d4b8d4", "#a3d3cd", "#e2e9ef",
		},
	}
}

func TestOscSequenceContainsAllEntries(t *testing.T) {
	seq := oscSequence(testPalette())

	for _, want := range []string{
		"\x1b]10;#c9d3dc\x1b\\", // foreground
		"\x1b]11;#161c22\x1b\\", // background
		"\x1b]12;#8fb8d0\x1b\\", // cursor
		"\x1b]4;0;#10151a\x1b\\",
		"\x1b]4;4;#8fb8d0\x1b\\",
		"\x1b]4;15;#e2e9ef\x1b\\",
	} {
		if !strings.Contains(seq, want) {
			t.Errorf("sequence missing %q, got %q", want, seq)
		}
	}
}

func TestOscSequenceHasPaletteEntriesInOrder(t *testing.T) {
	p := testPalette()
	seq := oscSequence(p)

	for i, c := range p.Colors {
		want := "\x1b]4;" + strconv.Itoa(i) + ";" + c + "\x1b\\"
		if !strings.Contains(seq, want) {
			t.Errorf("palette entry %d missing %q", i, want)
		}
	}
}

func TestIsDigits(t *testing.T) {
	cases := map[string]bool{
		"0":     true,
		"123":   true,
		"":      false,
		"ptmx":  false,
		"1a":    false,
		"-1":    false,
	}
	for in, want := range cases {
		if got := isDigits(in); got != want {
			t.Errorf("isDigits(%q) = %v, want %v", in, got, want)
		}
	}
}
