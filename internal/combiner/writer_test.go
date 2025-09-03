package combiner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteCombinedFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		RootDir: dir,
		Output:  filepath.Join(dir, "out.txt"),
		Force:   true,
	}

	empty := filepath.Join(dir, "empty.txt")
	full := filepath.Join(dir, "full.txt")
	os.WriteFile(empty, []byte(""), 0o644)
	os.WriteFile(full, []byte("line1\n\nline2\n"), 0o644)

	files := []string{"empty.txt", "full.txt"}
	if err := writeCombinedFiles(cfg, files); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(cfg.Output)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	if !strings.Contains(content, "[empty file]") {
		t.Errorf("expected [empty file] marker, got:\n%s", content)
	}
	if !strings.Contains(content, "line1") || !strings.Contains(content, "line2") {
		t.Errorf("missing file contents, got:\n%s", content)
	}
}

func TestWriteCombinedFilesStripEmpty(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		RootDir:    dir,
		Output:     filepath.Join(dir, "out.txt"),
		Force:      true,
		StripEmpty: true,
	}
	full := filepath.Join(dir, "f.txt")
	os.WriteFile(full, []byte("a\n\nb\n"), 0o644)

	files := []string{"f.txt"}
	if err := writeCombinedFiles(cfg, files); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(cfg.Output)
	content := string(data)

	if strings.Contains(content, "a\n\nb") {
		t.Errorf("expected empty lines stripped, got:\n%s", content)
	}
	if !strings.Contains(content, "a\nb") {
		t.Errorf("expected 'a' и 'b' подряд, got:\n%s", content)
	}
}
