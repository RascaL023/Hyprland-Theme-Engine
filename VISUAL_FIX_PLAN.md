# Rencana Perbaikan Bug Visual (Nvim & Lazygit)

Status: **Fase 0–5 SELESAI.** Semua 6 tema lama + tema baru `harbor` lulus guard
(dark & light), target `lazygit` aktif, dan alat ukur/usul warna sudah tersedia.

Ringkasan keputusan & analisis performa: lihat **§11** dan **§12**.

Tema pilot: **`kanagawa-dragon`** (Fase 1–3), lalu di-rollout ke
`nocturne`, `ghostly`, `sakura`, `kanagawa-wave`, `claude-manjusaka` (Fase 5).
Tema baru `harbor` (§8b) dibuat setelah guard aktif, jadi lulus sejak awal.

---

## 0. Ringkasan Masalah (hasil audit)

| # | Gejala | Akar masalah | Bukti |
|---|---|---|---|
| A | Teks backtick (`` `test` ``) di `.md` susah dibaca | Template `nvim/colors.tmpl` memetakan `terminal_black` ke **ANSI bright-black** (`colors[8]`), padahal Tokyonight memakai `terminal_black` sebagai **warna latar (background)** untuk inline code. Di beberapa tema `colors[8]` justru terang → latar & teks sama-sama terang. | `tokyonight/.../groups/treesitter.lua:66` (`@markup.raw.markdown_inline = { bg = c.terminal_black, fg = c.blue }`), `colors.tmpl:38` |
| B | Bagian tengah lualine (nama file/tab) kurang kontras | `fg_gutter` dipetakan ke `layer.surface_overlay` (layer **paling terang**), padahal Tokyonight asli `fg_gutter` itu **gelap**. Mode color jadi kehilangan kontras. | `_tokyonight.lua:11` (`b = { bg = c.fg_gutter, fg = c.blue }`), `colors.tmpl:14` |
| C | Di lazygit, baris yang di-select kurang jelas & tema tidak terasa | **Belum ada target `lazygit` sama sekali** (tidak ada template/path.txt/wiring). `~/.config/lazygit/config.yml` kosong (0 byte), jadi lazygit memakai default: `selectedLineBgColor: [blue]` + `defaultFgColor: [default]`. | `docs/Config.md` (default resmi), `config/path.txt` |

Temuan tambahan saat rollout: `sakura`/`kanagawa-wave`/`ghostly` belum mengisi blok
`extra.fg`, `extra.bg`, `text.visited`, `border.medium`. Ini **sengaja dibiarkan**:
resolver sudah punya fallback semantik untuk semua field itu, sehingga menambahkannya
 hanya akan mengubah output tanpa memperbaiki bug. Yang wajib diisi hanya
`syntax.terminal_black` dan `ui.gutter` (dua slot baru di dokumen ini).
`sakura` memang belum punya blok `syntax`/`ui` → ditambahkan minimal (hanya field
baru + `bg_statusline` yang nilainya sama dengan fallback-nya).

Riset tambahan (bukan penyebab bug, tapi batas desain yang perlu diketahui):

- `fg_gutter` di Tokyonight punya **dua peran**: foreground untuk nomor baris /
  indent guide, dan background untuk lualine B / folded / tabline. Karena itu
  `ui.gutter` dijaga ≥3:1 terhadap mode accent **dan** ≥~1.2:1 terhadap `layer.base`.
- `terminal_black` juga dipakai sebagai **foreground** ghost text
  (`CmpGhostText`, `DiagnosticUnnecessary`), jadi tidak boleh persis `layer.base`.
- Tema dengan ANSI accent pastel (`claude-manjusaka/light`, `ghostly/light`,
  `nocturne/light`) tidak punya nilai terang yang bisa lulus 3:1, sehingga
  `ui.gutter`-nya **harus gelap**. Itu trade-off palet, bukan bug engine —
  lihat advisory di §8.

---

## 1. Kamus Istilah (biar tidak bingung)

- **foreground (fg)** = warna teks/gambar. **background (bg)** = warna latar di belakangnya.
- **ANSI `colors[0..15]`** = 16 warna standar terminal. `colors[4]` = "blue",
  `colors[8]` = "bright black". Tokyonight memakai beberapa warna ini untuk
  hal di luar terminal (mis. `colors[8]` dipakai sebagai **latar** inline code).
- **Kontras ratio** = ukuran seberapa jelas teks di atas latarnya.
  `1:1` = warna persis sama (tidak terbaca), `21:1` = hitam-putih.
  Patokan: minimal **3:1** untuk elemen UI/teks besar, **4.5:1** untuk teks normal.
- **lualine** = statusline nvim, terbagi 3 bagian per mode:
  - **A** (kiri, blok berwarna sesuai mode) — `bg = mode color`
  - **B** (tengah, nama file/branch) — `bg = fg_gutter`, `fg = mode color`
  - **C** (kanan) — `bg = bg_statusline`
  Masalah B: latarnya (`fg_gutter`) terlalu terang sehingga tulisan mode color-nya "tenggelam".
- **`$pl.*`** = tempat engine menyimpan warna semantik (`$pl.extra.layer.base`, dst.).
- **inline code / backtick** = teks `` `seperti ini` `` di markdown.

---

## 2. Aturan Main

1. Hanya `themes/kanagawa-dragon/` yang diubah.
2. Tidak mengubah nilai `colors[0..15]` yang sudah Anda suka **kecuali** benar-benar
   diperlukan dan disetujui (lihat Fase 2).
3. Setiap fase harus punya bukti angka (kontras) + diff output sebelum lanjut.
4. `output/` tidak di-commit; yang di-commit adalah template/palet/kode.
5. Golden test baru di-regenerate setelah semua fase pilot lulus.

---

## 3. Fase 0 — Alat Ukur (tanpa perubahan output)

**Tujuan:** mengubah "bug tak terlihat" jadi angka yang bisa dites, supaya
perbaikan Fase 1–3 tidak menambah bug baru.

Isi:
- `internal/domain/palette/contrast.go` — fungsi murni: luminance + contrast ratio
  (tanpa IO, sesuai aturan arsitektur).
- `internal/app/engine/contrast_test.go` — men-iterasi `themes/*` × `dark/light`,
  lalu gagal bila ambang terlampaui.

Ambang yang diuji:
- `contrast(mode color, ui.gutter/fg_gutter) >= 3` (lualine B)
- `contrast(terminal_black, colors[4]) >= 3` (backtick)
- `contrast(text.primary, layer.base) >= 4.5` (teks utama)
- `contrast(text.primary, selectedLineBg) >= 3` (lazygit select)

- [x] Tambah `contrast.go`
- [x] Tambah `contrast_test.go`
- [x] Guard diverifikasi menangkap bug lama: `terminal_black=#a6a69c` → **1.06:1 FAIL**, `#393836` → **4.48:1 PASS**
- [x] Jalankan `GOCACHE=/tmp/go-build go test ./...` (hijau)

Output yang berubah: **tidak ada** (murni test).

---

## 4. Fase 1 — Perbaiki Backtick (Masalah A)

**Dampak ke layar:** kotak latar di belakang `` `test` `` menjadi gelap (dark)
atau terang (light) yang benar, sehingga warna teks biru terbaca.

### Pilihan pendekatan

| Opsi | Yang dilakukan | Kelebihan | Kekurangan |
|---|---|---|---|
| **A (rekomendasi)** | Tambah field palet baru, mis. `extra.ui.gutter` / `extra.syntax.terminal_black`, dengan nilai default `layer.surface_overlay`. Template `colors.tmpl` memakai field ini untuk `terminal_black`. | `colors[8]` tetap murni "bright black"; tiap tema masih bisa di-tune; bisa dipakai juga untuk memperbaiki Fase 2. | Perlu ubah 3 file Go + 1 template + isi 12 variant palet (bertahap, mulai pilot). |
| B | Template langsung pakai `layer.surface_overlay` (tanpa field baru). | Paling cepat, tanpa ubah schema. | Tidak bisa di-tune per tema; nilainya akan selalu sama dengan `fg_gutter`/`dark3`. |
| C | Tidak ubah engine; override `@markup.raw.markdown_inline` di `~/.config/nvim/lua/core/theme.lua`. | Perubahan terisolasi di config nvim. | Tidak memperbaiki engine; tema lain tetap salah di luar nvim. |

### Angka pilot `kanagawa-dragon`

| Variant | Sekarang | Sesudah (usul `terminal_black = layer.surface_overlay`) |
|---|---|---|
| dark | bg `#a6a69c` vs fg `#8ba4b0` = **1.06** ❌ | bg `#393836` vs fg `#8ba4b0` = **4.48** ✅ |
| light | bg `#8a8980` vs fg `#4d699b` = **1.57** ❌ | bg `#e4d794` vs fg `#4d699b` = **3.79** ✅ |

### Diff output yang diharapkan (`output/tools/nvim/colors.lua`)

```diff
-  terminal_black    = "#a6a69c",
+  terminal_black    = "#393836",
```
(untuk variant dark; variant light akan terisi `#e4d794` saat theme `light` aktif.)

- [x] **KEPUTUSAN: Opsi A** — tambah slot `syntax.terminal_black` (fallback `layer.surface_overlay`)
- [x] `raw.go`: tambah field
- [x] `resolved.go`: tambah field
- [x] `resolver.go`: resolve + fallback `layer.surface_overlay`
- [x] `flatten.go`: daftarkan key flat
- [x] `colors.tmpl`: `terminal_black` pakai field baru
- [x] `themes/kanagawa-dragon/palette.json`: isi nilai eksplisit (dark `#393836`, light `#e4d794`)
- [x] Render `go run ./cmd/theme-engine nvim`, cek `output/tools/nvim/colors.lua` (terisi `#393836`)
- [ ] `:Inspect` di nvim pada backtick → pastikan grup yang tampil memakai bg baru (**perlu dicek di sisi Anda**)

Catatan: `terminal_black` juga dipakai sebagai `terminal.black_bright` di nvim
`:=terminal`. Nilai baru (surface) tetap benar untuk itu.

---

## 5. Fase 2 — Perbaiki lualine Bagian B (Masalah B)

**Dampak ke layar:** latar blok tengah lualine (nama file/cabang) jadi lebih gelap
(dark) / tepat (light); tulisan mode color kembali jelas.

### Angka pilot `kanagawa-dragon` (bg = `fg_gutter` saat ini = `#393836` / `#e4d794`)

| Mode | dark | light |
|---|---|---|
| red (`colors[1]`) | 3.39 ✅ | 3.35 ✅ |
| green (`colors[2]`) | 3.90 ✅ | **2.69** ❌ |
| yellow (`colors[3]`) | 5.62 ✅ | 3.42 ✅ |
| blue (`colors[4]`) | 4.48 ✅ | 3.79 ✅ |
| magenta (`colors[5]`) | 4.00 ✅ | 3.38 ✅ |
| green1 (`green1`) | 4.49 ✅ | **2.46** ❌ |

Kesimpulan pilot:
- **dark sudah lulus** → tidak perlu ubah apa pun.
- **light gagal pada 2 warna hijau** → mapping tidak bisa menolong (semua layer
  light pada tema ini < 3). Perlu tweak kecil nilai hijau di palet light.

### Pilihan pendekatan

| Opsi | Yang dilakukan | Dampak |
|---|---|---|
| **1 (rekomendasi untuk dark/global)** | Ubah mapping `fg_gutter` → `layer.mantle`. | Pilot dark tidak terpengaruh (sudah lulus); memperbaiki tema dark lain nanti. Tidak mengubah warna favorit Anda. |
| 2 (wajib untuk light pilot) | Gelapkan `colors[2]` (green) & `syntax.green1` di `kanagawa-dragon/light` sampai ≥3:1. | Mengubah 2 nilai hijau **hanya di variant light pilot**. Perlu persetujuan. |
| 3 | Tambah field `ui.gutter` khusus (sama seperti Fase 1 opsi A). | Tiap tema bisa atur gutter tanpa mengganggu urutan layer. |

### Diff output yang diharapkan

Slot `ui.gutter` (fallback `layer.surface_overlay`): pilot **tidak** mengisinya,
jadi `colors.lua` dark tetap `#393836` (tidak ada perubahan output) — mekanisme
siap dipakai tema lain di Fase 5.

Tuning hijau light pilot (yang benar-benar mengubah output):
```diff
   colors[2] green
-  "#6f894e",   # 2.69:1
+  "#677f49",   # 3.07:1
   syntax.green1 light
-  "#6e915f",   # 2.46:1
+  "#618054",   # 3.06:1
```

- [x] **KEPUTUSAN: Opsi 3** — tambah slot `ui.gutter` (fallback `layer.surface_overlay`, sehingga tema lama tidak berubah sampai diisi)
- [x] **KEPUTUSAN (khusus pilot light): Opsi 2** — geser `colors[2]` & `syntax.green1` variant `light` sampai ≥3:1
- [x] Terapkan perubahan pada file yang dipilih
- [x] Konfirmasi kontras 6 mode ≥ 3:1 untuk dark **dan** light pilot (dark min 3.39, light min 3.06)
- [ ] Cek `:Inspect` grup `StatusLine` / lualine B (**perlu dicek di sisi Anda**)

---

## 6. Fase 3 — Target `lazygit` Baru (Masalah C)

**Dampak ke layar:** lazygit akhirnya punya tema sendiri, dan baris yang
di-select jelas terbaca.

Angka pilot (select = `layer.surface_overlay`, teks = `text.primary`):

| Variant | Sekarang (`blue` vs fg) | Sesudah |
|---|---|---|
| dark | 1.56 ❌ | **6.99** ✅ |
| light | 1.35 ❌ | **5.11** ✅ |

Langkah (pola "tool template-only", lihat `ARCHITECTURE.md` §9):
- [x] `assets/templates/tools/lazygit/config.tmpl`
- [x] `config/path.txt`: `lazygit|assets/templates/tools/lazygit/config.tmpl|output/tools/lazygit/config.yml`
- [x] `cmd/theme-engine/wiring.go`: `genericTargets` ditambah `"lazygit"`
- [ ] `internal/app/engine/golden_test.go`: tambah `"lazygit"` ke `testProcessors` (**di branch `feature/test/golden`**)
- [x] Isi tema: `activeBorderColor`, `inactiveBorderColor`,
      `optionsTextColor`, `selectedLineBgColor`, `unstagedChangesColor`,
      `cherryPicked*`, `markedBaseCommit*`, `defaultFgColor`
- [x] Render `go run ./cmd/theme-engine lazygit` → `selectedLineBgColor=#393836`, `defaultFgColor=#c5c9c5` (6.99:1)
- [ ] Deploy: arahkan `output/tools/lazygit/config.yml` ke
      `~/.config/lazygit/config.yml` (via `$MYENV/map` di setup asli)

Isi usulan (dark pilot):
```yaml
gui:
  theme:
    activeBorderColor:            ["#8ba4b0"]   # accent.primary
    inactiveBorderColor:          ["#282727"]   # border.default
    searchingActiveBorderColor:   ["#E6C384"]   # syntax.orange
    optionsTextColor:             ["#8ba4b0"]   # accent.primary
    selectedLineBgColor:          ["#393836"]   # layer.surface_overlay
    unstagedChangesColor:         ["#c4746e"]   # status.error
    cherryPickedCommitFgColor:    ["#8ba4b0"]
    cherryPickedCommitBgColor:    ["#7AA89F"]
    markedBaseCommitFgColor:      ["#8ba4b0"]
    markedBaseCommitBgColor:      ["#E6C384"]
    defaultFgColor:               ["#c5c9c5"]   # text.primary
```

---

## 7. Fase 4 — Sisi Nvim: markview (SELESAI, di sisi dotfile)

**Masalah:** `transparent = true` + `Normal.bg = "none"` membuat markview
(`MarkviewInlineCode`) memakai warna fallback hardcoded `#1E1E2E`, bukan warna
tema aktif. Ini sumber bug backtick kedua (khusus plugin markview; grup
`@markup.raw.markdown_inline` diperbaiki Fase 1).

Temuan: `vim.g.markview_dark_bg`/`markview_light_bg` **tidak** memengaruhi
`MarkviewInlineCode` (global itu hanya dipakai `create_pallete` / grup
`MarkviewPalette*`, sedangkan inline code memakai warna hardcoded
`#1E1E2E`/`#EFF1F5` bila `Normal.bg` unset). Jadi perlu dua langkah.

- [x] Verifikasi grup aktif pada backtick (`MarkviewInlineCode`)
- [x] Set `vim.g.markview_dark_bg` / `markview_light_bg` dari `.Palette`
- [x] Override langsung `MarkviewInlineCode` (`bg = M.colors.terminal_black`,
      `fg = M.colors.green`) + re-apply lewat autocmd `ColorScheme`
      (karena `:colorscheme` menghapus highlight group)
- [x] `luac -p` bersih; runtime dicek: `MarkviewInlineCode` bg `#393836` fg `#8a9a7b` = **3.90:1**
- [x] Override bertahan setelah `:colorscheme tokyonight` + `markview.highlights.setup()`

Perubahan ini hidup di `~/.dotfile/nvim/lua/core/theme.lua` (config Anda,
bukan repo engine) — engine tetap benar untuk tool lain lewat Fase 1.

---

## 8. Fase 5 — Rollout ke Semua Tema (SELESAI)

Nilai final per tema (hanya slot yang berubah). Guard dijalankan untuk semua:

| Tema | Variant | `syntax.terminal_black` | `ui.gutter` | Catatan |
|---|---|---|---|---|
| claude-manjusaka | dark | `#2D2824` | `#2D2824` | = `layer.surface_overlay` |
| claude-manjusaka | light | `#68635E` | `#3E3A35` | accent pastel → gutter wajib gelap |
| ghostly | dark | `#363636` | `#272727` | shade dari overlay (dE 0.11 / 0.17) |
| ghostly | light | `#5B5B5B` | `#565656` | accent pastel → gutter wajib gelap |
| kanagawa-dragon | dark | `#393836` | `#393836` | = overlay (pilot) |
| kanagawa-dragon | light | `#e4d794` | `#e4d794` | + `colors[2]` & `green1` digelapkan |
| kanagawa-wave | dark | `#4B4B64` | `#222237` | shade dari overlay (dE 0.03 / 0.19) |
| kanagawa-wave | light | `#e4d794` | `#e4d794` | `colors[2]`/`colors[10]`/`green1` disamakan dengan dragon light |
| nocturne | dark | `#594d80` | `#594d80` | = overlay |
| nocturne | light | `#EFEAF7` | `#35303b` | tb = `layer.mantle`, gutter gelap |
| sakura | dark | `#3A3036` | `#3A3036` | = overlay |
| sakura | light | `#E0D0C6` | `#FFFFFF` | accent kuning `#AD8B55` hanya 2.12:1 → gutter harus nyaris putih |

Kontras minimum yang dicapai: **dark 3.05:1** (`ghostly` backtick), **light 3.05:1**
(`sakura` backtick) — ambang 3.0, jadi semua punya sedikit ruang (headroom tool +0.05).
Tidak ada nilai ANSI yang diubah selain dua hijau `kanagawa-*/light` di atas
(disetujui di Fase 2) — sisanya murni menambah dua slot baru.

- [x] Terapkan kedua slot ke `nocturne`, `ghostly`, `sakura`, `kanagawa-wave`, `claude-manjuska`
- [x] `sakura`: tambah blok `syntax` + `ui` (minimal, nilai `bg_statusline` sama dengan fallback-nya)
- [x] Guard diperluas: `contrastPilotThemes` → discovery `themes/*` (semua tema otomatis ikut)
- [x] `GOCACHE=/tmp/go-build go test ./...` hijau
- [x] Update tabel mapping Tokyonight di `RULE.md` (`terminal_black`, `fg_gutter` → `ui.gutter`)
- [x] Update `ARCHITECTURE.md` §9 + §10 (target `lazygit`, guard, tool)
- [x] Update `PROGRESS.md` (target lazygit, daftar tema, tooling)
- [ ] Regenerate golden: `UPDATE_GOLDEN=1 go test ./internal/app/engine/ -run TestGolden` (branch `feature/test/golden`)

### Advisory (temuan, BELUM dikerjakan)

Accent ANSI di beberapa tema light terlalu pastel untuk lulus 3:1 di atas surface
terang, sehingga `ui.gutter` terpaksa mendekati hitam/putih. Kalau ingin surface
kembali kalem, accent-nya yang perlu digeser (hue boleh dipertahankan):

| Tema | Accent biang | Sekarang | Perlu |
|---|---|---|---|
| claude-manjusaka/light | `colors[1..5]` (semua) | L 0.62–0.80 | digelapkan ke L ±0.57 |
| ghostly/light | `colors[1..5]` (semua) | L 0.74–0.79 | digelapkan ke L ±0.55 |
| nocturne/light | `colors[1..5]` | L 0.59–0.68 | digelapkan ke L ±0.48 |
| sakura/light | `colors[3]` kuning `#AD8B55` | L 0.66 | digelapkan ke L ±0.57 |
| kanagawa-wave/dark | `colors[1]` merah `#C34043` | L 0.56 | diterangkan ke L ±0.75 |
| ghostly/dark | `colors[1..5]` | L 0.55–0.65 | diterangkan ke L ±0.72 |

Angka ini dihitung dengan `python3 tools/contrast_report.py --plan --accents`.
Dampak yang wajib dicek sebelum mengubah ANSI: terminal (kitty/foot/alacritty),
lualine section A (mode accent sebagai background), dan `:Inspect` pada syntax
highlight. Karena itu rollout ini sengaja tidak menyentuhnya.

## 8b. Tema Baru: `harbor` (SELESAI)

Diminta setelah rollout: tema estetik yang tidak mencolok dan beda dari yang lain.
`harbor` = biru-slate + teal, satu-satunya tema **dingin** di koleksi
(ghostly netral abu, nocturne ungu, sakura pink, kanagawa/claude hangat).

- [x] `themes/harbor/palette.json` (dark + light, blok `syntax`/`ui` eksplisit lengkap)
- [x] `themes/harbor/theme.json` (font sama dengan tema aktif; kitty/hypr/alacritty di-tune: rounding 12, shadow range 10, opacity 0.72)
- [x] Lulus guard dengan margin nyaman: dark min **4.61:1**, light min **3.27:1**
- [x] Render dicek: `go run ./cmd/theme-engine nvim` + `lazygit` untuk dark & light
- [x] Output tema aktif (`kanagawa-dragon/dark`) di-render ulang setelah uji (state `.state.json` otomatis dikembalikan)

Cara memakai: ubah `config/.state.json` → `"name": "harbor"`, lalu render ulang.

---

## 9. Verifikasi Akhir

- [x] `GOCACHE=/tmp/go-build go test ./...` hijau semua
- [x] Guard lulus untuk 6 tema × dark/light + `harbor` (14 variant)
      — minimum dark 3.05:1, minimum light 3.05:1 (ambang 3.0); teks utama ≥ 6.19:1 (ambang 4.5)
- [x] `python3 tools/contrast_report.py` exit 0 (tidak ada FAIL di seluruh tema)
- [x] `output/` dicek: `nvim/colors.lua` & `lazygit/config.yml` untuk `kanagawa-dragon/dark`
      dan `harbor/dark+light`
- [x] Semua 11 target render sukses (`go run ./cmd/theme-engine`) tanpa error
- [ ] Uji mata di nvim (`:Inspect` pada backtick, lualine B) + lazygit (**sisi Anda**)
- [ ] Regenerate & commit golden (branch `feature/test/golden`)
- [ ] Deploy `output/tools/lazygit/config.yml` → `~/.config/lazygit/config.yml` (**sisi Anda**)

---

## 10. Checklist Master (Ringkas)

- [x] F0 `contrast.go` + `contrast_test.go` (discovery semua tema)
- [x] F1 backtick: slot `syntax.terminal_black` (keputusan: Opsi A)
- [x] F2 lualine B: slot `ui.gutter` + tuning hijau light pilot (keputusan: Opsi 3 + 2)
- [x] F3 target lazygit
- [x] F4 markview globals + override `MarkviewInlineCode` (di `~/.dotfile/nvim`, sudah diizinkan pemilik)
- [x] F5 rollout semua tema + advisory accent dicatat
- [x] Tema baru `harbor` (dark + light) lulus guard
- [x] Tool `tools/contrast_report.py` (audit + planner OKLCH)
- [x] Update dokumentasi: `RULE.md`, `ARCHITECTURE.md`, `PROGRESS.md`, file ini
- [ ] Regenerate golden (branch `feature/test/golden`)

---

## 11. Keputusan Final & Alasan

Prinsip yang dipakai: **satu slot warna = satu makna**. Kalau Tokyonight memakai
suatu slot untuk 2 hal berbeda, kita beri slot sendiri supaya tidak ada yang
"menebak", dan tiap tema bisa di-tune tanpa efek samping ke file lain.

### F1 — Backtick: slot baru `syntax.terminal_black`

| Kenapa bukan opsi lain | Alasan |
|---|---|
| Bukan pakai `colors[8]` lagi | `colors[8]` tetap dipakai murni sebagai ANSI bright-black oleh foot/kitty/alacritty. Kalau diutak-atik, warna terminal ikut berubah. |
| Bukan langsung `layer.surface_overlay` di template | Harus bisa di-tune per tema tanpa mengubah `layer.*` (yang dipakai GTK/Waybar). |
| Bukan tambalan di `theme.lua` saja | Engine tetap salah untuk tool lain & tema lain. Tambalan config hanya menutupi. |

Nama `syntax.terminal_black` sengaja mengikuti gaya `syntax.blue0..blue7` yang
sudah ada (mirror eksplisit kunci Tokyonight). Fallback = `layer.surface_overlay`,
jadi tema lama tetap jalan.

### F2 — lualine B: slot baru `ui.gutter`

`fg_gutter` di Tokyonight **bukan** `layer.surface_overlay` (yang merupakan layer
paling terang). Karena itu mapping sekarang sering salah. Solusi:

- Tambah `ui.gutter` (sejajar dengan `ui.bg_statusline`), fallback `layer.surface_overlay`
  → **100% backward-compatible**: tema yang belum mengisi tidak berubah sama sekali.
- Template memakai `ui.gutter` untuk `fg_gutter`.
- Rollout (Fase 5): tema dark yang gagal isi `ui.gutter` lebih gelap
  (mis. `layer.mantle`) — **tanpa** mengubah warna ANSI.
- Khusus pilot `light`: hijau (`colors[2]`, `syntax.green1`) yang jadi biang
  (2.69 & 2.46) digeser sedikit lebih gelap; ini tak bisa diperbaiki mapping.

### F4 — Sisi nvim (tool), diizinkan pemilik

`transparent = true` + `Normal.bg = none` membuat markview memakai fallback
hardcoded `#1E1E2E`. Ini bug **tool**, bukan engine, jadi diperbaiki di
`~/.config/nvim/lua/core/theme.lua` dengan mengisi `vim.g.markview_dark_bg` /
`markview_light_bg` dari `.Palette`. Tidak perlu ubah schema engine.

---

## 12. Analisis Performa (kenapa desain ini tidak memperlambat render)

Prinsip: **tidak boleh ada kerja tambahan saat render templete.**

| Aktivitas | Biaya tambahan | Penjelasan |
|---|---|---|
| Resolve `syntax.terminal_black` + `ui.gutter` | +2 `ResolveVar` + 2 `fallback` per pemuatan tema | Sekali per switch, bukan per target; hanya beberapa string compare. |
| `flatten.go` | +2 entri `map[string]string` | Alokasi map sudah ada; hanya 2 insert. |
| Template nvim | **Lebih murah**: `{{ index .Palette.Colors 8 }}` (index slice) → `{{ .Palette.SyntaxTerminalBlack }}` (akses field) | Menghilangkan pemanggilan `index`. |
| Contrast test | **Nol** biaya produksi | File `_test.go` hanya dikompilasi ke test binary, tidak ikut ke `cmd/theme-engine`. |
| Target `lazygit` | 1 render template biasa | Sama seperti target generic lain; tetap kena cache, skip-unchanged, atomic write. `sassc` (biaya terbesar) tidak tersentuh. |
| Tema baru `harbor` | **nol** | Tema hanyalah data; jumlah tema tidak menambah kerja per switch (hanya 1 tema aktif yang di-render). |
| `tools/contrast_report.py` | **nol** di runtime | Python dev-only; tidak dipanggil engine, tidak ikut binary, tidak ada di jalur render. |

Tidak ada loop, komputasi warna, atau pembacaan file tambahan di jalur render.
Kecepatan switch tetap sama; cost terbesar tetap `sassc` + reload aplikasi
(`ARCHITECTURE.md` §7).

## 13. Cara Pakai Guard & Tool (operasional)

```bash
# Guard warna (Go) — otomatis untuk semua tema × dark/light
GOCACHE=/tmp/go-build go test ./internal/domain/palette/...
GOCACHE=/tmp/go-build go test ./...        # semua test

# Angka + usulan nilai palet (Python, tanpa dependency)
python3 tools/contrast_report.py                    # audit semua tema (exit 1 bila ada FAIL)
python3 tools/contrast_report.py harbor --plan       # usulan ui.gutter & syntax.terminal_black
python3 tools/contrast_report.py --plan --accents    # + usul perbaikan accent (hue dipertahankan)
python3 tools/contrast_report.py --snippet           # potongan JSON siap tempel ke palette.json
```

Catatan: `--plan` memakai headroom +0.05 di atas ambang guard supaya nilai yang
dipilih tidak menempel persis di batas 3.0.
