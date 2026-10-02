# Hot Reload Plan — Theme Engine

Rencana implementasi hot reload: setelah engine me-render config tiap tool,
perubahan langsung berlaku di tool terkuit tanpa restart manual.

> **Status**: ✅ terimplementasi (lihat `README.md` §Hot Reload).
> Dokumen ini adalah dasar implementasi yang sudah diwujudkan; disimpan
> sebagai referensi riset (tabel mekanisme per tool, temuan inotify,
> trade-off write mode).

---

## 1. Ringkasan Temuan

Engine sudah menulis file config dengan benar (atomic write + skip-unchanged).
Yang belum ada:

1. **Fase apply** setelah render — memberi tahu tool yang sedang berjalan
   bahwa config-nya berubah.
2. **Kompatibilitas write mode** — atomic write (`temp + rename`) memutus
   inotify watch Hyprland/Alacritty (lihat §3).
3. **CLI** — `set-theme <nama>` dan mode `watch` (Prioritas 3 di PROGRESS.md).

---

## 2. Hasil Riset Hot Reload Per Tool

| Tool | Target | Mekanisme resmi | Otomatis? | Catatan |
| --- | --- | --- | --- | --- |
| **Hyprland** | `hypr` → `source.conf` | inotify: config auto-reload saat save; `hyprctl reload` manual | Ya (jika ditulis in-place) | Watcher mengikuti inode; write via rename **tidak** terdeteksi (§3) |
| **Alacritty** | `alacritty` → `ui.toml` | `live_config_reload = true` (default) — watch file config | Ya (jika ditulis in-place) | Beberapa opsi (`window.dimensions`, `startup_mode`) tetap butuh restart |
| **Kitty** | `kitty` → `ui.conf` | Remote control: `kitty @ set-colors -a <file>` | Tidak — perlu perintah | Butuh `allow_remote_control yes` di `kitty.conf`; hanya warna yang live, font/opacity butuh restart |
| **Cava** | `cava` → `cava_extra` | `SIGUSR2` → reload warna saja (sama dengan tombol `c`) | Tidak — perlu sinyal | Output kita memang warna-only (gradient 1–5), pas persis |
| **Foot** | `foot` → `ui.ini` | Tidak ada live reload (issue #1653: "restart the server") | Tidak | Instance baru / `footserver` restart baca config baru |
| **Nvim** | `nvim` → `colors.lua` | `:colorscheme <nama>` di instance berjalan; `nvim --server <addr> --remote-send` | Instance baru: ya | Instance berjalan perlu `--listen <addr>`; nama colorscheme = nama file lua |
| **Starship** | `starship` → `starship.toml` | Stateless — baca config setiap render prompt | **Ya, inherent** | Tidak perlu tindakan apa pun |
| **Yazi** | `yazi` → `theme.toml` | Dibaca saat startup; yazi TUI berumur pendek | **Ya** (per jalankan) | Tidak perlu tindakan apa pun |
| **Lazygit** | `lazygit` → `config.yml` | Dibaca saat startup saja (issue #1158) | Tidak | Restart instance diperlukan; tidak ada IPC reload |
| **GTK** | `gtk` → `source.css` + `source.rasi` | `gtk.css` dibaca saat app start; tidak ada live reload bawaan (GNOME/gtk#3409) | Tidak | Cara pakai: restart GTK apps. **Waybar** (yang `@import` css ini): set `"reload_style_on_change": true` di waybar config → auto-watch css + imports |
| **Rofi** | via `source.rasi` | Dibaca saat launch (on-demand) | **Ya** (per launch) | Tidak perlu tindakan apa pun |
| **System (dconf)** | `system` → `apply.sh` | `dconf write` langsung berlaku untuk app gsettings-aware | **Ya** (saat script dijalankan) | Menjalankan `apply.sh` sudah menjadi apply-nya sendiri |

Sumber: wiki.hyprland.org (Configuring), alacritty.org/config-alacritty.html,
sw.kovidgoyal.net/kitty/remote-control, github.com/karlstav/cava (README),
codeberg.org/dnkl/foot/issues/1653, github.com/jesseduffield/lazygit/issues/1158,
yazi-rs.github.io/docs, starship.rs/faq, man waybar.5 (`reload_style_on_change`).

---

## 3. Temuan Kritis

### 3.1 Atomic rename memutus inotify watcher (Hyprland & Alacritty)

Renderer sekarang menulis via `temp file + rename` (`renderer.go:writeIfChanged`).
Hyprland mem-watch config-nya dengan inotify **per inode**. Bukti langsung:
hyprwm/Hyprland discussion #11848 — save ala neovim (rename-based backup)
membuat reload pertama tidak terdeteksi, dan save kedua justru memuat config
default (reference inode putus). Alacritty punya kelas masalah yang sama
(issue #5355: swap symlink target tidak terdeteksi).

**Konsekuensi**: untuk target yang di-watch tool-nya, engine harus menulis
**in-place** (`O_TRUNC` + write) agar watcher tetap hidup.

### 3.2 Template cache tidak invalidasi (mode watch)

`renderer.cache` mem-parse template sekali per path per proses. Di mode
`watch` (proses panjang), template yang diubah **tidak** akan di-parse ulang.
Perlu `renderer.ClearCache()` (atau cek mtime) tiap siklus watch.

### 3.3 Reload failure harus non-fatal

Konvensi project: unknown processor = Warn + skip, exit 0. Fase apply harus
mengikuti prinsip yang sama: tool tidak terinstal / remote control mati /
tidak ada proses cava → Warn, render tetap sukses.

---

## 4. Desain Implementasi

### 4.1 Port baru: `Reloader` (opsional)

`internal/app/ports/ports.go`:

```go
// Reloader menerapkan output yang baru di-render ke tool yang sedang
// berjalan. Opsional — processor yang tidak mengimplementasikannya
// dilewati tanpa error (konvensi unknown-processor).
type Reloader interface {
    Reload(outputPath string, ctx *renderctx.Context) error
}
```

`ctx` dibawa supaya `nvim` tahu nama colorscheme aktif (= nama theme);
tambah `ThemeName string` ke `renderctx.Context` (diisi dari
`st.Theme.Name` di `engine.New`).

### 4.2 Fase apply di engine

`internal/app/engine/engine.go` — `Run(name)`, setelah `Render` sukses:

```go
if e.applyEnabled {
    if r, ok := proc.(ports.Reloader); ok {
        if err := r.Reload(paths.OutputPath, &e.ctx); err != nil {
            log.Warn("reload %s: %v", name, err) // non-fatal
        }
    }
}
```

`applyEnabled` dari env `THEME_ENGINE_APPLY` (nilai `auto|1|0`, default
`auto`): mati otomatis saat `DefaultMapDir()` memetakan ke `config/` lokal
(dev/test), aktif untuk map eksternal (`$MYENV/map`). Pola ini meniru
`THEME_ENGINE_NOTIFY=auto` yang sudah ada di `runner.sh`.

### 4.3 Write mode per target

| Target | Write mode | Apply |
| --- | --- | --- |
| `hypr` | atomic write (tetap) | `Reload` → `hyprctl reload` (idempoten, mekanisme manual resmi wiki) |
| `alacritty` | **in-place** (`renderer.RenderInPlace`) | tidak ada Reloader — `live_config_reload` yang bekerja |
| `kitty` | atomic | `Reload` → filter baris warna dari output → temp file → `kitty @ set-colors -a <tmp>` |
| `cava` | atomic | `Reload` → `pkill -SIGUSR2 cava` (no-match = no-op) |
| `nvim` | atomic | `Reload` → jika `$NVIM_LISTEN_ADDRESS` ada: `nvim --server <addr> --remote-send "<cmd>colorscheme <nama><cr>"` |
| lainnya | atomic | tanpa Reloader (perilaku §2) |

`RenderInPlace` baru di `renderer.go`: reuse cache + buffer + compare, tapi
tulis langsung ke file (`os.OpenFile(path, O_WRONLY|O_CREATE|O_TRUNC, 0644)`).
Trade-off: tidak atomic (crash mid-write bisa korup) — risiko kecil untuk file
config sekecil ini, dan ini satu-satunya cara agar watcher alacritty tetap
hidup. Dokumendasikan di README.

Kitty: `set-colors` hanya menerima entri warna, sementara `kitty.tmpl`
juga berisi font/opacity/layout — `Reload` harus mengekstrak hanya baris
warna (`foreground`, `background`, `cursor`, `colorN`, `selection_*`,
`*_tab_*`, `*_border_color`) ke file temp sebelum memanggil `kitty @`.

### 4.4 `set-theme <nama>`

`cmd/theme-engine/settheme.go`:

1. Validasi `themes/<nama>/palette.json` ada (error `ExitLoad` kalau tidak).
2. Load `.state.json`, set `theme.name`, tulis atomic (pola temp+rename).
3. `RunAll()` — sekaligus memicu fase apply ke semua tool.

Opsional: flag `--type dark|light` untuk ikut mengganti varian.

### 4.5 Mode `watch` — polling stdlib, bukan fsnotify

PROGRESS.md menyarankan fsnotify, **tapi `go.mod` project ini sengaja
stdlib-only** (nol dependensi). Usulan: **polling** `os.Stat` (mtime+size)
agar properti stdlib-only tetap terjaga:

- Akar pantauan: `assets/templates/`, `themes/<aktif>/`, map dir
  (`path.txt` + `.state.json`).
- Interval 300–500ms; debounce 300ms (coalesce save beruntun).
- Pada perubahan: `renderer.ClearCache()` (§3.2) → `engine.New` → `RunAll`.
- `engine.New` gagal (misal theme.json sedang di-edit, JSON invalid) →
  log error, **tetap watching**, jangan crash.
- Handle SIGINT/SIGTERM untuk keluar bersih (`os/signal`, stdlib).

Jika nanti mau fsnotify, itu keputusan dependensi — diskusi terpisah.

### 4.6 CLI

`main.go` saat ini: argumen opsional = nama target. Diperluas (tetap
mundur-kompatibel):

```txt
theme-engine              # render semua target
theme-engine <target>     # render satu target (perilaku sekarang)
theme-engine set-theme <nama> [--type dark|light]
theme-engine watch
```

`set-theme`/`watch` jadi kata reserved (target bernama sama tidak mungkin
karena bukan key di `path.txt`).

---

## 5. Checklist Perubahan (implementasi)

| # | File | Perubahan |
| --- | --- | --- |
| 1 | `internal/infra/renderer/renderer.go` | + `RenderInPlace()` (write in-place setelah compare), + `ClearCache()` |
| 2 | `internal/app/ports/ports.go` | + interface `Reloader` |
| 3 | `internal/domain/renderctx/` | + field `ThemeName` |
| 4 | `internal/app/engine/engine.go` | + fase apply (type-assert `Reloader`, Warn non-fatal), + `applyEnabled`, isi `ThemeName` |
| 5 | `internal/adapters/tools/hypr/adapter.go` | + `Reload` → `hyprctl reload` |
| 6 | `internal/adapters/tools/alacritty/adapter.go` | `Render` → pakai `renderer.RenderInPlace` |
| 7 | `internal/adapters/tools/kitty/adapter.go` | + `Reload` → ekstrak warna → `kitty @ set-colors -a` |
| 8 | `internal/adapters/tools/cava/adapter.go` | + `Reload` → `pkill -SIGUSR2 cava` |
| 9 | `internal/adapters/tools/nvim/adapter.go` (generic) | + `Reload` → `nvim --server $NVIM_LISTEN_ADDRESS --remote-send` |
| 10 | `cmd/theme-engine/main.go` | parsing `set-theme` / `watch` |
| 11 | `cmd/theme-engine/settheme.go` (baru) | edit `.state.json` atomic + `RunAll` |
| 12 | `cmd/theme-engine/watch.go` (baru) | polling loop + debounce + signal handling |
| 13 | `internal/infra/renderer/renderer_test.go` | test `RenderInPlace` + `ClearCache` |
| 14 | test `set-theme` | round-trip `.state.json` |
| 15 | `README.md` + `PROGRESS.md` | dokumentasi hot reload, Prioritas 3 → selesai |

`wiring.go` **tidak berubah** — Reloader melekat pada adapter masing-masing.

---

## 6. Yang Harus Diaktifkan Pengguna (per tool)

- **Kitty**: `allow_remote_control yes` (atau `socket`/`socket-only`) di
  `kitty.conf`, kalau tidak `kitty @` error → Warn (render tetap sukses).
- **Waybar**: `"reload_style_on_change": true` di `config` waybar agar
  perubahan `source.css` (hasil GTK) langsung di-watch.
- **Nvim**: mulai nvim dengan `--listen <addr>` (atau `$NVIM_LISTEN_ADDRESS`)
  supaya instance berjalan bisa di-*send* `colorscheme`.
- **Alacritty**: pastikan `live_config_reload = true` (sudah default) dan
  output engine benar-benar file yang di-watch (atau `import`-nya — watch
  untuk file import tergantung versi, issue #5852).
- **Foot / Lazygit / GTK apps**: restart wajib — tidak ada mekanisme live
  reload di tool tersebut (bukan keterbatasan engine).

---

## 7. Risiko & Trade-off

| Risiko | Mitigasi |
| --- | --- |
| `RenderInPlace` kehilangan jaminan atomic write (alacritty) | File kecil + single `Write`; satu-satunya cara watcher alacritty hidup |
| Reload command tidak ada di sistem (kitty/hyprctl/pkill) | Warn + exit 0, sesuai konvensi unknown-processor |
| Dev mode (`config/` lokal) memicu apply ke tool asli | `THEME_ENGINE_APPLY=auto` matikan apply saat MapDir lokal |
| Watch loop error saat tema sedang diedit | Error log, loop tetap jalan |
| nvim multi-instance tidak ter-cover | Kirim ke `$NVIM_LISTEN_ADDRESS` saja; dokumentasikan |

---

## 8. Urutan Implementasi yang Disarankan

1. `ports.Reloader` + fase apply di engine + `RenderInPlace`/`ClearCache`
   (fondasi, belum ada adapter yang reload).
2. Adapter `cava` + `hypr` (paling sederhana: sinyal & `hyprctl`).
3. `alacritty` in-place + `kitty` set-colors (perlu ekstraksi warna).
4. `nvim` remote-send.
5. `set-theme` CLI.
6. `watch` mode polling.
7. Test + dokumentasi.
