package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunAdd(t *testing.T) {
	tests := []struct {
		component string
		wantFile  string
	}{
		{"calendar", "calendar.templ"},
		{"navigator", "navigator.templ"},
		{"jumper", "jumper.templ"},
	}
	for _, tt := range tests {
		t.Run(tt.component, func(t *testing.T) {
			dest := t.TempDir()
			if err := runAdd(tt.component, dest); err != nil {
				t.Fatalf("runAdd(%q, dest) error: %v", tt.component, err)
			}
			path := filepath.Join(dest, tt.wantFile)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("expected file %s not found: %v", path, err)
			}
			if len(data) == 0 {
				t.Errorf("file %s is empty", path)
			}
		})
	}
}

func TestRunAddUnknownComponent(t *testing.T) {
	err := runAdd("nonexistent", t.TempDir())
	if err == nil {
		t.Fatal("expected error for unknown component, got nil")
	}
}

func TestRunAddCreatesDestination(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "deeply", "nested", "dir")
	if err := runAdd("calendar", dest); err != nil {
		t.Fatalf("runAdd error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "calendar.templ")); err != nil {
		t.Errorf("expected file in created destination: %v", err)
	}
}
