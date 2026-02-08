package main

import (
	"os"
	"path/filepath"

	"theme-engine/internal/core/context"
	"theme-engine/internal/core/themes/palette"
	"theme-engine/internal/core/themes/state"
	"theme-engine/internal/core/themes/theme"
	"theme-engine/internal/helper"
	"theme-engine/internal/loader"
)


const (
	ExitOK        = 0
	ExitGeneral   = 1
	ExitLoad	    = 2
	ExitIOError   = 3
	ExitParseFail = 4
	ExitResolve		= 5
	ExitRender		= 6
)

func main() {
	option := os.Args[1];

	// ============================= RESOURCE INPUT =============================

	assetsMap := os.ExpandEnv("$MYENV/map");

	state, err := loader.LoadJSON[state.State](assetsMap + "/.state.json");
	helper.CheckErr(err, ExitLoad, "Error on loading state!");

	apps, err := loader.LoadToolMap(assetsMap + "/path.txt", state);
	helper.CheckErr(err, ExitLoad, "Error on loading tools map");

	rawPalette, err := loader.LoadJSON[palette.Raw](filepath.Join("themes", state.Theme.Name, "palette.json"));
	helper.CheckErr(err, ExitLoad, "Error on loading palette");

	palette := rawPalette.ResolveSelected(state.Theme.Type);
	theme.BuildFlattenPalette(palette);

	rawTheme, err := loader.LoadJSON[theme.Theme](filepath.Join("themes", state.Theme.Name, "theme.json"));
	helper.CheckErr(err, ExitLoad, "Error on loading theme %s", state.Theme.Name);

	// ============================= RESOURCE INPUT =============================

	ctx := context.Context{
		Palette: palette,
		Theme: rawTheme,
	}

	switch option {
		case "": runAll(apps, rawTheme, ctx);
		default:
			app, ok := apps[option]; 
			if !ok {
				helper.ExitWrapper(ExitGeneral, "Unknown option %s", option);
			}
			
			runSpecific(option, app, rawTheme, ctx);
	}
}
