package nilerr_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/gostaticanalysis/nilerr"
)

func Test(t *testing.T) {
	cases := []string{"a", "b", "c", "d", "e"}

	testdata := analysistest.TestData()
	for _, pkg := range cases {
		t.Run(pkg, func(t *testing.T) {
			analysistest.Run(t, testdata, nilerr.Analyzer, pkg)
		})
	}
}
