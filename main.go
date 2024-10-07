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

func main() {
	// project directory to scan
	dir := ""
	// FunctionDupe(dir)
	VariableDupe(dir)
}

func FunctionDupe(dir string) {
	// Create a new FileSet to track file positions
	fset := token.NewFileSet()

	// Store all variable and function declarations
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
			funcName := funcDecl.Name.Name

			if ok && funcName != "EXPECT" {

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
					fmt.Printf("Duplicate function content:'%s' \n%s\n%s\n =================================== \n", funcName, prevFile, path)
				} else {
					functionBodies[funcBodyStr] = path
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

func VariableDupe(dir string) {
	// Create a new FileSet to track file positions
	fset := token.NewFileSet()

	// Store variables (name + value) and their corresponding file
	variables := make(map[string]string)

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

		// Inspect the AST for variable declarations
		ast.Inspect(node, func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			if !ok {
				return true
			}

			for _, spec := range decl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				for _, name := range valueSpec.Names {
					if name.Name == "err" ||
						name.Name == "_" ||
						name.Name == "id" ||
						name.Name == "res" ||
						name.Name == "resp" ||
						name.Name == "secret" ||
						name.Name == "respObj" {
						continue
					}
					// Convert the variable value to a string for comparison
					var valueBuf bytes.Buffer
					if len(valueSpec.Values) > 0 {
						err := printer.Fprint(&valueBuf, fset, valueSpec.Values[0])
						if err != nil {
							fmt.Printf("Error printing variable value in %s: %v\n", path, err)
							continue
						}
					}

					varName := name.Name
					varValue := strings.TrimSpace(valueBuf.String())

					// Key for comparison: variable name + value
					key := varName + ":" + varValue

					if prevFile, exists := variables[key]; exists {
						// Duplicate variable name with the same value
						fmt.Printf("Duplicate variable: '%s' \n%s\n%s\n======================================\n", varName, prevFile, path)
					} else {
						// Store the variable (name + value) and its location
						variables[key] = path
					}
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
