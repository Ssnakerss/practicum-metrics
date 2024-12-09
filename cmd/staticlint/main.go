// This package defines the main function for an analysis driver
// with several analyzers from packages go/analysis and staticcheck.io
// Usage:
// `go run cmd/staticlint/main.go <analyzers> <files>`

// Analyzers from package golang.org/x/tools/go/analysis/passes:
// printf - enable printf analysis
// structtag - enable structtag analysis
// unreachable - enable unreachable analysis for unreachable code
// shadow - enable shadow analysisfor variables declared in
// for loops, switch statements, and if statements
// unusedresult - enable unusedresult analysis for unused results

// Analyzers from package staticcheck.io:
// SA* - see docs (https://staticcheck.dev/docs/checks/#SA)

// Custom analyzers
// OSExitCheck - check for os.Exit from main functions of package main

// Examples:
// Checking all files in the current folder with all analyzers
// go run cmd/staticlint/main.go ./...

package main

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/staticcheck"

	"github.com/Ssnakerss/practicum-metrics/internal/lint"
)

// setupAnalyzers - функция для конфигурации анализаторов.
func setupAnalyzers() []*analysis.Analyzer {
	mychecks := []*analysis.Analyzer{
		lint.OSExitCheck,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		unreachable.Analyzer,
		unusedresult.Analyzer,
	}

	//добавляем все SA анализаторы из staticcheck
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	return mychecks
}

func main() {

	multichecker.Main(
		setupAnalyzers()...,
	)
}
