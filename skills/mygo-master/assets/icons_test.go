package main

import (
	"os"
	"strings"
	"testing"
)

func TestIconsParseAndTheirLicenseTravelsWithThem(t *testing.T) {
	files, err := iconFiles.ReadDir("icons")
	if err != nil || len(files) == 0 {
		t.Fatalf("no icons embedded: %v", err)
	}
	for _, file := range files {
		data, err := iconFiles.ReadFile("icons/" + file.Name())
		if err != nil || len(data) == 0 {
			t.Errorf("%s is empty: %v", file.Name(), err)
		}
	}
	license, err := os.ReadFile("icons/LICENSE")
	if err != nil || !strings.Contains(string(license), "ISC License") {
		t.Errorf("icons/LICENSE is missing or is not the ISC License: %v", err)
	}
}
