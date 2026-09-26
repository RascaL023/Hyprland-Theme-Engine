package pathenv

import (
	"strings"

	"theme-engine/internal/domain/vars"
)

// ResolveVar resolves $pl.* palette references against the given sources.
// Non-$pl strings are returned unchanged. Unknown keys are returned
// unchanged to preserve the historical lenient behaviour.
func ResolveVar(s string, sources ...vars.VarSource) string {
	if !strings.HasPrefix(s, "$pl.") {
		return s
	}

	key := s[4:]

	for _, src := range sources {
		if v, ok := src.Get(key); ok {
			return v
		}
	}

	return s
}
