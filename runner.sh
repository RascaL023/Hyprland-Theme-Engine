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

error() {
  echo "Error: $@" >&2
  notify-send -t 3000 "Error!" "$@"
}


# Path
BINARY="./cmd/bin/gtr"
SRC="./cmd/gtr"

# Args
mode=$1
opt=$2

mkdir -p "./cmd/bin"

usage() {
  log "Usage: $0 [mode] [theme|waybar]"
  log "Mode:"
  log "  run-bin   -> Run binary only"
  log "  run-raw   -> Run with go"
  log "  build-bin -> Build binary"
  log "  build-run -> Build and run binary"
  log "  clean     -> Delete binary"

  error "script error ecounter!"
  exit 1
}

if [ $# -eq 0 ]; then
  usage
fi

case "$mode" in
  run-bin)
    log "[❯] Running.."
    "$BINARY" "$opt"
    log "[❯] Done."
    ;;  
  run-raw)
    log "[❯] Running.."
    go run "$SRC" "$opt"
    log "[❯] Done."
    ;;  
  build-bin)
    log "[❯] Building..."
    go build -o "$BINARY" "$SRC"
    log "[❯] Done."
    ;;  
  build-run)
    log "[❯] Building..."
    go build -o "$BINARY" "$SRC"
    log "[❯] Running binary"
    "./$BINARY" "$opt"
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

exit 0

