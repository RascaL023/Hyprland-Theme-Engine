// Package lazygit extends the generic template processor.
// Research conclusion (lazygit issues #1158, closed 2021,
// and #4602, 2025): lazygit reads config.yml only at
// startup and exposes no IPC or signal to re-read it.
// SIGHUP is not handled by the TUI — sending it would
// kill the process, losing the user's in-flight work.
// Reload is therefore a documented no-op: freshly
// started instances pick up the rendered config, while
// running ones need a manual restart.
package lazygit

import (
	generic "theme-engine/internal/adapters/tools/generic"
	"theme-engine/internal/domain/renderctx"
)

type Processor struct {
	generic.Processor
}

func New() Processor { return Processor{} }

// Reload is a deliberate no-op. There is no safe way to
// make a running lazygit re-read its config (issues #1158,
// #4602), so returning nil keeps the engine's apply phase
// quiet instead of warning on every render.
func (Processor) Reload(_ string, _ *renderctx.Context) error {
	return nil
}
