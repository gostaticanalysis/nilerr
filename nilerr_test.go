package nilerr_test

import (
	"testing"

	"github.com/gostaticanalysis/nilerr"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nilerr.Analyzer, "a")
}
