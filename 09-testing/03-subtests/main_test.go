package subtests

import (
	"errors"
	"testing"
)

/*
============================================================
MODULE 09 — TESTING IN GO
STEP 09.3 — SUBTESTS (t.Run) — COMPLETE DEEP DIVE
============================================================

This file is EXECUTABLE DOCUMENTATION.

Goal:
Understand WHAT subtests are,
WHY they exist,
HOW they work internally,
and WHY CNCF projects rely on them everywhere.

If you understand THIS file,
Kubernetes test output will finally make sense.
*/

/*
------------------------------------------------------------
1. THE PROBLEM WITH PLAIN TABLE-DRIVEN TESTS
------------------------------------------------------------

In 09.2, we had table-driven tests like:

	for _, tt := range tests {
		if got != want {
			t.Errorf(...)
		}
	}

Problem:
- All failures belong to ONE test function
- Hard to isolate failures
- Hard to rerun a single case
- Output becomes noisy with large tables

Subtests SOLVE this.
*/

// ==========================================================
// 2. CODE UNDER TEST
// ==========================================================

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

/*
------------------------------------------------------------
3. WHAT IS A SUBTEST?
------------------------------------------------------------

A subtest is:
- A test inside a test
- Created using t.Run(name, func(t *testing.T))

Each subtest:
- Has its OWN *testing.T
- Can PASS or FAIL independently
- Shows up separately in test output
*/

/*
------------------------------------------------------------
4. BASIC SUBTEST STRUCTURE
------------------------------------------------------------
*/

func TestDivide_WithSubtests(t *testing.T) {
	tests := []struct {
		name        string
		a, b        int
		expected    int
		expectError bool
	}{
		{
			name:        "normal division",
			a:           10,
			b:           2,
			expected:    5,
			expectError: false,
		},
		{
			name:        "division by zero",
			a:           10,
			b:           0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		// IMPORTANT:
		// We pass tt into t.Run to avoid closure issues.
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Fatalf(
					"Divide(%d,%d) = %d; want %d",
					tt.a, tt.b, result, tt.expected,
				)
			}
		})
	}
}

/*
------------------------------------------------------------
5. WHY t.Run CHANGES EVERYTHING
------------------------------------------------------------

With subtests:
- Each test case is isolated
- Failures are clearly labeled
- You can run a SINGLE test case
- Output is structured hierarchically

Example:

	go test -v

Output:

	TestDivide_WithSubtests
	 ├── normal division
	 └── division by zero
*/

/*
------------------------------------------------------------
6. FAILING A SUBTEST vs FAILING THE PARENT
------------------------------------------------------------

- t.Fatalf inside subtest:
	→ fails ONLY that subtest

- t.Fatalf in parent test:
	→ stops the entire test function

This granularity is critical at scale.
*/

/*
------------------------------------------------------------
7. SUBTESTS ARE REQUIRED FOR PARALLELISM
------------------------------------------------------------

You CANNOT safely parallelize plain table loops.

Subtests allow:

	t.Run(name, func(t *testing.T) {
		t.Parallel()
	})

We will deep dive this in 09.6.
*/

/*
------------------------------------------------------------
8. COMMON FOOTGUN: LOOP VARIABLE CAPTURE
------------------------------------------------------------

This line is NOT optional:

	tt := tt

Why?
- Closures capture variables by reference
- Loop variable changes each iteration

Without this:
- Tests may read wrong data
- Especially dangerous with t.Parallel

Even though Go 1.22 fixes some cases,
CNCF codebases STILL REQUIRE this pattern.
*/

/*
------------------------------------------------------------
9. WHY KUBERNETES USES SUBTESTS EVERYWHERE
------------------------------------------------------------

In Kubernetes:
- Tests often have 50–200 cases
- Subtests allow:
	- selective execution
	- readable failures
	- parallelism
	- structured logs

Example usage:
- API validation tests
- Admission webhook tests
- Scheduler predicate tests
*/

/*
------------------------------------------------------------
10. RUNNING SUBTESTS SELECTIVELY
------------------------------------------------------------

You can run ONE subtest:

	go test -run TestDivide_WithSubtests/normal

This is incredibly powerful when debugging.
*/

/*
------------------------------------------------------------
11. SUBTESTS vs MULTIPLE TEST FUNCTIONS
------------------------------------------------------------

Multiple TestX functions:
- Good for different behaviors

Subtests:
- Good for variations of SAME behavior

Kubernetes uses BOTH.
*/

/*
------------------------------------------------------------
12. COMMON BEGINNER MISTAKES
------------------------------------------------------------

❌ Forgetting tt := tt
❌ Using t.Errorf instead of t.Fatalf blindly
❌ Nesting too deeply
❌ Overusing subtests for unrelated behavior
❌ Parallelizing without understanding races
*/

/*
------------------------------------------------------------
13. WHAT YOU SHOULD BE COMFORTABLE WITH NOW
------------------------------------------------------------

After this file, you should:
- Recognize t.Run instantly
- Understand subtest output
- Know when to use subtests
- Avoid loop capture bugs
- Read Kubernetes test files confidently
*/
