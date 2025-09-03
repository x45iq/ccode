package combiner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func collectIgnorePatterns(root string) ([]string, error) {
	var patterns []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !info.IsDir() {
			return nil
		}

		ignPath := filepath.Join(path, ".ccodeignore")
		if _, err := os.Stat(ignPath); err != nil {
			return nil
		}

		dirRel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		dirRel = filepath.ToSlash(dirRel)
		if dirRel == "." {
			dirRel = ""
		}

		file, err := os.Open(ignPath)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", ignPath, err)
		}
		defer file.Close()

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			line := sc.Text()
			trimLeft := strings.TrimLeft(line, " \t")
			if trimLeft == "" {
				continue
			}
			if strings.HasPrefix(trimLeft, `\#`) {
				trimLeft = strings.TrimPrefix(trimLeft, `\`)
			} else if strings.HasPrefix(trimLeft, "#") {
				continue
			}

			neg := false
			if strings.HasPrefix(trimLeft, `\!`) {
				trimLeft = strings.TrimPrefix(trimLeft, `\`)
			} else if strings.HasPrefix(trimLeft, "!") {
				neg = true
				trimLeft = strings.TrimPrefix(trimLeft, "!")
			}

			p := strings.TrimSpace(trimLeft)
			if p == "" {
				continue
			}

			leadingSlash := strings.HasPrefix(p, "/")
			if leadingSlash {
				p = strings.TrimPrefix(p, "/")
			}

			prefix := dirRel
			if prefix != "" {
				prefix += "/"
			}

			final := p
			if prefix != "" || leadingSlash {
				final = "/" + prefix + p
			}
			if neg {
				final = "!" + final
			}
			patterns = append(patterns, final)
		}

		if err := sc.Err(); err != nil {
			return fmt.Errorf("failed to read %s: %w", ignPath, err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return patterns, nil
}
