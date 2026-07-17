# Theme Engine - Implementation Summary

## Overview
A dynamic theme engine that generates configuration files for various terminal applications based on a unified palette and theme definition system.

## Key Components

### 1. Core Architecture
- **resolvertion**: Dynamic variable resolution using `$pl.extra.*` syntax
- **Path-based Loading**: Generates different outputs based on the tool type
- **Template Rendering**: Uses Go templates for dynamic file generation
- **Caching**: Smart caching to avoid redundant file operations

### 2. Palette Structure
**Old Schema** (no longer used):
- Surface[3], Overlay[3], Accent(2), Text(3), Status(3)
- Limited semantic color naming

**New Schema (Material/Catppuccin-inspired)**:
- **accent**: { primary, secondary, on_accent }
- **text**: { primary, secondary, muted, link }
- **layer**: { base, mantle, crust, surface, surface_raised, surface_overlay }
- **border**: { default, active }
- **status**: { success, warning, error, critical, info }

### 3. Supported Tools

#### Domain Tools (Full Processing)
- **GTK-CSS**: Generates GTK CSS and Rasi theme files
  - Input: `theme.json` + palette
  - Output: CSS for GTK applications
  - Location: `/output/domain/gtk/`

#### Application Tools (Direct Variable Replacement)
- **Kitty**: Direct color injection into configuration
- **Foot**: Terminal emulator configuration
- **Hypr**: Window manager configuration
- **Yazi**: File manager theme
- **Cava**: Visualization tool configuration
- **Nvim**: Neovim integration via `colors.lua`
- **Starship**: Cross-shell prompt with theme-aware colors

#### Neovim Integration (Lua-based)
- Generates `colors.lua` with theme colors
- Used by Tokyo Night theme plugin
- Allows dynamic color injection without modifying nvim configuration

#### Starship Integration
- Static template generates theme-adaptive TOML prompt config
- Uses palette accent/status colors instead of hardcoded ANSI names
- Prompt layout remains user-defined; only colors adapt to active theme

### 4. Design Patterns

#### Template Structure
- Tool-specific templates
- Consistent variable naming pattern using {{ .Palette.Field }}
- Automatic variable resolution

#### Path Configuration
`/config/path.txt` maps tool names to templates:
```
tool|template_path|output_dir
```

#### Progressive Enhancement
**Old Targets (Removed in Refactor)**:
- Waybar - Now uses GTK CSS import

**Maintained Compatibility**:
- Existing themes (nocturne) work with new schema via template updates

### 5. Best Practices Implemented

#### 1. Semantic Naming
**Before:** `{primarysurface}`, `{secondarysurface}`, etc.
**After:** `{layer.surface}`, `{layer.surface_raised}` etc.

#### 2. Type Safety
- Strong typing in Go structs
- Clear interface definitions for processors

#### 3. Atomic Operations
- Use temp files for writing, then rename
- Skip unchanged files (reduce I/O)

#### 4. Documentation
- Comprehensive README with examples
- FLOW.md explains architecture
- PROGRESS.md tracks development

### 6. Color Management

#### Dynamic Injection
**Old Way (static):**
```ini
warning = #ff80ff
critical = #cc33ff
```

**New Way (context-aware):**
```ini
warning = $pl.extra.status.warning
critical = $pl.extra.status.critical
```

#### Semantic Tokens
- Colors are no longer tied to specific hex codes
- Maintain semantic meaning across different themes
- Easy to adapt to different color schemes

### 7. Usage Examples

#### Basic Usage
```bash
# Render all targets
./runner.sh dev

# Render specific tool
./runner.sh dev nvim

# Build binary for production
./runner.sh build

# Run the tool
./runner.sh run
```

#### Theme Selection
```bash
# Modify state.json to switch themes
# Example: change from nocturne to ghostly
{
  "theme": {
    "name": "ghostly"
  }
}
```

### 8. Key Features

1. **Multiple Tools**: Supports 10+ different applications
2. **Dynamic Palette**: Palette loaded dynamically from JSON files
3. **Smart Caching**: Only processes files that have changed
4. **Atomic Writes**: Prevents corruption from partial writes
5. **Flexible File Structure**: Organized with clear separation of concerns
6. **Extensible Design**: Easy to add new processors and templates
7. **Cross-Platform**: Works on Linux, macOS, Windows
8. **Comprehensive Documentation**: Clear guidance for users

### 9. Testing and Validation

- Unit tests for core functionality
- Test suite validates palette processing
- Integration tests for end-to-end workflows
- Smoke tests to ensure theme compatibility

### 10. Development Environment

#### Prerequisites
- Go 1.22 or later
- Node.js (for optional development tools)
- Nix (optional for isolated development)

#### Development Commands
```bash
# Run tests
./runner.sh test

# Render themes in development mode
./runner.sh dev

# Build binary
./runner.sh build

# Clean up
./runner.sh clean
```

## Recent Changes

### Schema Refactor
Completely reworked the palette schema to follow design system best practices:
- Replaced numeric-based surfaces with semantic `layer` definitions
- Added semantic border and status color tokens
- Maintained backward compatibility through template refactoring

### Starship Integration
Added a new template-based starship integration:
- Generates `starship.toml` with theme-adaptive colors
- Uses `{{ .Palette.AccentPrimary }}`, `{{ .Palette.StatusWarning }}`, etc.
- Symlink: `ln -sf …/output/tools/starship/starship.toml ~/.config/starship.toml`

### New Themes (Kanagawa, Dracula, Sakura)
Added several hand-crafted themes with correct ANSI luminance ordering:
- **Kanagawa Wave/Dragon**: Fixed reversed brightness pairs (dark cyan, white)
- **Dracula**: Official spec with Dracula Classic + Alucard light
- **Sakura**: Original cherry blossom themed palette

### Starship Integration
New static processor generating theme-adaptive `starship.toml`:
- Uses palette `extra.*` colors (accent, status) instead of hardcoded ANSI names
- Prompt layout stays user-defined; colors adapt per theme

### Neovim Integration
Added a new Lua-based integration for the Tokyo Night theme plugin:
- Generates `colors.lua` from palette
- Maintains plugin compatibility
- Allows dynamic color injection

### Tool Removal
Removed the Waybar template to reduce redundancy:
- Now imports GTK CSS directly
- Centralizes styling in GTK theme
- Reduces maintenance overhead

## Migration Guide

### From Old Schema to New

#### Palette Changes
**Old Keys**: `{primarysurface}`, `{secondarysurface}`, `{primaryaccent}`, `{secondaryaccent}`, etc.

**New Keys**: `{layer.surface}`, `{layer.surface_raised}`, `{accent.primary}`, `{accent.secondary}`, etc.

```lua
-- Old template snippet
{{- index .Palette "primarysurface" -}}

-- New template snippet  
{{- .Palette.LayerSurface -}}
```

#### Tool Configuration
No tool-specific changes needed! The processor automatically handles the semantic token mapping.

## Extensibility

### Adding New Tools
1. Create processor in `internal/core/tools/`.new_tool/
2. Add template in `assets/templates/tools/new_tool/`
3. Register in `config/path.txt`

### Modifying Palette Schema
1. Update `Palette` structs in `palette/raw.go`
2. Update `ResolvedPalette` in `palette/palette.go`
3. Update resolver in `palette/resolver.go`
4. Update flatten mapping in `theme/flatter.go`
5. Update all templates referencing new fields

## Best Practices for Contributors

### 1. Theme Selection
Choose a theme that matches your needs:
- **Nocturne**: Dark mode with purple/magenta accents
- **Ghostly**: Modern theme with extensive light/dark support
- **Kanagawa Wave**: Classic Kanagawa wave dark + lotus light
- **Kanagawa Dragon**: Muted Kanagawa dragon variant
- **Dracula**: Official Dracula Classic + Alucard light
- **Sakura**: Cherry blossom inspired, warm rose tones

### 2. File Organization
- Maintain existing project structure
- Follow Go package conventions
- Keep templates consistent with existing patterns

### 3. Testing
- Always run tests before committing
- Verify theme output manually
- Test both dark and light modes

### 4. Documentation
- Update README or relevant documentation
- Document new features or breaking changes
- Keep examples up to date

## Conclusion

The theme engine provides a consistent, maintainable way to manage multiple application themes. By using a unified palette schema and template-based generation, it ensures consistency across all supported tools while maintaining flexibility for customization.

The refactor represents a significant improvement in code organization, semantic naming, and extensibility. The engine is now well-positioned for future development and maintenance.
