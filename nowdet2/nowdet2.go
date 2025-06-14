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
		walkInstructions(pass, timeNow)
	}
	return nil, nil
}

// posTimeNow returns the instructions that call time.Now
func posTimeNow(pass *analysis.Pass) []ssa.Instruction {
	timeNows := make([]ssa.Instruction, 0)

	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, f := range funcs {
		for _, block := range f.Blocks {
			for _, instr := range block.Instrs {
				// TODO: It may be better to handle the case of *ssa.MakeClosure.
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

// walkInstructions walks through the SSA graph to detect Spanner functions that use a value from time.Now
func walkInstructions(pass *analysis.Pass, instr ssa.Instruction) {
	// Check that the instruction has not been checked yet to avoid infinite recursion
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
		}

		// Walk through the referrers
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}

	// Walk pointed value to pointer when the instructions relate to a pointer, in other words, walk to the definition of right-hand side.
	case *ssa.Store:
		// TODO: Type assertion or else branch may be changed when we analyze across packages.
		if addr, ok := v.Addr.(ssa.Instruction); ok {
			walkInstructions(pass, addr)
		}
	case *ssa.FieldAddr:
		for _, referrer := range *v.X.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.IndexAddr:
		for _, referrer := range *v.X.Referrers() {
			walkInstructions(pass, referrer)
		}

	// Following cases are trivial
	case *ssa.Phi:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.BinOp:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.UnOp:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.ChangeType:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Convert:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.MultiConvert:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.ChangeInterface:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.SliceToArrayPointer:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.MakeInterface:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Field:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Lookup:
		// Of course, this is a conservative checking.
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Select:
		// TODO: Now, we cannot ensure to detect time.Now().
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Range:
		// Of course, this is a conservative checking.
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Next:
		// TODO: Now, we cannot ensure to detect time.Now().
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.TypeAssert:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Extract:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Send:
		for _, referrer := range *v.Chan.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.MapUpdate:
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	}

	// Do nothing in the case of *ssa.Alloc, MakeClosure, MakeMap, Return, RunDefers, Panic, Go, Defer, and DebugRef.
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
