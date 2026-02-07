package foot

import (
	"encoding/json"
	"theme-engine/internal/core/context"
	"theme-engine/internal/processor"
	"theme-engine/internal/renderer"
)

type FootProcessor struct{}

func init() { processor.RegisterProcessor(FootProcessor{}); }

func (FootProcessor) Name() string { return "foot"; }

func (FootProcessor) Parse(in any) (any, error) {
	var cfg Raw;
	err := json.Unmarshal(in.(json.RawMessage), &cfg);
	return cfg, err;
}

func (FootProcessor) Resolve(in any, ctx *context.Context) (any, error) {
	inp := in.(Raw);

	cfg := Foot{
		Font: ctx.Theme.Theme.Fonts.Primary,
		FontSize: inp.FontSize,
		PaddingX: inp.Padding.X,
		PaddingY: inp.Padding.Y,

		// Foreground: ctx.Palette.Foreground,
		// Background: ctx.Palette.Background,
		//
		// Colors: [16]string(ctx.Palette.Colors),
	}
	var temp string;

	temp = ctx.Palette.Foreground;
	cfg.Foreground = temp[1:];

	temp = ctx.Palette.Background;
	cfg.Background = temp[1:];

	for idx, temp := range ctx.Palette.Colors {
		cfg.Colors[idx] = temp[1:];
	}

	return cfg, nil;
}

func (FootProcessor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data);
}

