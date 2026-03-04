package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Simple minifier using Regex
func minify(js string) string {
	// 1. Remove multi-line comments /* ... */
	reMulti := regexp.MustCompile(`(?s)/\*.*?\*/`)
	js = reMulti.ReplaceAllString(js, "")

	// 2. Remove single-line comments // ...
	reSingle := regexp.MustCompile(`(?m)//.*$`)
	js = reSingle.ReplaceAllString(js, "")

	// 3. Remove unnecessary whitespace (tabs, newlines, extra spaces)
	// This is a "safe" minify: it replaces sequences of whitespace with a single space
	reSpace := regexp.MustCompile(`\s+`)
	js = reSpace.ReplaceAllString(js, " ")

	return strings.TrimSpace(js)
}

func main() {
	// Define Flags
	srcPtr := flag.String("src", "./src/js", "Directory containing source JS files")
	outPtr := flag.String("out", "./dist/bundle.js", "Path for the output bundled file")
	minifyPtr := flag.Bool("minify", false, "Enable basic minification")

	flag.Parse()

	srcDir := filepath.Clean(*srcPtr)
	outputFile := filepath.Clean(*outPtr)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer out.Close()

	fmt.Printf("🚀 Bundling: %s -> %s (Minify: %v)\n", srcDir, outputFile, *minifyPtr)

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".js") && path != outputFile {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var finalContent string
			if *minifyPtr {
				finalContent = minify(string(content))
			} else {
				// Add header for easier debugging when not minified
				finalContent = fmt.Sprintf("\n/* Source: %s */\n%s\n", path, string(content))
			}

			if _, err := out.WriteString(finalContent); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("✅ Bundle complete!")
	}
}
