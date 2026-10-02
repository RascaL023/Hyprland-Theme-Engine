// Package ports defines the inbound contracts (use-case boundaries)
// implemented by tool/platform adapters.
package ports

import "theme-engine/internal/domain/renderctx"

type Parser interface {
	Parse(raw any) (any, error)
}

type Resolver interface {
	Resolve(in any, ctx *renderctx.Context) (any, error)
}

type Renderer interface {
	Render(templatePath, outputPath string, data any) error
}

// Reloader applies freshly rendered output to a tool that is
// already running. Optional: processors that don't implement it
// are skipped by the engine's apply phase (unknown-processor
// convention — absence is not an error).
type Reloader interface {
	Reload(outputPath string, ctx *renderctx.Context) error
}

type Processor interface {
	Parser
	Resolver
	Renderer
}
