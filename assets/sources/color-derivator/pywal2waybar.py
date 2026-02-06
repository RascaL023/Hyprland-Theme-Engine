#!/usr/bin/env python3
"""
Pywal to Waybar - FULLY ADAPTIVE VERSION
Extracts brightness range and color tones from entire pywal palette
"""

import json
import colorsys
from pathlib import Path
from typing import Dict, Tuple, List


class ColorUtils:
    """Color conversion utilities"""
    
    @staticmethod
    def hex_to_rgb(hex_color: str) -> Tuple[int, int, int]:
        hex_color = hex_color.lstrip('#')
        return tuple(int(hex_color[i:i+2], 16) for i in (0, 2, 4))
    
    @staticmethod
    def rgb_to_hex(rgb: Tuple[float, float, float]) -> str:
        return '#{:02x}{:02x}{:02x}'.format(
            int(max(0, min(255, rgb[0]))),
            int(max(0, min(255, rgb[1]))),
            int(max(0, min(255, rgb[2])))
        )
    
    @staticmethod
    def rgb_to_hsl(rgb: Tuple[int, int, int]) -> Tuple[float, float, float]:
        r, g, b = [x / 255.0 for x in rgb]
        h, l, s = colorsys.rgb_to_hls(r, g, b)
        return (h * 360, s * 100, l * 100)
    
    @staticmethod
    def hsl_to_rgb(hsl: Tuple[float, float, float]) -> Tuple[float, float, float]:
        h, s, l = hsl[0] / 360.0, hsl[1] / 100.0, hsl[2] / 100.0
        r, g, b = colorsys.hls_to_rgb(h, l, s)
        return (r * 255, g * 255, b * 255)
    
    @classmethod
    def create_color(cls, hue: float, saturation: float, lightness: float) -> str:
        """Create color from HSL values"""
        rgb = cls.hsl_to_rgb((hue, saturation, lightness))
        return cls.rgb_to_hex(rgb)
    
    @classmethod
    def interpolate_hsl(cls, color1: str, color2: str, factor: float) -> str:
        """Interpolate between two colors in HSL space"""
        h1, s1, l1 = cls.rgb_to_hsl(cls.hex_to_rgb(color1))
        h2, s2, l2 = cls.rgb_to_hsl(cls.hex_to_rgb(color2))
        
        # Handle hue wraparound
        if abs(h1 - h2) > 180:
            if h1 > h2:
                h2 += 360
            else:
                h1 += 360
        
        h = (h1 + (h2 - h1) * factor) % 360
        s = s1 + (s2 - s1) * factor
        l = l1 + (l2 - l1) * factor
        
        return cls.create_color(h, s, l)


class PywalAnalyzer:
    """Analyze pywal palette to extract color characteristics"""
    
    def __init__(self, pywal_data: dict):
        self.colors = pywal_data['colors']
        self.special = pywal_data['special']
        self.utils = ColorUtils()
        
    def get_brightness_range(self) -> Tuple[float, float]:
        """Get min/max lightness from pywal palette"""
        lightness_values = []
        
        # Analyze all colors
        for i in range(16):
            key = f'color{i}'
            if key in self.colors:
                rgb = self.utils.hex_to_rgb(self.colors[key])
                h, s, l = self.utils.rgb_to_hsl(rgb)
                lightness_values.append(l)
        
        # Also check special colors
        for color in [self.special['background'], self.special['foreground']]:
            rgb = self.utils.hex_to_rgb(color)
            h, s, l = self.utils.rgb_to_hsl(rgb)
            lightness_values.append(l)
        
        return min(lightness_values), max(lightness_values)
    
    def get_accent_color(self) -> str:
        """Find best accent color (most saturated mid-tone)"""
        candidates = []
        
        for i in range(2, 7):
            key = f'color{i}'
            if key not in self.colors:
                continue
            
            color = self.colors[key]
            rgb = self.utils.hex_to_rgb(color)
            h, s, l = self.utils.rgb_to_hsl(rgb)
            
            # Score: prefer high saturation + mid lightness
            if 20 < l < 80:
                score = s * (1 - abs(l - 50) / 50)  # Prefer mid-tones
                candidates.append((score, color, h, s, l))
        
        if candidates:
            candidates.sort(reverse=True)
            return candidates[0][1]
        
        return self.colors.get('color2', self.colors['color1'])
    
    def get_base_hue(self) -> float:
        """Get dominant hue from accent color"""
        accent = self.get_accent_color()
        rgb = self.utils.hex_to_rgb(accent)
        h, s, l = self.utils.rgb_to_hsl(rgb)
        return h
    
    def get_average_saturation(self) -> float:
        """Get average saturation of mid-tone colors"""
        saturations = []
        
        for i in range(1, 7):
            key = f'color{i}'
            if key in self.colors:
                rgb = self.utils.hex_to_rgb(self.colors[key])
                h, s, l = self.utils.rgb_to_hsl(rgb)
                if 20 < l < 80:  # Mid-tones only
                    saturations.append(s)
        
        return sum(saturations) / len(saturations) if saturations else 10


class AdaptiveWaybarGenerator:
    """Generate waybar theme that adapts to pywal palette characteristics"""
    
    # Relative lightness positions (0.0 to 1.0 of the range)
    LIGHTNESS_POSITIONS = {
        'text': 0.95,        # Near top of range
        'subtext-1': 0.88,
        'subtext-0': 0.80,
        'lavender': 0.72,
        'overlay2': 0.70,
        'mauve': 0.65,
        'overlay1': 0.60,
        'overlay0': 0.52,
        'surface2': 0.45,
        'surface1': 0.35,
        'surface0': 0.25,
        'base': 0.18,
        'mantle': 0.10,
        'crust': 0.02,       # Near bottom of range
    }
    
    # Saturation multipliers (relative to palette average)
    SATURATION_MULTIPLIERS = {
        'mauve': 2.5,        # More saturated than average
        'lavender': 2.0,
        'text': 0.5,         # Less saturated
        'subtext-1': 0.4,
        'subtext-0': 0.3,
        'overlay2': 0.2,
        'overlay1': 0.15,
        'overlay0': 0.15,
        'surface2': 0.15,
        'surface1': 0.2,
        'surface0': 0.25,
        'base': 0.3,
        'mantle': 0.3,
        'crust': 0.0,        # Pure grayscale
    }
    
    # Hue shift amounts (in degrees)
    HUE_SHIFTS = {
        'mauve': 0,
        'lavender': 5,
        'text': -20,
        'subtext-1': -20,
        'subtext-0': -15,
        'overlay2': -5,
        'overlay1': -10,
        'overlay0': -10,
        'surface2': -15,
        'surface1': -25,
        'surface0': -35,
        'base': -45,
        'mantle': -50,
        'crust': 0,          # No shift for darkest
    }
    
    def __init__(self, pywal_data: dict):
        self.analyzer = PywalAnalyzer(pywal_data)
        self.utils = ColorUtils()
        
        # Extract palette characteristics
        self.min_lightness, self.max_lightness = self.analyzer.get_brightness_range()
        self.lightness_range = self.max_lightness - self.min_lightness
        self.base_hue = self.analyzer.get_base_hue()
        self.avg_saturation = self.analyzer.get_average_saturation()
        
    def calculate_lightness(self, position: float) -> float:
        """
        Calculate actual lightness from position (0.0-1.0) in the palette range
        
        This makes colors adaptive to wallpaper brightness!
        """
        # Map position to actual lightness range from pywal
        lightness = self.min_lightness + (position * self.lightness_range)
        
        # Ensure minimum contrast
        if self.lightness_range < 50:  # Low contrast palette
            # Expand the range slightly
            expansion = (50 - self.lightness_range) * position * 0.5
            lightness += expansion
        
        return max(5, min(95, lightness))  # Clamp to sane values
    
    def calculate_saturation(self, multiplier: float, lightness: float) -> float:
        """
        Calculate saturation based on:
        1. Palette average saturation
        2. Multiplier for this color
        3. Lightness (reduce saturation in very dark/light colors)
        """
        base_sat = self.avg_saturation * multiplier
        
        # Reduce saturation for very dark or very light colors
        if lightness < 15:
            base_sat *= 0.5
        elif lightness > 85:
            base_sat *= 0.7
        
        return max(0, min(100, base_sat))
    
    def calculate_hue(self, shift: float) -> float:
        """Calculate hue with shift"""
        return (self.base_hue + shift) % 360
    
    def generate_palette(self) -> Dict[str, str]:
        """Generate complete adaptive waybar palette"""
        print(f"🎨 Palette Analysis:")
        print(f"  Base Hue: {self.base_hue:.1f}°")
        print(f"  Lightness Range: {self.min_lightness:.1f}% - {self.max_lightness:.1f}% (Δ{self.lightness_range:.1f}%)")
        print(f"  Average Saturation: {self.avg_saturation:.1f}%")
        print(f"  Strategy: Fully adaptive to wallpaper characteristics")
        print()
        
        palette = {}
        
        for color_name in self.LIGHTNESS_POSITIONS.keys():
            # Calculate values
            position = self.LIGHTNESS_POSITIONS[color_name]
            lightness = self.calculate_lightness(position)
            
            sat_mult = self.SATURATION_MULTIPLIERS[color_name]
            saturation = self.calculate_saturation(sat_mult, lightness)
            
            hue_shift = self.HUE_SHIFTS[color_name]
            hue = self.calculate_hue(hue_shift)
            
            # Create color
            palette[color_name] = self.utils.create_color(hue, saturation, lightness)
        
        return palette
    
    def format_css(self, palette: Dict[str, str]) -> str:
        """Format as waybar CSS"""
        lines = ["/* Auto-generated Adaptive Waybar Theme */"]
        
        rgb_colors = {'base', 'mantle', 'crust'}
        
        order = [
            'mauve', 'lavender',
            'text', 'subtext-1', 'subtext-0',
            'overlay2', 'overlay1', 'overlay0',
            'surface2', 'surface1', 'surface0',
            'base', 'mantle', 'crust',
        ]
        
        for name in order:
            color = palette[name]
            if name in rgb_colors:
                rgb = self.utils.hex_to_rgb(color)
                lines.append(f"@define-color {name:12} rgb({rgb[0]}, {rgb[1]}, {rgb[2]});")
            else:
                lines.append(f"@define-color {name:12} {color};")
        
        return '\n'.join(lines)
    
    def print_debug(self, palette: Dict[str, str]):
        """Print detailed analysis"""
        print("\n" + "="*70)
        print("  GENERATED PALETTE ANALYSIS")
        print("="*70)
        print(f"{'Color':<15} {'Hex':<10} {'Hue':<8} {'Sat':<7} {'Light':<7} {'Position'}")
        print("-"*70)
        
        for name in self.LIGHTNESS_POSITIONS.keys():
            hex_color = palette[name]
            rgb = self.utils.hex_to_rgb(hex_color)
            h, s, l = self.utils.rgb_to_hsl(rgb)
            pos = self.LIGHTNESS_POSITIONS[name]
            
            print(f"{name:<15} {hex_color:<10} {h:>6.1f}° {s:>5.1f}% {l:>5.1f}% ({pos:.0%})")


def main():
    import argparse
    
    parser = argparse.ArgumentParser(
        description='Fully adaptive pywal to waybar theme converter'
    )
    parser.add_argument(
        'input',
        nargs='?',
        default='~/Documents/Temp/Colors/merged.json',
        help='Path to pywal colors.json'
    )
    parser.add_argument(
        '-o', '--output',
        default='~/Documents/Temp/Colors/output.css',
        help='Output CSS file'
    )
    parser.add_argument(
        '--stdout',
        action='store_true',
        help='Print to stdout only'
    )
    parser.add_argument(
        '--debug',
        action='store_true',
        help='Show detailed color analysis'
    )
    
    args = parser.parse_args()
    
    input_path = Path(args.input).expanduser()
    output_path = Path(args.output).expanduser()
    
    if not input_path.exists():
        print(f"❌ File not found: {input_path}")
        return 1
    
    # Load pywal colors
    with open(input_path) as f:
        pywal_data = json.load(f)
    
    print(f"📖 Reading: {input_path}\n")
    
    # Generate palette
    generator = AdaptiveWaybarGenerator(pywal_data)
    palette = generator.generate_palette()
    
    if args.debug:
        generator.print_debug(palette)
        print()
    
    # Format CSS
    css_output = generator.format_css(palette)
    
    print(css_output)
    print()
    
    # Save
    if args.stdout:
        print("✅ Output printed to stdout")
    else:
        output_path.parent.mkdir(parents=True, exist_ok=True)
        with open(output_path, 'w') as f:
            f.write(css_output + '\n')
        
        print(f"✅ Saved to: {output_path}")
        print()
        print("💡 Add to waybar style.css:")
        print(f"   @import '{output_path.name}';")
    
    return 0


if __name__ == '__main__':
    exit(main())
