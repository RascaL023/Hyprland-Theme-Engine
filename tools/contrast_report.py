#!/usr/bin/env python3
"""Contrast auditor + palette-slot planner for Theme-Engine.

The audit half mirrors ``internal/domain/palette/contrast_test.go``: the numbers
printed here are the numbers the Go guard asserts.

    nvim inline code   colors[4]      on syntax.terminal_black    >= 3.0
    editor main text   text.primary   on layer.base               >= 4.5
    lazygit selection  text.primary   on layer.surface_overlay    >= 3.0
    lualine section B  mode accents   on ui.gutter                >= 3.0

The planning half answers "which colour should this slot be?". It walks shades
of ``layer.surface_overlay`` in OKLCH, so only lightness moves: hue and chroma
survive the nudge instead of drifting to a washed-out white. Candidates are
ranked by perceptual distance from the colour the theme uses today, with a
small bonus for reusing a value the palette already declares.

Usage:
    python3 tools/contrast_report.py                  # audit every theme
    python3 tools/contrast_report.py harbor --plan    # one theme, with proposals
    python3 tools/contrast_report.py --variant light --plan
    python3 tools/contrast_report.py --json
"""

from __future__ import annotations

import argparse
import json
import math
import os
import sys
from dataclasses import dataclass, field

REPO_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
THEMES_DIR = os.path.join(REPO_ROOT, "themes")

# Guard thresholds, identical to the Go guard.
MIN_UI = 3.0
MIN_TEXT = 4.5
# Extra headroom applied when planning so a value never sits exactly on the edge.
DEFAULT_HEADROOM = 0.05

VARIANTS = ("dark", "light")


# --------------------------------------------------------------- colour maths


def parse_hex(value: str) -> tuple[float, float, float]:
    """``#rgb`` / ``#rrggbb`` -> linear-in-gamma sRGB floats (0..1)."""
    h = value.strip().lstrip("#")
    if len(h) == 3:
        h = "".join(c * 2 for c in h)
    if len(h) != 6:
        raise ValueError(f"not a hex colour: {value!r}")
    return tuple(int(h[i : i + 2], 16) / 255 for i in (0, 2, 4))  # type: ignore[return-value]


def to_hex(rgb: tuple[float, float, float]) -> str:
    channels = [max(0, min(255, round(c * 255))) for c in rgb]
    return "#%02X%02X%02X" % tuple(channels)


def _linearize(c: float) -> float:
    return c / 12.92 if c <= 0.04045 else ((c + 0.055) / 1.055) ** 2.4


def _delinearize(c: float) -> float:
    return 12.92 * c if c <= 0.0031308 else 1.055 * c ** (1 / 2.4) - 0.055


def relative_luminance(hex_color: str) -> float:
    """WCAG relative luminance. Returns 0.0 for unparseable input."""
    try:
        r, g, b = (max(0.0, min(1.0, _linearize(c))) for c in parse_hex(hex_color))
    except ValueError:
        return 0.0
    return 0.2126 * r + 0.7152 * g + 0.0722 * b


def contrast(fg: str, bg: str) -> float:
    """WCAG contrast ratio, 1.0..21.0. 0.0 when either colour is invalid."""
    if not fg or not bg:
        return 0.0
    try:
        parse_hex(fg)
        parse_hex(bg)
    except ValueError:
        return 0.0
    a, b = relative_luminance(fg), relative_luminance(bg)
    if a < b:
        a, b = b, a
    return (a + 0.05) / (b + 0.05)


# OKLab / OKLCH (Bjorn Ottosson). Used only for planning; the audit above is
# plain WCAG so it stays bit-identical to the Go guard.


def hex_to_oklab(hex_color: str) -> tuple[float, float, float]:
    r, g, b = (_linearize(c) for c in parse_hex(hex_color))
    l = 0.4122214708 * r + 0.5363325363 * g + 0.0514459929 * b
    m = 0.2119034982 * r + 0.6806995451 * g + 0.1073969566 * b
    s = 0.0883024619 * r + 0.2817188376 * g + 0.6299787005 * b
    l_, m_, s_ = (math.copysign(abs(v) ** (1 / 3), v) for v in (l, m, s))
    return (
        0.2104542553 * l_ + 0.7936177850 * m_ - 0.0040720468 * s_,
        1.9779984951 * l_ - 2.4285922050 * m_ + 0.4505937099 * s_,
        0.0259040371 * l_ + 0.7827717662 * m_ - 0.8086757660 * s_,
    )


def oklab_to_rgb(lab: tuple[float, float, float]) -> tuple[float, float, float]:
    L, a, b = lab
    l_ = L + 0.3963377774 * a + 0.2158037573 * b
    m_ = L - 0.1055613458 * a - 0.0638541728 * b
    s_ = L - 0.0894841775 * a - 1.2914855480 * b
    l, m, s = (v**3 for v in (l_, m_, s_))
    return (
        _delinearize(4.0767416621 * l - 3.3077115913 * m + 0.2309699292 * s),
        _delinearize(-1.2684380046 * l + 2.6097574011 * m - 0.3413193965 * s),
        _delinearize(-0.0041960863 * l - 0.7034186147 * m + 1.7076147010 * s),
    )


def hex_to_oklch(hex_color: str) -> tuple[float, float, float]:
    L, a, b = hex_to_oklab(hex_color)
    c = math.hypot(a, b)
    h = math.degrees(math.atan2(b, a)) % 360
    return L, c, h


def _in_gamut(rgb: tuple[float, float, float], eps: float = 1e-4) -> bool:
    return all(-eps <= c <= 1 + eps for c in rgb)


def oklch_to_hex(L: float, C: float, h: float) -> str:
    """OKLCH -> hex, reducing chroma (never hue) until the colour fits sRGB."""
    rad = math.radians(h)
    L = max(0.0, min(1.0, L))
    rgb = oklab_to_rgb((L, C * math.cos(rad), C * math.sin(rad)))
    if not _in_gamut(rgb):
        lo, hi = 0.0, C
        for _ in range(40):
            mid = (lo + hi) / 2
            if _in_gamut(oklab_to_rgb((L, mid * math.cos(rad), mid * math.sin(rad)))):
                lo = mid
            else:
                hi = mid
        rgb = oklab_to_rgb((L, lo * math.cos(rad), lo * math.sin(rad)))
    return to_hex(tuple(max(0.0, min(1.0, c)) for c in rgb))


def shade(hex_color: str, dL: float = 0.0, chroma_scale: float = 1.0) -> str:
    """Shift a colour's OKLCH lightness (and optionally mute its chroma)."""
    L, C, h = hex_to_oklch(hex_color)
    return oklch_to_hex(L + dL, C * chroma_scale, h)


def delta_oklab(a: str, b: str) -> float:
    """Perceptual distance (OKLab euclidean) between two colours."""
    la, lb = hex_to_oklab(a), hex_to_oklab(b)
    return math.sqrt(sum((x - y) ** 2 for x, y in zip(la, lb)))


# ------------------------------------------------------------- palette model


def resolve_var(value: str, raw: dict) -> str:
    """Single-pass ``$pl.`` resolution, same contract as the Go resolver.

    ``palette.json`` may only reference the base palette (foreground,
    background, cursor, color0..color15) from inside ``extra.*``.
    """
    if not isinstance(value, str) or not value.startswith("$pl."):
        return value
    key = value[4:]
    if key == "foreground":
        return raw.get("foreground", "")
    if key == "background":
        return raw.get("background", "")
    if key == "cursor":
        return raw.get("cursor", "")
    if key.startswith("color"):
        index = int(key[5:])
        colors = raw.get("colors", [])
        if 0 <= index < len(colors):
            return resolve_var(colors[index], raw)
    return value


@dataclass
class Palette:
    """One resolved variant, keeping the raw tree around for provenance."""

    theme: str
    variant: str
    raw: dict
    slots: dict[str, str] = field(default_factory=dict)
    colors: list[str] = field(default_factory=list)

    def get(self, name: str) -> str:
        return self.slots.get(name, "")

    def pool(self) -> dict[str, str]:
        """Every declared/resolved value, for ``reuse an existing colour``."""
        out: dict[str, str] = {}
        for name, value in self.slots.items():
            out.setdefault(value, name)
        for index, value in enumerate(self.colors):
            out.setdefault(value, f"colors[{index}]")
        for layer, value in self.raw.get("extra", {}).get("layer", {}).items():
            out.setdefault(resolve_var(value, self.raw), f"layer.{layer}")
        return out


def _fallback(value: str, fallback: str) -> str:
    return value if value else fallback


def resolve_palette(theme: str, variant: str, raw: dict) -> Palette:
    """Mirror of ``resolver.go``: same fields, same fallback chain."""
    extra = raw.get("extra", {})
    accent = extra.get("accent", {})
    text = extra.get("text", {})
    layer = extra.get("layer", {})
    border = extra.get("border", {})
    status = extra.get("status", {})
    syntax = extra.get("syntax", {})
    ui = extra.get("ui", {})
    fg = extra.get("fg", {})
    bg = extra.get("bg", {})

    def rv(value: str) -> str:
        return resolve_var(value, raw)

    colors = [rv(c) for c in raw.get("colors", [])]

    slots = {
        "foreground": rv(raw.get("foreground", "")),
        "background": rv(raw.get("background", "")),
        "cursor": rv(raw.get("cursor", "")),
        "accent.primary": rv(accent.get("primary", "")),
        "accent.secondary": rv(accent.get("secondary", "")),
        "accent.on_accent": rv(accent.get("on_accent", "")),
        "text.primary": rv(text.get("primary", "")),
        "text.secondary": rv(text.get("secondary", "")),
        "text.muted": rv(text.get("muted", "")),
        "text.link": rv(text.get("link", "")),
        "layer.base": rv(layer.get("base", "")),
        "layer.mantle": rv(layer.get("mantle", "")),
        "layer.crust": rv(layer.get("crust", "")),
        "layer.surface": rv(layer.get("surface", "")),
        "layer.surface_raised": rv(layer.get("surface_raised", "")),
        "layer.surface_overlay": rv(layer.get("surface_overlay", "")),
        "border.default": rv(border.get("default", "")),
        "border.active": rv(border.get("active", "")),
        "status.success": rv(status.get("success", "")),
        "status.warning": rv(status.get("warning", "")),
        "status.error": rv(status.get("error", "")),
        "status.critical": rv(status.get("critical", "")),
        "status.info": rv(status.get("info", "")),
        "ui.bg_statusline": rv(ui.get("bg_statusline", "")),
        "ui.gutter": rv(ui.get("gutter", "")),
        "syntax.terminal_black": rv(syntax.get("terminal_black", "")),
    }

    palette = Palette(theme=theme, variant=variant, raw=raw, slots=slots, colors=colors)

    base = palette.get
    syntax_fallbacks = {
        "syntax.purple": base("accent.primary"),
        "syntax.magenta2": base("accent.secondary"),
        "syntax.blue0": base("layer.surface"),
        "syntax.blue1": base("text.link"),
        "syntax.blue5": colors[4] if len(colors) > 4 else "",
        "syntax.blue6": colors[6] if len(colors) > 6 else "",
        "syntax.blue7": base("border.default"),
        "syntax.green1": base("status.success"),
        "syntax.green2": colors[2] if len(colors) > 2 else "",
        "syntax.orange": base("status.warning"),
        "syntax.red1": base("status.critical"),
        "syntax.teal": colors[6] if len(colors) > 6 else "",
    }
    for name, fallback_value in syntax_fallbacks.items():
        slots[name] = _fallback(rv(syntax.get(name.split(".")[1], "")), fallback_value)

    slots["syntax.terminal_black"] = _fallback(
        slots["syntax.terminal_black"], base("layer.surface_overlay")
    )
    slots["ui.bg_statusline"] = _fallback(base("ui.bg_statusline"), base("layer.surface_raised"))
    slots["ui.gutter"] = _fallback(base("ui.gutter"), base("layer.surface_overlay"))
    slots["text.visited"] = _fallback(rv(text.get("visited", "")), base("accent.secondary"))
    slots["border.medium"] = _fallback(rv(border.get("medium", "")), base("border.default"))
    slots["fg.dim"] = _fallback(rv(fg.get("dim", "")), base("text.muted"))
    slots["fg.dim_muted"] = _fallback(rv(fg.get("dim_muted", "")), base("text.muted"))
    slots["fg.disabled"] = _fallback(rv(fg.get("disabled", "")), base("text.muted"))
    slots["bg.conflict"] = _fallback(rv(bg.get("conflict", "")), base("status.warning"))
    slots["bg.disk_usage"] = _fallback(rv(bg.get("disk_usage", "")), base("layer.surface_overlay"))

    return palette


def load_theme(theme: str, variant: str) -> Palette:
    path = os.path.join(THEMES_DIR, theme, "palette.json")
    with open(path, encoding="utf-8") as handle:
        data = json.load(handle)
    if variant not in data.get("palettes", {}):
        raise KeyError(f"{theme} has no '{variant}' palette")
    return resolve_palette(theme, variant, data["palettes"][variant])


def discover_themes() -> list[str]:
    names = []
    for entry in sorted(os.listdir(THEMES_DIR)):
        if entry.endswith(".bak"):
            continue
        if os.path.isfile(os.path.join(THEMES_DIR, entry, "palette.json")):
            names.append(entry)
    return names


# ------------------------------------------------------------ checks & slots


@dataclass
class Constraint:
    """A foreground that must stay readable on the slot's colour."""

    label: str
    hex: str
    minimum: float


@dataclass
class SideRole:
    """A second job the same colour has to do, with a soft floor."""

    label: str
    other: str
    minimum: float
    weight: float


@dataclass
class SlotPlan:
    """A palette slot the planner is allowed to move."""

    slot: str
    title: str
    why: str
    constraints: list[Constraint]
    baseline: str
    side_roles: list[SideRole]


def mode_accents(palette: Palette) -> list[tuple[str, str]]:
    """Tokyonight's lualine mode accents: the ANSI set plus syntax.green1."""
    names = ["red", "green", "yellow", "blue", "magenta"]
    out = [(names[i], palette.colors[i + 1]) for i in range(5)]
    out.append(("green1", palette.get("syntax.green1")))
    return out


def plan_slots(palette: Palette) -> list[SlotPlan]:
    base = palette.get
    inline_fg = palette.colors[4] if len(palette.colors) > 4 else ""

    return [
        SlotPlan(
            slot="ui.gutter",
            title="lualine section B background (fg_gutter)",
            why="tokyonight lualine b = { bg = fg_gutter, fg = mode colour }",
            constraints=[
                Constraint(label=name, hex=hex_value, minimum=MIN_UI)
                for name, hex_value in mode_accents(palette)
            ],
            # The effective value: the declared ui.gutter, else the injected
            # fallback (layer.surface_overlay).
            baseline=base("ui.gutter"),
            side_roles=[
                # fg_gutter is also the LineNr / indent-guide foreground: it must
                # stay faintly visible on the editor background, not vanish.
                SideRole("line numbers / indent guides on base", base("layer.base"), 1.2, 0.5),
                # ...and a background behind text.primary (folded lines, tabline).
                # Low floor on purpose: it only has to avoid a collision.
                SideRole("text on gutter (folded / tabline)", base("text.primary"), 1.5, 0.15),
            ],
        ),
        SlotPlan(
            slot="syntax.terminal_black",
            title="inline-code background (markdown backticks)",
            why="tokyonight @markup.raw.markdown_inline = { bg = terminal_black, fg = blue }",
            constraints=[Constraint(label="blue (colors[4])", hex=inline_fg, minimum=MIN_UI)],
            baseline=base("syntax.terminal_black"),
            side_roles=[
                SideRole("ghost text / diagnostics on base", base("layer.base"), 1.2, 0.4),
            ],
        ),
    ]


# ------------------------------------------------------------------ planner


def constraint_report(hex_value: str, constraints: list[Constraint]) -> list[tuple[str, float]]:
    return [(c.label, contrast(c.hex, hex_value)) for c in constraints]


def satisfies(hex_value: str, constraints: list[Constraint], target: float) -> bool:
    return all(contrast(c.hex, hex_value) >= target for c in constraints)


def side_penalty(hex_value: str, roles: list[SideRole]) -> float:
    penalty = 0.0
    for role in roles:
        if not role.other:
            continue
        got = contrast(role.other, hex_value)
        if got < role.minimum:
            penalty += (role.minimum - got) * role.weight
    return penalty


def generate_candidates(baseline: str, pool: dict[str, str]) -> list[tuple[str, str]]:
    """Shades of the baseline plus every colour the palette already declares.

    Lightness moves in fine steps; chroma is also tried at 70% and 40% so a
    value that cannot stay saturated in sRGB lands on a tasteful muted shade
    rather than on a clipped one.
    """
    candidates: dict[str, str] = {baseline: "current value (unchanged)"}

    L0, _, _ = hex_to_oklch(baseline)
    step = 0.002
    for chroma_scale in (1.0, 0.7, 0.4):
        dL = step
        while dL <= 0.75:
            for sign in (1, -1):
                new_L = L0 + sign * dL
                if 0.0 <= new_L <= 1.0:
                    value = shade(baseline, dL=sign * dL, chroma_scale=chroma_scale)
                    candidates.setdefault(value, f"shade of current value at L {new_L:.2f}")
            dL += step

    for value, name in pool.items():
        if len(value) == 7 and value.startswith("#"):
            candidates.setdefault(value, f"existing {name}")
    return list(candidates.items())


def plan_slot(slot: SlotPlan, pool: dict[str, str], headroom: float) -> dict:
    """Pick the best value for a slot: closest look that still passes."""
    target = max(c.minimum for c in slot.constraints) + headroom

    current = constraint_report(slot.baseline, slot.constraints)
    current_ok = satisfies(slot.baseline, slot.constraints, target)

    best: dict | None = None
    for value, source in generate_candidates(slot.baseline, pool):
        if not value or len(value) != 7:
            continue
        if not satisfies(value, slot.constraints, target):
            continue
        score = delta_oklab(slot.baseline, value)
        if value in pool:
            # Tiny nudge only: a value the palette already declares wins ties,
            # but never beats a noticeably closer shade.
            score -= 0.02
        score += side_penalty(value, slot.side_roles)
        if best is None or score < best["score"]:
            best = {
                "value": value,
                "source": source,
                "score": score,
                "delta": delta_oklab(slot.baseline, value),
                "ratios": constraint_report(value, slot.constraints),
            }

    advisories = []
    if best and best["delta"] > 0.12 and not current_ok:
        advisories.append(
            f"the slot had to move {best['delta']:.2f} dE away from "
            f"{slot.baseline} — this palette's accents are the real limit"
        )
    for constraint in slot.constraints:
        if contrast(constraint.hex, slot.baseline) >= target:
            continue
        needed = lightness_for_contrast(constraint.hex, slot.baseline, target)
        if needed is None:
            continue
        have = hex_to_oklch(constraint.hex)[0]
        verb = "darken" if needed < have else "lighten"
        advisories.append(
            f"{constraint.label} {constraint.hex} sits at L {have:.2f}; {verb} it to "
            f"L {needed:.2f} and the slot can stay at {slot.baseline}"
        )

    return {
        "current": current,
        "current_ok": current_ok,
        "best": best,
        "advisories": advisories,
    }


def lightness_for_contrast(fg: str, bg: str, target: float) -> float | None:
    """OKLCH lightness the foreground needs (hue/chroma kept) to reach target.

    Contrast is monotonic in relative luminance, so a bisection lands on the
    smallest lightness that clears the threshold. ``None`` means the colour
    cannot get there without changing its hue.
    """
    L0 = hex_to_oklch(fg)[0]
    at = lambda L: shade(fg, dL=L - L0)  # noqa: E731 - tiny local helper

    if relative_luminance(fg) >= relative_luminance(bg):
        lo, hi = L0, 1.0
        if contrast(at(hi), bg) < target:
            return None
        for _ in range(50):
            mid = (lo + hi) / 2
            if contrast(at(mid), bg) >= target:
                hi = mid
            else:
                lo = mid
        return hi

    lo, hi = 0.0, L0
    if contrast(at(lo), bg) < target:
        return None
    for _ in range(50):
        mid = (lo + hi) / 2
        if contrast(at(mid), bg) >= target:
            lo = mid
        else:
            hi = mid
    return lo


# ------------------------------------------------------------------ audit


@dataclass
class AuditRow:
    label: str
    fg: str
    bg: str
    minimum: float
    ratio: float
    note: str

    @property
    def ok(self) -> bool:
        return self.ratio >= self.minimum


def audit(palette: Palette) -> list[AuditRow]:
    base = palette.get
    rows = [
        AuditRow(
            "editor main text",
            base("text.primary"),
            base("layer.base"),
            MIN_TEXT,
            contrast(base("text.primary"), base("layer.base")),
            "text.primary on layer.base",
        ),
        AuditRow(
            "lazygit selected line",
            base("text.primary"),
            base("layer.surface_overlay"),
            MIN_UI,
            contrast(base("text.primary"), base("layer.surface_overlay")),
            "defaultFgColor on selectedLineBgColor",
        ),
    ]

    plans = plan_slots(palette)
    gutter = plans[0]
    terminal_black = plans[1]
    for name, hex_value in mode_accents(palette):
        rows.append(
            AuditRow(
                f"lualine B {name}",
                hex_value,
                base("ui.gutter"),
                MIN_UI,
                contrast(hex_value, base("ui.gutter")),
                gutter.why,
            )
        )
    for constraint in terminal_black.constraints:
        rows.append(
            AuditRow(
                "inline code (backtick)",
                constraint.hex,
                base("syntax.terminal_black"),
                MIN_UI,
                contrast(constraint.hex, base("syntax.terminal_black")),
                terminal_black.why,
            )
        )
    return rows


# ------------------------------------------------------------------ output

RESET = "\033[0m"
RED = "\033[31m"
GREEN = "\033[32m"
DIM = "\033[2m"
BOLD = "\033[1m"


def use_color() -> bool:
    return sys.stdout.isatty() and os.environ.get("NO_COLOR") is None


def paint(text: str, code: str) -> str:
    return f"{code}{text}{RESET}" if use_color() else text


def render_text(theme: str, variant: str, palette: Palette, with_plan: bool,
                headroom: float, with_accents: bool = False) -> None:
    header = f"{theme}/{variant}"
    print(paint(header, BOLD))
    print(paint("-" * len(header), DIM))

    for row in audit(palette):
        mark = paint("ok  ", GREEN) if row.ok else paint("FAIL", RED)
        print(f"  {mark} {row.label:<22} {row.ratio:5.2f}:1 (min {row.minimum:.1f}) "
              f"{row.fg} on {row.bg}")
        if not row.ok:
            print(paint(f"       {row.note}", DIM))

    if not with_plan:
        print()
        return

    pool = palette.pool()
    for slot in plan_slots(palette):
        result = plan_slot(slot, pool, headroom)
        current = ", ".join(f"{name} {ratio:.2f}" for name, ratio in result["current"])
        if result["current_ok"]:
            print(f"  {paint('keep ', GREEN)} {slot.slot} = {slot.baseline} ({current})")
        elif result["best"]:
            best = result["best"]
            ratios = ", ".join(f"{name} {ratio:.2f}" for name, ratio in best["ratios"])
            print(f"  {paint('set  ', RED)} {slot.slot} = {best['value']} "
                  f"[{best['source']}, dE {best['delta']:.3f}] ({ratios})")
            print(paint(f"       was {slot.baseline} → {current}", DIM))
        else:
            print(f"  {paint('none ', RED)} {slot.slot}: no value in this palette can reach "
                  f"{slot.constraints[0].minimum + headroom:.2f}:1")
        for advisory in result["advisories"]:
            print(paint(f"       advisory: {advisory}", DIM))

    if with_accents:
        for suggestion in accent_suggestions(palette, headroom):
            print(f"  {paint('accent', RED)} {suggestion['label']}: {suggestion['current']} → "
                  f"{suggestion['suggested']} ({suggestion['was']:.2f} → {suggestion['ratio']:.2f}:1 "
                  f"on {suggestion['against']}, hue kept)")
    print()


def accent_suggestions(palette: Palette, headroom: float) -> list[dict]:
    """Accents that block a slot, and the lightness that would unblock it.

    Raising a washed-out accent is the better fix than dragging the surface to
    near black: it keeps the theme's own surface vocabulary intact. Hue and
    chroma are preserved, only lightness moves.
    """
    out: list[dict] = []
    for slot in plan_slots(palette):
        target = max(c.minimum for c in slot.constraints) + headroom
        for constraint in slot.constraints:
            if contrast(constraint.hex, slot.baseline) >= target:
                continue
            needed = lightness_for_contrast(constraint.hex, slot.baseline, target)
            if needed is None:
                continue
            have = hex_to_oklch(constraint.hex)[0]
            fixed = shade(constraint.hex, dL=needed - have)
            out.append(
                {
                    "slot": slot.slot,
                    "label": constraint.label,
                    "current": constraint.hex,
                    "suggested": fixed,
                    "ratio": contrast(fixed, slot.baseline),
                    "was": contrast(constraint.hex, slot.baseline),
                    "against": slot.baseline,
                }
            )
    return out


def snippet(theme: str, variant: str, palette: Palette, headroom: float) -> dict:
    pool = palette.pool()
    slots: dict[str, str] = {}
    for slot in plan_slots(palette):
        result = plan_slot(slot, pool, headroom)
        if result["current_ok"]:
            slots[slot.slot] = slot.baseline
        elif result["best"]:
            slots[slot.slot] = result["best"]["value"]
    return {"theme": theme, "variant": variant, **slots}


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("themes", nargs="*", help="theme names (default: every active theme)")
    parser.add_argument("--variant", choices=VARIANTS, help="only this variant")
    parser.add_argument("--plan", action="store_true", help="propose values for the guard slots")
    parser.add_argument("--json", action="store_true", help="machine-readable output")
    parser.add_argument("--snippet", action="store_true", help="print palette.json fragments")
    parser.add_argument("--accents", action="store_true",
                        help="also propose accent fixes so a slot can stay put")
    parser.add_argument("--headroom", type=float, default=DEFAULT_HEADROOM,
                        help="extra ratio added to the threshold when planning")
    args = parser.parse_args(argv)

    themes = args.themes or discover_themes()
    variants = [args.variant] if args.variant else list(VARIANTS)

    exit_code = 0
    results = []
    for theme in themes:
        for variant in variants:
            try:
                palette = load_theme(theme, variant)
            except (OSError, KeyError) as error:
                print(f"{theme}/{variant}: {error}", file=sys.stderr)
                exit_code = 1
                continue

            rows = audit(palette)
            if any(not row.ok for row in rows):
                exit_code = 1

            if args.json:
                results.append(
                    {
                        "theme": theme,
                        "variant": variant,
                        "audit": [
                            {
                                "check": row.label,
                                "fg": row.fg,
                                "bg": row.bg,
                                "ratio": round(row.ratio, 3),
                                "min": row.minimum,
                                "ok": row.ok,
                            }
                            for row in rows
                        ],
                        "plan": snippet(theme, variant, palette, args.headroom),
                        "accents": accent_suggestions(palette, args.headroom),
                    }
                )
            elif args.snippet:
                data = snippet(theme, variant, palette, args.headroom)
                print(f"{theme}/{variant}")
                print(f'  "syntax": {{ "terminal_black": "{data["syntax.terminal_black"]}" }},')
                print(f'  "ui":     {{ "gutter": "{data["ui.gutter"]}" }}')
            else:
                render_text(theme, variant, palette, args.plan, args.headroom, args.accents)

    if args.json:
        print(json.dumps(results, indent=2))
    return exit_code


if __name__ == "__main__":
    sys.exit(main())
