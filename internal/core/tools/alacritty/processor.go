package alacritty

import (
	"encoding/json"
	"theme-engine/internal/core/context"
	"theme-engine/internal/processor"
	"theme-engine/internal/renderer"
)

type AlacrittyProcessor struct{}

func init() { processor.RegisterProcessor(AlacrittyProcessor{}); }

func (AlacrittyProcessor) Name() string { return "alacritty"; }

func (AlacrittyProcessor) Parse(in any) (any, error) {
	var cfg Raw
	if in == nil {
		return cfg, nil
	}
	err := json.Unmarshal(in.(json.RawMessage), &cfg)
	return cfg, err
}

func (AlacrittyProcessor) Resolve(in any, ctx *context.Context) (any, error) {
	inp, ok := in.(Raw)
	if !ok {
		inp = Raw{}
	}

	opacity := inp.Opacity
	if opacity == 0 {
		opacity = 1.0
	}

	cursorShape := inp.Cursor.Style.Shape
	if cursorShape == "" {
		cursorShape = "Beam"
	}

	return Alacritty{
		Palette: ctx.Palette,

		Opacity:  opacity,
		PaddingX: inp.Padding.X,
		PaddingY: inp.Padding.Y,

		CursorShape: cursorShape,

		FontSize:   ctx.Theme.Theme.Fonts.Terminal.Size,
		FontFamily: ctx.Theme.Theme.Fonts.Terminal.Family,
	}, nil
}

func (AlacrittyProcessor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data)
}
