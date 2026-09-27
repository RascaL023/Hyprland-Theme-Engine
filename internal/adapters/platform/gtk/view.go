package gtk

import "theme-engine/internal/domain/renderctx"

// Gtk carries the global render context so the template accesses
// .Palette and .Theme directly, like every other adapter.
type Gtk struct {
	*renderctx.Context
}
