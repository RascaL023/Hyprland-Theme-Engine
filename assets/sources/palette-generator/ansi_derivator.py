#!/usr/bin/env python3
"""
Adaptive ANSI Terminal Color Generator
Generates functional 16-color ANSI palettes that adapt to wallpaper complexity.
- Monochrome wallpapers: Desaturated, grayish functional colors.
- Colorful wallpapers: Vibrant, hue-tinted functional colors.
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
    sats = []
    best_sat = -1
    base_hue = 0
    
    for i in range(1, 8):
        c = colors.get(f'color{i}')
        if not c: continue
        h, s, l = rgb_to_hsl(hex_to_rgb(c))
        sats.append(s)
        # Find the most saturated color (even if it's low) to act as base hue
        if s > best_sat:
            best_sat = s
            base_hue = h
            
    avg_sat = sum(sats) / len(sats) if sats else 0
    # Determine if monochrome (avg saturation < 15%)
    is_mono = avg_sat < 15
    return base_hue, avg_sat, best_sat, is_mono

def shift_hue_towards(target_hue, base_hue, blend_factor):
    diff = (base_hue - target_hue + 180) % 360 - 180
    new_hue = (target_hue + diff * blend_factor) % 360
    return new_hue

def generate_ansi_palette(pywal_data, is_dark=True):
    base_hue, avg_sat, max_sat, is_mono = analyze_vibe(pywal_data)
    
    STANDARD_HUES = {
        'red': 0,
        'green': 120,
        'yellow': 40,
        'blue': 220,
        'magenta': 300,
        'cyan': 180
    }
    
    if is_mono:
        # Monochrome logic: Desaturated colors, low tint factor
        target_sat_normal = 25  # Highly desaturated (grayish)
        target_sat_bright = 35
        tint_factor = 0.1       # Barely tint to base hue to keep standard hues distinct
    else:
        # Colorful logic: Vibrant colors, higher tint factor
        target_sat_normal = min(80, max(40, avg_sat))
        target_sat_bright = min(100, target_sat_normal * 1.2)
        tint_factor = 0.25      # Blend towards wallpaper vibe
        
    if is_dark:
        l_bg = 10; l_bg_br = 30
        l_fg = 85; l_fg_br = 95
        l_norm = 65; l_br = 75
    else:
        l_bg = 95; l_bg_br = 80
        l_fg = 20; l_fg_br = 10
        l_norm = 40; l_br = 30

    ansi = [None] * 16
    
    # Black
    ansi[0] = rgb_to_hex(hsl_to_rgb((base_hue, target_sat_normal * 0.2, l_bg)))
    ansi[8] = rgb_to_hex(hsl_to_rgb((base_hue, target_sat_normal * 0.3, l_bg_br)))
    
    color_map = [
        (1, 9, 'red'),
        (2, 10, 'green'),
        (3, 11, 'yellow'),
        (4, 12, 'blue'),
        (5, 13, 'magenta'),
        (6, 14, 'cyan')
    ]
    
    for norm_idx, br_idx, name in color_map:
        hue = shift_hue_towards(STANDARD_HUES[name], base_hue, tint_factor)
        
        ansi[norm_idx] = rgb_to_hex(hsl_to_rgb((hue, target_sat_normal, l_norm)))
        ansi[br_idx] = rgb_to_hex(hsl_to_rgb((hue, target_sat_bright, l_br)))
        
    # White
    ansi[7] = rgb_to_hex(hsl_to_rgb((base_hue, target_sat_normal * 0.1, l_fg)))
    ansi[15] = rgb_to_hex(hsl_to_rgb((base_hue, target_sat_normal * 0.15, l_fg_br)))
    
    return ansi, is_mono

def main():
    parser = argparse.ArgumentParser(description="Adaptive ANSI Terminal Color Generator")
    parser.add_argument('input', help="Path to Pywal colors.json")
    parser.add_argument('-o', '--output', help="Output JSON path", default="new_ansi_colors.json")
    parser.add_argument('--light', action='store_true', help="Generate for light background")
    
    args = parser.parse_args()
    
    with open(args.input) as f:
        pywal = json.load(f)
        
    colors, is_mono = generate_ansi_palette(pywal, not args.light)
    mode = "Monochrome" if is_mono else "Colorful"
    
    output = {
        "description": f"Adaptive ANSI Palette ({mode} Mode)",
        "colors": colors
    }
    
    with open(args.output, 'w') as f:
        json.dump(output, f, indent=2, ensure_ascii=False)
        
    print(f"[*] Generated {mode} ANSI palette saved to {args.output}")

if __name__ == "__main__":
    main()
