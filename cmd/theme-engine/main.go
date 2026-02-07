package main

import (
	"fmt"
	"os"
	"path/filepath"

	"theme-engine/internal/core/context"
	"theme-engine/internal/core/themes/palette"
	"theme-engine/internal/core/themes/state"
	"theme-engine/internal/core/themes/theme"
	"theme-engine/internal/helper"
	"theme-engine/internal/loader"
	"theme-engine/internal/processor"
)

func main() {
	// ============================= RESOURCE INPUT =============================

	assetsMap := os.ExpandEnv("$MYENV/map");

	state, err := loader.LoadJSON[state.State](assetsMap + "/.state.json");
	helper.CheckErr(err, "Error on loading state!");

	apps, err := loader.LoadToolMap(assetsMap + "/path.txt", state);
	helper.CheckErr(err, "Error on loading tools map");

	rawPalette, err := loader.LoadJSON[palette.Raw](filepath.Join("themes", state.Theme.Name, "palette.json"));
	helper.CheckErr(err, "Error on loading palette");

	palette := rawPalette.ResolveSelected(state.Theme.Type);
	theme.BuildFlattenPalette(palette);

	rawTheme, err := loader.LoadJSON[theme.Theme](filepath.Join("themes", state.Theme.Name, "theme.json"));
	helper.CheckErr(err, "Error on loading theme");

	// ============================= RESOURCE INPUT =============================

	ctx := context.Context{
		Palette: palette,
		Theme: rawTheme,
	}

	// ============================= RUN =============================

	var parsed any = nil;
	var errParse error = nil;

	for name, path := range apps {
		processor, ok := processor.GetProcessor(name);
		if !ok {
			fmt.Println("Unknown tool", name);
			continue;
		}

		raw, hasConfig := rawTheme.Tools[name];
		if hasConfig {
			parsed, errParse = processor.Parse(raw);
			helper.CheckErr(errParse, "Error on parsing tool");
		}

		resolved, err := processor.Resolve(parsed, &ctx);
		helper.CheckErr(err, "Error on resolving tool");

		helper.CheckErr(processor.Render(path.TemplatePath, path.OutputPath, resolved), "Error on rendering tool");
		fmt.Println("Success rendering", name);
	}

	// ============================= RUN =============================
}
