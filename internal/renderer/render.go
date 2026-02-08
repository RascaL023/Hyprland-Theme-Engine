package renderer

import (
	"os"
	"path/filepath"
	"text/template"
)

var funcMap = template.FuncMap{
	"hex": func(s string) string {
		if len(s) > 0 && s[0] == '#' {
			return s[1:];
		}
		return s;
	},
}

func Render(templatePath, outputPath string, tool any) error {
	tmpl, err := template.New("").
		Funcs(funcMap).
		ParseFiles(templatePath);
	if err != nil {
		return err;
	}

	out, err := os.Create(outputPath);
	if err != nil {
		return err;
	}
	defer out.Close();

	name := filepath.Base(templatePath);
	return tmpl.ExecuteTemplate(out, name, tool);
}
