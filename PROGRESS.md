# PROGRESS & ROADMAP DOKUMENTASI

Dokumen ini melacak status rilis, fitur yang telah diimplementasikan, serta memberikan panduan terarah bagi pengembang saat ingin melanjutkan proyek **Theme Engine**.

---

## 1. Ringkasan Status Proyek

- **Status Saat Ini**: Aktif / Siap Dilanjutkan (Stabil & Teruji)
- **Arsitektur Inti (Core Engine)**: **100% Selesai & Stabil** (Mendukung pemrosesan state, pemetaan dinamis, caching, penulisan aman/atomik, dan skip-unchanged).
- **Cakupan Pengujian**: Unit test mencakup modul krusial seperti `loader` dan `renderer`.
- **[BARU] Schema Rebranding & Expansion**: Kategori `extra` di-*refactor* total dari skema
  `surface[3]` + `overlay[3]` + `base/mantle/crust` + `accent[2]` + `text[3]` menjadi
  terstruktur seperti *design system* modern:
  - `accent` → ditambah `on_accent`
  - `text` → `teritary` diganti `muted`, ditambah `link`
  - `layer` → pengganti `base/mantle/crust` + `surface[3]` + `overlay[3]`
  - `border` → kategori baru (default & active)
  - `status` → kategori baru (success, warning, error, critical, info)
- **[BARU] Ghostly Theme**: Tema ghostly kini sudah memiliki `palette.json` (dark & light) dan `theme.json`.

---

## 2. Peta Kemajuan (Progress Matrix)

### A. Infrastruktur Inti (Core Engine)
- [x] **CLI Parsing**: Tipis, cepat, dan menerima target tunggal atau semua target (`main.go`).
- [x] **State Management**: Membaca aktif tema dan variannya (`.state.json`).
- [x] **Dynamic Path Mapping**: Membaca, mengabaikan komentar, dan mengekspansi variabel seperti `$WAYBAR` dan `$THEME` (`path.txt`).
- [x] **Dynamic Palette Resolver**: Mengurai palette mentah, menyubstitusi referensi internal (seperti `$pl.color5`), dan meratakan (*flattening*) warna agar cepat diakses.
- [x] **Template Caching**: Parsing template hanya sekali per path eksekusi.
- [x] **Skip-Unchanged Writes**: Hanya menulis file jika isi hasil render berbeda dengan file di disk (mengurangi SSD wear dan meminimalkan trigger hot-reload eksternal).
- [x] **Atomic Disk Write**: Menulis ke file temp terlebih dahulu lalu melakukan rename agar terhindar dari file konfigurasi korup atau terpotong jika proses mati tengah jalan.

### B. Implementasi Target & Processor
| Target | Tipe Processor | Status | File Template | Lokasi Output |
| :--- | :--- | :---: | :--- | :--- |
| **GTK-CSS** | Domain | [x] Selesai | `assets/templates/domain/gtk/source.tmpl` | `output/domain/gtk/` (menghasilkan CSS & Rasi via `sassc`) |
| **Cava** | Custom Tool | [x] Selesai | `assets/templates/tools/cava/cava.tmpl` | `output/tools/cava/cava_extra` |
| **Foot** | Custom Tool | [x] Selesai | `assets/templates/tools/foot/foot.tmpl` | `output/tools/foot/ui.ini` |
| **Kitty** | Custom Tool | [x] Selesai | `assets/templates/tools/kitty/kitty.tmpl` | `output/tools/kitty/ui.conf` |
| **Hypr** | Custom Tool | [x] Selesai | `assets/templates/tools/hypr/hypr.tmpl` | `output/tools/hypr/source.conf` |
| **Yazi** | Static | [x] Selesai | `assets/templates/tools/yazi/theme.tmpl` | `output/tools/yazi/theme.toml` |
| **Nvim** | Static | [x] Selesai | `assets/templates/tools/nvim/colors.tmpl` | `output/tools/nvim/colors.lua` |
| **Starship** | Static | [x] Selesai | `assets/templates/tools/starship/starship.tmpl` | `output/tools/starship/starship.toml` |
| **Waybar** | Migrated to GTK | [x] Selesai | via `output/domain/gtk/css/source.css` | Config langsung `@import` source.css dari GTK |
| **Ncmcpp** | Static / Custom | [ ] Rencana | - | - |
| **Btop** | Static / Custom | [ ] Rencana | - | - |
| **Micro** | Static / Custom | [ ] Rencana | - | - |
| **Swaync** | Static / Custom | [ ] Rencana | - | - |
| **Nemo** | Static / Custom | [ ] Rencana | - | - |
| **Zathura** | Static / Custom | [ ] Rencana | - | - |

### C. Daftar Tema
| Tema | Dark | Light | Theme JSON |
| :--- | :---: | :---: | :---: |
| **Nocturne** | [x] | [x] | [x] |
| **Ghostly** | [x] | [x] | [x] |
| **Kanagawa Wave** | [x] | [x] | [x] |
| **Kanagawa Dragon** | [x] | [x] | [x] |
| **Custom Kanagawa Wave** | [x] | [x] | [x] |
| **Custom Kanagawa Dragon** | [x] | [x] | [x] |
| **Kanagawa Dragon Original** | [x] | [x] | [x] |
| **Dracula** | [x] | [x] | [x] |
| **Sakura** | [x] | [x] | [x] |

### D. Perkakas Pengembangan (Developer Tooling)
- [x] **Runner Script (`runner.sh`)**: Orkestrator otomatis untuk kompilasi, pengujian, pembersihan, dan penanganan log senyap dengan notifikasi desktop (`notify-send`).
- [x] **Unit Testing**: Pengujian otomatis untuk penulisan file, render, caching, dan pemetaan path.

---

## 3. Rencana Langkah Selanjutnya (Next Action Items)

Jika Anda ingin melanjutkan coding sekarang, pilih salah satu dari tugas terarah berikut:

### 🎯 Prioritas 1: Menambah Target Baru (Ncmcpp & Nvim)
- **Tujuan**: Membawa integrasi music player (Ncmcpp) dan editor (Nvim) ke dalam sistem tema terpadu.
- **Langkah-Langkah**:
  1. Tentukan apakah target baru ini butuh opsi unik di `theme.json` (Gunakan **Custom Tool Processor**) atau cukup data warna global saja (Gunakan **Generic Static Processor**).
  2. Rujuk panduan penambahan target di **`FLOW.md` Bab 7** untuk cara pembuatan kodenya.
  3. Buat template `.tmpl`-nya dan daftarkan di `config/path.txt`.

### 🎯 Prioritas 2: Validasi Skema JSON & Error Handling yang Lebih Ramah
- **Tujuan**: Mencegah aplikasi crash jika pengguna melakukan salah ketik warna atau properti di `theme.json` atau `palette.json`.
- **Langkah-Langkah**:
  1. Tambahkan pengecekan format hex warna (misal, memastikan warna diawali dengan `#` dan memiliki 6 digit hex yang valid).
  2. Implementasikan parser error reporter yang menunjukkan letak file dan baris yang bermasalah.

### 🎯 Prioritas 3: CLI Flag Interaktif & Mode Watch (Hot Reload)
- **Tujuan**: Mempermudah penggantian tema lewat terminal dan melakukan render ulang secara otomatis saat file template diubah.
- **Langkah-Langkah**:
  1. Tambahkan argument parsing yang lebih kaya di `cmd/theme-engine/main.go` (misal, command `theme-engine set-theme <nama_tema>` untuk mengedit `.state.json` otomatis).
  2. Gunakan library pembaca perubahan file filesystem (seperti `fsnotify`) untuk mengimplementasikan mode `--watch` atau `watch` command.

---

## 4. Cheat Sheet Pengoperasian (Cara Cepat Melanjutkan)

Gunakan perintah-perintah praktis berikut di terminal untuk berinteraksi dengan proyek:

### Menjalankan Kode (Mode Development)
Jalankan file tanpa melakukan kompilasi terlebih dahulu (sangat cocok saat mengedit kode Go):
```bash
./runner.sh dev          # Render semua target
./runner.sh dev kitty    # Hanya render target kitty
```

### Menjalankan Tes Unit (Unit Testing)
Pastikan tidak ada kode yang rusak setelah Anda melakukan modifikasi:
```bash
./runner.sh test
```

### Membangun File Binary Kompilasi (Production-ready)
Lakukan kompilasi ke file biner lokal untuk performa maksimal (biasanya dipanggil oleh keybinding desktop):
```bash
./runner.sh build
```

### Menjalankan Biner Kompilasi
Gunakan mode `run` jika binary sudah dikompilasi (sangat cepat):
```bash
./runner.sh run          # Jalankan render semua target
./runner.sh run hyprland # Hanya render hyprland
```

### Tips Pengujian Mandiri:
1. File hasil kompilasi dan rendering diletakkan di direktori `output/`.
2. Periksa isi file di `output/` setelah menjalankan renderer untuk memverifikasi apakah perubahan Anda pada template sudah diaplikasikan dengan benar.
3. Gunakan `config/.state.json` untuk mengganti tema aktif (misalnya ubah `"name": "nocturne"` menjadi tema lain di folder `themes/`).
