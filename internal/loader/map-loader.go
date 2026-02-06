package loader

import (
	"bufio"
	"os"
	"strings"
	"theme-engine/internal/core/themes/state"
	"theme-engine/internal/helper"
)

type ToolMap struct {
	TemplatePath	string
	OutputPath		string
}

func LoadToolMap(path string, state *state.State) (map[string]*ToolMap, error) {
	path = helper.ExpandPath(path, state);

	fileMap, err := os.Open(path);
	if err != nil {
		return nil, err;
	}
	defer fileMap.Close();

	res := make(map[string]*ToolMap);
	scanner := bufio.NewScanner(fileMap);

	for scanner.Scan() {
		line := helper.ExpandPath(scanner.Text(), state);
		parts := strings.Split(line, "|");

		if len(parts) != 3 {
			continue;
		}
		
		res[parts[0]] = &ToolMap{
			TemplatePath: parts[1],
			OutputPath: parts[2],
		}
	}

	return res, scanner.Err();
}

