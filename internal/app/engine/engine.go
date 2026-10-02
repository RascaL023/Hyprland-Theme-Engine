package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"theme-engine/internal/app/ports"
	"theme-engine/internal/domain/palette"
	"theme-engine/internal/domain/renderctx"
	"theme-engine/internal/domain/state"
	"theme-engine/internal/domain/theme"
	"theme-engine/internal/infra/loader"
	"theme-engine/internal/infra/log"
	"theme-engine/internal/infra/pathenv"
)

type Config struct {
	MapDir     string
	ThemeDir   string
	Processors map[string]ports.Processor
}

type Engine struct {
	apps        map[string]*loader.ToolMap
	theme       *theme.Theme
	ctx         renderctx.Context
	processors  map[string]ports.Processor
	applyEnabled bool
}

func New(cfg Config) (*Engine, error) {
	if cfg.ThemeDir == "" {
		cfg.ThemeDir = "themes"
	}
	if cfg.Processors == nil {
		return nil, fmt.Errorf("engine: no processors wired")
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

	resolvedPalette, err := rawPalette.ResolveSelected(st.Theme.Type)
	if err != nil {
		return nil, fmt.Errorf("resolve palette: %w", err)
	}
	palette.BuildFlattenPalette(resolvedPalette)

	rawTheme, err := loader.LoadJSON[theme.Theme](filepath.Join(cfg.ThemeDir, st.Theme.Name, "theme.json"))
	if err != nil {
		return nil, fmt.Errorf("load theme %s: %w", st.Theme.Name, err)
	}

	rawTheme.Theme.Fonts.ResolveDefaults()

	return &Engine{
		apps:         apps,
		theme:        rawTheme,
		processors:   cfg.Processors,
		applyEnabled: applyEnabled(cfg.MapDir),
		ctx: renderctx.Context{
			Palette:   resolvedPalette,
			Theme:     rawTheme,
			ThemeType: st.Theme.Type,
			ThemeName: st.Theme.Name,
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

	proc, ok := e.processors[name]
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

	if e.applyEnabled {
		if r, ok := proc.(ports.Reloader); ok {
			if err := r.Reload(paths.OutputPath, &e.ctx); err != nil {
				log.Warn("reload %s: %v", name, err)
			}
		}
	}

	log.Info("Success rendering %s", name)
	return nil
}

func DefaultMapDir() string {
	if path := pathenv.Env("THEME_ENGINE_MAP"); path != "" {
		return path
	}

	if _, err := os.Stat("config/path.txt"); err == nil {
		return "config"
	}

	if path := pathenv.Env("MYENV"); path != "" {
		return filepath.Join(path, "map")
	}

	return "config"
}

// applyEnabled decides whether the post-render apply phase runs.
// THEME_ENGINE_APPLY=1 forces it on, =0 forces it off, and the
// default (auto) enables it only when the map directory is an
// external deployment ($THEME_ENGINE_MAP / $MYENV), not the
// local config/ used for dev and tests.
func applyEnabled(mapDir string) bool {
	switch pathenv.Env("THEME_ENGINE_APPLY") {
	case "1":
		return true
	case "0":
		return false
	default: // "auto"
		return !isLocalMapDir(mapDir)
	}
}

func isLocalMapDir(mapDir string) bool {
	abs, err := filepath.Abs(mapDir)
	if err != nil {
		return false
	}
	local, err := filepath.Abs("config")
	if err != nil {
		return false
	}
	return abs == local
}
