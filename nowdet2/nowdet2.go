package nowdet2

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"

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
	// SSA dumping
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, f := range funcs {
		f.WriteTo(os.Stdout)
	}

	timeNows := posTimeNow(pass)
	for _, timeNow := range timeNows {
		// Print the register of time.Now() found
		fmt.Println("time.Now() found at:", timeNow.Name())

		walkInstructions(pass, timeNow)
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
		fmt.Printf("Already checked %s: %s\n", reflect.TypeOf(instr), instr)
		return
	}
	checked = append(checked, instr)

	// Check whether the instruction is a Spanner function that uses a value of time.Now as an argument and walk through the referrers
	switch v := instr.(type) {
	case *ssa.Call:
		fmt.Printf("Checking %s: %s\n", reflect.TypeOf(v), v)
		if fn, ok := v.Call.Value.(*ssa.Function); ok {
			// Check if this is a Spanner function
			if isSpannerFunction(fn) {
				fmt.Println("Spanner function found:", fn.Name())
				pass.Reportf(v.Pos(), "%s may use an argument that is a value from time.Now()", fn.String())
			}
		}

		// Walk through the referrers
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}

	// Walk pointed value to pointer when the instructions relate to a pointer, in other words, walk to the definition of right-hand side.
	case *ssa.Store:
		fmt.Printf("Checking %s: %s\n", reflect.TypeOf(v), v)
		// TODO: Type assertion or else branch may be changed when we analyze across packages.
		if addr, ok := v.Addr.(ssa.Instruction); ok {
			walkInstructions(pass, addr)
		}
	case *ssa.FieldAddr:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.X.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.IndexAddr:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.X.Referrers() {
			walkInstructions(pass, referrer)
		}

	// Following cases are trivial
	case *ssa.Phi:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.BinOp:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.UnOp:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.ChangeType:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Convert:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.MultiConvert:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.ChangeInterface:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.SliceToArrayPointer:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.MakeInterface:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Slice:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Field:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Lookup:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		// Of course, this is a conservative checking.
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Select:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		// TODO: Now, we cannot ensure to detect time.Now().
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Range:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		// Of course, this is a conservative checking.
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Next:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		// TODO: Now, we cannot ensure to detect time.Now().
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.TypeAssert:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Extract:
		fmt.Printf("Checking %s: %s = %s\n", reflect.TypeOf(v), v.Name(), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.Send:
		fmt.Printf("Checking %s: %s\n", reflect.TypeOf(v), v)
		for _, referrer := range *v.Chan.Referrers() {
			walkInstructions(pass, referrer)
		}
	case *ssa.MapUpdate:
		fmt.Printf("Checking %s: %s\n", reflect.TypeOf(v), v)
		for _, referrer := range *v.Referrers() {
			walkInstructions(pass, referrer)
		}
	}

	// Do nothing in the case of *ssa.Alloc, MakeClosure, MakeMap, MakeChan, MakeSlice, Return, RunDefers, Panic, Go, Defer, and DebugRef.
}

// isSpannerFunction checks if the given function is a Spanner-related function
func isSpannerFunction(fn *ssa.Function) bool {
	if fn.Pkg == nil || fn.Pkg.Pkg == nil {
		return false
	}

	pkgPath := fn.Pkg.Pkg.Path()
	return strings.Contains(pkgPath, "cloud.google.com/go/spanner")
}
