#!/usr/bin/env bash

# ==================================== #
# |  Wrapper script...               | #
# ==================================== #

# Settings
QUIET=0
if ! [[ -t 1 ]]; then
  QUIET=1
else
  clear
fi
[[ -n "$VERBOSE" ]] && QUIET=0

log() {
  [[ $QUIET -eq 0 ]] && echo "$@"
}

mkdir -p "./cmd/bin"


# Path
BINARY="./cmd/bin/theme-engine"
SRC="./cmd/theme-engine"

# Args
mode=$1
opt=$2



notify_error() {
  case "$1" in
    1) notify-send -u critical "Theme Engine" "Script call error" ;;
    2) notify-send -u critical "Theme Engine" "Fail to load resource" ;;
    3) notify-send -u critical "Theme Engine" "I/O error" ;;
    4) notify-send -u critical "Theme Engine" "Config / parse error" ;;
    5) notify-send -u critical "Theme Engine" "Resolve error" ;;
    6) notify-send -u critical "Theme Engine" "Render error" ;;
    *) notify-send -u critical "Theme Engine" "Unknown error ($1)" ;;
  esac
}

usage() {
  log "Usage: $0 [mode] [theme|waybar]"
  log "Mode:"
  log "  run-bin   -> Run binary only"
  log "  run-raw   -> Run with go"
  log "  build-bin -> Build binary"
  log "  build-run -> Build and run binary"
  log "  clean     -> Delete binary"

  notify_error 1
}

if [ $# -eq 0 ]; then
  usage
  exit 1
fi

run() {
  "$@"
  code=$?

  if [[ $code -eq 0 ]]; then
    [[ $QUIET -eq 0 ]] && notify-send -t 2000 "Theme" "Success Rendering"
  else
    notify_error "$code"
  fi

  return $code
}


exitCode=0

case "$mode" in
  run-bin)
    log "[❯] Running.."
    run "$BINARY" "$opt"
    exitCode=$?
    log "[❯] Done."
    ;; 
  run-raw)
    log "[❯] Running.."
    run go run "$SRC" "$opt"
    exitCode=$?
    log "[❯] Done."
    ;;
  build-bin)
    log "[❯] Building..."
    go build -o "$BINARY" "$SRC" || {
      notify-send -u critical "Build failed"
      exitCode=2
    }

    log "[❯] Done."
    ;;
  build-run)
    log "[❯] Building..."
    go build -o "$BINARY" "$SRC"|| {
      notify-send -u critical "Build failed"
      exit 2
    }
 
    log "[❯] Running binary"
    run "$BINARY" "$opt"
    exitCode=$?
    log "[❯] Done."
    ;;
  clean)
    log "[❯] Cleaning binary..."
    rm -f "$BINARY"
    log "[❯] Done."
    ;;
  *)
    log "[!] Invalid mode!"
    usage
    ;;
esac

exit $exitCode

