// the OSExitCheck package is a static analyzer that
// detects the use of os.Exit in the main function.

package lint

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// OSExitCheck - exported var for analuzer usage
var OSExitCheck = &analysis.Analyzer{
	Name: "mainexitcheck",
	Doc:  "check for using os.Exit in main function",
	Run:  run,
}

// isOsExitCalling - checking for  os.Exit calls .
func isOsExitCalling(pass *analysis.Pass, call *ast.CallExpr) bool {
	// checking function to contain 2 parts - os and Exit
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		// checking mathod name  — "Exit"
		if sel.Sel.Name == "Exit" {
			// checking identifier — "os"
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "os" {
				// checking os package imported.
				for _, imp := range pass.Pkg.Imports() {
					if imp.Path() == "os" {
						return true
					}
				}
			}
		}
	}
	return false
}

// run -main anamlyze fubnc to call.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// пропускаю файлы кэша, чтобы анализировать только исходные файлы
		filename := pass.Fset.Position(file.Pos()).Filename
		if strings.Contains(filename, "/.cache/go-build") {
			continue
		}

		// checking AST nodes.
		ast.Inspect(file, func(node ast.Node) bool {
			// checking node type is function.
			if fn, ok := node.(*ast.FuncDecl); ok {
				// Проверяем, что это функция main
				if fn.Name.Name == "main" {
					// checking function body.
					for _, stmt := range fn.Body.List {
						// checking for  os.Exit call.
						if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
							if call, ok := exprStmt.X.(*ast.CallExpr); ok && isOsExitCalling(pass, call) {
								pass.Reportf(call.Pos(), "using os.Exit in main function is prohibited")
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
