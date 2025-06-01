package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/nowdet2/nowdet2"
)

func main() {
	singlechecker.Main(
		nowdet2.Analyzer,
	)
}
