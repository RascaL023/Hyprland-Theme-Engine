package helper

import (
	"os"
	"theme-engine/internal/core/log"
)

func CheckErr(err error, code int, msg string, args ...any) {
	if err != nil {
		log.Error(msg + ": %v", append(args, err)...);
		os.Exit(code);
	}
}

func ExitWrapper(code int, msg string, args ...any) {
	log.Error(msg, args...);
	os.Exit(code);
}
