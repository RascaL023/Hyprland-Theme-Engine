package gtk

import (
	"bytes"
	"os/exec"
	"theme-engine/internal/core/context"
	"theme-engine/internal/processor"
	"theme-engine/internal/renderer"
)

type GtkProcessor struct {}

func (GtkProcessor) Name() string { return "gtk"; }

func init() { processor.RegisterProcessor(GtkProcessor{}); }

func (GtkProcessor) Parse(_ any) (any, error) {
	return nil, nil;
}

func (GtkProcessor) Resolve(_ any, ctx *context.Context) (any, error) {
	return Gtk{
		Config: ctx,
	}, nil;
}

func compileSCSS(templatePath, outputPath string) error {
	cmd := exec.Command("sassc", templatePath, outputPath);

	var stderr bytes.Buffer;
	cmd.Stderr = &stderr;

	err := cmd.Run();
	if err != nil {
		return err;
	}

	return nil;
}

func (GtkProcessor) Render(
	templatePath, 
	outputPath string,
	data any,
) error {
	err := renderer.Render(
		templatePath, 
		outputPath + "/scss/source/_source.scss", 
		data,
	);
	if err != nil { return err; }
	
	err = compileSCSS(outputPath + "/scss/base/base.scss", outputPath + "/css/source.css");
	if err != nil { return err; }

	return compileSCSS(outputPath + "/scss/base/rofi-base.scss", outputPath + "/rasi/source.rasi");
}
