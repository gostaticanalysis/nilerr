package main

import (
	"github.com/gostaticanalysis/nilerr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(nilerr.Analyzer) }
