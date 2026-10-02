package main

import (
	"theme-engine/internal/adapters/platform/gtk"
	"theme-engine/internal/adapters/platform/system"
	"theme-engine/internal/adapters/tools/alacritty"
	"theme-engine/internal/adapters/tools/cava"
	"theme-engine/internal/adapters/tools/foot"
	"theme-engine/internal/adapters/tools/generic"
	"theme-engine/internal/adapters/tools/hypr"
	"theme-engine/internal/adapters/tools/kitty"
	"theme-engine/internal/adapters/tools/nvim"
	"theme-engine/internal/app/ports"
)

func defaultProcessors() map[string]ports.Processor {
	genericTargets := []string{"yazi", "starship", "lazygit"}

	procs := make(map[string]ports.Processor, 10)
	procs["gtk"] = gtk.New()
	procs["system"] = system.New()
	procs["kitty"] = kitty.New()
	procs["foot"] = foot.New()
	procs["alacritty"] = alacritty.New()
	procs["cava"] = cava.New()
	procs["hypr"] = hypr.New()
	procs["nvim"] = nvim.New()
	for _, name := range genericTargets {
		procs[name] = generic.New()
	}

	return procs
}
