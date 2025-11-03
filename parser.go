package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

func loadRandomFunction(cfg config) (string, string, error) {
	folderPath := cfg.FolderPath

	// Expand home directory if needed
	if strings.HasPrefix(folderPath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", "", err
		}
		folderPath = filepath.Join(homeDir, folderPath[1:])
	}

	var goFiles []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})

	if err != nil {
		return "", "", err
	}

	if len(goFiles) == 0 {
		return "", "", fmt.Errorf("no Go files found in %s", folderPath)
	}

	// Try to find a suitable function
	maxAttempts := len(goFiles) * 3
	for attempt := 0; attempt < maxAttempts; attempt++ {
		randomFile := goFiles[rand.Intn(len(goFiles))]
		functions, err := extractFunctions(randomFile)
		if err != nil {
			continue
		}

		// Filter functions based on line count
		var validFunctions []string
		for _, fn := range functions {
			lines := countLines(fn)
			if lines >= cfg.MinLines && lines <= cfg.MaxLines {
				validFunctions = append(validFunctions, fn)
			}
		}

		if len(validFunctions) > 0 {
			return validFunctions[rand.Intn(len(validFunctions))], randomFile, nil
		}
	}

	return "", "", fmt.Errorf("no functions between %d and %d lines found after %d attempts", cfg.MinLines, cfg.MaxLines, maxAttempts)
}

func extractFunctions(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var functions []string

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			start := fset.Position(fn.Pos())
			end := fset.Position(fn.End())

			// Read the source file
			content, err := os.ReadFile(filePath)
			if err != nil {
				return true
			}

			lines := strings.Split(string(content), "\n")
			if start.Line <= len(lines) && end.Line <= len(lines) {
				funcText := strings.Join(lines[start.Line-1:end.Line], "\n")
				functions = append(functions, funcText)
			}
		}
		return true
	})

	return functions, nil
}
