package gtk

import "theme-engine/internal/domain/renderctx"

// Gtk keeps the historical .Config wrapper for the current template
// while also embedding the context so new templates can use .Palette
// and .Theme directly.
type Gtk struct {
	*renderctx.Context
	Config *renderctx.Context
}
