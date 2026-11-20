package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/gostaticanalysis/nilerr"
)

func main() { singlechecker.Main(nilerr.Analyzer) }
