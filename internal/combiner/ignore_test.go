package combiner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectIgnorePatterns(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatal(err)
	}

	rootIgnore := filepath.Join(dir, ".ccodeignore")
	if err := os.WriteFile(rootIgnore, []byte("# comment\nfile1.txt\n\\#notcomment\n!keep.txt\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	subIgnore := filepath.Join(subdir, ".ccodeignore")
	if err := os.WriteFile(subIgnore, []byte("subfile.txt\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	patterns, err := collectIgnorePatterns(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]bool{
		"file1.txt":        false,
		"#notcomment":      false,
		"!keep.txt":        false,
		"/sub/subfile.txt": false,
	}
	for _, p := range patterns {
		if _, ok := expected[p]; !ok {
			t.Errorf("unexpected pattern: %s", p)
		} else {
			expected[p] = true
		}
	}
	for p, seen := range expected {
		if !seen {
			t.Errorf("missing expected pattern: %s", p)
		}
	}
}
