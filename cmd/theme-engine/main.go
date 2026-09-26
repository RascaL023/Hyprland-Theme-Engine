package main

import (
	"os"

	"theme-engine/internal/app/engine"
	"theme-engine/internal/infra/log"
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

func fail(code int, msg string, args ...any) {
	log.Error(msg, args...)
	os.Exit(code)
}

func main() {
	target := ""
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	app, err := engine.New(engine.Config{
		MapDir:     engine.DefaultMapDir(),
		Processors: defaultProcessors(),
	})
	if err != nil {
		fail(ExitLoad, "Error on initializing engine: %v", err)
	}

	if target == "" {
		if err := app.RunAll(); err != nil {
			fail(ExitGeneral, "Error on rendering all targets: %v", err)
		}
		return
	}

	if err := app.Run(target); err != nil {
		fail(ExitGeneral, "Error on rendering %s: %v", target, err)
	}
}
