package tabledriventests

import (
	"errors"
	"testing"
)

/*
============================================================
MODULE 09 — TESTING IN GO
STEP 09.2 — TABLE-DRIVEN TESTS (COMPLETE DEEP DIVE)
============================================================

This file is EXECUTABLE DOCUMENTATION.

Goal:
Understand WHY table-driven tests exist,
HOW they are structured,
and WHY CNCF projects rely on them heavily.

If you master THIS pattern,
you can read almost any Kubernetes test.
*/

/*
------------------------------------------------------------
1. THE PROBLEM WITH NAIVE TESTS
------------------------------------------------------------

Without table-driven tests, you often see:

	func TestAdd1(t *testing.T) { ... }
	func TestAdd2(t *testing.T) { ... }
	func TestAdd3(t *testing.T) { ... }

Problems:
- Repetition
- Hard to extend
- Easy to miss edge cases
- No single "source of truth"

Table-driven tests SOLVE this.
*/

// ==========================================================
// 2. CODE UNDER TEST
// ==========================================================

// Add is intentionally simple.
func Add(a, b int) int {
	return a + b
}

// Divide demonstrates returning errors.
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

/*
------------------------------------------------------------
3. BASIC TABLE-DRIVEN TEST STRUCTURE
------------------------------------------------------------

A table-driven test has THREE parts:

1. A slice of test cases (the table)
2. A loop over the table
3. Assertions inside the loop

This pattern is UNIVERSAL in Go.
*/

func TestAdd_TableDriven(t *testing.T) {
	// Step 1: Define test cases
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "both positive",
			a:        2,
			b:        3,
			expected: 5,
		},
		{
			name:     "with zero",
			a:        0,
			b:        5,
			expected: 5,
		},
		{
			name:     "both negative",
			a:        -2,
			b:        -3,
			expected: -5,
		},
	}

	// Step 2: Iterate over test cases
	for _, tt := range tests {
		// Step 3: Run assertions
		result := Add(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf(
				"%s: Add(%d,%d) = %d; want %d",
				tt.name,
				tt.a,
				tt.b,
				result,
				tt.expected,
			)
		}
	}
}

/*
------------------------------------------------------------
4. WHY THE 'name' FIELD MATTERS
------------------------------------------------------------

The 'name' field:
- Documents intent
- Makes failures readable
- Is REQUIRED for large test tables

In CNCF code reviews:
Tests WITHOUT names are often rejected.
*/

/*
------------------------------------------------------------
5. TABLE-DRIVEN TESTS WITH ERRORS
------------------------------------------------------------

Errors MUST be tested explicitly.

We test:
- success cases
- failure cases
- expected error presence
*/

func TestDivide_TableDriven(t *testing.T) {
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
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		result, err := Divide(tt.a, tt.b)

		if tt.expectError {
			if err == nil {
				t.Errorf("%s: expected error, got nil", tt.name)
			}
			continue
		}

		if err != nil {
			t.Errorf("%s: unexpected error: %v", tt.name, err)
			continue
		}

		if result != tt.expected {
			t.Errorf(
				"%s: Divide(%d,%d) = %d; want %d",
				tt.name,
				tt.a,
				tt.b,
				result,
				tt.expected,
			)
		}
	}
}

/*
------------------------------------------------------------
6. WHY THIS SCALES SO WELL
------------------------------------------------------------

Adding a new test case is trivial:

	tests = append(tests, {...})

No new functions.
No copy-paste.
No boilerplate.

This is why Kubernetes tests often have:
- 20+
- 50+
- 100+ cases in one test
*/

/*
------------------------------------------------------------
7. TABLE-DRIVEN TESTS DOCUMENT BEHAVIOR
------------------------------------------------------------

The table:
- describes the contract
- lists edge cases
- shows intended behavior

Many maintainers READ tables
before reading implementation.
*/

/*
------------------------------------------------------------
8. COMMON BEGINNER MISTAKES
------------------------------------------------------------

❌ Writing separate tests instead of tables
❌ Forgetting descriptive names
❌ Not testing error cases
❌ Testing implementation instead of behavior
❌ Hardcoding logic inside loops
*/

/*
------------------------------------------------------------
9. CNCF EXPECTATIONS
------------------------------------------------------------

In Kubernetes:
- Table-driven tests are the DEFAULT
- Single-case tests are rare
- Tables grow over time
- Bugs get fixed by adding a new row

This makes regressions obvious.
*/

/*
------------------------------------------------------------
10. WHAT YOU SHOULD BE COMFORTABLE WITH NOW
------------------------------------------------------------

After this file, you should:
- Recognize table-driven tests instantly
- Be able to add new cases safely
- Understand test intent from tables
- Read Kubernetes test files confidently
*/
