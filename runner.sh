#!/usr/bin/env bash

set -u

APP_NAME="Theme Engine"
BIN_DIR="./cmd/bin"
BINARY="$BIN_DIR/theme-engine"
SRC="./cmd/theme-engine"

EXIT_USAGE=1
EXIT_BUILD=10

is_tty=0
if [[ -t 1 ]]; then
  is_tty=1
fi

quiet=0
if [[ "$is_tty" -eq 0 ]]; then
  quiet=1
fi
if [[ "${VERBOSE:-}" != "" ]]; then
  quiet=0
fi

notify_enabled=0
if command -v notify-send >/dev/null 2>&1; then
  if [[ "${THEME_ENGINE_NOTIFY:-auto}" == "1" ]]; then
    notify_enabled=1
  elif [[ "${THEME_ENGINE_NOTIFY:-auto}" == "auto" && "$is_tty" -eq 0 ]]; then
    notify_enabled=1
  fi
fi

log() {
  [[ "$quiet" -eq 0 ]] && printf '%s\n' "$*"
}

notify_info() {
  [[ "$notify_enabled" -eq 1 ]] || return 0
  notify-send -t 1800 "$APP_NAME" "$1" >/dev/null 2>&1 || true
}

notify_error() {
  [[ "$notify_enabled" -eq 1 ]] || return 0

  local code="${1:-1}"
  local msg
  case "$code" in
    1) msg="Script call error" ;;
    2) msg="Fail to load resource" ;;
    3) msg="I/O error" ;;
    4) msg="Config / parse error" ;;
    5) msg="Resolve error" ;;
    6) msg="Render error" ;;
    7) msg="Usage error" ;;
    10) msg="Build failed" ;;
    127) msg="Binary not found. Run: ./runner.sh build" ;;
    *) msg="Unknown error ($code)" ;;
  esac

  notify-send -u critical "$APP_NAME" "$msg" >/dev/null 2>&1 || true
}

usage() {
  cat <<'EOF'
Usage:
  ./runner.sh [target]
  ./runner.sh run [target]
  ./runner.sh dev [target]
  ./runner.sh build
  ./runner.sh build-run [target]
  ./runner.sh test
  ./runner.sh clean
  ./runner.sh set-term [theme]

Modes:
  run        Run compiled binary. Best for keybinds/buttons.
  dev        Run with go run. Best while editing code.
  build      Build ./cmd/bin/theme-engine.
  build-run  Build then run compiled binary.
  test       Run go test ./...
  clean      Delete compiled binary.
  set-term   Apply the active palette to running terminals
             (tools/set_term — OSC sequences, no engine run).

The binary also accepts subcommands (passed through by run/dev):
  theme-engine set-theme <name> [--type dark|light]
  theme-engine watch

Compatibility aliases:
  run-bin, run-raw, build-bin

Environment:
  THEME_ENGINE_MAP      Override config map directory.
  THEME_ENGINE_APPLY    auto, 1, or 0. Post-render apply phase. Default: auto.
  THEME_ENGINE_NOTIFY   auto, 1, or 0. Default: auto.
  VERBOSE               Print logs even without a terminal.
EOF
}

run_and_report() {
  "$@"
  local code=$?

  if [[ "$code" -eq 0 ]]; then
    notify_info "Render success"
  else
    notify_error "$code"
  fi

  return "$code"
}

build_binary() {
  mkdir -p "$BIN_DIR"
  log "[build] $BINARY"
  if ! go build -o "$BINARY" "$SRC"; then
    notify_error "$EXIT_BUILD"
    return "$EXIT_BUILD"
  fi
}

run_binary() {
  if [[ ! -x "$BINARY" ]]; then
    log "[error] Binary not found: $BINARY"
    log "        Run: ./runner.sh build"
    notify_error 127
    return 127
  fi

  log "[run] $BINARY $*"
  run_and_report "$BINARY" "$@"
}

run_dev() {
  log "[dev] go run $SRC $*"
  run_and_report go run "$SRC" "$@"
}

run_tests() {
  log "[test] go test ./..."
  go test ./...
}

run_set_term() {
  log "[set-term] tools/set_term $*"
  python3 tools/set_term "$@"
}

mode="${1:-run}"
if [[ "$#" -gt 0 ]]; then
  shift
fi

case "$mode" in
  run|run-bin)
    run_binary "$@"
    ;;
  dev|run-raw)
    run_dev "$@"
    ;;
  build|build-bin)
    build_binary
    ;;
  build-run)
    build_binary && run_binary "$@"
    ;;
  test)
    run_tests
    ;;
  set-term)
    run_set_term "$@"
    ;;
  clean)
    log "[clean] $BINARY"
    rm -f "$BINARY"
    ;;
  help|-h|--help)
    usage
    ;;
  *)
    run_binary "$mode" "$@"
    ;;
esac

