package main

import (
	"theme-engine/internal/core/context"
	"theme-engine/internal/core/log"
	"theme-engine/internal/core/themes/theme"
	"theme-engine/internal/helper"
	"theme-engine/internal/loader"
	"theme-engine/internal/processor"
)

func runAll(
	apps map[string]*loader.ToolMap,
	rawTheme *theme.Theme,
	ctx context.Context,
) {
	var parsed any = nil;
	var errParse error = nil;

	for name, path := range apps {
		processor, ok := processor.GetProcessor(name);
		if !ok {
			log.Warn("Unknown processor %s", name);
			continue;
		}

		raw, hasConfig := rawTheme.Tools[name];
		if hasConfig {
			parsed, errParse = processor.Parse(raw);
			helper.CheckErr(errParse, ExitParseFail, "Error on parsing %s", name);
		}

		resolved, err := processor.Resolve(parsed, &ctx);
		helper.CheckErr(err, ExitResolve, "Error on resolving %s", name);

		helper.CheckErr(processor.Render(path.TemplatePath, path.OutputPath, resolved), ExitRender, "Error on rendering %s", name);
		log.Info("Success rendering %s", name);
	}
}

