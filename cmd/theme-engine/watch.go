package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"theme-engine/internal/app/engine"
	"theme-engine/internal/domain/state"
	"theme-engine/internal/infra/loader"
	"theme-engine/internal/infra/log"
	"theme-engine/internal/infra/renderer"
)

const (
	watchInterval = 500 * time.Millisecond
	watchDebounce = 300 * time.Millisecond
)

// runWatch polls the template tree, the active theme's
// directory and the map directory for changes, then
// re-renders everything. Polling (not fsnotify) keeps
// the project dependency-free; the roots are small so
// os.Stat per file every 500ms is cheap.
func runWatch() {
	mapDir := engine.DefaultMapDir()

	log.Info("watch: templates, active theme and %s (interval %s)", mapDir, watchInterval)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(watchInterval)
	defer ticker.Stop()

	var prev map[string]fileStamp
	first := true

	for {
		select {
		case <-sig:
			log.Info("watch: stopped")
			return
		case <-ticker.C:
		}

		roots, excludes, err := watchRoots(mapDir)
		if err != nil {
			log.Error("watch: %v (still watching)", err)
			continue
		}

		cur := snapshot(roots, excludes)
		if first {
			prev = cur
			first = false
			continue
		}

		changed := diff(prev, cur)
		if len(changed) == 0 {
			continue
		}

		// Debounce: editors save via temp+rename, which can
		// land as several events within milliseconds. Sleep
		// once, then absorb whatever landed in that window
		// into a single reload.
		time.Sleep(watchDebounce)
		prev = snapshot(roots, excludes)

		log.Info("watch: changed: %s", strings.Join(changed, ", "))
		renderer.ClearCache()
		reload(mapDir)
	}
}

// watchRoots returns the directories to poll: the template
// tree, the active theme's directory (re-evaluated every
// cycle so editing .state.json to switch themes works) and
// the map directory. excludes holds absolute output paths
// from the tool map so a re-render can never feed back
// into the watcher.
func watchRoots(mapDir string) ([]string, map[string]struct{}, error) {
	st, err := loader.LoadJSON[state.State](filepath.Join(mapDir, ".state.json"))
	if err != nil {
		return nil, nil, fmt.Errorf("load state: %w", err)
	}

	excludes := make(map[string]struct{})
	apps, err := loader.LoadToolMap(filepath.Join(mapDir, "path.txt"), st)
	if err != nil {
		return nil, nil, fmt.Errorf("load tool map: %w", err)
	}
	for _, m := range apps {
		if abs, err := filepath.Abs(m.OutputPath); err == nil {
			excludes[abs] = struct{}{}
		}
	}

	roots := []string{
		"assets/templates",
		filepath.Join("themes", st.Theme.Name),
		mapDir,
	}
	return roots, excludes, nil
}

type fileStamp struct {
	size  int64
	mtime int64
}

func snapshot(roots []string, excludes map[string]struct{}) map[string]fileStamp {
	sn := make(map[string]fileStamp)
	for _, root := range roots {
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // missing or unreadable root — skip
			}
			if d.IsDir() {
				return nil
			}
			if abs, err := filepath.Abs(path); err == nil {
				if _, ok := excludes[abs]; ok {
					return nil
				}
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			sn[path] = fileStamp{size: info.Size(), mtime: info.ModTime().UnixNano()}
			return nil
		})
	}
	return sn
}

func diff(a, b map[string]fileStamp) []string {
	var changed []string
	for path, stamp := range b {
		if old, ok := a[path]; !ok || old != stamp {
			changed = append(changed, path)
		}
	}
	for path := range a {
		if _, ok := b[path]; !ok {
			changed = append(changed, path+" (removed)")
		}
	}
	sort.Strings(changed)
	return changed
}

func reload(mapDir string) {
	app, err := engine.New(engine.Config{
		MapDir:     mapDir,
		Processors: defaultProcessors(),
	})
	if err != nil {
		// e.g. theme.json mid-edit: keep watching, the
		// next cycle renders once the file is valid.
		log.Error("watch: reload failed (still watching): %v", err)
		return
	}
	if err := app.RunAll(); err != nil {
		log.Error("watch: render failed: %v", err)
	}
}
