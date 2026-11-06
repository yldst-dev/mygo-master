package main

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const mygoModule = "github.com/egoist/mygo"

func modulePath(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if name, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(name)
		}
	}
	t.Fatal("go.mod names no module")
	return ""
}

func TestDependenciesPointInward(t *testing.T) {
	module := modulePath(t)
	internal := module + "/internal/"
	layers := map[string][]string{
		"domain":  {},
		"usecase": {internal + "domain"},
		"adapter": {internal + "domain", internal + "usecase", mygoModule},
		"view":    {internal + "domain", internal + "usecase", mygoModule + "/ui"},
	}
	allowed := func(layer, path string) bool {
		if !strings.HasPrefix(path, module+"/") && !strings.Contains(strings.Split(path, "/")[0], ".") {
			return true
		}
		for _, prefix := range layers[layer] {
			if path == prefix || strings.HasPrefix(path, prefix+"/") {
				return true
			}
		}
		return false
	}
	err := filepath.WalkDir("internal", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		layer := strings.Split(filepath.ToSlash(path), "/")[1]
		if _, known := layers[layer]; !known {
			t.Errorf("%s: declare the layer %q in this test", path, layer)
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			imported, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			if !allowed(layer, imported) {
				t.Errorf("%s: the %s layer must not import %s", path, layer, imported)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
