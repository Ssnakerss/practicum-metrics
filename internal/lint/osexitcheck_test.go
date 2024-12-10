package lint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMyAnalyzer(t *testing.T) {
	// // the analysistest function.Run applies the OSExitCheck analyzer under test
	// to packages from the testdata folder and checks expectations
	// ./... — checking all subdirectories in testdata
	// you can specify ./pkg1 to check only pkg1
	analysistest.Run(t, analysistest.TestData(), OSExitCheck, "./...")
}
