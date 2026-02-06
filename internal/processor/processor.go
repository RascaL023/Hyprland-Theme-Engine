package processor

import "theme-engine/internal/register"

var Registered = map[string]Processor{}

type Processor interface {
    Name() string
		register.Parser
		register.Resolver
		register.Renderer
}

func RegisterProcessor(p Processor) {
    Registered[p.Name()] = p;
}

func GetProcessor(name string) (Processor, bool) {
    p, ok := Registered[name];
    return p, ok;
}

