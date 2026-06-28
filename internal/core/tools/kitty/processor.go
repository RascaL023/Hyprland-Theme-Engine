package kitty

import (
	"encoding/json"
	"theme-engine/internal/core/context"
	"theme-engine/internal/core/themes/palette"
	"theme-engine/internal/processor"
	"theme-engine/internal/renderer"
	"theme-engine/internal/resolver"
)

type KittyProcessor struct{}

func init() { processor.RegisterProcessor(KittyProcessor{}); }

func (KittyProcessor) Name() string { return "kitty"; }

func (KittyProcessor) Parse(in any) (any, error) {
	var cfg Raw;
	err := json.Unmarshal(in.(json.RawMessage), &cfg);
	return cfg, err;
}

func (KittyProcessor) Resolve(in any, ctx *context.Context) (any, error) {
	inp := in.(Raw);
	return Kitty{
		Palette: ctx.Palette,
		SelectionBackground: ctx.Palette.AccentPrimary,
		SelectionForeground: ctx.Palette.LayerBase,

		CursorShape: inp.CursorShape,
		Opacity: inp.Opacity,

		TabBar: inp.Tab.Bar,
		TabPowerline: inp.Tab.Style,
		ActiveTabBg: resolver.ResolveVar(inp.Tab.Active.Background, palette.ResolvedPaletteVars{P: ctx.Palette}),
		ActiveTabFg: resolver.ResolveVar(inp.Tab.Active.Foreground, palette.ResolvedPaletteVars{P: ctx.Palette}),
		InActiveTabBg: resolver.ResolveVar(inp.Tab.InActive.Background, palette.ResolvedPaletteVars{P: ctx.Palette}),
		InActiveTabFg: resolver.ResolveVar(inp.Tab.InActive.Foreground, palette.ResolvedPaletteVars{P: ctx.Palette}),

		ActiveBorder: resolver.ResolveVar(inp.Window.Border.Active, palette.ResolvedPaletteVars{P: ctx.Palette}),
		InActiveBorder: resolver.ResolveVar(inp.Window.Border.InActive, palette.ResolvedPaletteVars{P: ctx.Palette}),

		WindowBorder: inp.Window.Border.Width,
		WindowPadding: inp.Window.Padding,
		WindowMargin: inp.Window.Margin,

		FontSize: inp.FontSize,
		FontFamily: ctx.Theme.Theme.Fonts.Primary,
	}, nil;
}

func (KittyProcessor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data);
}


