# Theme Engine

Renderer tema cepat berbasis JSON untuk setup desktop Hyprland.

Theme Engine mengambil satu sumber tema, resolve variable palette, lalu render
file config untuk beberapa target seperti GTK, Cava, Foot, Kitty, Yazi,
Starship, dan Hyprland. Ide utamanya sederhana: data tema cukup satu sumber,
kerja runtime dibuat sekecil mungkin, dan tool baru gampang ditambahkan.

## Highlight

- Single source of truth dari `themes/<name>/theme.json` dan `palette.json`
- Render path cepat dengan cached template dan skip-write kalau output tidak berubah
- Mapping path sederhana lewat `config/path.txt` atau `$MYENV/map`
- Processor khusus untuk tool yang butuh config sendiri
- Processor template-only untuk target sederhana seperti Starship, Yazi, dan Nvim
- Atomic write supaya file config tidak pernah setengah tertulis
- Hot reload: fase apply setelah render mengirim config baru ke tool yang sedang berjalan
- CLI `set-theme <nama>` dan mode `watch` (polling stdlib, tanpa dependensi)
- Config lokal tersedia, jadi repo bisa dites tanpa file eksternal

## Cara Kerja

```txt
config/.state.json
        |
        v
themes/<theme>/palette.json + themes/<theme>/theme.json
        |
        v
internal/app/engine
        |
        +--> load path map
        +--> resolve palette variable
        +--> pilih processor
        +--> render template
        v
output/<target config>
```

CLI sengaja dibuat tipis. Orchestration utama ada di `internal/app/engine`,
sementara setiap target punya parsing dan resolving logic sendiri di
`internal/adapters/tools/<tool>` atau `internal/adapters/platform/<domain>`.
Wiring processor eksplisit ada di `cmd/theme-engine/wiring.go` (tanpa `init()`).

## Struktur Project

```txt
cmd/theme-engine/              Entrypoint CLI dan wiring processor eksplisit
internal/app/engine/           Pipeline utama load -> resolve -> render
internal/app/ports/            Interface Processor (Parser/Resolver/Renderer)
internal/domain/               Model murni: palette, theme, state, vars, renderctx
internal/adapters/tools/       Adapter per tool (kitty, foot, alacritty, cava, hypr)
internal/adapters/tools/generic/ Target template-only (nvim, yazi, starship)
internal/adapters/platform/    Domain besar: gtk (sassc), system (dconf apply.sh)
internal/infra/loader/         Loader JSON dan path map
internal/infra/renderer/       Template cache, atomic write, skip unchanged write
internal/infra/pathenv/        Expand $WAYBAR/$THEME/$ENV + resolve $pl.*
internal/infra/log/            Logger

assets/templates/              Template output
config/                        Config lokal untuk runtime/test
themes/                        Definisi theme dan palette
output/                        Hasil render config (golden, jangan diedit manual)
```

Tiap adapter tool berisi pola seragam:

```txt
internal/adapters/tools/<tool>/
  model.go    # bentuk input JSON dari theme.json
  view.go     # data final yang dikirim ke template
  adapter.go  # Parse(), Resolve(), Render() + New()
```

Helper bersama: `internal/adapters/tools/jsonx/jsonx.go` (`Decode`).
Kontrak: `Parse(nil)` -> zero `Raw` tanpa error, tipe non-`json.RawMessage`
-> error (tidak panic). `Resolve` missing/bukan `Raw` -> zero `Raw`
(`cava` tetap error bersih jika gradients < 5).

## Target Yang Didukung

| Target | Tipe | Adapter | Template |
| --- | --- | --- | --- |
| `gtk` | platform | `internal/adapters/platform/gtk` | `assets/templates/domain/gtk/source.tmpl` |
| `system` | platform | `internal/adapters/platform/system` | `assets/templates/domain/system/apply.tmpl` |
| `cava` | tool | `internal/adapters/tools/cava` | `assets/templates/tools/cava/cava.tmpl` |
| `foot` | tool | `internal/adapters/tools/foot` | `assets/templates/tools/foot/foot.tmpl` |
| `kitty` | tool | `internal/adapters/tools/kitty` | `assets/templates/tools/kitty/kitty.tmpl` |
| `alacritty` | tool | `internal/adapters/tools/alacritty` | `assets/templates/tools/alacritty/alacritty.tmpl` |
| `hypr` | tool | `internal/adapters/tools/hypr` | `assets/templates/tools/hypr/hypr.tmpl` |
| `yazi` | template-only | `internal/adapters/tools/generic` | `assets/templates/tools/yazi/theme.tmpl` |
| `nvim` | template-only + reload | `internal/adapters/tools/nvim` | `assets/templates/tools/nvim/colors.tmpl` |
| `starship` | template-only | `internal/adapters/tools/generic` | `assets/templates/tools/starship/starship.tmpl` |

## Quick Start

Render semua target memakai config lokal `config/`:

```bash
go run ./cmd/theme-engine
```

Render satu target:

```bash
go run ./cmd/theme-engine kitty
go run ./cmd/theme-engine hyprland
```

Ganti tema aktif (update `.state.json` + render semua target):

```bash
go run ./cmd/theme-engine set-theme harbor
go run ./cmd/theme-engine set-theme harbor --type light
```

Watch template/tema/map, render ulang otomatis saat berubah:

```bash
go run ./cmd/theme-engine watch
```

Atau pakai wrapper:

```bash
./runner.sh build
./runner.sh run kitty
./runner.sh hyprland
```

Jalankan test:

```bash
go test ./...
```

Kalau Go cache tidak bisa ditulis karena environment sandbox:

```bash
GOCACHE=/tmp/go-build go test ./...
GOCACHE=/tmp/go-build go run ./cmd/theme-engine
```

## Konfigurasi

Engine mencari direktori config dengan urutan:

1. `THEME_ENGINE_MAP`
2. `config/` lokal
3. `$MYENV/map`

Artinya development bisa pakai file lokal di repo, sementara setup desktop asli
tetap bisa menunjuk ke map eksternal.

Contoh:

```bash
THEME_ENGINE_MAP="$MYENV/map" go run ./cmd/theme-engine
```

## Hot Reload

Setelah render sukses, engine menjalankan fase **apply** yang memberitahu
tool yang sedang berjalan bahwa config-nya baru saja berubah. Fase ini
opsional dan non-fatal: tool tidak terinstal atau mekanismenya tidak
aktif hanya menghasilkan warning, render tetap sukses.

Fase apply aktif otomatis untuk map eksternal (`$THEME_ENGINE_MAP` /
`$MYENV/map`) dan mati untuk `config/` lokal (dev/test). Atur manual:

```bash
THEME_ENGINE_APPLY=1 go run ./cmd/theme-engine   # selalu apply
THEME_ENGINE_APPLY=0 go run ./cmd/theme-engine   # selalu matikan
```

Mekanisme per tool:

| Tool | Cara apply | Otomatis? |
| --- | --- | --- |
| `hypr` | `hyprctl reload` | Ya (juga inotify saat save in-place) |
| `alacritty` | live config reload — engine menulis **in-place** (bukan rename) | Ya |
| `kitty` | `kitty @ set-colors -a` (remote control) | Tidak — butuh `allow_remote_control yes` |
| `cava` | `pkill -SIGUSR2 cava` (reload warna saja) | Tidak — sinyal |
| `foot` | OSC 10/11/12 + 4;N ke semua `/dev/pts/*` (set_term_colors) | Ya — foot tidak punya live reload (issue #1653) |
| `nvim` | `nvim --server $NVIM_LISTEN_ADDRESS --remote-send :colorscheme` | Tidak — butuh nvim `--listen` |
| `starship`, `yazi`, `rofi` | baca config saat startup/prompt | Ya, inherent |
| `lazygit` | restart instance | Tidak ada IPC |
| GTK apps / waybar | restart app; waybar: `"reload_style_on_change": true` | Tidak |

Trade-off: `alacritty` ditulis in-place (`RenderInPlace`) karena
Alacritty mem-watch config per inode via inotify — write temp+rename
memutus watcher (alacritty#5355). Risiko crash mid-write diterima untuk
file config sekecil ini.

### `tools/set_term` — set_term_colors mandiri

Script standalone (`tools/set_term`, Python 3, nol dependensi) yang
meng-apply palette ke semua terminal berjalan via OSC — berguna setelah
mengedit `palette.json` secara manual atau dari keybind, tanpa jalan
engine. Memakai lookup map-dir yang sama (`$THEME_ENGINE_MAP` →
`config/` → `$MYENV/map`) dan resolve `$pl.` yang sama dengan Go —
`internal/adapters/tools/foot/set_term_test.go` men-assert output
`--dry-run`-nya identik dengan `oscSequence` engine untuk semua tema ×
varian. `.state.json` hanya dibaca untuk default; tema + varian
eksplisit jalan tanpa map dir:

```bash
python3 tools/set_term                     # tema aktif dari .state.json
python3 tools/set_term harbor --variant dark
python3 tools/set_term --dry-run           # print sequence, tidak menulis
```

`tools/set_term.sh` adalah twin bash-nya — palette di-parse dengan
awk/sed (tanpa python/jq), API sama. Berisi fungsi `set_term_colors`
(memakai `FG`, `BG`, `CURSOR`, `COLOR0..COLOR15`) yang bisa di-`source`
ke `.bashrc`; menjalankannya (bukan sourcing) memakai CLI yang sama:

```bash
tools/set_term.sh harbor --variant dark    # API sama dengan versi python
source tools/set_term.sh                   # hanya definisikan fungsi
set_term_colors                            # terapkan FG/BG/CURSOR/COLOR0..15
```

Untuk tool baru, tambahkan metode ke adapter-nya:

```go
func (Processor) Reload(outputPath string, ctx *renderctx.Context) error {
	// kirim config baru ke instance yang berjalan
}
```

Engine otomatis memanggilnya lewat interface `ports.Reloader` —
adapter tanpa `Reload` dilewati tanpa error.

## Watch Mode

`theme-engine watch` polling (stdlib `os.Stat`, tanpa fsnotify) direktori
`assets/templates/`, `themes/<aktif>/`, dan map directory setiap 500ms.
Perubahan di-debounce 300ms (save beruntun editor di-coalesce), template
cache di-clear, lalu semua target di-render ulang. Error saat tema sedang
diedit (JSON invalid) hanya di-log — watching terus jalan. SIGINT/SIGTERM
keluar bersih.

```bash
go run ./cmd/theme-engine watch
```

## Runner Script

`runner.sh` adalah wrapper tipis untuk pemakaian harian, terminal, dan tombol
desktop/keybind.

Mode utama:

```bash
./runner.sh build             # build binary ke cmd/bin/theme-engine
./runner.sh run [target]      # jalankan binary, paling cocok untuk tombol/keybind
./runner.sh dev [target]      # jalankan go run, cocok saat development
./runner.sh build-run [target]
./runner.sh test
./runner.sh clean
./runner.sh set-term [theme]  # apply palette ke terminal berjalan (tools/set_term)
```

Shorthand ini juga valid:

```bash
./runner.sh kitty
./runner.sh hyprland
```

Saat dijalankan dari terminal, runner menulis output ke terminal. Saat dijalankan
tanpa terminal, runner otomatis diam dan memakai `notify-send` kalau tersedia.
Ini cocok untuk tombol switch theme.

Untuk memaksa log tetap keluar walau tanpa terminal:

```bash
VERBOSE=1 ./runner.sh run kitty
```

Untuk mengatur notifikasi:

```bash
THEME_ENGINE_NOTIFY=auto ./runner.sh run kitty
THEME_ENGINE_NOTIFY=1 ./runner.sh run kitty
THEME_ENGINE_NOTIFY=0 ./runner.sh run kitty
```

Untuk path cepat di tombol/keybind, build dulu sekali lalu panggil mode `run`:

```bash
./runner.sh build
./runner.sh run hyprland
```

Mode `set-term` tidak perlu build — langsung pakai Python:

```bash
./runner.sh set-term              # tema aktif
./runner.sh set-term harbor       # tema tertentu
./runner.sh set-term --dry-run    # print sequence, tidak menulis
```

### `config/.state.json`

```json
{
  "version": 1,
  "theme": {
    "name": "nocturne",
    "type": "dark"
  },
  "waybar": "default"
}
```

`theme.name` memilih folder `themes/<name>/`.

`theme.type` memilih variant palette, misalnya `dark` atau `light`.

`waybar` dipakai untuk mengganti `$WAYBAR` di path map.

### `config/path.txt`

Setiap baris memetakan satu target ke satu template dan satu output:

```txt
target|template_path|output_path
```

Map lokal saat ini:

```txt
gtk|assets/templates/domain/gtk/source.tmpl|output/domain/gtk
cava|assets/templates/tools/cava/cava.tmpl|output/tools/cava/cava_extra
foot|assets/templates/tools/foot/foot.tmpl|output/tools/foot/ui.ini
kitty|assets/templates/tools/kitty/kitty.tmpl|output/tools/kitty/ui.conf
alacritty|assets/templates/tools/alacritty/alacritty.tmpl|output/tools/alacritty/ui.toml
hypr|assets/templates/tools/hypr/hypr.tmpl|output/tools/hypr/source.conf
yazi|assets/templates/tools/yazi/theme.tmpl|output/tools/yazi/theme.toml
nvim|assets/templates/tools/nvim/colors.tmpl|output/tools/nvim/colors.lua
starship|assets/templates/tools/starship/starship.tmpl|output/tools/starship/starship.toml
```

Baris kosong dan komentar yang diawali `#` akan diabaikan.

## File Theme

Data theme ada di:

```txt
themes/<theme-name>/theme.json
themes/<theme-name>/palette.json
```

`palette.json` berisi warna dasar dan warna turunan.

`theme.json` berisi metadata global dan config opsional per tool:

```json
{
  "theme": {
    "name": "nocturne",
    "fonts": {
      "primary": "FiraCode Nerd Font",
      "secondary": "Maple Mono Nerd Font",
      "size": 11.5
    }
  },
  "tools": {
    "kitty": {
      "cursorShape": "beam",
      "opacity": 0.7
    }
  }
}
```

Variable palette memakai prefix `$pl.`:

```json
"active": "$pl.extra.accent.primary"
```

## Template

Template memakai Go `text/template`.

Data umum yang bisa dipakai di template:

```gotemplate
{{ .Palette.Foreground }}
{{ .Palette.AccentPrimary }}
{{ .Palette.Background }}
{{ .Theme.Theme.Fonts.Primary }}
{{ .Theme.Theme.Fonts.Size }}
```

Renderer juga menyediakan helper `hex` untuk menghapus awalan `#`:

```gotemplate
rgb({{ hex .Palette.AccentPrimary }})
```

## Menambah Tool Baru

Ada dua jalur.

### 1. Tool Template-Only

Pakai ini kalau target cuma butuh data global dari theme dan palette.

Contoh: menambah `dunst`.

Buat template:

```txt
assets/templates/tools/dunst/dunstrc.tmpl
```

Tambah entry path map:

```txt
dunst|assets/templates/tools/dunst/dunstrc.tmpl|output/tools/dunst/dunstrc
```

Register target di `cmd/theme-engine/wiring.go`:

```go
procs["dunst"] = generic.New()
```

Itu sudah cukup untuk template yang hanya butuh `.Palette` dan `.Theme`.

### 2. Tool Dengan Config Khusus

Pakai ini kalau target butuh config khusus dari `theme.json`.

Buat folder:

```txt
internal/adapters/tools/dunst/
  model.go
  view.go
  adapter.go
```

Tanggung jawab setiap file:

| File | Fungsi |
| --- | --- |
| `model.go` | Bentuk input JSON dari `theme.json` |
| `view.go` | Data final yang dikirim ke template |
| `adapter.go` | Logic parse, resolve, dan render + `New()` (pakai `jsonx.Decode`, nil-safe) |

Register adapter di `cmd/theme-engine/wiring.go`:

```go
procs["dunst"] = dunst.New()
```

Tambah config tool ke `themes/<theme>/theme.json`:

```json
"dunst": {
  "radius": 8,
  "border": "$pl.extra.accent.primary"
}
```

Tambah entry path map:

```txt
dunst|assets/templates/tools/dunst/dunstrc.tmpl|output/tools/dunst/dunstrc
```

## Catatan Performa

Render path dioptimalkan untuk repeated run dan future hot reload:

- Template diparse sekali per path lalu disimpan di cache.
- Hasil render masuk buffer dulu sebelum ditulis.
- Output tidak ditulis ulang kalau isinya sama.
- Write dilakukan atomic lewat temp file lalu rename.
- Lookup variable palette memakai flattened map.
- Path map tetap plain text supaya ringan dan gampang diedit.

Cost runtime terbesar biasanya datang dari tool eksternal atau command reload
setelah render, bukan dari parsing JSON atau loading path map.

## Catatan Naming

- `app/engine` berarti pipeline utama load -> resolve -> render.
- `app/ports` berarti kontrak Processor (Parser/Resolver/Renderer).
- `domain` berarti model murni tanpa IO: palette, theme, state, vars, renderctx.
- `adapters/tools` berarti target spesifik aplikasi seperti Kitty atau Cava.
- `adapters/platform` berarti domain config yang lebih luas (gtk, system).
- `adapters/tools/generic` berarti target template-only tanpa custom JSON parser (ex-`static`).
- `adapters/platform/system` sengaja pisah dari `generic` walau sama-sama
  template-only: `generic` untuk `output/tools/*`, `system` untuk
  `output/domain/system/apply.sh` (script dconf). Duplikasi 3 method
  disengaja demi kepemilikan yang jelas.
- `infra/loader` berarti loading filesystem, JSON, dan path map.
- `infra/renderer` berarti eksekusi template dan penulisan output.
- `infra/pathenv` berarti expand path + resolve variable `$pl.*`.
- `model.go` = input JSON, `view.go` = data final template, `adapter.go` = Parse/Resolve/Render.

## Status Saat Ini

- Render semua target: jalan
- Render target spesifik: jalan
- Hot reload (fase apply per tool): jalan
- `set-theme` dan mode `watch`: jalan
- Config lokal: tersedia di `config/`
- Test: coverage fokus untuk loader, renderer, foot OSC, kitty colour extraction, set-theme
- Theme: nocturne, ghostly, kanagawa-wave, kanagawa-dragon,
  kanagawa-dragon-original, dracula, sakura
- Tools: gtk, cava, foot, kitty, alacritty, hypr, yazi, nvim, starship
