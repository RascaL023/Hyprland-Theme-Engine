# Arsitektur Sistem Font pada Theme Engine

## Konsep Utama: 3 Kategori Semantis
Theme-Engine memakai **arsitektur font berbasis semantik (fungsi)**. Jadi, alih-alih cuma pakai istilah umum kayak font "utama" (*primary*) atau "sekunder" (*secondary*), font di sini dikelompokkan berdasarkan **tujuan dan fungsinya** di dalam sistem.

### 1. `system`
*   **Tujuan:** Untuk UI di tingkat OS (GTK, Firefox, Nemo, Desktop Environment).
*   **Target:** `gsettings`, `dconf`, `~/.config/gtk-3.0/settings.ini`.
*   **Karakteristik:** Biasanya pakai font *sans-serif* proporsional (contohnya: *Inter*, *Roboto*, *Cantarell*).

### 2. `widget`
*   **Tujuan:** Untuk widget desktop (seperti Waybar, Eww, Rofi/Wofi).
*   **Target:** Di-render lewat variabel SCSS GTK (`$font-widget`).
*   **Karakteristik:** Bisa pakai font proporsional atau *monospace*, dan sering kali butuh dukungan ikon *Nerd Font* yang bagus.
*   **Fallback (Cadangan):** Kalau dikosongkan, otomatis bakal ngikutin font `system`.

### 3. `terminal`
*   **Tujuan:** Untuk emulator terminal dan *code editor* (seperti Kitty, Foot, Alacritty, Neovim).
*   **Target:** Langsung dibaca/dipakai oleh konfigurasi masing-masing aplikasi.
*   **Karakteristik:** **WAJIB** pakai font *monospace* (contohnya: *JetBrainsMono Nerd Font*, *Maple Mono Nerd Font*).
*   **Fallback (Cadangan):** Kalau dikosongkan, default-nya bakal ngikut ke `system` (tapi sangat disarankan buat tetap ngatur font *monospace* sendiri).

---

## Konfigurasi (`theme.json`)
Semua pengaturan font terpusat (terpusat secara ketat) di dalam file JSON milik tema.

```json
{
  "theme": {
    "name": "my-theme",
    "fonts": {
      "system": {
        "family": "Inter",
        "size": 11.0
      },
      "widget": {
        "family": "Inter",
        "size": 11.0
      },
      "terminal": {
        "family": "Maple Mono Nerd Font",
        "size": 11.0
      }
    }
  }
}
```

> **Catatan:** Pengaturan ukuran font khusus per aplikasi (contohnya: `tools.kitty.fontSize`) sudah **DIHAPUS** demi menjaga prinsip *Single Source of Truth* (satu sumber acuan tunggal). Semua emulator terminal sekarang bakal tunduk pada nilai `fonts.terminal.size`.

---

## Integrasi dengan GTK
Di dalam file `assets/templates/domain/gtk/source.tmpl`, font akan diekspor sebagai variabel SCSS:
```scss
$font-system   : "Inter";
$font-widget   : "Inter";
$font-terminal : "Maple Mono Nerd Font";
```
Pas kamu lagi bikin widget (kayak Waybar), tinggal *import* file `source.scss` lalu pakai `font-family: $font-widget;`. Cara ini mencegah terjadinya bentrok font akibat *hardcoding* (menulis nama font secara manual di banyak tempat).

---

## Menerapkan Font Sistem
Karena Theme-Engine ini sifatnya adalah sebuah *generator* (pembuat file konfigurasi), cara paling ampuh buat menerapkan font `system` secara langsung (*live*) ke desktop GTK kamu adalah lewat *shell script* (misalnya `runner.sh`) yang membaca hasil *output*-nya lalu menjalankan perintah `gsettings`.

Contoh perintah Bash untuk menerapkan font secara *live*:
```bash
gsettings set org.gnome.desktop.interface font-name "Inter 11"
gsettings set org.gnome.desktop.interface document-font-name "Inter 11"
```
