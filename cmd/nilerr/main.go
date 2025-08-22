package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/gostaticanalysis/nilerr"
)

func main() { unitchecker.Main(nilerr.Analyzer) }
