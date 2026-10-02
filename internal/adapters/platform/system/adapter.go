// Package system renders the platform apply script (dconf).
// Template-only like generic, kept separate because its template
// and output live under domain/, not tools/.
package system

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"theme-engine/internal/domain/renderctx"
	"theme-engine/internal/infra/renderer"
)

type Processor struct{}

func New() Processor { return Processor{} }

func (Processor) Parse(_ any) (any, error) {
	return nil, nil
}

func (Processor) Resolve(_ any, ctx *renderctx.Context) (any, error) {
	return ctx, nil
}

func (Processor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data)
}

// Reload runs the freshly rendered apply script. The
// script's `dconf write` calls propagate live to GTK
// applications that are already running: the settings
// daemon picks the change up over D-Bus and pushes
// XSETTINGS updates, so fonts, cursor size and the
// prefer-dark/prefer-light scheme switch without a
// restart. The script is idempotent, so running it on
// every render is safe. A non-zero exit (e.g. dconf
// not installed, or a key rejected) is returned as an
// error — the engine logs it as a non-fatal warning.
// The renderer writes mode 0600, so the script is
// invoked through bash instead of executed directly.
func (Processor) Reload(outputPath string, _ *renderctx.Context) error {
	if _, err := os.Stat(outputPath); err != nil {
		return err
	}

	var stderr bytes.Buffer
	cmd := exec.Command("bash", outputPath)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("apply %s: %w: %s",
			outputPath, err, strings.TrimSpace(stderr.String()))
	}
	return nil
}
