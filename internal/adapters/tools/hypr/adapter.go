package hypr

import (
	"os/exec"
	"strings"

	"theme-engine/internal/adapters/tools/jsonx"
	"theme-engine/internal/domain/palette"
	"theme-engine/internal/domain/renderctx"
	"theme-engine/internal/infra/pathenv"
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

// Reload asks Hyprland to re-read its config. Hyprland
// inotify-watches the file it was launched with, so a plain
// save is usually enough; this also covers atomic (rename)
// writes and is the documented manual mechanism (hyprctl
// reload). Idempotent — safe to call on every render.
func (Processor) Reload(_ string, _ *renderctx.Context) error {
	return exec.Command("hyprctl", "reload").Run()
}
