package register

import "theme-engine/internal/core/context"

type Parser interface {
	Parse(json any) (any, error)
}

type Resolver interface {
	Resolve(in any, ctx *context.Context) (any, error)
}

type Renderer interface {
	Render(templatePath, outputPath string, data any) error
}
