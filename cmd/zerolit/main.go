package main

import (
	"github.com/gostaticanalysis/zerolit"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(zerolit.Analyzer) }
