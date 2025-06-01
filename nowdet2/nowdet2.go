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
	posTimeNow(pass)
	return nil, nil
}

func posTimeNow(pass *analysis.Pass) {
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, f := range funcs {
		for _, block := range f.Blocks {
			for _, instr := range block.Instrs {
				if call, ok := instr.(*ssa.Call); ok {
					fnName := call.Call.Value.Name()
					if fnName == "Now" {
						pass.Reportf(call.Pos(), "time.Now() should not be used")
					}
				}
			}
		}
	}
}
