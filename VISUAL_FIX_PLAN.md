# Rencana Perbaikan Bug Visual (Nvim & Lazygit)

Status: **Fase 0, 1, 2 & 3 SELESAI untuk pilot `kanagawa-dragon`. Fase 4-5 belum dikerjakan.**

Ringkasan keputusan & analisis performa: lihat **§11** dan **§12**.

Tema pilot: **`kanagawa-dragon` SAJA.**
Tema lain (`nocturne`, `ghostly`, `sakura`, `kanagawa-wave`, `claude-manjuska`)
**tidak disentuh** sampai pilot ini lulus semua checklist.

---

## 0. Ringkasan Masalah (hasil audit)

| # | Gejala | Akar masalah | Bukti |
|---|---|---|---|
| A | Teks backtick (`` `test` ``) di `.md` susah dibaca | Template `nvim/colors.tmpl` memetakan `terminal_black` ke **ANSI bright-black** (`colors[8]`), padahal Tokyonight memakai `terminal_black` sebagai **warna latar (background)** untuk inline code. Di beberapa tema `colors[8]` justru terang → latar & teks sama-sama terang. | `tokyonight/.../groups/treesitter.lua:66` (`@markup.raw.markdown_inline = { bg = c.terminal_black, fg = c.blue }`), `colors.tmpl:38` |
| B | Bagian tengah lualine (nama file/tab) kurang kontras | `fg_gutter` dipetakan ke `layer.surface_overlay` (layer **paling terang**), padahal Tokyonight asli `fg_gutter` itu **gelap**. Mode color jadi kehilangan kontras. | `_tokyonight.lua:11` (`b = { bg = c.fg_gutter, fg = c.blue }`), `colors.tmpl:14` |
| C | Di lazygit, baris yang di-select kurang jelas & tema tidak terasa | **Belum ada target `lazygit` sama sekali** (tidak ada template/path.txt/wiring). `~/.config/lazygit/config.yml` kosong (0 byte), jadi lazygit memakai default: `selectedLineBgColor: [blue]` + `defaultFgColor: [default]`. | `docs/Config.md` (default resmi), `config/path.txt` |

Temuan tambahan (dicatat, belum dikerjakan): `sakura`/`kanagawa-wave`/`ghostly`
belum mengisi blok `extra.fg`, `extra.bg`, `text.visited`, `border.medium`;
`sakura` bahkan belum punya `syntax` & `ui`.

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

## 7. Fase 4 — Sisi Nvim: markview (opsional, setelah pilot)

**Masalah:** `transparent = true` + `Normal.bg = "none"` membuat markview
(`MarkviewInlineCode`) memakai warna fallback hardcoded `#1E1E2E`, bukan warna
tema aktif. Ini bisa jadi sumber bug backtick kedua.

- [ ] Verifikasi dulu via `:Inspect` pada backtick: grup mana yang aktif?
- [ ] Jika `MarkviewInlineCode`: set `vim.g.markview_dark_bg` /
      `vim.g.markview_light_bg` dari `.Palette` di `lua/core/theme.lua`
- [ ] Ulangi cek kontras

---

## 8. Fase 5 — Rollout ke Tema Lain (setelah pilot lulus, JANGAN sekarang)

- [ ] Terapkan field/nilai yang sama ke `nocturne`, `ghostly`, `sakura`,
      `kanagawa-wave`, `claude-manjuska`
- [ ] Lengkapi blok `extra.fg`/`extra.bg`/`text.visited`/`border.medium` yang kurang
- [ ] Update tabel mapping Tokyonight di `RULE.md` (`terminal_black`, `dark3`,
      `dark5`, `fg_gutter`, `blue0/blue7`)
- [ ] Update `ARCHITECTURE.md` §9 dengan target `lazygit`
- [ ] Update `PROGRESS.md` (daftar target + status tema)
- [ ] Regenerate golden: `UPDATE_GOLDEN=1 go test ./internal/app/engine/ -run TestGolden`

---

## 9. Verifikasi Akhir

- [x] `GOCACHE=/tmp/go-build go test ./...` hijau semua
- [x] Contrast test lulus untuk pilot (dark 4.48:1 & light 3.79:1)
- [ ] `output/` pilot di-review manual
- [ ] Uji mata di nvim (`colorscheme tokyonight`) + lazygit
- [ ] Regenerate & commit golden (branch `feature/test/golden`)

---

## 10. Checklist Master (Ringkas)

- [x] F0 `contrast.go` + `contrast_test.go`
- [x] F1 backtick: slot `syntax.terminal_black` (keputusan: Opsi A)
- [x] F2 lualine B: slot `ui.gutter` + tuning hijau light pilot (keputusan: Opsi 3 + 2)
- [x] F3 target lazygit
- [ ] F4 markview globals (di `~/.config/nvim`, sudah diizinkan pemilik)
- [ ] F5 rollout tema lain (setelah pilot)
- [ ] Regenerate golden + update dokumentasi

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

Tidak ada loop, komputasi warna, atau pembacaan file tambahan di jalur render.
Kecepatan switch tetap sama; cost terbesar tetap `sassc` + reload aplikasi
(`ARCHITECTURE.md` §7).
