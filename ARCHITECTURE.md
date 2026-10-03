# Arsitektur Theme Engine

Renderer tema berbasis JSON untuk setup desktop Hyprland: satu sumber tema
(`themes/<nama>/`) di-resolve lalu di-render menjadi file config per target
(GTK, Kitty, Foot, Hyprland, dll.) via Go `text/template`.

## 1. Prinsip Lapisan

Dependensi hanya mengarah ke dalam:

```txt
adapters → app → domain        (aturan utama)
infra → domain                 (hanya untuk tipe data)
cmd → semuanya                 (wiring saja)
```

* `domain` murni: struct + logika resolve, tanpa IO.
* `app` orkestrasi: tidak tahu JSON, path file, atau `sassc`.
* `adapters` implementasi per target terhadap kontrak `app/ports`.
* `infra` hal teknis reusable: file, template, env, log.
* `cmd` tipis: parse argumen + wiring eksplisit, tanpa `init()`/blank-import.

## 2. Struktur Direktori

```txt
cmd/theme-engine/
  main.go          # CLI: 1 argumen opsional (nama target), exit code 0-6
  wiring.go        # defaultProcessors(): map nama → adapter (otoritas identitas)

internal/app/
  engine/engine.go # Pipeline New → RunAll → Run; DefaultMapDir()
  ports/ports.go   # Kontrak Processor = Parser + Resolver + Renderer

internal/domain/
  palette/         # Raw, ResolvedPalette, ResolvePalette, flatten, RawPaletteVars
  theme/           # Theme (fonts + Tools map), ResolveDefaults
  state/           # State (.state.json)
  vars/            # VarSource{ Get(key) }
  renderctx/       # Context{ Palette, Theme, ThemeType } — data render global

internal/adapters/
  tools/           # Adapter per tool: kitty, foot, alacritty, cava, hypr
    <tool>/          # model.go (input JSON) + view.go (data template) + adapter.go
    generic/         # Template-only: nvim, yazi, starship, lazygit (tanpa schema khusus)
    jsonx/           # Decode: Parse nil-safe, anti-panic
  platform/
    gtk/             # view.go (Gtk{*Context}) + adapter.go (render + sassc)
    system/          # Template-only apply.sh (dconf); pisah dari generic (lihat §8)

internal/infra/
  loader/          # json.go (LoadJSON generik), toolmap.go (parse path.txt)
  renderer/        # Cache template, skip-unchanged, atomic write, helper hex
  pathenv/         # ExpandPath ($WAYBAR/$THEME/$ENV) + ResolveVar ($pl.*)
  log/             # Info/Warn/Error, verbose global

assets/templates/  # Sumber .tmpl (tools/*, domain/gtk/*, domain/system/*)
config/            # path.txt + .state.json lokal untuk dev/test
themes/<nama>/     # palette.json + theme.json per theme
output/            # Hasil render (tidak di-commit, lihat .gitignore)
```

## 3. Alur Runtime

`cmd/theme-engine/main.go:25` → `engine.New` → `RunAll`/`Run`:

```txt
.state.json ──► theme aktif (name/type) + waybar
path.txt    ──► map target → {template, output}, dengan ExpandPath
palette.json──► ResolveSelected(type) → ResolvePalette → BuildFlattenPalette
theme.json  ──► fonts.ResolveDefaults + Tools map (mentah per tool)
                        │
                        ▼
              per target: Parse → Resolve → Render (sorted by name)
```

* `engine.New` (`internal/app/engine/engine.go:32`): load + validasi semua input,
  gagal di sini berarti fail-fast sebelum ada file tersentuh.
* `Run(name)` (`engine.go:96`): target tak dikenal → error; processor tak dikenal
  → `Warn` + skip (return nil, exit tetap 0 — disengaja, lihat §8); tiap tahap
  dibungkus error `parse/resolve/render <name>`.
* `DefaultMapDir` (`engine.go:131`): `THEME_ENGINE_MAP` → `config/` lokal (dev)
  → `$MYENV/map` (state live desktop) → `config`.

## 4. Kontrak Adapter (`app/ports`)

```go
type Processor interface {
    Parser    // Parse(raw any) (any, error) — json mentah → struct model
    Resolver  // Resolve(in any, ctx *renderctx.Context) (any, error) — model → view
    Renderer  // Render(templatePath, outputPath string, data any) error
}
```

Pola seragam per tool (`internal/adapters/tools/<tool>/`):

| File | Isi |
|---|---|
| `model.go` | Struct `Raw`: cerminan `tools.<nama>` di `theme.json` |
| `view.go` | Struct view: data final siap template (+ `Palette` bila perlu warna global) |
| `adapter.go | `Processor{}` + `New()` + 3 method kontrak |

Aturan defensif (semua adapter):

* `Parse` via `jsonx.Decode`: input `nil` → zero `Raw`; tipe salah → error
  (tidak pernah panic atas type assertion).
* `Resolve`: `inp, _ := in.(Raw)` — config absen → zero value, bukan panic.
  Pengecualian: `cava` error bersih bila gradients < 5.
* Identitas target = key map di `wiring.go`, bukan method — adapter tidak punya
  `Name()` agar tidak ada dua sumber kebenaran.

Dua jenis adapter khusus:

* `generic` (`adapters/tools/generic/adapter.go`): template-only, teruskan
  `*renderctx.Context` mentah ke template. Dipakai `nvim`, `yazi`, `starship`, `lazygit`.
* `platform/gtk` (`adapters/platform/gtk/adapter.go`): satu-satunya dengan
  `Render` kustom — render `source.tmpl` → `_source.scss`, lalu `sassc`
  compile ke `css/source.css` + `rasi/source.rasi` (dengan `MkdirAll` +
  error berisi stderr sassc). `outputPath`-nya direktori, bukan file.

## 5. Sistem Variabel

Dua sistem terpisah, jangan dicampur:

**a. `$pl.*` — warna palette** (`infra/pathenv/resolve.go:12`)

* Di `palette.json`: hanya base palette — `foreground`, `background`,
  `cursor`, `color0..15`. Resolver single-pass: referensi ke sesama
  `extra.*` **tidak** ter-resolve dan bocor sebagai literal (aturan
  didokumentasikan di `RULE.md`, bab Referensi `$pl.`).
* Di `theme.json` (level tool): namespace hasil flatten — `extra.*`,
  `fg.*`, `bg.*`, `syntax.*`, `ui.*`, `colorN` (lihat `palette/flatten.go`).
* Miss → return input apa adanya (lenient historis); `syntax`/`ui`/`visited`/
  `border.medium`/`fg`/`bg` punya fallback semantik di `resolver.go`.
* Validasi ringan: `ResolvePalette` error bila `len(Colors) < 7` (dibutuhkan
  `Colors[2,4,6]` untuk fallback syntax) — satu int-compare, tanpa cost render.

**b. `$VAR` — path filesystem** (`infra/pathenv/expand.go:31`)

* `$WAYBAR` → `state.Waybar`, `$THEME` → nama theme aktif, sisanya env OS
  (upper-cased). Dipakai di `path.txt` dan path map itu sendiri.

## 6. Config & Data

* `config/path.txt`: `target|template|output` per baris; kosong/`#` diabaikan;
  error menyebut `file:line` (`loader/toolmap.go:41`).
* `config/.state.json`: `{version, theme{name,type}, waybar}`.
* `themes/<nama>/{palette.json,theme.json}`: palette berisi base + `extra`
  semantik (aturan warna di `RULE.md`); theme berisi `fonts` (SSOT ukuran font
  terminal di `fonts.terminal`, detail di `FONT_SYSTEM.md`) + `tools` opsional
  per tool.
* Folder `themes/*.bak` diabaikan (konvensi arsip lokal).

## 7. Renderer & Jaminan Performa Switch

`infra/renderer/renderer.go`:

* Template di-parse sekali per path, disimpan di cache `sync.RWMutex`.
* Hasil di-execute ke buffer dulu, lalu `writeIfChanged`: isi sama → skip
  tulis (hemat SSD + tidak memicu hot-reload eksternal).
* Tulis atomik: temp file di direktori tujuan + `rename`.
* Helper template `hex`: strip awalan `#` (`rgb({{ hex .Palette.AccentPrimary }})`).

Cost terbesar ada di tool eksternal (`sassc`) dan reload aplikasi setelah
render, bukan di parsing/loading — jadi penambahan validasi ringan (satu
compare, nil-check) tidak berpengaruh ke kecepatan switch.

## 8. Keputusan Disengaja (jangan "diperbaiki" tanpa diskusi)

* **Unknown processor = Warn + skip, exit 0** (`engine.go:102`). Terlihat
  seperti bug (notifikasi runner bilang sukses), tapi dipertahankan sementara;
  pengubahnya wajib mengubah juga ekspektasi `runner.sh`.
* **`system` pisah dari `generic`** walau kode identik: `generic` untuk
  `output/tools/*`, `system` untuk `output/domain/system/apply.sh`.
  Duplikasi 3 method demi kepemilikan yang jelas.
* **Miss `$pl.*` lenient** (bukan error) — kompatibilitas historis.
* **`output/` tidak di-commit** (`.gitignore`: `**/output/**`).

## 9. Menambah Target Baru

Template-only (cukup warna global):

1. `assets/templates/tools/<nama>/config.tmpl`
2. `config/path.txt`: `<nama>|assets/...|output/...`
3. `wiring.go`: `procs["<nama>"] = generic.New()`

Contoh yang sudah lewat jalur ini: `lazygit`
(`assets/templates/tools/lazygit/config.tmpl` → `output/tools/lazygit/config.yml`).
Untuk nvim, warna `fg_gutter` dan `terminal_black` di `colors.tmpl` adalah target
kontras yang dijaga test — lihat `RULE.md` bab mapping Tokyonight.

Butuh config khusus:

1. `internal/adapters/tools/<nama>/{model.go,view.go,adapter.go}` + `New()`
2. Template + `path.txt` + `wiring.go`: `procs["<nama>"] = <nama>.New()`
3. Default di `themes/<theme>/theme.json` bila perlu

## 10. Test

```bash
go test ./...                                   # semua
GOCACHE=/tmp/go-build go test ./...            # bila cache Go tak bisa ditulis
./runner.sh test
```

Guard khusus warna: `internal/domain/palette/contrast_test.go` meng-assert ambang
WCAG (3:1 UI / 4.5:1 teks) untuk setiap tema × variant. Angka yang sama bisa
dilihat tanpa menjalankan Go lewat `python3 tools/contrast_report.py`
(tool dev-only di `tools/`; helper `--plan`/`--accents` untuk memilih nilai palet).

Coverage saat ini fokus pada `infra/loader` (expand path map, komentar,
malformed line) dan `infra/renderer` (buat parent dir, skip-unchanged via
mtime). Adapter/domain belum ber-unit-test — verifikasi utama masih render
penuh + `diff output/`.
