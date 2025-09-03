package combiner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDryRun(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	os.WriteFile(f1, []byte("abc"), 0o644)

	cfg := Config{
		RootDir: dir,
		Output:  filepath.Join(dir, "out.txt"),
		DryRun:  true,
	}
	if err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWrite(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	os.WriteFile(f1, []byte("abc"), 0o644)

	cfg := Config{
		RootDir: dir,
		Output:  filepath.Join(dir, "out.txt"),
		Force:   true,
	}
	if err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(cfg.Output)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "abc") {
		t.Errorf("expected combined content, got: %s", string(data))
	}
}
