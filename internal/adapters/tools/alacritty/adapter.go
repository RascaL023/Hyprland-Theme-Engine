package alacritty

import (
	"theme-engine/internal/adapters/tools/jsonx"
	"theme-engine/internal/domain/renderctx"
	"theme-engine/internal/infra/renderer"
)

type Processor struct{}

func New() Processor { return Processor{} }

func (Processor) Parse(in any) (any, error) {
	var cfg Raw
	if err := jsonx.Decode(in, &cfg); err != nil {
		return Raw{}, err
	}
	return cfg, nil
}

func (Processor) Resolve(in any, ctx *renderctx.Context) (any, error) {
	inp, _ := in.(Raw)

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

// Render writes the config in place: Alacritty live-reloads
// its config (live_config_reload, on by default) via inotify
// on the file's inode, which a temp+rename write breaks.
func (Processor) Render(templatePath, outputPath string, data any) error {
	return renderer.RenderInPlace(templatePath, outputPath, data)
}
