package renderer

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

var funcMap = template.FuncMap{
	"hex": func(s string) string {
		if len(s) > 0 && s[0] == '#' {
			return s[1:]
		}
		return s
	},
}

var cache = struct {
	sync.RWMutex
	templates map[string]*template.Template
}{
	templates: make(map[string]*template.Template),
}

func Render(templatePath, outputPath string, tool any) error {
	tmpl, err := loadTemplate(templatePath)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	name := filepath.Base(templatePath)
	if err := tmpl.ExecuteTemplate(&buf, name, tool); err != nil {
		return err
	}

	return writeIfChanged(outputPath, buf.Bytes())
}

func loadTemplate(templatePath string) (*template.Template, error) {
	cache.RLock()
	tmpl, ok := cache.templates[templatePath]
	cache.RUnlock()
	if ok {
		return tmpl, nil
	}

	cache.Lock()
	defer cache.Unlock()

	if tmpl, ok = cache.templates[templatePath]; ok {
		return tmpl, nil
	}

	tmpl, err := template.New("").
		Funcs(funcMap).
		ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	cache.templates[templatePath] = tmpl
	return tmpl, nil
}

func writeIfChanged(path string, data []byte) error {
	old, err := os.ReadFile(path)
	if err == nil && bytes.Equal(old, data) {
		return nil
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}

	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}

	return os.Rename(tmpName, path)
}
