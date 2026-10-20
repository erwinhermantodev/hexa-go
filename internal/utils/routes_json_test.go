package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddRoutes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "routes.json")
	if err := os.WriteFile(path, []byte(`{"prefix":"/api/v1","routes":[]}`), 0644); err != nil {
		t.Fatal(err)
	}
	r := []Route{{Method: "GET", Path: "/items", Handler: "ItemHandler.GetAllItems"}}
	for i := 0; i < 2; i++ { // the second call must not duplicate
		if err := AddRoutes(path, r); err != nil {
			t.Fatal(err)
		}
	}
	got, _ := os.ReadFile(path)
	if n := strings.Count(string(got), "ItemHandler.GetAllItems"); n != 1 {
		t.Errorf("route present %d times, want 1:\n%s", n, got)
	}
	if !strings.Contains(string(got), `"prefix": "/api/v1"`) {
		t.Errorf("prefix lost:\n%s", got)
	}
}

func TestAddRoutesRejectsInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "routes.json")
	os.WriteFile(path, []byte("{ nope"), 0644)
	if err := AddRoutes(path, nil); err == nil {
		t.Fatal("expected an error for invalid JSON")
	}
}
