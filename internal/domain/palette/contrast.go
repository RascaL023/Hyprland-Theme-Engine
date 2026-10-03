package palette

import (
	"math"
	"strconv"
	"strings"
)

// ContrastRatio returns the WCAG 2.x contrast ratio between two colours
// expressed as "#rgb" or "#rrggbb". The result is in [1, 21]: 1 means the
// colours are identical and higher means more readable. It returns 0 when
// either colour cannot be parsed.
//
// Domain helper only; it is never called on the render path (see
// VISUAL_FIX_PLAN.md §12).
func ContrastRatio(a, b string) float64 {
	la, okA := RelativeLuminance(a)
	lb, okB := RelativeLuminance(b)
	if !okA || !okB {
		return 0
	}
	if lb > la {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

// RelativeLuminance returns the WCAG relative luminance of a colour expressed
// as "#rgb" or "#rrggbb". ok is false when the input cannot be parsed.
func RelativeLuminance(hex string) (luminance float64, ok bool) {
	r, g, b, ok := parseHex(hex)
	if !ok {
		return 0, false
	}
	return 0.2126*linearize(r) + 0.7152*linearize(g) + 0.0722*linearize(b), true
}

func linearize(c float64) float64 {
	if c <= 0.03928 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func parseHex(hex string) (r, g, b float64, ok bool) {
	s := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	switch len(s) {
	case 3:
		s = string(s[0]) + string(s[0]) + string(s[1]) + string(s[1]) + string(s[2]) + string(s[2])
	case 6:
	default:
		return 0, 0, 0, false
	}
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	return float64((v>>16)&0xff) / 255, float64((v>>8)&0xff) / 255, float64(v&0xff) / 255, true
}
