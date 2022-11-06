package lint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// CheckOSExit Check call os.Exit from main file
//
//	func main() {
//		os.Exit(0) // Linting error
//	}
var CheckOSExit = &analysis.Analyzer{
	Name: "CheckOSExit",
	Doc:  "Doc",
	Run:  checkOSExitRun,
}

func checkOSExitRun(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if file.Name.String() == "main" {
			ast.Inspect(file, func(node ast.Node) bool {
				if fd, ok := node.(*ast.FuncDecl); ok {
					if fd.Name.String() != "main" {
						return false
					}
				}

				if expr, ok := node.(*ast.ExprStmt); ok {
					if call, ok := expr.X.(*ast.CallExpr); ok {
						if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
							if ident, ok := selector.X.(*ast.Ident); ok {
								if ident.Name == "os" && selector.Sel.Name == "Exit" {
									pass.Reportf(call.Pos(), "call os.Exit in main")
									return false
								}
							}
						}
					}
				}

				return true
			})
		}
	}

	return nil, nil
}
