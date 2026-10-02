// Package nvim extends the generic template processor with a
// Reload: a running nvim (started with --listen) is asked to
// switch colourscheme so the freshly rendered colours file
// takes effect without restarting it.
package nvim

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	generic "theme-engine/internal/adapters/tools/generic"
	"theme-engine/internal/domain/renderctx"
)

type Processor struct {
	generic.Processor
}

func New() Processor { return Processor{} }

// Reload sends ":colorscheme <name>" to nvim listening on
// $NVIM_LISTEN_ADDRESS. The colourscheme name is the base
// name of the rendered file — nvim registers a colourscheme
// per lua file under colors/. Without a listening nvim this
// is a no-op: new instances read the file at startup anyway.
func (Processor) Reload(outputPath string, _ *renderctx.Context) error {
	addr := os.Getenv("NVIM_LISTEN_ADDRESS")
	if addr == "" {
		return nil
	}

	name := strings.TrimSuffix(filepath.Base(outputPath), ".lua")
	if name == "" {
		return nil
	}

	// <C-\><C-n> guarantees normal mode before the command
	// runs, so typing in an insert-mode buffer is safe.
	return exec.Command("nvim", "--server", addr, "--remote-send",
		"<C-\\><C-n>:colorscheme "+name+"<CR>").Run()
}
