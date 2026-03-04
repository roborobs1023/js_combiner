package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func minify(js string) string {
	reMulti := regexp.MustCompile(`(?s)/\*.*?\*/`)
	js = reMulti.ReplaceAllString(js, "")
	reSingle := regexp.MustCompile(`(?m)^[ \t]*//.*$`)
	js = reSingle.ReplaceAllString(js, "")
	reSpace := regexp.MustCompile(`\s+`)
	js = reSpace.ReplaceAllString(js, " ")
	return strings.TrimSpace(js)
}

func main() {
	srcPtr := flag.String("src", "./src/js", "Source directory")
	outPtr := flag.String("out", "./dist/bundle.js", "Output file path")
	doMinify := flag.Bool("minify", false, "Minify the output")

	flag.Parse()

	// Clean paths to avoid trailing slash issues
	srcDir := filepath.Clean(*srcPtr)
	outFile := filepath.Clean(*outPtr)

	fmt.Printf("🔍 Scanning: %s\n", srcDir)

	// Verify source directory exists
	info, err := os.Stat(srcDir)
	if os.IsNotExist(err) {
		fmt.Printf("❌ Error: Source directory '%s' not found.\n", srcDir)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Printf("❌ Error: '%s' is a file, not a directory.\n", srcDir)
		os.Exit(1)
	}

	// Prepare output
	if err := os.MkdirAll(filepath.Dir(outFile), 0755); err != nil {
		fmt.Printf("❌ Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	out, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("❌ Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	fileCount := 0

	// Using WalkDir for O(n) efficiency with DirEntry
	err = filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Case-insensitive extension check
		isJS := strings.ToLower(filepath.Ext(path)) == ".js"

		if !d.IsDir() && isJS && path != outFile {
			fmt.Printf("  -> Adding: %s\n", path)
			fileCount++

			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", path, err)
			}

			data := string(content)
			if *doMinify {
				data = minify(data)
			} else {
				data = fmt.Sprintf("\n// Source: %s\n%s\n", path, data)
			}

			if _, err := out.WriteString(data); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Walk failed: %v\n", err)
	} else if fileCount == 0 {
		fmt.Println("⚠️  Warning: No .js files found in the source directory.")
	} else {
		fmt.Printf("✅ Success! Bundled %d files into %s\n", fileCount, outFile)
	}
}
