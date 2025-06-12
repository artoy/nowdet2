package nowdet2

import (
	"slices"

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

func run(pass *analysis.Pass) (any, error) {
	timeNows := posTimeNow(pass)
	for _, timeNow := range timeNows {
		walkToDetectSpannerFunc(pass, timeNow)
	}
	return nil, nil
}

// posTimeNow returns the instructions that call time.Now
func posTimeNow(pass *analysis.Pass) []*ssa.Call {
	timeNows := make([]*ssa.Call, 0)

	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, f := range funcs {
		for _, block := range f.Blocks {
			for _, instr := range block.Instrs {
				if call, ok := instr.(*ssa.Call); ok {
					if fn, ok := call.Call.Value.(*ssa.Function); ok {
						// Detect time.Now()
						if fn.Pkg != nil && fn.Pkg.Pkg.Path() == "time" && fn.Name() == "Now" {
							// Accumulate the variable that have value from time.Now()
							timeNows = append(timeNows, call)
						}
					}
				}
			}
		}
	}

	return timeNows
}

var checked []ssa.Instruction

// walkToDetectSpannerFunc walks through the SSA graph to detect Spanner functions that use a value from time.Now
func walkToDetectSpannerFunc(pass *analysis.Pass, instr ssa.Instruction) {
	if slices.Contains(checked, instr) {
		return
	}
	checked = append(checked, instr)

	// Check whether the instruction is a Spanner function that uses a value of time.Now as an argument and walk through the referrers
	switch v := instr.(type) {
	case *ssa.Call:
		if fn, ok := v.Call.Value.(*ssa.Function); ok {
			// Check if this is a Spanner function
			if isSpannerFunction(fn) {
				pass.Reportf(v.Pos(), "%s may use an argument that is a value from time.Now()", fn.String())
			}

			// Walk through the referrers
			for _, referrer := range *v.Referrers() {
				walkToDetectSpannerFunc(pass, referrer)
			}
		}
	case *ssa.Phi:
		// In case of Phi, it is not certain that the value is used in the function but may be used.
		for _, referrer := range *v.Referrers() {
			walkToDetectSpannerFunc(pass, referrer)
		}
	case *ssa.BinOp:
		// In case of BinOp, it is not certain that the value is used in the function but may be used.
		for _, referrer := range *v.Referrers() {
			walkToDetectSpannerFunc(pass, referrer)
		}
	// Walk pointed value to pointer when the instruction relates to a pointer
	case *ssa.Store:
		// TODO: Type assertion or else branch may be changed when we analyze across packages.
		if addr, ok := v.Addr.(ssa.Instruction); ok {
			walkToDetectSpannerFunc(pass, addr)
		} else {
			return
		}
	case *ssa.FieldAddr:
		for _, referrer := range *v.X.Referrers() {
			walkToDetectSpannerFunc(pass, referrer)
		}
	case *ssa.IndexAddr:
		for _, referrer := range *v.X.Referrers() {
			walkToDetectSpannerFunc(pass, referrer)
		}
	}
}

// isSpannerFunction checks if the given function is a Spanner-related function
func isSpannerFunction(fn *ssa.Function) bool {
	if fn.Pkg == nil {
		return false
	}

	pkgPath := fn.Pkg.Pkg.Path()

	// Check for common Spanner package paths
	spannerPkgs := []string{
		"cloud.google.com/go/spanner",
	}

	return slices.Contains(spannerPkgs, pkgPath)
}
