package combiner

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/sabhiram/go-gitignore"
)

func Run(ctx context.Context, cfg Config) error {
	patterns, err := collectIgnorePatterns(cfg.RootDir)
	if err != nil {
		return fmt.Errorf("failed to collect .ccodeignore patterns: %w", err)
	}

	var ign *ignore.GitIgnore
	if len(patterns) > 0 {
		ign = ignore.CompileIgnoreLines(patterns...)
	}

	outAbs, err := filepath.Abs(cfg.Output)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	files, err := collectFiles(cfg, ign, outAbs)
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	if cfg.DryRun {
		fmt.Println("[Dry Run] Files to be combined:")
		for _, f := range files {
			fmt.Println(" -", f)
		}
		return nil
	}

	return writeCombinedFiles(cfg, files)
}
