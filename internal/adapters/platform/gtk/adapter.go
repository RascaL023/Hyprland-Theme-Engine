package gtk

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
		Context: ctx,
		Config:  ctx,
	}, nil
}

func compileSCSS(src, dst string, loadPaths []string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	args := append([]string{src, dst}, "-I", src)
	for _, lp := range loadPaths {
		args = append(args, "-I", lp)
	}
	cmd := exec.Command("sassc", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sassc %s -> %s: %w: %s", src, dst, err, strings.TrimSpace(stderr.String()))
	}

	return nil
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
