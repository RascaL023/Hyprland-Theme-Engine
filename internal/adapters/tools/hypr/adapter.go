package hypr

import (
	"encoding/json"
	"strings"

	"theme-engine/internal/domain/palette"
	"theme-engine/internal/domain/renderctx"
	"theme-engine/internal/infra/pathenv"
	"theme-engine/internal/infra/renderer"
)

type Processor struct{}

func New() Processor { return Processor{} }

func (Processor) Name() string { return "hypr" }

func (Processor) Parse(in any) (any, error) {
	var cfg Raw
	if in == nil {
		return cfg, nil
	}
	err := json.Unmarshal(in.(json.RawMessage), &cfg)
	return cfg, err
}

func resolveColor(s string, ctx *renderctx.Context) string {
	resolved := pathenv.ResolveVar(s, palette.ResolvedPaletteVars{P: ctx.Palette})
	if strings.HasPrefix(resolved, "#") {
		hex := resolved[1:]
		if len(hex) == 6 {
			return "rgba(" + hex + "ee)"
		}
	}
	return resolved
}

func (Processor) Resolve(in any, ctx *renderctx.Context) (any, error) {
	inp, _ := in.(Raw)
	return Hypr{
		GapsIn:  inp.Gaps.In,
		GapsOut: inp.Gaps.Out,

		BorderSize:            inp.Border.Size,
		BorderActiveDegree:    inp.Border.Active.Degree,
		BorderActivePrimary:   resolveColor(inp.Border.Active.Primary, ctx),
		BorderActiveSecondary: resolveColor(inp.Border.Active.Secondary, ctx),
		BorderInActivePrimary: resolveColor(inp.Border.Inactive.Primary, ctx),

		BlurOption: inp.Blur.Option,
		BlurSize:   inp.Blur.Size,
		BlurPasses: inp.Blur.Passes,

		Rounding: inp.Rounding,

		ShadowOption:      inp.Shadow.Option,
		ShadowRange:       inp.Shadow.Range,
		ShadowRenderPower: inp.Shadow.RenderPower,
		ShadowColor:       resolveColor(inp.Shadow.Color, ctx),

		WindowRuleFootOpacity: inp.WindowRule.Foot.Opacity,
	}, nil
}

func (Processor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data)
}
