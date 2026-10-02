package firefox

import (
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

func (Processor) Resolve(in any, ctx *renderctx.Context) (any, error) {
	inp, _ := in.(Raw)

	return Firefox{
		Background: pathenv.ResolveVar(inp.Special.Background, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Foreground: pathenv.ResolveVar(inp.Special.Foreground, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Cursor: pathenv.ResolveVar(inp.Special.Cursor, palette.ResolvedPaletteVars{P: ctx.Palette}),

		Color0: pathenv.ResolveVar(inp.Colors.Color0, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color1: pathenv.ResolveVar(inp.Colors.Color1, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color2: pathenv.ResolveVar(inp.Colors.Color2, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color3: pathenv.ResolveVar(inp.Colors.Color3, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color4: pathenv.ResolveVar(inp.Colors.Color4, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color5: pathenv.ResolveVar(inp.Colors.Color5, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color6: pathenv.ResolveVar(inp.Colors.Color6, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color7: pathenv.ResolveVar(inp.Colors.Color7, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color8: pathenv.ResolveVar(inp.Colors.Color8, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color9: pathenv.ResolveVar(inp.Colors.Color9, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color10: pathenv.ResolveVar(inp.Colors.Color10, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color11: pathenv.ResolveVar(inp.Colors.Color11, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color12: pathenv.ResolveVar(inp.Colors.Color12, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color13: pathenv.ResolveVar(inp.Colors.Color13, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color14: pathenv.ResolveVar(inp.Colors.Color14, palette.ResolvedPaletteVars{P: ctx.Palette}),
		Color15: pathenv.ResolveVar(inp.Colors.Color15, palette.ResolvedPaletteVars{P: ctx.Palette}),
	}, nil
}

func (Processor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data)
}
