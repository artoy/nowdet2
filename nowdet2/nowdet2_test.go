package nowdet2

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"

	"github.com/stretchr/testify/require"
)

func TestPosTimeNow(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected int // number of time.Now() calls expected
	}{
		{
			name: "single time.Now call",
			source: `package main
import "time"
func main() {
	now := time.Now()
	_ = now
}`,
			expected: 1,
		},
		{
			name: "multiple time.Now calls",
			source: `package main
import "time"
func main() {
	now1 := time.Now()
	now2 := time.Now()
	_ = now1
	_ = now2
}`,
			expected: 2,
		},
		{
			name: "no time.Now calls",
			source: `package main
func main() {
	x := 5
	_ = x
}`,
			expected: 0,
		},
		{
			name: "other Now function (not time.Now)",
			source: `package main
func Now() int { return 42 }
func main() {
	x := Now()
	_ = x
}`,
			expected: 0,
		},
		{
			name: "time.Now in function argument",
			source: `package main
import "time"
func doSomething(t time.Time) {}
func main() {
	doSomething(time.Now())
}`,
			expected: 1,
		},
		{
			name: "time package alias",
			source: `package main
import t "time"
func main() {
	now := t.Now()
	_ = now
}`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse source
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.source, parser.ParseComments)
			require.NoError(t, err)

			pkg := types.NewPackage("test", "")

			// Build SSA
			ssaPkg, _, err := ssautil.BuildPackage(
				&types.Config{Importer: importer.Default()},
				fset,
				pkg,
				[]*ast.File{file},
				ssa.SanityCheckFunctions,
			)
			require.NoError(t, err)

			// Get source functions
			var srcFuncs []*ssa.Function
			for _, member := range ssaPkg.Members {
				if fn, ok := member.(*ssa.Function); ok {
					srcFuncs = append(srcFuncs, fn)
				}
			}

			// Create analysis pass
			pass := &analysis.Pass{
				Analyzer: Analyzer,
				Fset:     fset,
				Files:    []*ast.File{file},
				ResultOf: map[*analysis.Analyzer]interface{}{
					buildssa.Analyzer: &buildssa.SSA{
						Pkg:      ssaPkg,
						SrcFuncs: srcFuncs,
					},
				},
			}

			// Test PosTimeNow function
			timeNows := posTimeNow(pass)
			require.Equal(t, tt.expected, len(timeNows))
		})
	}
}
