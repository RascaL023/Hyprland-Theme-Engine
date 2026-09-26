// Package system renders the platform apply script (dconf).
// Template-only like generic, kept separate because its template
// and output live under domain/, not tools/.
package system

import (
	"theme-engine/internal/domain/renderctx"
	"theme-engine/internal/infra/renderer"
)

type Processor struct{}

func New() Processor { return Processor{} }

func (Processor) Name() string { return "system" }

func (Processor) Parse(_ any) (any, error) {
	return nil, nil
}

func (Processor) Resolve(_ any, ctx *renderctx.Context) (any, error) {
	return ctx, nil
}

func (Processor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data)
}
