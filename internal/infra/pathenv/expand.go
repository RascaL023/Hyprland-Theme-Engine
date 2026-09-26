package pathenv

import (
	"os"
	"strings"

	"theme-engine/internal/domain/state"
)

func Env(key string) string {
	return os.Getenv(key)
}

func parseKeyword(keyword string, st *state.State) string {
	res := strings.ToUpper(keyword)

	switch res {
	case "WAYBAR":
		res = st.Waybar
	case "THEME":
		res = st.Theme.Name
	default:
		if res = os.Getenv(res); res == "" {
			res = "$" + keyword
		}
	}

	return res
}

// ExpandPath substitutes $VAR segments in filesystem paths.
// Supported: $WAYBAR and $THEME from state, otherwise OS env (upper-cased).
func ExpandPath(path string, st *state.State) string {
	buf := make([]byte, 0, len(path)+32)

	for i := 0; i < len(path); {
		if path[i] != '$' {
			buf = append(buf, path[i])
			i++
			continue
		}

		start := i + 1
		end := start

		for end < len(path) {
			ch := path[end]
			if ch == '/' || ch == '|' || ch == '\\' || ch == '.' {
				break
			}
			end++
		}

		parsed := parseKeyword(path[start:end], st)

		buf = append(buf, parsed...)

		i = end
	}

	return string(buf)
}
