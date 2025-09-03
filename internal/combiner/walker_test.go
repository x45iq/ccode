package combiner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sabhiram/go-gitignore"
)

func TestCollectFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{RootDir: dir}

	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	if err := os.WriteFile(f1, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte("world"), 0o644); err != nil {
		t.Fatal(err)
	}

	ign := ignore.CompileIgnoreLines("b.txt")

	files, err := collectFiles(cfg, ign, filepath.Join(dir, "out.txt"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 1 || files[0] != "a.txt" {
		t.Errorf("expected only a.txt, got %v", files)
	}
}
