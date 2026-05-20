// Command check_test_placement fails if any *_test.go exists outside backend/tests/.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "check_test_placement: %v\n", err)
		os.Exit(1)
	}

	var violations []string
	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			name := d.Name()
			if name == "vendor" || name == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		rel, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, "tests/") || rel == "tests" {
			return nil
		}
		violations = append(violations, rel)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "check_test_placement: %v\n", err)
		os.Exit(1)
	}
	if len(violations) == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "check_test_placement: *_test.go must live under tests/ (mirror source path):\n")
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  - %s\n", v)
	}
	fmt.Fprintf(os.Stderr, "See .cursor/rules/_lang/go.md and backend/tests/README.md\n")
	os.Exit(1)
}
