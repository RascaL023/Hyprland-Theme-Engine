# FLOW & ARSITEKTUR THEME ENGINE

Dokumen ini menjelaskan secara mendalam bagaimana **Theme Engine** bekerja dari awal eksekusi hingga menghasilkan file konfigurasi akhir yang digunakan oleh aplikasi-aplikasi desktop (Hyprland, Waybar, Kitty, Foot, Cava, GTK, dll.).

---

## 1. Konsep Utama: Single Source of Truth (SSOT)

Sebelum memahami kodenya, pahami filosofi desainnya:
- **Satu Sumber Data**: Pengaturan font, ukuran, warna, dan dekorasi berada di satu tempat (`themes/<nama-tema>/`).
- **Separasi Data & Struktur**: 
  - `palette.json` menyimpan definisi warna dasar (Hex) dan warna turunan.
  - `theme.json` menyimpan preferensi global (font, size) dan konfigurasi spesifik aplikasi (opacity, border).
- **Zero-cost Runtime**: Mengurangi penulisan disk jika tidak diperlukan dan menggunakan caching template agar proses rendering sangat instan.

---

## 2. Struktur Direktori Utama

Berikut peta modul internal agar Anda bisa bernavigasi dengan mudah:

```txt
cmd/theme-engine/              # Titik masuk aplikasi (CLI parser)
internal/
 ├── engine/                   # Pipeline utama: orchestrator flow aplikasi
 ├── loader/                   # I/O Loader untuk JSON, State, dan Path Map (path.txt)
 ├── resolver/                 # Engine penyelesai variabel seperti "$pl.extra.layer.base"
 ├── renderer/                 # Template parser, cache, & atomic write dengan skip-unchanged
 ├── processor/                # Pendaftaran (Registry) dan pemanggilan Processor target
 ├── register/                 # Interface kontrak (Parser, Resolver, Renderer) bagi Processor
 └── core/
      ├── context/             # Context global yang menampung Palette & Theme ter-resolve
      ├── log/                 # Log helper (Info, Warn, Error)
      ├── themes/              # Struktur data model (Palette, Theme, State)
      └── tools/               # Implementasi Processor spesifik per aplikasi / domain
```

---

## 3. Siklus Hidup Eksekusi (Runtime Lifecycle)

Ketika Anda menjalankan `./runner.sh run` atau `go run ./cmd/theme-engine [target]`, urutan kejadian di bawah kap mesin adalah sebagai berikut:

```txt
[1. CLI Entrypoint]
       │
       ▼
[2. Load State (.state.json)] ──► Ambil nama tema & tipe (dark/light)
       │
       ▼
[3. Load Path Map (path.txt)] ──► Cari lokasi template & output, expand variabel ($WAYBAR)
       │
       ▼
[4. Load & Resolve Palette]   ──► Baca palette.json, selesaikan variabel internal, lalu flatten
       │
       ▼
[5. Load Theme (theme.json)]  ──► Baca metadata global & config tools mentah
       │
       ▼
[6. Build Engine & Context]   ──► Buat Context global (Palette + Theme ter-resolve)
       │
       ▼
[7. Filter Target]            ──► Jika ada argumen, render target itu saja; jika tidak, render semua
       │
       ▼
[8. Loop & Execute Processors]
       │
       ├─► a. Dapatkan processor teregistrasi (misal "kitty")
       ├─► b. Parse: Ekstrak sub-JSON config tool dari theme.json
       ├─► c. Resolve: Selesaikan variabel warna "$pl.xxxx" ke Hex nyata
       └─► d. Render: Eksekusi Go template -> Bandingkan isi -> Tulis atomik ke tujuan
```

### Detail Langkah demi Langkah:

#### Langkah 1: Membaca State (`config/.state.json`)
Engine memanggil `loader.LoadJSON[state.State]` untuk mengetahui tema aktif saat ini.
Contoh isi `.state.json`:
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

#### Langkah 2: Membaca Path Map (`config/path.txt`)
`LoadToolMap` membaca file pemetaan path. Sebelum diproses, fungsi `ExpandPath` dipanggil untuk menyubstitusi variabel seperti `$WAYBAR` dengan nilai dari `.state.json` atau `$THEME` dengan nama tema aktif.
Contoh baris map:
```txt
waybar|assets/templates/waybar/$WAYBAR.tmpl|output/waybar/sources.css
```
Jika `state.Waybar` adalah `"default"`, ini dikembangkan menjadi:
```txt
waybar|assets/templates/waybar/default.tmpl|output/waybar/sources.css
```

#### Langkah 3: Load & Resolve Palette (`themes/<tema>/palette.json`)
1. Membaca data palette kasar (`palette.Raw`).
2. Menentukan subset varian sesuai tipe aktif (`dark` atau `light`) via `ResolveSelected(type)`.
3. Menyelesaikan warna referensi internal (seperti `$pl.color5`) ke nilai warna Hex asli menggunakan `resolver.ResolveVar`.
4. Memanggil `themeconfig.BuildFlattenPalette` untuk meratakan seluruh warna ke dalam `map[string]string` dengan format nama kunci flat (contoh: `"extra.primaryaccent"`). Ini bertujuan agar pencarian warna berikutnya sangat cepat.

#### Langkah 4: Load Theme Config (`themes/<tema>/theme.json`)
Membaca file konfigurasi utama tema yang berisi properti font global dan konfigurasi spesifik masing-masing alat/aplikasi dalam representasi `map[string]json.RawMessage`.

#### Langkah 5: Pencarian & Inisialisasi Processor
Setiap aplikasi (seperti Kitty, Foot) mendaftarkan dirinya secara otomatis ke `processor.Registered` melalui blok `init()` di package-nya masing-masing. Engine mencocokkan nama target dengan processor yang terdaftar.

#### Langkah 6: Pemrosesan Data oleh Processor
Setiap processor memenuhi interface kontrak `Processor` yang terdiri dari tiga metode utama:
1. **`Parse(in any) (any, error)`**: Mengonversi sub-JSON mentah (`json.RawMessage`) dari `theme.json` khusus untuk tool tersebut menjadi Go struct representatif tool tersebut (misal `kitty.Raw`).
2. **`Resolve(in any, ctx *context.Context) (any, error)`**: Menerima struct mentah dan mengembalikan struct siap render. Di langkah ini, nilai bertipe variabel (seperti `"$pl.extra.accent.primary"`) diselesaikan menjadi warna Hex asli dengan memanggil `resolver.ResolveVar` terhadap map warna yang sudah di-*flatten* di Langkah 3.
3. **`Render(templatePath, outputPath string, data any) error`**: Mengirim data yang sudah matang ke mesin pembuat file.

#### Langkah 7: Template Rendering & Atomic/Skip Write
Modul `internal/renderer/render.go` menangani penulisan file ke filesystem dengan optimasi tinggi:
1. **Template Caching**: Template hanya diparse dari disk sekali per path eksekusi, lalu disimpan di cache memori. Eksekusi berikutnya untuk template yang sama langsung menggunakan cache.
2. **Skip-Unchanged**: Membaca isi file tujuan jika sudah ada. Jika hasil render baru sama persis dengan isi file lama, penulisan dilewati sepenuhnya untuk menghindari SSD/HDD wear-and-tear serta tidak mengganggu utilitas watch-reload eksternal.
3. **Atomic Write**: Jika konten berubah, engine menulis data ke file temp di folder tujuan, lalu melakukan rename sistem operasi secara instan. Ini menjamin file konfigurasi tidak pernah terputus atau rusak di tengah jalan jika proses tiba-tiba mati.

---

## 4. Mekanisme Resolving Variabel `$pl.`

Bagaimana `$pl.extra.accent.primary` berubah menjadi kode warna hex nyata?

1. Di `internal/resolver/var-resolver.go`, fungsi `ResolveVar(s string, sources ...vars.VarSource)` bertugas menyaring string input.
2. Jika string **tidak** diawali `$pl.`, fungsi langsung mengembalikan string tersebut (misal, string `#ffffff` atau string biasa tetap utuh).
3. Jika string diawali `$pl.`, prefix tersebut dipotong dan sisanya dicari di daftar penyedia variabel (`VarSource`).
4. Selama resolusi palette mentah, penyedia variabel adalah `RawPaletteVars` yang membaca warna ANSI dari indeks array `colors` (seperti `color0`, `color15`).
5. Selama resolusi config tool, penyedia variabel adalah `ResolvedPaletteVars` yang membaca langsung dari map `Flat` yang berisi pasangan key-value warna ter-resolve lengkap.

---

## 5. Klasifikasi Target / Processor

Di dalam Theme Engine, terdapat tiga jenis target pemrosesan tergantung kompleksitas kebutuhan aplikasi target:

### A. Generic / Template-Only (Static Processor)
- **Karakteristik**: Aplikasi yang konfigurasinya hanya memerlukan data palette global dan metadata global theme tanpa memerlukan parser konfigurasi khusus sendiri di `theme.json`.
- **Lokasi**: `internal/core/tools/static/processor.go`
- **Contoh target**: `waybar`, `hyprland`.
- **Cara Kerja**: Processor ini langsung melewatkan objek Context global (`*context.Context` yang berisi `.Palette` dan `.Theme`) ke dalam template. Semua kustomisasi dilakukan langsung di file `.tmpl` menggunakan sintaks Go template biasa.

### B. Custom Tool Processor
- **Karakteristik**: Aplikasi yang memerlukan struktur konfigurasi unik di dalam `theme.json` untuk mengontrol perilakunya (seperti `cursorShape` di Kitty atau `gradients` di Cava).
- **Lokasi**: `internal/core/tools/<nama_tool>/`
- **Contoh target**: `kitty`, `foot`, `cava`.
- **Cara Kerja**: Memiliki Go struct sendiri di file `raw.go` dan `tool.go`. Processor mem-parse sub-JSON milik dirinya dari `theme.json`, me-resolve warna, lalu mengirim objek instansi struct khususnya sendiri ke template.

### C. Domain Processor
- **Karakteristik**: Mengatur cakupan (domain) konfigurasi yang sangat besar dan dinamis, sering kali menghasilkan beberapa file keluaran dan melibatkan proses kompilasi eksternal.
- **Lokasi**: `internal/core/domain/<nama_domain>/`
- **Contoh target**: `gtk`.
- **Cara Kerja**: GTK processor me-render template scss (`_source.scss`), kemudian secara otomatis memanggil command eksternal `sassc` untuk mengompilasi SCSS dasar tersebut menjadi CSS siap pakai (`source.css` untuk GTK/Waybar dan `source.rasi` untuk launcher Rofi).

---

## 6. Alur Data di Template (Go text/template)

Di dalam file `.tmpl`, Anda memiliki akses penuh ke data yang dikirim oleh Processor. 

### Helper Khusus: `hex`
Secara bawaan, warna di palette diawali dengan tanda `#` (contoh `#1e1e2e`). Beberapa konfigurasi aplikasi (seperti CSS rgb, atau beberapa format tool) membutuhkan warna tanpa tanda `#`.
Gunakan helper `hex` di template Anda:
```gotemplate
/* Mengembalikan "1e1e2e" */
color = {{ hex .Palette.Background }}
```

### Akses Data pada Static Processor (Waybar/Hyprland)
Karena Static Processor meneruskan Context global, Anda mengaksesnya lewat properti kapital:
```gotemplate
# Mengakses warna palette
foreground_color = {{ .Palette.Foreground }}

# Mengakses font dari tema global
font_name = {{ .Theme.Theme.Fonts.Primary }}
```

### Akses Data pada Custom Tool Processor (Kitty/Foot/Cava)
Data diakses langsung lewat properti struct matang yang dibuat oleh processor terkait:
```gotemplate
# Mengakses warna palette (karena disalin ke structnya)
cursor = {{ .Palette.Cursor }}

# Mengakses properti unik yang sudah di-resolve
cursor_shape = {{ .CursorShape }}
font_size = {{ .FontSize }}
```

---

## 7. Panduan Cepat Menambah Target Baru

Gunakan checklist ini jika Anda ingin memperluas Theme Engine di masa depan:

### Opsi A: Jika target baru cukup dengan Template-Only
1. Buat template baru di `assets/templates/tools/<nama>/config.tmpl`.
2. Daftarkan target baru tersebut di `internal/core/tools/static/processor.go` menggunakan fungsi `processor.RegisterProcessor(Processor{name: "nama_tool_baru"})`.
3. Tambahkan baris pemetaan baru di `config/path.txt`:
   ```txt
   nama_tool_baru|assets/templates/tools/nama/config.tmpl|output/tools/nama/config
   ```

### Opsi B: Jika target baru butuh konfigurasi khusus di `theme.json`
1. Buat folder baru di `internal/core/tools/<nama_tool>/`.
2. Buat tiga file utama:
   - `raw.go`: Berisi struct JSON pencerminan opsi konfigurasi di `theme.json`.
   - `<nama_tool>.go`: Berisi struct final siap pakai di template.
   - `processor.go`: Implementasikan interface `Processor` (metode `Name`, `Parse`, `Resolve`, `Render`) dan panggil registrasinya di fungsi `init()`.
3. Daftarkan package baru Anda di `cmd/theme-engine/import.go` agar di-import secara blank/side-effect:
   ```go
   _ "theme-engine/internal/core/tools/<nama_tool>"
   ```
4. Tambahkan konfigurasi default tool tersebut di `themes/<tema>/theme.json`.
5. Buat template di `assets/templates/tools/<nama_tool>/` dan daftarkan path pemetaannya di `config/path.txt`.
