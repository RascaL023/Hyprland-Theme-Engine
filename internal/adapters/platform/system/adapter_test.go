package system

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestReloadRunsRenderedScript(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "applied")
	script := filepath.Join(dir, "apply.sh")

	// A stand-in for the rendered dconf script: if
	// Reload executes it, the marker file appears.
	content := "#!/usr/bin/env bash\nset -e\ntouch " + strconv.Quote(marker) + "\n"
	if err := os.WriteFile(script, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := (Processor{}).Reload(script, nil); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("rendered script was not executed: %v", err)
	}
}

func TestReloadMissingOutput(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "apply.sh")

	if err := (Processor{}).Reload(missing, nil); err == nil {
		t.Fatal("want error when rendered script does not exist")
	}
}
