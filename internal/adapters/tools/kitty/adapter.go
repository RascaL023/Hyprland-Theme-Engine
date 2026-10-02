package kitty

import (
	"bufio"
	"os"
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

func (Processor) Resolve(in any, ctx *renderctx.Context) (any, error) {
	inp, _ := in.(Raw)
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

// Reload applies the colour section of the freshly
// rendered config to a running kitty via remote
// control: only colour entries are extracted (kitty
// @ set-colors accepts colours, not fonts/opacity),
// written to a temp file, then applied to all
// windows. Requires remote control enabled in
// kitty.conf (allow_remote_control yes) — without
// it kitty @ fails and the engine logs a warning.
func (Processor) Reload(outputPath string, _ *renderctx.Context) error {
	if _, err := exec.LookPath("kitty"); err != nil {
		return err
	}

	colours, err := extractColourLines(outputPath)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp("", "kitty-colours-*.conf")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.WriteString(colours); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	return exec.Command("kitty", "@", "set-colors", "-a", tmpName).Run()
}

// extractColourLines copies the colour entries from a
// rendered kitty config: foreground/background/cursor,
// selection_*, color0-15, *_tab_*, *_border_color.
// Layout and font entries are excluded because
// kitty @ set-colors only accepts colours.
func extractColourLines(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var out strings.Builder
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key := line
		if i := strings.IndexAny(line, " \t"); i >= 0 {
			key = line[:i]
		}
		if !isColourKey(key) {
			continue
		}
		out.WriteString(line)
		out.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func isColourKey(key string) bool {
	switch key {
	case "foreground", "background", "cursor",
		"cursor_text_color",
		"selection_foreground", "selection_background",
		"active_tab_foreground", "active_tab_background",
		"inactive_tab_foreground", "inactive_tab_background",
		"active_border_color", "inactive_border_color":
		return true
	}
	return strings.HasPrefix(key, "color") && isDigits(key[len("color"):])
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
