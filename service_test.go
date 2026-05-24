package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSafeSlice(t *testing.T) {
	tests := []struct {
		input  string
		n      int
		expect string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hello"},
		{"", 5, ""},
		{"hi", 0, ""},
		{"abcdef", 3, "abc"},
	}

	for _, tt := range tests {
		result := safeSlice(tt.input, tt.n)
		if result != tt.expect {
			t.Errorf("safeSlice(%q, %d) = %q, want %q", tt.input, tt.n, result, tt.expect)
		}
	}
}

func TestLoadEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	// No file
	env := loadEnv(envPath)
	if len(env) != 0 {
		t.Error("loadEnv should return empty map for missing file")
	}

	// With content
	os.WriteFile(envPath, []byte("KEY1=val1\n# comment\nKEY2=val2\n"), 0644)
	env = loadEnv(envPath)
	if env["KEY1"] != "val1" {
		t.Errorf("KEY1 = %q, want %q", env["KEY1"], "val1")
	}
	if env["KEY2"] != "val2" {
		t.Errorf("KEY2 = %q, want %q", env["KEY2"], "val2")
	}
	if _, ok := env["# comment"]; ok {
		t.Error("comments should be skipped")
	}
}

func TestLoadEnv_EmptyLines(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	os.WriteFile(envPath, []byte("\n\nKEY=val\n\n"), 0644)
	env := loadEnv(envPath)
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d: %v", len(env), env)
	}
}

func TestSafeSlice_Panics(t *testing.T) {
	// Verify no panic on edge cases
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("safeSlice panicked: %v", r)
		}
	}()

	_ = safeSlice("", -1)  // negative n
	_ = safeSlice("a", 0)  // zero
	_ = safeSlice("a", 1)  // exact
	_ = safeSlice("", 0)   // both empty
}
