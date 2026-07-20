package gtk

import (
	"bytes"
	"os/exec"
	"path/filepath"

	"theme-engine/internal/core/context"
	"theme-engine/internal/processor"
	"theme-engine/internal/renderer"
)

type GtkProcessor struct{}

func (GtkProcessor) Name() string { return "gtk" }

func init() { processor.RegisterProcessor(GtkProcessor{}) }

func (GtkProcessor) Parse(_ any) (any, error) {
	return nil, nil
}

func (GtkProcessor) Resolve(_ any, ctx *context.Context) (any, error) {
	return Gtk{
		Config: ctx,
	}, nil
}

func compileSCSS(src, dst string, loadPaths []string) error {
	args := append([]string{src, dst}, "-I", src)
	for _, lp := range loadPaths {
		args = append(args, "-I", lp)
	}
	cmd := exec.Command("sassc", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil { return err }

	return nil
}

func (GtkProcessor) Render(
	templatePath,
	outputPath string,
	data any,
) error {
	assetsDir := filepath.Dir(templatePath)

	err := renderer.Render(
		templatePath,
		outputPath+"/scss/source/_source.scss",
		data,
	); if err != nil { return err }

	loadPaths := []string{
		outputPath + "/scss/source",
		assetsDir,
	}

	err = compileSCSS(assetsDir+"/base.scss", outputPath+"/css/source.css", loadPaths)
	if err != nil { return err }

	return compileSCSS(assetsDir+"/rofi-base.scss", outputPath+"/rasi/source.rasi", loadPaths)
}
