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

type Processor interface {
	Parser
	Resolver
	Renderer
}
