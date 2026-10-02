package foot

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"theme-engine/internal/domain/palette"
)

var (
	setTermScript      = filepath.Join("..", "..", "..", "..", "tools", "set_term")
	setTermShellScript = filepath.Join("..", "..", "..", "..", "tools", "set_term.sh")
	themesDir          = filepath.Join("..", "..", "..", "..", "themes")
)

// setTermThemes mirrors the contrast guard's discovery: any
// directory with a palette.json is a shipped theme and gets
// asserted automatically — a new theme cannot silently skip
// this guard.
func setTermThemes(t *testing.T) []string {
	t.Helper()

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		t.Fatalf("read themes dir: %v", err)
	}

	var themes []string
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasSuffix(entry.Name(), ".bak") {
			continue
		}
		if _, err := os.Stat(filepath.Join(themesDir, entry.Name(), "palette.json")); err != nil {
			continue
		}
		themes = append(themes, entry.Name())
	}

	if len(themes) < 2 {
		t.Fatalf("expected at least 2 active themes, found %d", len(themes))
	}
	return themes
}

// assertSetTermMatchesEngine runs a standalone set_term
// implementation in --dry-run mode for every theme and
// variant and asserts its output is byte-identical to
// oscSequence over the palette resolved through the same
// path the renderer uses. The explicit theme and variant
// mean no .state.json is involved.
func assertSetTermMatchesEngine(t *testing.T, argv []string) {
	t.Helper()

	for _, theme := range setTermThemes(t) {
		for _, variant := range []string{"dark", "light"} {
			t.Run(theme+"/"+variant, func(t *testing.T) {
				raw, err := os.ReadFile(filepath.Join(themesDir, theme, "palette.json"))
				if err != nil {
					t.Fatalf("read palette: %v", err)
				}

				var p palette.Raw
				if err := json.Unmarshal(raw, &p); err != nil {
					t.Fatalf("parse palette: %v", err)
				}

				resolved, err := p.ResolveSelected(variant)
				if err != nil {
					t.Fatalf("resolve palette: %v", err)
				}

				cmd := append(append([]string{}, argv...),
					theme, "--variant", variant, "--dry-run")
				out, err := exec.Command(cmd[0], cmd[1:]...).Output()
				if err != nil {
					t.Fatalf("run %v: %v", argv, err)
				}

				if want := oscSequence(resolved); string(out) != want {
					t.Errorf("%s diverges from engine\ngot:  %q\nwant: %q", argv[0], out, want)
				}
			})
		}
	}
}

// TestSetTermDryRunMatchesEngine keeps tools/set_term — the
// standalone form of Reload — byte-identical to the engine's
// OSC output.
func TestSetTermDryRunMatchesEngine(t *testing.T) {
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available: tools/set_term differential skipped")
	}
	assertSetTermMatchesEngine(t, []string{"python3", setTermScript})
}

// TestSetTermShellDryRunMatchesEngine does the same for the
// bash twin (tools/set_term.sh), which parses the palette
// with awk/sed instead of python.
func TestSetTermShellDryRunMatchesEngine(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available: tools/set_term.sh differential skipped")
	}
	assertSetTermMatchesEngine(t, []string{"bash", setTermShellScript})
}
