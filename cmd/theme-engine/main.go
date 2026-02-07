package main

import (
	"os"
	"path/filepath"

	"theme-engine/internal/core/context"
	"theme-engine/internal/core/log"
	"theme-engine/internal/core/themes/palette"
	"theme-engine/internal/core/themes/state"
	"theme-engine/internal/core/themes/theme"
	"theme-engine/internal/helper"
	"theme-engine/internal/loader"
	"theme-engine/internal/processor"
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

	// ============================= RUN ALL =============================

	var parsed any = nil;
	var errParse error = nil;

	for name, path := range apps {
		processor, ok := processor.GetProcessor(name);
		if !ok {
			log.Warn("Unknown tool %s", name);
			continue;
		}

		raw, hasConfig := rawTheme.Tools[name];
		if hasConfig {
			parsed, errParse = processor.Parse(raw);
			helper.CheckErr(errParse, ExitParseFail, "Error on parsing %s", name);
		}

		resolved, err := processor.Resolve(parsed, &ctx);
		helper.CheckErr(err, ExitResolve, "Error on resolving %s", name);

		helper.CheckErr(processor.Render(path.TemplatePath, path.OutputPath, resolved), ExitRender, "Error on rendering", name);
		log.Info("Success rendering %s", name);
	}

	// ============================= RUN ALL =============================
}
