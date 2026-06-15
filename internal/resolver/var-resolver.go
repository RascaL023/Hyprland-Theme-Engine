package resolver

import (
	"strings"
	"theme-engine/internal/core/themes/vars"
)

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
