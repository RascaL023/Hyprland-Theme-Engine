package main

import (
	"fmt"
	"os"
	"path/filepath"

	"theme-engine/internal/core/context"
	"theme-engine/internal/core/themes/palette"
	"theme-engine/internal/core/themes/state"
	"theme-engine/internal/helper"
	"theme-engine/internal/loader"
)

func main() {
	assetsMap := os.ExpandEnv("$MYENV/map");

	state, err := loader.LoadJSON[state.State](assetsMap + "/.state.json");
	helper.CheckErr(err, "Error on loading state!");

	tools, err := loader.LoadToolMap(assetsMap + "/path.txt", state);
	helper.CheckErr(err, "Error on loading tools map");

	rawPalette, err := loader.LoadJSON[palette.Raw](filepath.Join("themes", state.Theme.Name, "palette.json"));
	helper.CheckErr(err, "Error on loading palette");
	palette := rawPalette.ResolveSelected(state.Theme.Type);


	ctx := context.Context{
		Palette: palette,
	}

	fmt.Println(tools["cava"]);
	fmt.Println(ctx.Palette);
}
