package foot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"theme-engine/internal/adapters/tools/jsonx"
	"theme-engine/internal/domain/palette"
	"theme-engine/internal/domain/renderctx"
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

	return Foot{
		Font:     ctx.Theme.Theme.Fonts.Terminal.Family,
		FontSize: ctx.Theme.Theme.Fonts.Terminal.Size,
		PaddingX: inp.Padding.X,
		PaddingY: inp.Padding.Y,

		Palette: ctx.Palette,
	}, nil
}

func (Processor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data)
}

// Reload pushes the palette straight into every accessible
// terminal session via OSC sequences — the set_term_colors
// trick (tools/set_term is the standalone form of this,
// usable without running the engine). Foot has no live
// config reload (issue #1653: only a footserver restart
// picks up ui.ini), but its emulator honours OSC 4/10/11/12
// at runtime, so running foot instances switch colours
// without a restart. Sequences written to a pty are
// delivered to whichever emulator owns it, so
// alacritty/kitty/xterm sessions get the same palette too
// — consistent with what their configs get. Sessions we
// cannot write (another user's pty) are skipped; with no
// terminal sessions at all this is a no-op.
func (Processor) Reload(_ string, ctx *renderctx.Context) error {
	if ctx == nil || ctx.Palette == nil {
		return fmt.Errorf("reload foot: no palette in context")
	}

	seq := oscSequence(ctx.Palette)

	entries, err := os.ReadDir("/dev/pts")
	if err != nil {
		return fmt.Errorf("list /dev/pts: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !isDigits(entry.Name()) {
			continue
		}
		// EACCES on other users' ptys is expected — skip.
		_ = os.WriteFile(filepath.Join("/dev/pts", entry.Name()), []byte(seq), 0)
	}
	return nil
}

// oscSequence builds the runtime colour update: OSC 10
// (default foreground), OSC 11 (default background), OSC 12
// (cursor colour) and OSC 4;N (palette entry N) for all 16
// ANSI colours. ESC \ terminates every sequence (ST).
func oscSequence(p *palette.ResolvedPalette) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\x1b]10;%s\x1b\\", p.Foreground)
	fmt.Fprintf(&b, "\x1b]11;%s\x1b\\", p.Background)
	fmt.Fprintf(&b, "\x1b]12;%s\x1b\\", p.Cursor)
	for i, c := range p.Colors {
		fmt.Fprintf(&b, "\x1b]4;%d;%s\x1b\\", i, c)
	}
	return b.String()
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
