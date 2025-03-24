package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type ignoreRule struct {
	pattern  string
	regex    *regexp.Regexp
	negation bool
}

type IgnoreMatcher struct {
	rules []ignoreRule
}

func NewIgnoreMatcher(ignoreFile string) (*IgnoreMatcher, error) {
	file, err := os.Open(ignoreFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error opening %s: %w", ignoreFile, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var rules []ignoreRule
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		negation := false
		if strings.HasPrefix(line, "!") {
			negation = true
			line = strings.TrimSpace(line[1:])
		}
		regexStr := convertGitignorePattern(line)
		re, err := regexp.Compile(regexStr)
		if err != nil {
			log.Printf("Ignoring invalid pattern in %s at line %d: %s", ignoreFile, lineNum, line)
			continue
		}
		rules = append(rules, ignoreRule{pattern: line, regex: re, negation: negation})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading %s: %w", ignoreFile, err)
	}
	return &IgnoreMatcher{rules: rules}, nil
}

func convertGitignorePattern(pattern string) string {
	pattern = filepath.ToSlash(pattern)
	var regex strings.Builder
	regex.WriteString("^")
	if !strings.HasPrefix(pattern, "/") {
		regex.WriteString("(.*/)?")
	} else {
		pattern = pattern[1:]
	}
	for i := 0; i < len(pattern); i++ {
		if i+1 < len(pattern) && pattern[i] == '*' && pattern[i+1] == '*' {
			regex.WriteString(".*")
			i++
			continue
		}
		switch pattern[i] {
		case '*':
			regex.WriteString("[^/]*")
		case '?':
			regex.WriteString("[^/]")
		case '.':
			regex.WriteString("\\.")
		case '/':
			regex.WriteString("/")
		default:
			regex.WriteString(regexp.QuoteMeta(string(pattern[i])))
		}
	}
	regex.WriteString("$")
	return regex.String()
}

func (im *IgnoreMatcher) Match(relativePath string) bool {
	relativePath = filepath.ToSlash(relativePath)
	ignored := false
	for _, rule := range im.rules {
		if rule.regex.MatchString(relativePath) {
			ignored = !rule.negation
		}
	}
	return ignored
}

func loadIgnoreMatchers(baseDir string) (*IgnoreMatcher, error) {
	var combinedRules []ignoreRule
	currentDir := baseDir
	for {
		ignoreFile := filepath.Join(currentDir, ".ccodeignore")
		matcher, err := NewIgnoreMatcher(ignoreFile)
		if err != nil {
			log.Printf("Error loading %s: %v", ignoreFile, err)
		} else if matcher != nil {
			combinedRules = append(combinedRules, matcher.rules...)
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}
	return &IgnoreMatcher{rules: combinedRules}, nil
}

func isExcluded(matcher *IgnoreMatcher, baseDir, filePath string) bool {
	if matcher == nil {
		return false
	}
	relativePath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		return false
	}
	return matcher.Match(relativePath)
}

func save(file, output string, matcher *IgnoreMatcher, baseDir string) error {
	if isExcluded(matcher, baseDir, file) {
		return nil
	}
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", file, err)
	}
	f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return fmt.Errorf("error opening output file: %w", err)
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	header := fmt.Sprintf("\n// File: %s\n\n", file)
	if _, err := writer.WriteString(header); err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}
	if _, err := writer.Write(content); err != nil {
		return fmt.Errorf("error writing content: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing content: %w", err)
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "ccode",
	Short: "Collect your code in one file",
	Long:  "Collect your code in one file using .ccodeignore for exclusions",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputFile, _ := cmd.Flags().GetString("output")
		force, _ := cmd.Flags().GetBool("force")
		targetDir := args[0]
		absDir, err := filepath.Abs(targetDir)
		if err != nil {
			log.Fatalf("Error resolving path: %v", err)
		}
		matcher, err := loadIgnoreMatchers(absDir)
		if err != nil {
			log.Printf("Warning: %v", err)
		}
		absOutput, err := filepath.Abs(outputFile)
		if err != nil {
			log.Fatalf("Error resolving output path: %v", err)
		}
		if _, err := os.Stat(absOutput); err == nil && !force {
			log.Fatalf("Output file %s already exists. Use --force to overwrite.", absOutput)
		}
		err = os.WriteFile(absOutput, []byte{}, 0660)
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			return save(path, absOutput, matcher, absDir)
		})
		if err != nil {
			log.Fatalf("Error processing files: %v", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("output", "o", "combined.txt", "Output file")
	rootCmd.Flags().Bool("force", false, "Overwrite existing output file")
}
