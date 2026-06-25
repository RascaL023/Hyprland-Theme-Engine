package hypr

import (
	"encoding/json"
	"strings"
	"theme-engine/internal/core/context"
	"theme-engine/internal/core/themes/palette"
	"theme-engine/internal/processor"
	"theme-engine/internal/renderer"
	"theme-engine/internal/resolver"
)

type HyprProcessor struct{}

func init() { processor.RegisterProcessor(HyprProcessor{}); }

func (HyprProcessor) Name() string { return "hypr"; }

func (HyprProcessor) Parse(in any) (any, error) {
	var cfg Raw;
	if in == nil { return cfg, nil; }
	err := json.Unmarshal(in.(json.RawMessage), &cfg);
	return cfg, err;
}

func resolveColor(s string, ctx *context.Context) string {
	resolved := resolver.ResolveVar(s, palette.ResolvedPaletteVars{P: ctx.Palette})
	if strings.HasPrefix(resolved, "#") {
		hex := resolved[1:]
		if len(hex) == 6 {
			return "rgba(" + hex + "ee)"
		}
	}
	return resolved
}

func (HyprProcessor) Resolve(in any, ctx *context.Context) (any, error) {
	inp, _ := in.(Raw);
	return Hypr{
		GapsIn: inp.Gaps.In,
		GapsOut: inp.Gaps.Out,

		BorderSize: inp.Border.Size,
		BorderActiveDegree: inp.Border.Active.Degree,
		BorderActivePrimary: resolveColor(inp.Border.Active.Primary, ctx),
		BorderActiveSecondary: resolveColor(inp.Border.Active.Secondary, ctx),
		BorderInActivePrimary: resolveColor(inp.Border.Inactive.Primary, ctx),

		BlurOption: inp.Blur.Option,
		BlurSize: inp.Blur.Size,
		BlurPasses: inp.Blur.Passes,

		Rounding: inp.Rounding,

		ShadowOption: inp.Shadow.Option,
		ShadowRange: inp.Shadow.Range,
		ShadowRenderPower: inp.Shadow.RenderPower,
		ShadowColor: resolveColor(inp.Shadow.Color, ctx),

		WindowRuleFootOpacity: inp.WindowRule.Foot.Opacity,
	}, nil;
}

func (HyprProcessor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data);
}

