#!/usr/bin/env bash
# Apply a theme palette to every running terminal (set_term_colors).
#
# Bash twin of tools/set_term: foot has no live config reload
# (codeberg.org/dnkl/foot/issues/1653), but its emulator — like
# alacritty, kitty and xterm — honours OSC sequences at runtime.
# Data written to a pty slave is delivered to the emulator that
# owns it, so every accessible terminal session picks up the new
# colours without a restart.
#
# The palette is read from themes/<name>/palette.json. The default
# theme/variant comes from .state.json (same map-dir lookup as the
# engine) — but only when not given explicitly, so an explicit
# theme --variant pair works with no map dir at all. $pl. references
# are resolved the same way the Go resolver does. Parsing uses only
# awk/sed — no jq, no python.
#
# Usage:
#   tools/set_term.sh                     # active theme from .state.json
#   tools/set_term.sh harbor              # explicit theme
#   tools/set_term.sh --variant dark
#   tools/set_term.sh --map-dir "$MYENV/map"
#   tools/set_term.sh --dry-run           # print the sequence, write nothing
#
# Sessions owned by another user are skipped (EACCES). With no
# terminal sessions at all this is a no-op.
#
# Running the script executes the CLI below; sourcing it only
# defines set_term_colors (and its helpers) for a .bashrc.

REPO_ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
THEMES_DIR="$REPO_ROOT/themes"

die() {
  printf 'set_term: %s\n' "$*" >&2
  exit 1
}

# ------------------------------------------------------------------ sequence

# set_term_colors writes the foot set_term_colors sequence — OSC 10
# (foreground), OSC 11 (background), OSC 12 (cursor) and OSC 4;N for
# the 16 ANSI colours — to every terminal session we can write.
# Colours come from FG, BG, CURSOR and COLOR0..COLOR15.
#
# The printf calls are inlined (not shared with print_sequence) so the
# function stays pasteable into a .bashrc on its own — keep the two
# in sync.
set_term_colors() {
  local pts
  SET_TERM_UPDATED=0
  SET_TERM_SKIPPED=0
  for pts in /dev/pts/[0-9]*; do
    [[ -e "$pts" ]] || continue
    if {
      printf '\e]10;%s\e\\' "$FG"      # foreground
      printf '\e]11;%s\e\\' "$BG"      # background
      printf '\e]12;%s\e\\' "$CURSOR"  # cursor
      printf '\e]4;0;%s\e\\'  "$COLOR0"
      printf '\e]4;1;%s\e\\'  "$COLOR1"
      printf '\e]4;2;%s\e\\'  "$COLOR2"
      printf '\e]4;3;%s\e\\'  "$COLOR3"
      printf '\e]4;4;%s\e\\'  "$COLOR4"
      printf '\e]4;5;%s\e\\'  "$COLOR5"
      printf '\e]4;6;%s\e\\'  "$COLOR6"
      printf '\e]4;7;%s\e\\'  "$COLOR7"
      printf '\e]4;8;%s\e\\'  "$COLOR8"
      printf '\e]4;9;%s\e\\'  "$COLOR9"
      printf '\e]4;10;%s\e\\' "$COLOR10"
      printf '\e]4;11;%s\e\\' "$COLOR11"
      printf '\e]4;12;%s\e\\' "$COLOR12"
      printf '\e]4;13;%s\e\\' "$COLOR13"
      printf '\e]4;14;%s\e\\' "$COLOR14"
      printf '\e]4;15;%s\e\\' "$COLOR15"
    } 2>/dev/null > "$pts"; then
      (( SET_TERM_UPDATED++ ))
    else
      (( SET_TERM_SKIPPED++ ))
    fi
  done
}

# print_sequence emits the same bytes on stdout (--dry-run).
print_sequence() {
  printf '\e]10;%s\e\\' "$FG"      # foreground
  printf '\e]11;%s\e\\' "$BG"      # background
  printf '\e]12;%s\e\\' "$CURSOR"  # cursor
  printf '\e]4;0;%s\e\\'  "$COLOR0"
  printf '\e]4;1;%s\e\\'  "$COLOR1"
  printf '\e]4;2;%s\e\\'  "$COLOR2"
  printf '\e]4;3;%s\e\\'  "$COLOR3"
  printf '\e]4;4;%s\e\\'  "$COLOR4"
  printf '\e]4;5;%s\e\\'  "$COLOR5"
  printf '\e]4;6;%s\e\\'  "$COLOR6"
  printf '\e]4;7;%s\e\\'  "$COLOR7"
  printf '\e]4;8;%s\e\\'  "$COLOR8"
  printf '\e]4;9;%s\e\\'  "$COLOR9"
  printf '\e]4;10;%s\e\\' "$COLOR10"
  printf '\e]4;11;%s\e\\' "$COLOR11"
  printf '\e]4;12;%s\e\\' "$COLOR12"
  printf '\e]4;13;%s\e\\' "$COLOR13"
  printf '\e]4;14;%s\e\\' "$COLOR14"
  printf '\e]4;15;%s\e\\' "$COLOR15"
}

# ------------------------------------------------------------------ palette

# variant_block prints the "dark"/"light" object of a palette.json.
# Blocks sit at a 4-space indent and hold only 6-space-or-deeper
# lines, so the first 4-space close ends the block.
variant_block() {
  awk -v v="$1" '
    $0 ~ "^    \"" v "\": \\{" { inside = 1; next }
    inside && /^    \}/ { inside = 0 }
    inside { print }
  '
}

# colors_block prints the entries of the "colors" array.
colors_block() {
  awk '/"colors": *\[/ { inside = 1; next }
       inside && /\]/ { inside = 0 }
       inside { print }'
}

# field prints the first value of a "key": "value" line.
field() {
  sed -n "s/.*\"$1\": *\"\\([^\"]*\\)\".*/\\1/p" | head -n 1
}

# resolve_var: single-pass $pl. resolution, same contract as the Go
# resolver — palette.json may only reference the base palette
# (foreground, background, cursor, color0..color15).
resolve_var() {
  case "$1" in
    '$pl.foreground') printf '%s' "$RAW_FOREGROUND" ;;
    '$pl.background') printf '%s' "$RAW_BACKGROUND" ;;
    '$pl.cursor')     printf '%s' "$RAW_CURSOR" ;;
    '$pl.color'*)
      local idx="${1#\$pl.color}"
      if [[ "$idx" =~ ^[0-9]+$ ]] && (( idx < ${#RAW_COLORS[@]} )); then
        resolve_var "${RAW_COLORS[idx]}"
      else
        printf '%s' "$1"
      fi
      ;;
    *) printf '%s' "$1" ;;
  esac
}

load_palette() {  # $1 = theme, $2 = variant; sets FG BG CURSOR COLOR*
  local palette="$THEMES_DIR/$1/palette.json" block i
  [[ -r "$palette" ]] || die "no such theme: $1"

  block="$(variant_block "$2" < "$palette")"
  [[ -n "$block" ]] || die "$1 has no '$2' palette"

  RAW_FOREGROUND="$(printf '%s\n' "$block" | field foreground)"
  RAW_BACKGROUND="$(printf '%s\n' "$block" | field background)"
  RAW_CURSOR="$(printf '%s\n' "$block" | field cursor)"
  [[ -n "$RAW_FOREGROUND" && -n "$RAW_BACKGROUND" && -n "$RAW_CURSOR" ]] \
    || die "$1/$2: incomplete base palette"

  mapfile -t RAW_COLORS < <(
    printf '%s\n' "$block" | colors_block | sed -n 's/.*"\(#[^"]*\)".*/\1/p'
  )
  (( ${#RAW_COLORS[@]} >= 7 )) \
    || die "$1/$2: palette needs at least 7 colors, got ${#RAW_COLORS[@]}"

  FG="$(resolve_var "$RAW_FOREGROUND")"
  BG="$(resolve_var "$RAW_BACKGROUND")"
  CURSOR="$(resolve_var "$RAW_CURSOR")"
  for ((i = 0; i < ${#RAW_COLORS[@]}; i++)); do
    printf -v "COLOR$i" '%s' "$(resolve_var "${RAW_COLORS[i]}")"
  done
}

# ------------------------------------------------------------------ state

default_map_dir() {  # same lookup order as engine.DefaultMapDir()
  if [[ -n "${THEME_ENGINE_MAP:-}" ]]; then
    printf '%s' "$THEME_ENGINE_MAP"
  elif [[ -f config/path.txt ]]; then
    printf '%s' config
  elif [[ -n "${MYENV:-}" ]]; then
    printf '%s' "$MYENV/map"
  else
    printf '%s' config
  fi
}

state_field() {  # $1 = map dir, $2 = field of the "theme" object
  awk '/"theme": *\{/ { inside = 1; next }
       inside && /^  \}/ { inside = 0 }
       inside { print }' "$1/.state.json" | field "$2"
}

# ------------------------------------------------------------------ cli

usage() {
  cat <<'EOF'
Usage: tools/set_term.sh [theme] [options]

Apply a theme palette to every running terminal (set_term_colors).

Options:
  --variant dark|light   palette variant (default: active variant)
  --map-dir DIR          map directory holding .state.json
  --dry-run              print the OSC sequence, write nothing
EOF
}


# Executing (not sourcing) the script runs the CLI.
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then

set -u

THEME=""
VARIANT=""
MAP_DIR="$(default_map_dir)"
DRY_RUN=0

while (($#)); do
  case "$1" in
    --variant)   (($# >= 2)) || die "--variant needs a value"; VARIANT="$2"; shift 2 ;;
    --variant=*) VARIANT="${1#*=}"; shift ;;
    --map-dir)   (($# >= 2)) || die "--map-dir needs a value"; MAP_DIR="$2"; shift 2 ;;
    --map-dir=*) MAP_DIR="${1#*=}"; shift ;;
    --dry-run)   DRY_RUN=1; shift ;;
    -h|--help)   usage; exit 0 ;;
    -*)          die "unknown option: $1" ;;
    *)           [[ -z "$THEME" ]] || die "unexpected argument: $1"; THEME="$1"; shift ;;
  esac
done

[[ -z "$VARIANT" || "$VARIANT" == dark || "$VARIANT" == light ]] \
  || die "--variant must be dark or light"

# .state.json is only the source of the defaults: an explicit theme
# and variant make the script usable without any map dir.
if [[ -z "$THEME" || -z "$VARIANT" ]]; then
  [[ -r "$MAP_DIR/.state.json" ]] || die "no such file: $MAP_DIR/.state.json"
  [[ -n "$THEME" ]] || THEME="$(state_field "$MAP_DIR" name)"
  [[ -n "$VARIANT" ]] || VARIANT="$(state_field "$MAP_DIR" type)"
  [[ -n "$THEME" ]] || die "$MAP_DIR/.state.json: no active theme"
fi
VARIANT="${VARIANT:-dark}"

load_palette "$THEME" "$VARIANT"

if (( DRY_RUN )); then
  print_sequence
  exit 0
fi

set_term_colors
if (( SET_TERM_SKIPPED > 0 )); then
  printf 'set_term: %s/%s -> %d terminal(s) updated, %d skipped (permission)\n' \
    "$THEME" "$VARIANT" "$SET_TERM_UPDATED" "$SET_TERM_SKIPPED"
else
  printf 'set_term: %s/%s -> %d terminal(s) updated\n' \
    "$THEME" "$VARIANT" "$SET_TERM_UPDATED"
fi

fi
