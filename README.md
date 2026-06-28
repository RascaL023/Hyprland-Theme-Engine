# Theme Engine

Renderer tema cepat berbasis JSON untuk setup desktop Hyprland.

Theme Engine mengambil satu sumber tema, resolve variable palette, lalu render
file config untuk beberapa target seperti GTK, Cava, Foot, Kitty, Waybar, dan
Hyprland. Ide utamanya sederhana: data tema cukup satu sumber, kerja runtime
dibuat sekecil mungkin, dan tool baru gampang ditambahkan.

## Highlight

- Single source of truth dari `themes/<name>/theme.json` dan `palette.json`
- Render path cepat dengan cached template dan skip-write kalau output tidak berubah
- Mapping path sederhana lewat `config/path.txt` atau `$MYENV/map`
- Processor khusus untuk tool yang butuh config sendiri
- Processor template-only untuk target sederhana seperti Waybar dan Hyprland
- Atomic write supaya file config tidak pernah setengah tertulis
- Config lokal tersedia, jadi repo bisa dites tanpa file eksternal

## Cara Kerja

```txt
config/.state.json
        |
        v
themes/<theme>/palette.json + themes/<theme>/theme.json
        |
        v
internal/engine
        |
        +--> load path map
        +--> resolve palette variable
        +--> pilih processor
        +--> render template
        v
output/<target config>
```

CLI sengaja dibuat tipis. Orchestration utama ada di `internal/engine`,
sementara setiap target punya parsing dan resolving logic sendiri di
`internal/core/tools/<tool>`.

## Struktur Project

```txt
cmd/theme-engine/              Entrypoint CLI dan import processor
internal/engine/               Pipeline utama load -> resolve -> render
internal/loader/               Loader JSON dan path map
internal/renderer/             Template cache, atomic write, skip unchanged write
internal/resolver/             Resolver variable, contoh $pl.extra.layer.base
internal/processor/            Registry processor
internal/register/             Interface bersama untuk processor

internal/core/themes/          Model data theme, palette, dan state
internal/core/tools/           Processor spesifik per tool
internal/core/tools/static/    Target template-only
internal/core/domain/          Processor domain besar, saat ini GTK

assets/templates/              Template output
assets/sources/                Script helper dan data warna
config/                        Config lokal untuk runtime/test
themes/                        Definisi theme dan palette
output/                        Hasil render config
```

## Target Yang Didukung

| Target | Tipe | Processor | Template |
| --- | --- | --- | --- |
| `gtk` | domain processor | `internal/core/domain/gtk` | `assets/templates/domain/gtk/source.tmpl` |
| `cava` | tool processor | `internal/core/tools/cava` | `assets/templates/tools/cava/cava.tmpl` |
| `foot` | tool processor | `internal/core/tools/foot` | `assets/templates/tools/foot/foot.tmpl` |
| `kitty` | tool processor | `internal/core/tools/kitty` | `assets/templates/tools/kitty/kitty.tmpl` |
| `waybar` | template-only | `internal/core/tools/static` | `assets/templates/waybar/default.tmpl` |
| `hyprland` | template-only | `internal/core/tools/static` | `assets/templates/hypr/hyprland.tmpl` |

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
waybar|assets/templates/waybar/$WAYBAR.tmpl|output/waybar/sources.css
hyprland|assets/templates/hypr/hyprland.tmpl|output/hypr/hyprland.conf
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

Register target di `internal/core/tools/static/processor.go`:

```go
processor.RegisterProcessor(Processor{name: "dunst"})
```

Itu sudah cukup untuk template yang hanya butuh `.Palette` dan `.Theme`.

### 2. Tool Dengan Config Khusus

Pakai ini kalau target butuh config khusus dari `theme.json`.

Buat folder:

```txt
internal/core/tools/dunst/
  raw.go
  dunst.go
  processor.go
```

Tanggung jawab setiap file:

| File | Fungsi |
| --- | --- |
| `raw.go` | Bentuk input JSON dari `theme.json` |
| `dunst.go` | Data final yang dikirim ke template |
| `processor.go` | Logic parse, resolve, dan render |

Register package di `cmd/theme-engine/import.go`:

```go
_ "theme-engine/internal/core/tools/dunst"
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

Sebagian besar naming sengaja dibuat polos:

- `engine` berarti pipeline utama.
- `loader` berarti loading dari filesystem, JSON, dan path map.
- `renderer` berarti eksekusi template dan penulisan output.
- `processor` berarti implementasi renderer untuk satu target.
- `resolver` berarti resolve variable.
- `tools` berarti target spesifik aplikasi seperti Kitty atau Cava.
- `domain` berarti domain config yang lebih luas dan bisa menghasilkan beberapa file,
  saat ini GTK/SCSS/Rasi/CSS.
- `static` berarti target template-only tanpa custom JSON parser.

Kalau project ini makin besar, naming pertama yang layak dipertimbangkan adalah
rename `static` menjadi `templateonly` atau `generic`. Untuk sekarang, `static`
masih cukup aman karena konvensinya sudah terdokumentasi.

## Status Saat Ini

- Render semua target: jalan
- Render target spesifik: jalan
- Config lokal: tersedia di `config/`
- Test: coverage fokus untuk loader dan renderer
- Next step yang masuk akal: template Hyprland lebih lengkap, validasi schema,
  dan mode hot reload/watch
