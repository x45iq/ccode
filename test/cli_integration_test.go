package test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/x45iq/ccode/cmd"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestExecuteDryRun(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(file, []byte("abc"), 0o644); err != nil {
		t.Fatal(err)
	}

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(dir)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"ccode", "--dry-run", "--output", filepath.Join(dir, "out.txt")}
	out := captureOutput(func() {
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "a.txt") {
		t.Errorf("expected file in dry run output, got:\n%s", out)
	}
}

func TestExecuteForceWrite(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(file, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(dir)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	outFile := filepath.Join(dir, "out.txt")
	os.Args = []string{"ccode", "--force", "--output", outFile}
	out := captureOutput(func() {
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "Combined 1 files") {
		t.Errorf("expected success message, got:\n%s", out)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "hello") {
		t.Errorf("expected file content in output, got:\n%s", string(data))
	}
}

func TestExecuteStripEmpty(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(file, []byte("x\n\ny\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(dir)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	outFile := filepath.Join(dir, "out.txt")
	os.Args = []string{"ccode", "--force", "--strip-empty", "--output", outFile}
	captureOutput(func() {
		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "x\n\ny") {
		t.Errorf("expected empty lines stripped, got:\n%s", string(data))
	}
	if !strings.Contains(string(data), "x\ny") {
		t.Errorf("expected lines joined without empty, got:\n%s", string(data))
	}
}
