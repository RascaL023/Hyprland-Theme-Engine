package main

import (
	"os"

	"theme-engine/internal/engine"
	"theme-engine/internal/helper"
)

const (
	ExitOK        = 0
	ExitGeneral   = 1
	ExitLoad      = 2
	ExitIOError   = 3
	ExitParseFail = 4
	ExitResolve   = 5
	ExitRender    = 6
)

func main() {
	target := ""
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	app, err := engine.New(engine.Config{MapDir: engine.DefaultMapDir()})
	helper.CheckErr(err, ExitLoad, "Error on initializing engine")

	if target == "" {
		helper.CheckErr(app.RunAll(), ExitGeneral, "Error on rendering all targets")
		return
	}

	helper.CheckErr(app.Run(target), ExitGeneral, "Error on rendering %s", target)
}
