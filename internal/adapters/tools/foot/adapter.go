package foot

import (
	"encoding/json"

	"theme-engine/internal/domain/renderctx"
	"theme-engine/internal/infra/renderer"
)

type Processor struct{}

func New() Processor { return Processor{} }

func (Processor) Name() string { return "foot" }

func (Processor) Parse(in any) (any, error) {
	var cfg Raw
	err := json.Unmarshal(in.(json.RawMessage), &cfg)
	return cfg, err
}

func (Processor) Resolve(in any, ctx *renderctx.Context) (any, error) {
	inp := in.(Raw)

	return Foot{
		Font:     ctx.Theme.Theme.Fonts.Terminal.Family,
		FontSize: ctx.Theme.Theme.Fonts.Terminal.Size,
		PaddingX: inp.Padding.X,
		PaddingY: inp.Padding.Y,

		Palette: ctx.Palette,
	}, nil
}

func (Processor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data)
}
