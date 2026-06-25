package hypr

import (
	"encoding/json"
	"theme-engine/internal/core/context"
	"theme-engine/internal/processor"
	"theme-engine/internal/renderer"
)

type HyprProcessor struct{}

func init() { processor.RegisterProcessor(HyprProcessor{}); }

func (HyprProcessor) Name() string { return "hypr"; }

func (HyprProcessor) Parse(in any) (any, error) {
	var cfg Raw;
	err := json.Unmarshal(in.(json.RawMessage), &cfg);
	return cfg, err;
}

func (HyprProcessor) Resolve(in any, ctx *context.Context) (any, error) {
	inp := in.(Raw);
	return Hypr{
		GapsIn: inp.Gaps.In,
		GapsOut: inp.Gaps.Out,

		BorderSize: inp.Border.Size,
		BorderActiveDegree: inp.Border.Active.Degree,
		BorderActivePrimary: inp.Border.Active.Primary,
		BorderActiveSecondary: inp.Border.Active.Secondary,
		BorderInActivePrimary: inp.Border.Inactive.Primary,

		BlurOption: inp.Blur.Option,
		BlurSize: inp.Blur.Size,
		BlurPasses: inp.Blur.Passes,

		Rounding: inp.Rounding,

		ShadowOption: inp.Shadow.Option,
		ShadowRange: inp.Shadow.Range,
		ShadowRenderPower: inp.Shadow.RenderPower,
		ShadowColor: inp.Shadow.Color,

		WindowRuleFootOpacity: inp.WindowRule.Foot.Opacity,
	}, nil;
}

func (HyprProcessor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data);
}
