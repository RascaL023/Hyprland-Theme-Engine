package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"theme-engine/internal/core/context"
	"theme-engine/internal/core/log"
	"theme-engine/internal/core/themes/palette"
	"theme-engine/internal/core/themes/state"
	themeconfig "theme-engine/internal/core/themes/theme"
	"theme-engine/internal/helper"
	"theme-engine/internal/loader"
	"theme-engine/internal/processor"
)

type Config struct {
	MapDir   string
	ThemeDir string
}

type Engine struct {
	apps  map[string]*loader.ToolMap
	theme *themeconfig.Theme
	ctx   context.Context
}

func New(cfg Config) (*Engine, error) {
	if cfg.ThemeDir == "" {
		cfg.ThemeDir = "themes"
	}

	st, err := loader.LoadJSON[state.State](filepath.Join(cfg.MapDir, ".state.json"))
	if err != nil {
		return nil, fmt.Errorf("load state: %w", err)
	}

	apps, err := loader.LoadToolMap(filepath.Join(cfg.MapDir, "path.txt"), st)
	if err != nil {
		return nil, fmt.Errorf("load tool map: %w", err)
	}

	rawPalette, err := loader.LoadJSON[palette.Raw](filepath.Join(cfg.ThemeDir, st.Theme.Name, "palette.json"))
	if err != nil {
		return nil, fmt.Errorf("load palette: %w", err)
	}

	resolvedPalette := rawPalette.ResolveSelected(st.Theme.Type)
	themeconfig.BuildFlattenPalette(resolvedPalette)

	rawTheme, err := loader.LoadJSON[themeconfig.Theme](filepath.Join(cfg.ThemeDir, st.Theme.Name, "theme.json"))
	if err != nil {
		return nil, fmt.Errorf("load theme %s: %w", st.Theme.Name, err)
	}

	rawTheme.Theme.Fonts.ResolveDefaults()

	return &Engine{
		apps:  apps,
		theme: rawTheme,
		ctx: context.Context{
			Palette: resolvedPalette,
			Theme:   rawTheme,
		},
	}, nil
}

func (e *Engine) RunAll() error {
	names := make([]string, 0, len(e.apps))
	for name := range e.apps {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if err := e.Run(name); err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) Run(name string) error {
	paths, ok := e.apps[name]
	if !ok {
		return fmt.Errorf("unknown target %q", name)
	}

	proc, ok := processor.GetProcessor(name)
	if !ok {
		log.Warn("Unknown processor %s", name)
		return nil
	}

	var parsed any
	raw, hasConfig := e.theme.Tools[name]
	if hasConfig {
		var err error
		parsed, err = proc.Parse(raw)
		if err != nil {
			return fmt.Errorf("parse %s: %w", name, err)
		}
	}

	resolved, err := proc.Resolve(parsed, &e.ctx)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", name, err)
	}

	if err := proc.Render(paths.TemplatePath, paths.OutputPath, resolved); err != nil {
		return fmt.Errorf("render %s: %w", name, err)
	}

	log.Info("Success rendering %s", name)
	return nil
}

func DefaultMapDir() string {
	if path := helper.Env("THEME_ENGINE_MAP"); path != "" {
		return path
	}

	if _, err := os.Stat("config/path.txt"); err == nil {
		return "config"
	}

	if path := helper.Env("MYENV"); path != "" {
		return filepath.Join(path, "map")
	}

	return "config"
}
