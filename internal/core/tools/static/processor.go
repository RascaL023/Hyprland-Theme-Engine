package static

import (
	"theme-engine/internal/core/context"
	"theme-engine/internal/processor"
	"theme-engine/internal/renderer"
)

type Processor struct {
	name string
}

func init() {
	processor.RegisterProcessor(Processor{name: "waybar"})
	processor.RegisterProcessor(Processor{name: "hyprland"})
	processor.RegisterProcessor(Processor{name: "yazi"})
}

func (p Processor) Name() string { return p.name }

func (Processor) Parse(_ any) (any, error) {
	return nil, nil
}

func (Processor) Resolve(_ any, ctx *context.Context) (any, error) {
	return ctx, nil
}

func (Processor) Render(templatePath, outputPath string, data any) error {
	return renderer.Render(templatePath, outputPath, data)
}
