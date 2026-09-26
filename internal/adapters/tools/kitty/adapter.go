package kitty

import (
	"encoding/json"

	"theme-engine/internal/domain/palette"
	"theme-engine/internal/domain/renderctx"
	"theme-engine/internal/infra/pathenv"
	"theme-engine/internal/infra/renderer"
)

type Processor struct{}

func New() Processor { return Processor{} }

func (Processor) Name() string { return "kitty" }

func (Processor) Parse(in any) (any, error) {
	var cfg Raw
	err := json.Unmarshal(in.(json.RawMessage), &cfg)
	return cfg, err
}

func (Processor) Resolve(in any, ctx *renderctx.Context) (any, error) {
	inp := in.(Raw)
	return Kitty{
		Palette:             ctx.Palette,
		SelectionBackground: ctx.Palette.AccentPrimary,
		SelectionForeground: ctx.Palette.LayerBase,

		CursorShape: inp.CursorShape,
		Opacity:     inp.Opacity,

		TabBar:        inp.Tab.Bar,
		TabPowerline:  inp.Tab.Style,
		ActiveTabBg:   pathenv.ResolveVar(inp.Tab.Active.Background, palette.ResolvedPaletteVars{P: ctx.Palette}),
		ActiveTabFg:   pathenv.ResolveVar(inp.Tab.Active.Foreground, palette.ResolvedPaletteVars{P: ctx.Palette}),
		InActiveTabBg: pathenv.ResolveVar(inp.Tab.InActive.Background, palette.ResolvedPaletteVars{P: ctx.Palette}),
		InActiveTabFg: pathenv.ResolveVar(inp.Tab.InActive.Foreground, palette.ResolvedPaletteVars{P: ctx.Palette}),

		ActiveBorder:   pathenv.ResolveVar(inp.Window.Border.Active, palette.ResolvedPaletteVars{P: ctx.Palette}),
		InActiveBorder: pathenv.ResolveVar(inp.Window.Border.InActive, palette.ResolvedPaletteVars{P: ctx.Palette}),

		WindowBorder:  inp.Window.Border.Width,
		WindowPadding: inp.Window.Padding,
		WindowMargin:  inp.Window.Margin,

		FontSize:   ctx.Theme.Theme.Fonts.Terminal.Size,
		FontFamily: ctx.Theme.Theme.Fonts.Terminal.Family,
	}, nil
}

func (Processor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data)
}
