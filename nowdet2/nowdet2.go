package nowdet2

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var Analyzer = &analysis.Analyzer{
	Name: "nowdet2",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

const Doc = "nowdet2 is a static analysis tool that detects time.now() in arguments of functions about Spanner."

func run(pass *analysis.Pass) (interface{}, error) {
	timeNows := posTimeNow(pass)
	for _, instr := range timeNows {
		pass.Reportf(instr.Pos(), "time.Now() should not be used")
	}
	return nil, nil
}

// posTimeNow returns the instructions that call time.Now()
func posTimeNow(pass *analysis.Pass) []ssa.Call {
	timeNows := make([]ssa.Call, 0)

	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, f := range funcs {
		for _, block := range f.Blocks {
			for _, instr := range block.Instrs {
				if call, ok := instr.(*ssa.Call); ok {
					if fn, ok := call.Call.Value.(*ssa.Function); ok {
						// Detect time.Now()
						if fn.Pkg != nil && fn.Pkg.Pkg.Path() == "time" && fn.Name() == "Now" {
							// Accumulate the variable that have value from time.Now()
							timeNows = append(timeNows, *call)
						}
					}
				}
			}
		}
	}

	return timeNows
}
