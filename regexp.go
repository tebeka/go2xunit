package main

const (
	version = "1.1.1"

	// gotest regular expressions

	// === RUN TestAdd
	gt_startRE = "^=== RUN:? ([a-zA-Z_][^[:space:]]*)"

	// --- PASS: TestSub (0.00 seconds)
	// --- FAIL: TestSubFail (0.00 seconds)
	// --- SKIP: TestSubSkip (0.00 seconds)
	gt_endRE = "^--- (PASS|FAIL|SKIP): ([a-zA-Z_][^[:space:]]*) \\((\\d+(.\\d+)?)"

	// FAIL	_/home/miki/Projects/goroot/src/xunit	0.004s
	// ok  	_/home/miki/Projects/goroot/src/anotherTest	0.000s
	gt_suiteRE = "^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(\\d+.\\d+)"

	// ?       alipay  [no test files]
	gt_noFiles = "^\\?.*\\[no test files\\]$"
	// FAIL    node/config [build failed]
	gt_buildFailed = `^FAIL.*\[(build|setup) failed\]$`

	// gocheck regular expressions

	// START: mmath_test.go:16: MySuite.TestAdd
	gc_startRE = "START: [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)"
	// PASS: mmath_test.go:16: MySuite.TestAdd	0.000s
	// FAIL: mmath_test.go:35: MySuite.TestDiv
	gc_endRE = "(PASS|FAIL): [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)([[:space:]]+([0-9]+.[0-9]+))?"
)
