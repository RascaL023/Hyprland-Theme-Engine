# Theme Color Rules

Dokumen ini adalah acuan saat membuat tema baru. Tujuannya supaya setiap warna punya makna semantic yang konsisten, bukan sekadar warna yang terlihat bagus secara terpisah.

Theme engine ini memakai palette semantic sendiri, lalu beberapa output seperti Neovim/tokyonight akan memetakan warna tersebut ke nama variabel milik tool terkait. Karena itu, urutan terang-gelap dan fungsi tiap field harus dijaga.

## Prinsip Utama

- Jangan isi field berdasarkan nama warna saja. Isi berdasarkan fungsi visualnya.
- `text.*` adalah warna foreground/text, bukan warna background UI.
- `layer.*` adalah warna background/surface, bukan warna text.
- `border.*` adalah warna pemisah/outline structural.
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
| `layer.surface_raised` | Surface lebih terlihat | statusline bg, raised component |
| `layer.surface_overlay` | Structural muted UI | gutter-like color, inactive/subtle UI |

Rule penting:

```text
layer.crust < layer.mantle < layer.base < layer.surface < layer.surface_raised < layer.surface_overlay
```

Untuk dark theme, tanda `<` berarti lebih gelap secara luminance. Untuk light theme, prinsipnya dibalik secara visual: base harus tetap nyaman sebagai background utama, dan surface/raised/overlay harus memberi separasi yang cukup.

## Text Colors

Text colors dipakai untuk foreground. Jangan pakai warna ini sebagai background structural kecuali ada alasan khusus.

| Field | Makna | Rule |
|---|---|---|
| `text.primary` | Teks utama | kontras tinggi terhadap `layer.base` |
| `text.secondary` | Teks sekunder | masih readable, lebih kalem dari primary |
| `text.muted` | Komentar / disabled-ish text | readable tapi low-emphasis |
| `text.link` | Link / navigasi / reference | biasanya dekat dengan accent/info |

Untuk dark theme:

```text
text.primary > text.secondary > text.muted > layer.surface_overlay
```

Catatan penting:

- `text.muted` boleh terang selama masih cocok untuk komentar.
- `text.muted` jangan dipakai untuk `fg_gutter` Tokyonight karena `fg_gutter` sering dipakai sebagai background lualine section B.

## Border Colors

| Field | Makna | Rule |
|---|---|---|
| `border.default` | Border non-active | lebih subtle dari active, masih terlihat di atas surface |
| `border.active` | Border active/focused | biasanya sama atau dekat dengan accent utama |

Untuk dark theme, `border.default` idealnya berada di antara `layer.surface_raised` dan `layer.surface_overlay`, atau setidaknya tidak lebih terang dari text muted.

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

## Checklist Tema Baru

Sebelum tema dianggap siap:

- `layer.*` punya urutan terang-gelap yang konsisten.
- `text.primary`, `text.secondary`, `text.muted` readable di atas `layer.base`.
- `text.muted` tidak dipakai sebagai structural background.
- `border.default` terlihat di atas `layer.surface` / `layer.surface_raised`.
- `accent.on_accent` kontras terhadap `accent.primary`.
- ANSI `colors[0..15]` lengkap dan mengikuti urutan standar.
- Mode colors punya kontras cukup terhadap `layer.surface_overlay`.
- `syntax` block diisi explicit.
- `ui.bg_statusline` diisi explicit.
- Neovim output dicek khusus untuk `fg_gutter`, `dark3`, `dark5`, dan `bg_statusline`.

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
  "active": "#8ba4b0"
},
"text": {
  "primary": "#c5c9c5",
  "secondary": "#C8C093",
  "muted": "#a6a69c",
  "link": "#8ba4b0"
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
