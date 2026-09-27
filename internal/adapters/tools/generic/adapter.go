// Package generic implements template-only targets: no dedicated
// theme.json schema, the global render context is passed straight
// to the template. Target identity comes from the wiring map key,
// so this processor carries no name of its own.
package generic

import (
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
