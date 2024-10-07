package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func FunctionDupe1() {
	// Set the directory to scan (current directory in this case)
	dir := "/Users/aryadhira/Documents/Works/Brankas/direct/directsvc/services"

	// Create a new FileSet to track file positions
	fset := token.NewFileSet()

	// Store all variable and function declarations
	functions := make(map[string]string)      // Function name map (function name -> file)
	functionBodies := make(map[string]string) // Function content map (content -> file)

	// Walk through all Go files in the directory and its subdirectories
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Parse the Go file
		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			fmt.Printf("Error parsing file %s: %v\n", path, err)
			return err
		}

		// Inspect the AST for variable and function declarations
		ast.Inspect(node, func(n ast.Node) bool {

			// Check for function declarations
			funcDecl, ok := n.(*ast.FuncDecl)
			if ok {
				funcName := funcDecl.Name.Name

				// Convert the function body into a string for comparison
				var funcBodyBuf bytes.Buffer
				err := printer.Fprint(&funcBodyBuf, fset, funcDecl.Body)
				if err != nil {
					fmt.Printf("Error printing function body in %s: %v\n", path, err)
					return true
				}
				funcBodyStr := strings.TrimSpace(funcBodyBuf.String())

				// Check if a function with the same body has been encountered before
				if prevFile, exists := functionBodies[funcBodyStr]; exists {
					fmt.Printf("Duplicate function content found in %s and %s\n", prevFile, path)
				} else {
					functionBodies[funcBodyStr] = path
				}

				// Also check for duplicate function names
				if prevFile, exists := functions[funcName]; exists {
					fmt.Printf("Duplicate function name '%s' found in %s and %s\n", funcName, prevFile, path)
				} else {
					functions[funcName] = path
				}
			}

			return true
		})

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", dir, err)
		return
	}
}
