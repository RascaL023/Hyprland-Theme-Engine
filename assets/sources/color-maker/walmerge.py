#!/usr/bin/env python3
import json
import subprocess
import re
import colorsys
from math import sqrt

# ---------------------------------------
# Utilities
# ---------------------------------------

def hex_to_rgb(h):
    h = h.lstrip("#")
    return tuple(int(h[i:i+2], 16) for i in (0, 2, 4))

def luminance(rgb):
    r, g, b = [x/255 for x in rgb]
    return 0.2126*r + 0.7152*g + 0.0722*b

def delta_e(c1, c2):
    return sqrt(sum((a - b)**2 for a, b in zip(c1, c2)))

def extract_magick_colors(image):
    cmd = [
        "magick", image, "-resize", "600x600",
        "-colors", "16", "-unique-colors", "txt:-"
    ]
    out = subprocess.check_output(cmd).decode()

    hexes = re.findall(r"#([0-9A-Fa-f]{6})", out)
    unique = list(dict.fromkeys(hexes))  # remove duplicates
    return ["#" + h for h in unique]

# ---------------------------------------
# Main merge logic
# ---------------------------------------

def brighten(rgb, pct):
    r, g, b = rgb
    return (
        min(int(r + (255 - r) * pct), 255),
        min(int(g + (255 - g) * pct), 255),
        min(int(b + (255 - b) * pct), 255),
    )

def merge(pyal_json_path, out_path):
    py = json.load(open(pyal_json_path))

    # ambil color0–7
    base = [hex_to_rgb(py["colors"][f"color{i}"]) for i in range(8)]

    # generate bright from dark
    brights = []

    for rgb in base:
        # brighten 35%
        rb = brighten(rgb, 0.35)
        brights.append("#{0:02X}{1:02X}{2:02X}".format(*rb))

    # assign ke color8–15
    for i in range(8):
        py["colors"][f"color{i+8}"] = brights[i]

    json.dump(py, open(out_path, "w"), indent=4)

    print("[OK] Bright colors generated ANSI-style.")

# ---------------------------------------
# CLI
# ---------------------------------------
if __name__ == "__main__":
    import sys
    if len(sys.argv) != 3:
        print("Usage: walmerge.py source.json out.json")
        exit(1)

    merge(sys.argv[1], sys.argv[2])

