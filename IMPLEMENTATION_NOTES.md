# Theme Engine – Gambaran Umum Implementasi

## Ringkasan
*Theme engine* ini bekerja secara dinamis untuk menghasilkan file konfigurasi berbagai aplikasi (terutama aplikasi terminal) berdasarkan sistem palet dan definisi tema yang terpusat.

---

## Komponen Utama

### 1. Arsitektur Inti
- **Resolver Variabel**: Menyelesaikan variabel secara dinamis menggunakan sintaks `$pl.extra.*`.
- **Pemuatan Berbasis Path**: Menghasilkan output yang disesuaikan dengan jenis aplikasi/alat yang dituju.
- **Rendering Template**: Menggunakan Go Template untuk membuat file konfigurasi secara dinamis.
- **Sistem Cache**: Menggunakan mekanisme *caching* cerdas untuk menghindari proses penulisan file yang tidak perlu.

---

### 2. Struktur Palet

**Skema Lama** *(sudah tidak digunakan)*:
- Menggunakan skema kaku: `Surface[3]`, `Overlay[3]`, `Accent(2)`, `Text(3)`, `Status(3)`.
- Penamaan warna semantik sangat terbatas.

**Skema Baru** *(terinspirasi dari Material Design & Catppuccin)*:
- **`accent`**: `{ primary, secondary, on_accent }`
- **`text`**: `{ primary, secondary, muted, link }`
- **`layer`**: `{ base, mantle, crust, surface, surface_raised, surface_overlay }`
- **`border`**: `{ default, active }`
- **`status`**: `{ success, warning, error, critical, info }`

---

### 3. Aplikasi yang Didukung

#### Aplikasi Domain (Pemrosesan Penuh)
- **GTK-CSS**: Menghasilkan file CSS untuk GTK dan tema Rasi (Rofi).
  - **Input**: `theme.json` + palet
  - **Output**: CSS untuk aplikasi GTK
  - **Lokasi Output**: `/output/domain/gtk/`

#### Aplikasi Standar (Injeksi Variabel Langsung)
- **Kitty**: Menginjeksikan warna langsung ke file konfigurasi Kitty.
- **Foot**: File konfigurasi untuk emulator terminal Foot.
- **Hypr**: Konfigurasi *window manager* Hyprland.
- **Yazi**: Tema untuk *file manager* Yazi.
- **Cava**: Konfigurasi untuk audio visualizer Cava.
- **Nvim**: Integrasi Neovim melalui file `colors.lua`.
- **Starship**: Prompt terminal lintas-shell dengan warna yang menyesuaikan tema aktif.

#### Integrasi Neovim (Berdasarkan Lua)
- Menghasilkan file `colors.lua` berisi variabel warna tema.
- Terintegrasi langsung dengan plugin tema seperti **Tokyo Night**.
- Memungkinkan perubahan warna secara dinamis tanpa perlu merombak konfigurasi Neovim utama.

#### Integrasi Starship
- Menggunakan template statis untuk menghasilkan konfigurasi prompt dalam format TOML (`starship.toml`).
- Menggunakan warna dari skema `accent` dan `status`, bukan warna kaku (ANSI hardcoded).
- *Layout* prompt tetap mengikuti preferensi pengguna; hanya warnanya saja yang berubah sesuai tema yang aktif.

---

### 4. Pola Desain (Design Patterns)

#### Struktur Template
- Setiap aplikasi memiliki template khusus.
- Penamaan variabel dibuat konsisten menggunakan format `{{ .Palette.Field }}`.
- Variabel akan terisi/terproses secara otomatis.

#### Konfigurasi Path
File `/config/path.txt` memetakan nama aplikasi ke template masing-masing:
```text
tool|template_path|output_dir
```

#### Kompatibilitas & Refaktorisasi
- **Waybar**: Template Waybar lama telah dihapus. Sekarang Waybar langsung mengimpor file CSS GTK agar lebih terpusat.
- **Dukungan Tema Lama**: Tema lama (seperti *Nocturne*) tetap dapat digunakan pada skema baru dengan memperbarui bagian template-nya.

---

### 5. Praktik Terbaik yang Diterapkan

1. **Penamaan Semantik**
   - **Sebelumnya:** `{primarysurface}`, `{secondarysurface}`
   - **Sekarang:** `{layer.surface}`, `{layer.surface_raised}`

2. **Keamanan Tipe (Type Safety)**
   - Struktur data diatur secara ketat menggunakan *struct* Go.
   - Antarmuka (*interface*) didefinisikan dengan jelas untuk setiap pemroses (processor).

3. **Operasi File Atomik**
   - Proses penulisan dilakukan ke file sementara (*temporary file*) terlebih dahulu sebelum diubah namanya (*rename*).
   - Sistem akan melewati (*skip*) file yang tidak mengalami perubahan untuk menghemat I/O disk.

4. **Dokumentasi Lengkap**
   - File `README.md` memuat panduan lengkap beserta contoh.
   - File `FLOW.md` menjelaskan alur arsitektur secara detail.
   - File `PROGRESS.md` digunakan untuk memantau status pengembangan.

---

### 6. Pengelolaan Warna

#### Injeksi Warna Dinamis

**Pendekatan Lama (Statis):**
```ini
warning = #ff80ff
critical = #cc33ff
```

**Pendekatan Baru (Kontekstual):**
```ini
warning = $pl.extra.status.warning
critical = $pl.extra.status.critical
```

#### Token Semantik
- Kode warna berbasis hex tidak lagi ditulis secara langsung (*hardcoded*).
- Makna warna tetap konsisten meskipun tema berganti.
- Memudahkan adaptasi ke skema warna baru di masa mendatang.

---

### 7. Contoh Penggunaan

#### Perintah Dasar
```bash
# Render konfigurasi untuk semua aplikasi
./runner.sh dev

# Render konfigurasi untuk aplikasi tertentu (contoh: Neovim)
./runner.sh dev nvim

# Build binary untuk kebutuhan produksi
./runner.sh build

# Jalankan aplikasi
./runner.sh run
```

#### Mengganti Tema
Untuk berganti tema, cukup ubah nilai pada file `state.json`:
```json
{
  "theme": {
    "name": "ghostly"
  }
}
```

---

### 8. Fitur Utama

1. **Dukungan Banyak Aplikasi**: Mendukung lebih dari 10 aplikasi berbeda.
2. **Palet Dinamis**: Memuat file palet format JSON secara fleksibel.
3. **Mekanisme Cache Cerdas**: Hanya memproses file yang mengalami perubahan.
4. **Penulisan File Atomik**: Mencegah file rusak akibat kegagalan proses penulisan di tengah jalan.
5. **Struktur File Teratur**: Kode tersusun rapi dengan pemisahan tanggung jawab (*separation of concerns*) yang jelas.
6. **Mudah Dikembangkan**: Memudahkan penambahan template atau pemroses aplikasi baru.
7. **Cross-Platform**: Dapat dijalankan di Linux, macOS, dan Windows.

---

### 9. Pengujian & Validasi

- **Unit Test**: Memastikan fungsi-fungsi inti berjalan sesuai ekspektasi.
- **Palette Test**: Memvalidasi proses pembacaan dan pemetaan palet warna.
- **Integration Test**: Menguji alur kerja sistem secara menyeluruh dari awal hingga akhir (*end-to-end*).
- **Smoke Test**: Memastikan tema-tema yang ada dapat dimuat tanpa kendala.

---

### 10. Lingkungan Pengembangan

#### Prasyarat
- Go 1.22 atau versi yang lebih baru
- Node.js (opsional, untuk *development tools*)
- Nix (opsional, untuk lingkungan pengembangan terisolasi)

#### Perintah Pengembangan
```bash
# Menjalankan pengujian (tests)
./runner.sh test

# Render tema dalam mode pengembangan (dev)
./runner.sh dev

# Mengompilasi binary
./runner.sh build

# Membersihkan file hasil generate/build
./runner.sh clean
```

---

## Perubahan Terbaru

- **Refaktor Skema Palet**: Mengubah struktur palet agar mengikuti standar sistem desain modern (menggantikan angka variabel dengan penamaan semantik seperti `layer`, `accent`, `border`, dan `status`).
- **Integrasi Starship**: Menambahkan generator untuk file `starship.toml` agar warna prompt otomatis menyesuaikan dengan tema aktif.
- **Penambahan Tema Baru**: Menyediakan tema *Kanagawa (Wave & Dragon)*, *Dracula*, dan *Sakura* dengan tingkat kecerahan ANSI yang sudah disesuaikan.
- **Integrasi Neovim**: Menghasilkan file `colors.lua` secara otomatis untuk mendukung plugin seperti *Tokyo Night*.
- **Penyederhanaan Waybar**: Menghapus template khusus Waybar dan dialihkan untuk langsung menggunakan file CSS GTK.

---

## Panduan Migrasi

### Mengubah Skema Lama ke Skema Baru

#### Pemetaan Kunci Palet
- **Kunci Lama**: `{primarysurface}`, `{secondarysurface}`, `{primaryaccent}`, `{secondaryaccent}`
- **Kunci Baru**: `{layer.surface}`, `{layer.surface_raised}`, `{accent.primary}`, `{accent.secondary}`

**Contoh pada Template:**
```lua
-- Format Lama
{{- index .Palette "primarysurface" -}}

-- Format Baru
{{- .Palette.LayerSurface -}}
```

#### Konfigurasi Aplikasi
Tidak ada perubahan khusus yang perlu dilakukan pada konfigurasi aplikasi. Sistem pemroses akan memetakan token semantik secara otomatis.

---

## Cara Mengembangkan Sistem

### Menambahkan Aplikasi Baru
1. Buat pemroses (*processor*) baru di direktori `internal/core/tools/.aplikasi_baru/`.
2. Tambahkan file template di `assets/templates/tools/aplikasi_baru/`.
3. Daftarkan aplikasi tersebut ke dalam file `config/path.txt`.

### Mengubah Skema Palet
1. Perbarui *struct* `Palette` di `palette/raw.go`.
2. Perbarui `ResolvedPalette` di `palette/palette.go`.
3. Pasang logika pemetaan baru pada *resolver* di `palette/resolver.go`.
4. Sesuaikan fungsi perataan (*flattening*) di `theme/flatter.go`.
5. Perbarui seluruh template yang menggunakan variabel terkait.

---

## Panduan untuk Kontributor

1. **Pemilihan Tema**: Pastikan tema yang dibuat memiliki kontras yang baik dan mendukung mode gelap/terang sesuai standar tema yang ada (seperti Nocturne, Ghostly, Kanagawa, Dracula, atau Sakura).
2. **Kerapian Kode**: Pertahankan struktur proyek yang sudah ada dan gunakan konvensi penulisan standar bahasa Go.
3. **Pengujian**: Selalu jalankan perintah `./runner.sh test` sebelum melakukan *commit* dan pastikan hasil render warna sudah sesuai secara visual.
4. **Dokumentasi**: Perbarui dokumentasi jika Anda menambahkan fitur baru atau mengubah konfigurasi yang ada.

---

## Kesimpulan

*Theme engine* ini dirancang untuk mempermudah pengelolaan tema di berbagai aplikasi secara konsisten dan terpusat. Dengan memanfaatkan skema palet semantik dan sistem *template*, Anda dapat mengubah tampilan seluruh lingkungan sistem hanya dengan satu kali konfigurasi.
