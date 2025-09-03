package combiner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func writeCombinedFiles(cfg Config, files []string) error {
	if !cfg.Force {
		if _, err := os.Stat(cfg.Output); err == nil {
			return fmt.Errorf("output file %s already exists (use --force to overwrite)", cfg.Output)
		}
	}

	outFile, err := os.Create(cfg.Output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)

	for _, f := range files {
		if _, err := fmt.Fprintf(writer, "---\nFile: %s\n---\n", f); err != nil {
			return err
		}
		content, err := os.ReadFile(filepath.Join(cfg.RootDir, f))
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", f, err)
		}

		if len(content) == 0 {
			if _, err := writer.WriteString("[empty file]\n\n"); err != nil {
				return err
			}
			continue
		}

		text := string(content)
		if cfg.StripEmpty {
			var sb strings.Builder
			sc := bufio.NewScanner(strings.NewReader(text))
			for sc.Scan() {
				line := sc.Text()
				if strings.TrimSpace(line) == "" {
					continue
				}
				sb.WriteString(line)
				sb.WriteString("\n")
			}
			if err := sc.Err(); err != nil {
				return fmt.Errorf("failed to process file %s: %w", f, err)
			}
			text = sb.String()
		}

		if _, err := writer.WriteString(text + "\n"); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}

	fmt.Printf("Combined %d files into %s\n", len(files), cfg.Output)
	return nil
}
