package nilerr_test

import (
	"testing"

	"github.com/gostaticanalysis/nilerr"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	tests := []struct {
		name string
		file string
	}{
		{"basic", "a"},
		{"defers", "defers"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			analysistest.Run(t, testdata, nilerr.Analyzer, test.file)
		})
	}
}
