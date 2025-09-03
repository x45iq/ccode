package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/x45iq/ccode/internal/combiner"
)

func NewRootCmd() *cobra.Command {
	var (
		force      bool
		dryRun     bool
		output     string
		stripEmpty bool
	)

	rootCmd := &cobra.Command{
		Use:   "ccode",
		Short: "Lightweight text combiner",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			cfg := combiner.Config{
				RootDir:    cwd,
				Output:     output,
				Force:      force,
				DryRun:     dryRun,
				StripEmpty: stripEmpty,
			}
			return combiner.Run(ctx, cfg)
		},
	}

	rootCmd.Flags().BoolVar(&force, "force", false, "Overwrite output file if exists")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview without writing output")
	rootCmd.Flags().StringVarP(&output, "output", "o", "combined.txt", "Path to output file")
	rootCmd.Flags().BoolVar(&stripEmpty, "strip-empty", false, "Remove empty lines")

	return rootCmd
}

func Execute() error {
	return NewRootCmd().Execute()
}
