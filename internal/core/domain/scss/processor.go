package scss

import (
	"fmt"
	"theme-engine/internal/core/context"
	"theme-engine/internal/processor"
	"theme-engine/internal/renderer"
)

type ScssProcessor struct {}

func (ScssProcessor) Name() string { return "scss"; }

func init() { processor.RegisterProcessor(ScssProcessor{}); }

func (ScssProcessor) Parse(_ any) (any, error) {
	return nil, nil;
}

func (ScssProcessor) Resolve(_ any, ctx *context.Context) (any, error) {
	return Scss{
		Accent1: ctx.Palette.PrimaryAccent,
		Accent2: ctx.Palette.SecondaryAccent,

		Text: ctx.Palette.PrimaryText,
		Subtext1: ctx.Palette.SecondaryText,
		Subtext0: ctx.Palette.TeritaryText,

		Overlay2: ctx.Palette.PrimaryOverlay,
		Overlay1: ctx.Palette.SecondaryOverlay,
		Overlay0: ctx.Palette.TeritaryOverlay,

		Surface2: ctx.Palette.PrimarySurface,
		Surface1: ctx.Palette.SecondarySurface,
		Surface0: ctx.Palette.TeritarySurface,

		Base: ctx.Palette.Base,
		Crust: ctx.Palette.Crust,
		Mantle: ctx.Palette.Mantle,

		Warning: ctx.Palette.Warning,
		Critical: ctx.Palette.Critical,
		Charging: ctx.Palette.Charging,
	}, nil;
}

func (ScssProcessor) Render(
	templatePath, 
	outputPath string,
	data any,
) error {
	err := renderer.Render(
		templatePath, 
		outputPath, 
		data,
	);

	if err != nil {
		fmt.Println("[ERROR]")
	}

	return err;
}
