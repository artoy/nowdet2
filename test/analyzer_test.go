package test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/nowdet2/nowdet2"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nowdet2.Analyzer, "server")
}
