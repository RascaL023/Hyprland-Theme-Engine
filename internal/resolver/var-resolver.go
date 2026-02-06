package resolver

import "theme-engine/internal/core/themes/vars"

func ResolveVar(s string, sources ...vars.VarSource) string {
	if len(s) < 4 || s[0] != '$' {
		return s;
	}

	key := s[4:] // asumsi $th.xxx

	for _, src := range sources {
		if v, ok := src.Get(key); ok {
			return v;
		}
	}

	return s;
}

