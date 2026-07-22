# Theme Color Rules

Dokumen ini adalah acuan saat membuat tema baru. Tujuannya supaya setiap warna punya makna semantic yang konsisten, bukan sekadar warna yang terlihat bagus secara terpisah.

Theme engine ini memakai palette semantic sendiri, lalu output akan memetakan warna tersebut ke variabel milik tool terkait. Karena itu, urutan terang-gelap dan fungsi tiap field harus dijaga.

## Prinsip Utama

- Jangan isi field berdasarkan nama warna saja. Isi berdasarkan fungsi visualnya.
- `text.*` adalah warna foreground/text, bukan warna background UI.
- `layer.*` adalah warna background/surface, bukan warna text.
- `fg.*` adalah warna foreground UI variants, bukan text konten.
- `border.*` adalah warna pemisah/outline structural.
- `bg.*` adalah warna background UI tambahan (bukan surface structural).
- ANSI `colors[0..15]` dipakai terminal dan juga beberapa adapter UI seperti lualine.
- Untuk dark theme, warna layer harus bergerak dari gelap ke lebih terang secara bertahap.
- Jangan pakai `text.muted` sebagai background structural kecuali memang sengaja ingin UI low-contrast.

## Layer Colors

Layer adalah fondasi background UI. Untuk dark theme, urutan ideal dari tergelap ke terang:

```text
layer.crust
layer.mantle
layer.base
layer.surface
layer.surface_raised
border.default
layer.surface_overlay
```

Makna tiap field:

| Field | Makna | Contoh penggunaan |
|---|---|---|
| `layer.crust` | Background paling gelap | deep background, section fg on accent |
| `layer.mantle` | Background lebih gelap dari base | sidebar, popup, float bg |
| `layer.base` | Background utama | editor bg, terminal bg |
| `layer.surface` | Surface halus di atas base | cursorline, subtle highlight |
| `layer.surface_raised` | Surface lebih terlihat | statusline bg, raised component, dialog bg |
| `layer.surface_overlay` | Structural muted UI | gutter-like color, inactive/subtle UI, alt bg |

Rule penting:

```text
layer.crust < layer.mantle < layer.base < layer.surface < layer.surface_raised < layer.surface_overlay
```

Untuk dark theme, tanda `<` berarti lebih gelap secara luminance. Untuk light theme, prinsipnya dibalik secara visual: base harus tetap nyaman sebagai background utama, dan surface/raised/overlay harus memberi separasi yang cukup.

## Text Colors

Text colors dipakai untuk foreground konten. Jangan pakai warna ini sebagai background structural kecuali ada alasan khusus.

| Field | Makna | Rule |
|---|---|---|
| `text.primary` | Teks utama | kontras tinggi terhadap `layer.base` |
| `text.secondary` | Teks sekunder | masih readable, lebih kalem dari primary |
| `text.muted` | Komentar / disabled-ish text | readable tapi low-emphasis |
| `text.link` | Link / navigasi / reference | biasanya dekat dengan accent/info |
| `text.visited` | Visited link | lebih muted/different dari link |

Untuk dark theme:

```text
text.primary > text.secondary > text.muted > layer.surface_overlay
```

Catatan penting:

- `text.muted` boleh terang selama masih cocok untuk komentar.
- `text.muted` jangan dipakai untuk `fg_gutter` Tokyonight karena `fg_gutter` dipakai sebagai background lualine section B.

## Border Colors

| Field | Makna | Rule |
|---|---|---|
| `border.default` | Border non-active | lebih subtle dari active, masih terlihat di atas surface |
| `border.active` | Border active/focused | biasanya sama atau dekat dengan accent utama |
| `border.medium` | Border medium emphasis | antara default dan active |

Untuk dark theme, `border.default` idealnya berada di antara `layer.surface_raised` dan `layer.surface_overlay`, atau setidaknya tidak lebih terang dari text muted. `border.medium` bisa lebih terang dari default.

## Accent Colors

| Field | Makna |
|---|---|
| `accent.primary` | Accent utama theme |
| `accent.secondary` | Accent pendamping |
| `accent.on_accent` | Foreground saat berada di atas accent background |

Rule penting:

- `accent.on_accent` harus kontras tinggi terhadap `accent.primary`.
- Untuk dark theme, `accent.on_accent` biasanya `layer.crust` atau warna gelap lain.
- Untuk light theme, `accent.on_accent` biasanya putih atau warna sangat terang, tergantung accent-nya.

## ANSI Colors

`colors` wajib berisi 16 warna dengan urutan ANSI standar:

| Index | Makna |
|---|---|
| `0` | black |
| `1` | red |
| `2` | green |
| `3` | yellow |
| `4` | blue |
| `5` | magenta |
| `6` | cyan |
| `7` | white |
| `8` | bright black |
| `9` | bright red |
| `10` | bright green |
| `11` | bright yellow |
| `12` | bright blue |
| `13` | bright magenta |
| `14` | bright cyan |
| `15` | bright white |

Rule penting:

- `colors[0]` harus cocok sebagai terminal black.
- `colors[7]` harus cocok sebagai terminal white.
- `colors[8]` adalah bright black, bukan background UI. Jangan diasumsikan sebagai layer color.
- `colors[4]`, `colors[2]`, `colors[3]`, `colors[5]`, `colors[1]`, dan `syntax.green1` dipakai oleh Tokyonight lualine sebagai mode accent.
- Mode accent harus cukup kontras terhadap `layer.surface_overlay` karena lualine section B memakai `fg_gutter` sebagai background.

Target praktis untuk lualine section B dark theme:

```text
contrast(colors[4], layer.surface_overlay) >= 3:1
contrast(colors[2], layer.surface_overlay) >= 3:1
contrast(colors[3], layer.surface_overlay) >= 3:1
contrast(colors[5], layer.surface_overlay) >= 3:1
contrast(colors[1], layer.surface_overlay) >= 3:1
```

## Status Colors

| Field | Makna |
|---|---|
| `status.success` | success / ok / added |
| `status.warning` | warning / caution |
| `status.error` | error |
| `status.critical` | critical / stronger error |
| `status.info` | info / hint |

Rule penting:

- Status colors harus readable di atas `layer.base`, `layer.surface`, dan `layer.surface_overlay`.
- `status.critical` biasanya lebih kuat/terang daripada `status.error`.

## Fg Colors

Fg colors adalah foreground variants untuk UI elements, berbeda dengan `text.*` yang untuk konten.

| Field | Makna | Fallback |
|---|---|---|
| `fg.dim` | Dim/low-emphasis foreground | `text.muted` |
| `fg.dim_muted` | Dim muted foreground | `text.muted` |
| `fg.disabled` | Disabled element foreground | `text.muted` |

Untuk dark theme, nilai ideal:

```text
fg.dim < text.muted < text.secondary < text.primary
```

## Bg Colors

Bg colors adalah background variants tambahan untuk UI spesifik (bukan surface structural).

| Field | Makna | Fallback |
|---|---|---|
| `bg.conflict` | Conflict/diff marker background | `status.warning` |
| `bg.disk_usage` | Disk usage meter background | `layer.surface_overlay` |

## Syntax Colors

Syntax block dipakai terutama untuk Neovim/Tokyonight compatibility.

| Field | Makna umum |
|---|---|
| `syntax.purple` | keyword/special/purple accent |
| `syntax.magenta2` | bright magenta/red accent |
| `syntax.blue0` | dark blue-ish background selection/helper |
| `syntax.blue1` | brighter blue/link/type accent |
| `syntax.blue5` | base blue accent |
| `syntax.blue6` | cyan/teal-ish blue accent |
| `syntax.blue7` | dark blue/structural color |
| `syntax.green1` | bright green accent |
| `syntax.green2` | base green accent |
| `syntax.orange` | orange/warning/string-ish accent |
| `syntax.red1` | stronger red diagnostic/accent |
| `syntax.teal` | teal/info accent |

Fallback resolver memang ada, tapi tema baru sebaiknya tetap mengisi `syntax` secara explicit supaya hasil Neovim tidak bergantung pada fallback yang terlalu generik.

## UI Colors

| Field | Makna |
|---|---|
| `ui.bg_statusline` | Background statusline/lualine section C |

Default yang sehat untuk dark theme:

```text
ui.bg_statusline = layer.surface_raised
```

Pastikan `text.secondary` atau `text.primary` tetap kontras saat dipakai di atas `ui.bg_statusline`.

## Neovim / Tokyonight Mapping

Template Neovim menghasilkan `colors.lua` dengan mapping compatibility untuk Tokyonight.

Mapping penting:

| Tokyonight field | Theme engine source | Alasan |
|---|---|---|
| `bg` | `layer.base` | editor background |
| `bg_dark` | `layer.mantle` | popup/sidebar darker bg |
| `bg_dark1` | `layer.crust` | darkest bg / fg over accent |
| `bg_highlight` | `layer.surface` | subtle highlight bg |
| `fg` | `text.primary` | main foreground |
| `fg_dark` | `text.secondary` | secondary foreground |
| `comment` | `text.muted` | comments are muted text |
| `fg_gutter` | `layer.surface_overlay` | Tokyonight uses this as subtle UI and lualine B background |
| `dark3` | `layer.surface_overlay` | should match gutter-like low-emphasis UI |
| `dark5` | `border.default` | darker structural bg |
| `bg_statusline` | `ui.bg_statusline` | lualine/statusline bg |

Rule penting:

- Jangan map `fg_gutter` ke `text.muted`.
- Jangan map `dark5` ke text color.
- `dark3` harus lebih terang dari `dark5` untuk dark theme.

## Lualine Expectations

Tokyonight lualine theme memakai pola berikut:

```lua
normal.a = { bg = colors.blue,    fg = colors.black }
insert.a = { bg = colors.green,   fg = colors.black }
command.a = { bg = colors.yellow, fg = colors.black }
visual.a = { bg = colors.magenta, fg = colors.black }
replace.a = { bg = colors.red,    fg = colors.black }

normal.b = { bg = colors.fg_gutter, fg = colors.blue }
```

Implikasi untuk theme baru:

- `colors.black` di adapter Neovim berasal dari `bg_dark1` / `layer.crust`.
- Mode colors (`blue`, `green`, `yellow`, `magenta`, `red`, `green1`) harus readable dengan `layer.crust` sebagai foreground section A.
- Mode colors juga harus readable sebagai foreground di atas `layer.surface_overlay` untuk section B.
- Kalau section B tabrakan, cek dulu `layer.surface_overlay` dan ANSI mode colors, bukan langsung custom lualine.

## Font System

Theme-Engine menggunakan sistem font semantik dengan 3 kategori utama, terpusat di `theme.json` pada bagian `theme.fonts`:

1.  **`system`**: Font proporsional untuk UI OS (GTK, Firefox, Nemo).
2.  **`widget`**: Font untuk elemen desktop seperti Waybar, Eww, Rofi. (Fallback otomatis ke `system` jika kosong).
3.  **`terminal`**: Font monospace wajib untuk Kitty, Foot, Alacritty, dan Neovim. (Ukuran `size` di sini akan menjadi SSOT bagi semua emulator terminal).

*Override font-size di level masing-masing tool sudah dihapus untuk menjaga kebersihan arsitektur.*
Baca `FONT_SYSTEM.md` untuk detail lebih lanjut.

## GTK Integration

Theme engine menghasilkan CSS untuk GTK melalui pipeline:

```text
source.tmpl (Go template)
    → render
    → _source.scss (SCSS variables, generated di output/domain/gtk/scss/source/)
    → di-@import oleh base.scss (static, di assets/templates/domain/gtk/)
    → sassc compile
    → source.css (CSS @define-color rules)

_sfunction.scss (static, di assets/templates/domain/gtk/)
    → di-@import oleh base.scss
    → menyediakan utility functions: alpha(), shade(), tint(), blend()

rofi-base.scss (static)
    → sassc compile
    → source.rasi (Rasi config untuk Rofi)
```

### SCSS Utility Functions

| Function | Deskripsi | Contoh |
|---|---|---|
| `alpha($color, $a)` | rgba dari hex + opacity | `alpha($accent-primary, 0.2)` → `rgba(139, 164, 176, 0.2)` |
| `shade($color, $amount)` | Darken | `shade($accent-primary, 8%)` → `#73919f` |
| `tint($color, $amount)` | Lighten | `tint($accent-primary, 15%)` → `#b9c8cf` |
| `blend($c1, $c2, $p)` | Mix/warna campuran | `blend($accent-primary, $text-primary, 50%)` → `#a8b7bb` |

File SCSS static (`base.scss`, `rofi-base.scss`, `_function.scss`) terletak di `assets/templates/domain/gtk/` — bukan output. Hanya `_source.scss` yang di-generate di `output/domain/gtk/scss/source/`.

### CSS Output

`source.css` berisi 78+ `@define-color` rules, mencakup:

- **Accent** (9): direct + derived via alpha, shade, tint, blend
- **Background** (14): direct palette + derived (surface-active, sidebar-hover, menu-selected, etc.)
- **Foreground** (13): direct palette + derived (bright, alt, muted, invert, osd)
- **Selection** (2): derived dari accent
- **Status** (9): direct + derived (hover, strong)
- **Border** (3): direct + derived (alpha)
- **Shadow** (4): rgba(0,0,0, X%)
- **Legacy aliases** (17): untuk backward compatibility dengan waybar dan config lama (`text-primary`, `layer-crust`, `accent-primary`, `status-*`, dll.)

Beberapa CSS variable dihasilkan dari SCSS function:

```css
@define-color accent-hover      shade(accent-primary, 8%);
@define-color bg-surface-active tint(layer-surface-overlay, 10%);
@define-color fg-muted          alpha(text-primary, 0.5);
@define-color border            alpha(text-primary, 0.12);
@define-color bg-sidebar-hover  tint(layer-surface-overlay, 8%);
@define-color accent-focus      alpha(text-secondary, 0.3);
```

Variable yang tidak bisa di-derive (butuh nilai explicit di palette):
`fg-dim`, `fg-dim-muted`, `fg-disabled`, `bg-conflict`, `bg-disk-usage`, `border-medium`, `text-visited`

## Checklist Tema Baru

Sebelum tema dianggap siap:

- `layer.*` punya urutan terang-gelap yang konsisten.
- `text.primary`, `text.secondary`, `text.muted`, `text.link`, `text.visited` terisi.
- `text.muted` tidak dipakai sebagai structural background.
- `border.default`, `border.active`, `border.medium` terisi.
- `accent.on_accent` kontras terhadap `accent.primary`.
- ANSI `colors[0..15]` lengkap dan mengikuti urutan standar.
- Mode colors punya kontras cukup terhadap `layer.surface_overlay`.
- `status.*` terisi.
- `fg.dim`, `fg.dim_muted`, `fg.disabled` terisi (bisa fallback).
- `bg.conflict`, `bg.disk_usage` terisi (bisa fallback).
- `syntax` block diisi explicit.
- `ui.bg_statusline` diisi explicit.
- Neovim output dicek khusus untuk `fg_gutter`, `dark3`, `dark5`, dan `bg_statusline`.
- GTK output dicek: `source.css` dan `source.rasi` ter-generate.

## Contoh Dark Theme Yang Sehat

```json
"layer": {
  "base": "#181616",
  "mantle": "#12120f",
  "crust": "#0d0c0c",
  "surface": "#1D1C19",
  "surface_raised": "#282727",
  "surface_overlay": "#393836"
},
"border": {
  "default": "#282727",
  "active": "#8ba4b0",
  "medium": "#4a4a49"
},
"text": {
  "primary": "#c5c9c5",
  "secondary": "#C8C093",
  "muted": "#a6a69c",
  "link": "#8ba4b0",
  "visited": "#7b6f8c"
},
"fg": {
  "dim": "#63635e",
  "dim_muted": "#7a7a74",
  "disabled": "#5a5f63"
},
"bg": {
  "conflict": "#635127",
  "disk_usage": "#50504a"
},
"ui": {
  "bg_statusline": "#282727"
}
```

Hasil Neovim yang diharapkan:

```lua
fg_gutter = "#393836"
dark3     = "#393836"
dark5     = "#282727"
comment   = "#a6a69c"
```

Ini menjaga lualine section B tetap readable tanpa harus membuat custom lualine theme.
