package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"theme-engine/internal/core/themes/state"
	"theme-engine/internal/helper"
)

type ToolMap struct {
	TemplatePath string
	OutputPath   string
}

func LoadToolMap(path string, state *state.State) (map[string]*ToolMap, error) {
	path = helper.ExpandPath(path, state)

	fileMap, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fileMap.Close()

	res := make(map[string]*ToolMap)
	scanner := bufio.NewScanner(fileMap)

	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = helper.ExpandPath(line, state)
		parts := strings.Split(line, "|")

		if len(parts) != 3 {
			return nil, fmt.Errorf("%s:%d: expected name|template|output", path, lineNo)
		}

		name := strings.TrimSpace(parts[0])
		if name == "" {
			return nil, fmt.Errorf("%s:%d: empty target name", path, lineNo)
		}

		res[name] = &ToolMap{
			TemplatePath: strings.TrimSpace(parts[1]),
			OutputPath:   strings.TrimSpace(parts[2]),
		}
	}

	return res, scanner.Err()
}
