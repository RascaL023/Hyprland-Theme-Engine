package loader

import (
	"os"
	"path/filepath"
	"testing"

	"theme-engine/internal/core/themes/state"
)

func TestLoadToolMapExpandsStateAndSkipsComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "path.txt")

	err := os.WriteFile(path, []byte(`
# local targets
waybar|assets/templates/waybar/$WAYBAR.tmpl|output/waybar/sources.css
foot|assets/templates/tools/foot/foot.tmpl|output/tools/foot/ui.ini
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	st := &state.State{Waybar: "default"}
	got, err := LoadToolMap(path, st)
	if err != nil {
		t.Fatal(err)
	}

	if got["waybar"].TemplatePath != "assets/templates/waybar/default.tmpl" {
		t.Fatalf("unexpected waybar template: %q", got["waybar"].TemplatePath)
	}
	if got["foot"].OutputPath != "output/tools/foot/ui.ini" {
		t.Fatalf("unexpected foot output: %q", got["foot"].OutputPath)
	}
}

func TestLoadToolMapRejectsMalformedLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "path.txt")

	if err := os.WriteFile(path, []byte("broken\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadToolMap(path, &state.State{})
	if err == nil {
		t.Fatal("expected malformed path map error")
	}
}
