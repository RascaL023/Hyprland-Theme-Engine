package cava

import (
	"fmt"

	"theme-engine/internal/adapters/tools/jsonx"
	"theme-engine/internal/domain/palette"
	"theme-engine/internal/domain/renderctx"
	"theme-engine/internal/infra/pathenv"
	"theme-engine/internal/infra/renderer"
)

type Processor struct{}

func New() Processor { return Processor{} }

func (Processor) Name() string { return "cava" }

func (Processor) Parse(in any) (any, error) {
	var cfg Raw
	if err := jsonx.Decode(in, &cfg); err != nil {
		return Raw{}, err
	}
	return cfg, nil
}

func (Processor) Resolve(in any, ctx *renderctx.Context) (any, error) {
	inp, _ := in.(Raw)
	if len(inp.Gradients) < 5 {
		return nil, fmt.Errorf("expected at least 5 gradients, got %d", len(inp.Gradients))
	}

	return Cava{
		Gradient1: pathenv.ResolveVar(inp.Gradients[0], palette.ResolvedPaletteVars{P: ctx.Palette}),
		Gradient2: pathenv.ResolveVar(inp.Gradients[1], palette.ResolvedPaletteVars{P: ctx.Palette}),
		Gradient3: pathenv.ResolveVar(inp.Gradients[2], palette.ResolvedPaletteVars{P: ctx.Palette}),
		Gradient4: pathenv.ResolveVar(inp.Gradients[3], palette.ResolvedPaletteVars{P: ctx.Palette}),
		Gradient5: pathenv.ResolveVar(inp.Gradients[4], palette.ResolvedPaletteVars{P: ctx.Palette}),
	}, nil
}

func (Processor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data)
}
