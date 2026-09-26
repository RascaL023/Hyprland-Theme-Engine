package gtk

import (
	"bytes"
	"os/exec"
	"path/filepath"

	"theme-engine/internal/domain/renderctx"
	"theme-engine/internal/infra/renderer"
)

type Processor struct{}

func New() Processor { return Processor{} }

func (Processor) Name() string { return "gtk" }

func (Processor) Parse(_ any) (any, error) {
	return nil, nil
}

func (Processor) Resolve(_ any, ctx *renderctx.Context) (any, error) {
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

	return cmd.Run()
}

func (Processor) Render(
	templatePath,
	outputPath string,
	data any,
) error {
	assetsDir := filepath.Dir(templatePath)

	if err := renderer.Render(
		templatePath,
		outputPath+"/scss/source/_source.scss",
		data,
	); err != nil {
		return err
	}

	loadPaths := []string{
		outputPath + "/scss/source",
		assetsDir,
	}

	if err := compileSCSS(assetsDir+"/base.scss", outputPath+"/css/source.css", loadPaths); err != nil {
		return err
	}

	return compileSCSS(assetsDir+"/rofi-base.scss", outputPath+"/rasi/source.rasi", loadPaths)
}
