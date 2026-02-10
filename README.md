# Go Theme Renderer

## Description
Render JSON-based theme definition into multiple tools configuration using SSOT (Single Source of Truth) approach.

### Features:
- Parse and render theme JSON into internal structured data
- Render theme into tools-specific config formats
- Modular renderer design (easy to add new tools)
- Custom input and output path
- Deterministic output (same input → same config)
- Minimal runtime dependency
- CLI log and notification

### Supported Tools:
- Foot
- Cava
- Kitty
- GTK-CSS
    Member of Gtk-Css config: Waybar, Eww, Wlogout

### Planned Tools:
- Hyprland
- Ncmcpp
- Yazi

### Future Plans:
- Hot-reload support
- Theme validation & schema checking
- CLI flags for selective rendering
- Documentation & examples
- SSOT color

#### Pending        : Resolve gtk variants
#### Last Progress  : Rendered scss
