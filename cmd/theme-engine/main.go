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
	assetsMap := os.ExpandEnv("$MYENV/map");

	state, err := loader.LoadJSON[state.State](assetsMap + "/.state.json");
	helper.CheckErr(err, "Error on loading state!");

	tools, err := loader.LoadToolMap(assetsMap + "/path.txt", state);
	helper.CheckErr(err, "Error on loading tools map");

	domains, err := loader.LoadToolMap(assetsMap + "/domain.txt", state);
	helper.CheckErr(err, "Error on loading tools map");

	rawPalette, err := loader.LoadJSON[palette.Raw](filepath.Join("themes", state.Theme.Name, "palette.json"));
	helper.CheckErr(err, "Error on loading palette");
	palette := rawPalette.ResolveSelected(state.Theme.Type);
	theme.BuildFlattenPalette(palette);

	rawTheme, err := loader.LoadJSON[theme.Theme](filepath.Join("themes", state.Theme.Name, "theme.json"));
	helper.CheckErr(err, "Error on loading theme");

	ctx := context.Context{
		Palette: palette,
	}

	for domain, path := range domains {
		processor, ok := processor.GetProcessor(domain);
		if !ok {
			fmt.Println("Unknown domain", domain);
			continue;
		}

		resolved, err := processor.Resolve(nil, &ctx);
		helper.CheckErr(err, "Error on resolving");

		if err := processor.Render(path.TemplatePath, path.OutputPath, resolved); err != nil {
			helper.CheckErr(err, "Error on rendering");
		}

		fmt.Println("Success rendering", domain);
	}

// fmt.Println(palette.Flat)
	for toolName, toolPath := range tools {
		processor, ok := processor.GetProcessor(toolName);
		if !ok {
			fmt.Println("Unknown tool", toolName);
			continue;
		}

		parsed, err := processor.Parse(rawTheme.Tools[toolName]);
		helper.CheckErr(err, "Error on parsing tool");

		resolved, err := processor.Resolve(parsed, &ctx);
		helper.CheckErr(err, "Error on resolving tool");

		helper.CheckErr(processor.Render(toolPath.TemplatePath, toolPath.OutputPath, resolved), "Error on rendering tool");
	}
}
