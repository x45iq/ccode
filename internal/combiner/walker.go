package combiner

import (
	"os"
	"path/filepath"

	"github.com/sabhiram/go-gitignore"
)

func collectFiles(cfg Config, ign *ignore.GitIgnore, outAbs string) ([]string, error) {
	var files []string

	err := filepath.Walk(cfg.RootDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == ".ccodeignore" {
			return nil
		}
		pAbs, err := filepath.Abs(path)
		if err == nil && pAbs == outAbs {
			return nil
		}
		relPath, err := filepath.Rel(cfg.RootDir, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		if ign != nil && ign.MatchesPath(relPath) {
			return nil
		}
		files = append(files, relPath)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}
