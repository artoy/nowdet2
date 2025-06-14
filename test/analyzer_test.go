package test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/gostaticanalysis/testutil"
	"github.com/nowdet2/nowdet2"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	vers := testutil.LatestVersion(t, "cloud.google.com/go/spanner", 1)
	testutil.RunWithVersions(t, testdata, nowdet2.Analyzer, vers, "server")
}

func TestAnalyzerSimple(t *testing.T) {
	testdata := analysistest.TestData()
	vers := testutil.LatestVersion(t, "cloud.google.com/go/spanner", 1)
	testutil.RunWithVersions(t, testdata, nowdet2.Analyzer, vers, "simple")
}
