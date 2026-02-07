package log

import (
	"fmt"
	"os"
)

// var verb = false;
var verb = true;

func SetVerbose(v bool) {
	verb = v;
}

func Info(msg string, args ...any) {
	if !verb {
		return;
	}

	fmt.Fprintf(os.Stdout, "[INFO]: "+ msg + "\n", args...);
}

func Warn(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "[WARNING]: " + msg + "\n", args...);
}

func Error(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "[ERROR]: " + msg + "\n", args...);
}
