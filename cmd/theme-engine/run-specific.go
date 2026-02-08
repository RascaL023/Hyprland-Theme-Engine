package main

import (
	"theme-engine/internal/core/context"
	"theme-engine/internal/core/log"
	"theme-engine/internal/core/themes/theme"
	"theme-engine/internal/helper"
	"theme-engine/internal/loader"
	"theme-engine/internal/processor"
)

func runSpecific(
	appName string,
	path *loader.ToolMap,
	rawTheme *theme.Theme,
	ctx context.Context,
) {
	var parsed any = nil;
	var err error = nil;

	processor, ok := processor.GetProcessor(appName);
	if !ok {
		helper.ExitWrapper(ExitGeneral, "Unknown processor %s", appName);
	}

	raw, hasConfig := rawTheme.Tools[appName];
	if hasConfig {
		parsed, err = processor.Parse(raw);
		helper.CheckErr(err, ExitParseFail, "Error on parsing %s", appName);
	}

	resolved, err := processor.Resolve(parsed, &ctx);
	helper.CheckErr(err, ExitResolve, "Error on resolving %s", appName);

	helper.CheckErr(processor.Render(path.TemplatePath, path.OutputPath, resolved), ExitRender, "Error on rendering %s", appName);
	log.Info("Success rendering %s", appName);
}


