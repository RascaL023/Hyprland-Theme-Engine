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

func TestRenderInPlaceWritesContent(t *testing.T) {
	dir := t.TempDir()
	templatePath := filepath.Join(dir, "sample.tmpl")
	outputPath := filepath.Join(dir, "nested", "out.txt")

	if err := os.WriteFile(templatePath, []byte("hello {{ .Name }}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	data := struct{ Name string }{Name: "theme"}
	if err := RenderInPlace(templatePath, outputPath, data); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello theme\n" {
		t.Fatalf("RenderInPlace content = %q", got)
	}

	// Update the template and render in place again:
	// content must change (no stale skip) and the file
	// must be overwritten directly, not replaced.
	if err := os.WriteFile(templatePath, []byte("bye {{ .Name }}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	ClearCache()

	if err := RenderInPlace(templatePath, outputPath, data); err != nil {
		t.Fatal(err)
	}

	got, err = os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "bye theme\n" {
		t.Fatalf("RenderInPlace after template edit = %q", got)
	}
}

func TestRenderInPlaceSkipsUnchanged(t *testing.T) {
	dir := t.TempDir()
	templatePath := filepath.Join(dir, "sample.tmpl")
	outputPath := filepath.Join(dir, "out.txt")

	if err := os.WriteFile(templatePath, []byte("x\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := RenderInPlace(templatePath, outputPath, nil); err != nil {
		t.Fatal(err)
	}

	first, err := os.Stat(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	if err := RenderInPlace(templatePath, outputPath, nil); err != nil {
		t.Fatal(err)
	}

	second, err := os.Stat(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	if !first.ModTime().Equal(second.ModTime()) {
		t.Fatalf("unchanged RenderInPlace rewrote file: %s != %s", first.ModTime(), second.ModTime())
	}
}

func TestClearCacheReparsesTemplate(t *testing.T) {
	dir := t.TempDir()
	templatePath := filepath.Join(dir, "sample.tmpl")
	outputPath := filepath.Join(dir, "out.txt")

	if err := os.WriteFile(templatePath, []byte("one\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Render(templatePath, outputPath, nil); err != nil {
		t.Fatal(err)
	}

	// Without ClearCache the parsed template is reused,
	// so an edited template would still render "one".
	if err := os.WriteFile(templatePath, []byte("two\n"), 0644); err != nil {
		t.Fatal(err)
	}
	ClearCache()

	if err := Render(templatePath, outputPath, nil); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "two\n" {
		t.Fatalf("after ClearCache content = %q, want %q", got, "two\n")
	}
}
