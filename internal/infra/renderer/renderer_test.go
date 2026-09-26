package renderer

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRenderCreatesParentAndSkipsUnchangedWrite(t *testing.T) {
	dir := t.TempDir()
	templatePath := filepath.Join(dir, "sample.tmpl")
	outputPath := filepath.Join(dir, "nested", "out.txt")

	if err := os.WriteFile(templatePath, []byte("hello {{ .Name }}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	data := struct{ Name string }{Name: "theme"}
	if err := Render(templatePath, outputPath, data); err != nil {
		t.Fatal(err)
	}

	first, err := os.Stat(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	if err := Render(templatePath, outputPath, data); err != nil {
		t.Fatal(err)
	}

	second, err := os.Stat(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	if !first.ModTime().Equal(second.ModTime()) {
		t.Fatalf("unchanged render rewrote file: %s != %s", first.ModTime(), second.ModTime())
	}
}
