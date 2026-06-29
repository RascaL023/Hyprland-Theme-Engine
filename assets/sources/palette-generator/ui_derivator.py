#!/usr/bin/env python3
"""
Adaptive Semantic UI Color Generator
Analyzes pywal palette to generate Catppuccin-style semantic UI colors.
- Monochrome wallpapers: Grayscale backgrounds, boosts the most saturated color as Accent.
- Colorful wallpapers: Base-hue tinted backgrounds, natural accents.
"""

import json
import colorsys
import argparse

def hex_to_rgb(hex_color):
    hex_color = hex_color.lstrip('#')
    return tuple(int(hex_color[i:i+2], 16) for i in (0, 2, 4))

def rgb_to_hex(rgb):
    return '#{:02x}{:02x}{:02x}'.format(
        int(max(0, min(255, rgb[0]))),
        int(max(0, min(255, rgb[1]))),
        int(max(0, min(255, rgb[2])))
    )

def rgb_to_hsl(rgb):
    r, g, b = [x / 255.0 for x in rgb]
    h, l, s = colorsys.rgb_to_hls(r, g, b)
    return (h * 360, s * 100, l * 100)

def hsl_to_rgb(hsl):
    h, s, l = hsl[0] / 360.0, hsl[1] / 100.0, hsl[2] / 100.0
    r, g, b = colorsys.hls_to_rgb(h, l, s)
    return (r * 255, g * 255, b * 255)

def analyze_vibe(pywal_data):
    colors = pywal_data['colors']
    lightness_vals = []
    sats = []
    best_sat = -1
    base_hue = 0
    
    for i in range(16):
        c = colors.get(f'color{i}')
        if not c: continue
        h, s, l = rgb_to_hsl(hex_to_rgb(c))
        lightness_vals.append(l)
        if 1 < i < 8: # Check normal colors for accent
            sats.append(s)
            if s > best_sat:
                best_sat = s
                base_hue = h
                
    min_l, max_l = min(lightness_vals), max(lightness_vals)
    avg_sat = sum(sats) / len(sats) if sats else 0
    is_mono = avg_sat < 15
    
    return base_hue, avg_sat, best_sat, min_l, max_l, is_mono

def generate_ui_palette(pywal_data, is_dark=True):
    base_hue, avg_sat, max_sat, min_l, max_l, is_mono = analyze_vibe(pywal_data)
    l_range = max_l - min_l
    
    if is_mono:
        bg_sat = 0          # Pure grayscale backgrounds
        accent_sat = 100    # Force max saturation for the accent color
    else:
        bg_sat = avg_sat * 0.2  # Tint backgrounds slightly
        accent_sat = min(100, max_sat * 1.5)
        
    if is_dark:
        MAPPING = {
            'text_primary':   {'pos': 0.95, 'sat': bg_sat * 0.5, 'hue_s': 0},
            'text_secondary': {'pos': 0.85, 'sat': bg_sat * 0.8, 'hue_s': 0},
            'text_muted':     {'pos': 0.65, 'sat': bg_sat,       'hue_s': 0},
            
            'layer_surface_overlay': {'pos': 0.40, 'sat': bg_sat, 'hue_s': 0},
            'layer_surface_raised':  {'pos': 0.25, 'sat': bg_sat, 'hue_s': 0},
            'layer_surface':         {'pos': 0.18, 'sat': bg_sat, 'hue_s': 0},
            'layer_base':            {'pos': 0.10, 'sat': bg_sat, 'hue_s': 0},
            'layer_mantle':          {'pos': 0.05, 'sat': bg_sat, 'hue_s': 0},
            'layer_crust':           {'pos': 0.01, 'sat': bg_sat, 'hue_s': 0},
            
            'border_active':  {'pos': 0.75, 'sat': accent_sat, 'hue_s': 0},
            'border_default': {'pos': 0.30, 'sat': bg_sat,     'hue_s': 0},
            
            'accent_primary':   {'pos': 0.75, 'sat': accent_sat, 'hue_s': 0},
            'accent_secondary': {'pos': 0.65, 'sat': accent_sat * 0.8, 'hue_s': 15},
            'accent_on':        {'pos': 0.05, 'sat': 0, 'hue_s': 0},
        }
    else:
        MAPPING = {
            'text_primary':   {'pos': 0.10, 'sat': bg_sat * 0.5, 'hue_s': 0},
            'text_secondary': {'pos': 0.25, 'sat': bg_sat * 0.8, 'hue_s': 0},
            'text_muted':     {'pos': 0.45, 'sat': bg_sat,       'hue_s': 0},
            
            'layer_surface_overlay': {'pos': 0.70, 'sat': bg_sat, 'hue_s': 0},
            'layer_surface_raised':  {'pos': 0.85, 'sat': bg_sat, 'hue_s': 0},
            'layer_surface':         {'pos': 0.90, 'sat': bg_sat, 'hue_s': 0},
            'layer_base':            {'pos': 0.95, 'sat': bg_sat, 'hue_s': 0},
            'layer_mantle':          {'pos': 0.97, 'sat': bg_sat, 'hue_s': 0},
            'layer_crust':           {'pos': 0.99, 'sat': bg_sat, 'hue_s': 0},
            
            'border_active':  {'pos': 0.40, 'sat': accent_sat, 'hue_s': 0},
            'border_default': {'pos': 0.80, 'sat': bg_sat,     'hue_s': 0},
            
            'accent_primary':   {'pos': 0.40, 'sat': accent_sat, 'hue_s': 0},
            'accent_secondary': {'pos': 0.50, 'sat': accent_sat * 0.8, 'hue_s': 15},
            'accent_on':        {'pos': 0.95, 'sat': 0, 'hue_s': 0},
        }

    ui = {}
    for name, conf in MAPPING.items():
        l = min_l + (conf['pos'] * l_range)
        if l_range < 40: l += (50 - l_range) * (conf['pos'] - 0.5)
        l = max(3, min(97, l))
        
        h = (base_hue + conf['hue_s']) % 360
        ui[name] = rgb_to_hex(hsl_to_rgb((h, conf['sat'], l)))
        
    return ui, is_mono

def format_extra_block(ui):
    return {
        "accent": {
            "primary": ui["accent_primary"],
            "secondary": ui["accent_secondary"],
            "on_accent": ui["accent_on"]
        },
        "text": {
            "primary": ui["text_primary"],
            "secondary": ui["text_secondary"],
            "muted": ui["text_muted"],
            "link": ui["accent_primary"]
        },
        "layer": {
            "base": ui["layer_base"],
            "mantle": ui["layer_mantle"],
            "crust": ui["layer_crust"],
            "surface": ui["layer_surface"],
            "surface_raised": ui["layer_surface_raised"],
            "surface_overlay": ui["layer_surface_overlay"]
        },
        "border": {
            "default": ui["border_default"],
            "active": ui["border_active"]
        }
    }

def main():
    parser = argparse.ArgumentParser(description="Adaptive Semantic UI Color Generator")
    parser.add_argument('input', help="Path to Pywal colors.json")
    parser.add_argument('-o', '--output', help="Output JSON path", default="new_ui_colors.json")
    parser.add_argument('--light', action='store_true', help="Generate for light theme")
    
    args = parser.parse_args()
    
    with open(args.input) as f:
        pywal = json.load(f)
        
    ui_flat, is_mono = generate_ui_palette(pywal, not args.light)
    extra_block = format_extra_block(ui_flat)
    mode = "Monochrome" if is_mono else "Colorful"
    
    out_data = {
        "description": f"Adaptive Semantic UI Colors ({mode} Mode)",
        "extra": extra_block
    }
    
    with open(args.output, 'w') as f:
        json.dump(out_data, f, indent=2, ensure_ascii=False)
        
    print(f"[*] Generated {mode} UI extra block saved to {args.output}")

if __name__ == "__main__":
    main()
