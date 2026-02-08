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

	return Foot{
		Font: ctx.Theme.Theme.Fonts.Primary,
		FontSize: inp.FontSize,
		PaddingX: inp.Padding.X,
		PaddingY: inp.Padding.Y,

		Palette: ctx.Palette,
	}, nil;
}

func (FootProcessor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data);
}

