package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"theme-engine/internal/app/engine"
	"theme-engine/internal/domain/state"
	"theme-engine/internal/infra/loader"
	"theme-engine/internal/infra/renderer"
)

// runSetTheme switches the active theme in .state.json and
// re-renders every target (which also triggers the apply
// phase for running tool instances).
func runSetTheme(args []string) {
	name, themeType, err := parseSetThemeArgs(args)
	if err != nil {
		fail(ExitUsage, "%v", err)
	}

	mapDir := engine.DefaultMapDir()

	if _, err := os.Stat(filepath.Join("themes", name, "palette.json")); err != nil {
		fail(ExitLoad, "Error: unknown theme %q (themes/%s/palette.json: %v)", name, name, err)
	}

	st, err := loader.LoadJSON[state.State](filepath.Join(mapDir, ".state.json"))
	if err != nil {
		fail(ExitLoad, "Error on loading state: %v", err)
	}

	st.Theme.Name = name
	if themeType != "" {
		st.Theme.Type = themeType
	}

	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		fail(ExitGeneral, "Error on encoding state: %v", err)
	}
	data = append(data, '\n')

	if err := renderer.WriteIfChanged(filepath.Join(mapDir, ".state.json"), data); err != nil {
		fail(ExitIOError, "Error on writing state: %v", err)
	}

	app, err := engine.New(engine.Config{
		MapDir:     mapDir,
		Processors: defaultProcessors(),
	})
	if err != nil {
		fail(ExitLoad, "Error on initializing engine: %v", err)
	}

	if err := app.RunAll(); err != nil {
		fail(ExitGeneral, "Error on rendering all targets: %v", err)
	}
}

func parseSetThemeArgs(args []string) (name, themeType string, err error) {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return "", "", fmt.Errorf("usage: theme-engine set-theme <name> [--type dark|light]")
	}

	name = args[0]
	for i := 1; i < len(args); i++ {
		if args[i] != "--type" {
			return "", "", fmt.Errorf("unexpected argument %q", args[i])
		}
		if i+1 >= len(args) {
			return "", "", fmt.Errorf("--type needs a value (dark or light)")
		}
		i++
		switch args[i] {
		case "dark", "light":
			themeType = args[i]
		default:
			return "", "", fmt.Errorf("--type must be dark or light, got %q", args[i])
		}
	}

	return name, themeType, nil
}
