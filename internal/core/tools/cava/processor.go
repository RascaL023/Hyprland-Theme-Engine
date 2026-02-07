package cava

import (
	"encoding/json"
	"theme-engine/internal/core/context"
	"theme-engine/internal/core/themes/palette"
	"theme-engine/internal/processor"
	"theme-engine/internal/renderer"
	"theme-engine/internal/resolver"
)

type CavaProcessor struct{}

func init() { processor.RegisterProcessor(CavaProcessor{}); }

func (CavaProcessor) Name() string { return "cava"; }

func (CavaProcessor) Parse(in any) (any, error) {
	var cfg Raw;
	err := json.Unmarshal(in.(json.RawMessage), &cfg);
	return cfg, err;
}

func (CavaProcessor) Resolve(in any, ctx *context.Context) (any, error) {
	inp := in.(Raw);

	return Cava{
		Gradient1: resolver.ResolveVar(inp.Gradients[0], palette.ResolvedPaletteVars{P: ctx.Palette}),
		Gradient2: resolver.ResolveVar(inp.Gradients[1], palette.ResolvedPaletteVars{P: ctx.Palette}),
		Gradient3: resolver.ResolveVar(inp.Gradients[2], palette.ResolvedPaletteVars{P: ctx.Palette}),
		Gradient4: resolver.ResolveVar(inp.Gradients[3], palette.ResolvedPaletteVars{P: ctx.Palette}),
		Gradient5: resolver.ResolveVar(inp.Gradients[4], palette.ResolvedPaletteVars{P: ctx.Palette}),
	}, nil;
}

func (CavaProcessor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data);
}

