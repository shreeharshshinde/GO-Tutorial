package testingbasics

import (
	"errors"
	"testing"
)

/*
============================================================
MODULE 09 — TESTING IN GO
STEP 09.1 — TESTING BASICS (COMPLETE DEEP DIVE)
============================================================

This file is EXECUTABLE DOCUMENTATION.

Goal:
Understand how Go testing REALLY works,
how CNCF projects structure tests,
and how to read/write tests without fear.

If you understand THIS file,
you can read Kubernetes tests confidently.
*/

/*
------------------------------------------------------------
1. WHAT IS A GO TEST?
------------------------------------------------------------

A Go test is:
- a normal Go function
- whose name starts with Test
- takes *testing.T as input
- lives in a *_test.go file

The Go tool discovers and runs it automatically.
*/

// ==========================================================
// 2. CODE UNDER TEST
// ==========================================================

// Add is a simple function we will test.
func Add(a, b int) int {
	return a + b
}

// Divide demonstrates error handling in tests.
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

/*
------------------------------------------------------------
3. BASIC TEST FUNCTION
------------------------------------------------------------

- Use t.Errorf / t.Fatalf to fail tests
- A failed test does NOT crash the program
- Each test is isolated
*/

func TestAdd(t *testing.T) {
	result := Add(2, 3)

	if result != 5 {
		t.Errorf("Add(2,3) = %d; want 5", result)
	}
}

/*
------------------------------------------------------------
4. TESTING ERRORS
------------------------------------------------------------

Testing errors is EXPLICIT.
You must check both:
- the error value
- the returned result
*/

func TestDivideSuccess(t *testing.T) {
	result, err := Divide(10, 2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 5 {
		t.Errorf("Divide(10,2) = %d; want 5", result)
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Divide(10, 0)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

/*
------------------------------------------------------------
5. FAIL FAST vs CONTINUE
------------------------------------------------------------

t.Fatalf:
- stops THIS test immediately

t.Errorf:
- records failure
- continues execution

Rule of thumb:
- Fatal when future checks depend on this
- Error when independent
*/

func TestFailFastExample(t *testing.T) {
	result := Add(1, 1)

	if result != 2 {
		t.Fatalf("unexpected result")
	}

	// Safe to continue only if above passed
	if result < 0 {
		t.Errorf("result should not be negative")
	}
}

/*
------------------------------------------------------------
6. TEST OUTPUT & VERBOSITY
------------------------------------------------------------

- Tests are silent by default
- Use -v to see output
- t.Log prints ONLY in verbose mode

Example:
	go test -v
*/

func TestLogging(t *testing.T) {
	t.Log("This log is visible only with -v")
}

/*
------------------------------------------------------------
7. PACKAGE NAMING IN TESTS
------------------------------------------------------------

Two valid patterns:

1. Same package:
	package testingbasics

   - Access unexported identifiers
   - Used heavily in CNCF internal tests

2. External package:
	package testingbasics_test

   - Tests public API only
   - Used for API guarantees

Kubernetes uses BOTH.
*/

/*
------------------------------------------------------------
8. HOW CNCF PROJECTS THINK ABOUT TESTS
------------------------------------------------------------

In CNCF:
- Tests are FIRST-CLASS code
- Tests document behavior
- Tests prevent regressions
- Tests are required for PRs

A PR without tests is usually rejected.
*/

/*
------------------------------------------------------------
9. COMMON BEGINNER MISTAKES
------------------------------------------------------------

❌ Using panic instead of t.Fatal
❌ Printing instead of asserting
❌ Ignoring errors
❌ Writing tests after code
❌ Testing implementation instead of behavior
*/

/*
------------------------------------------------------------
10. WHAT YOU SHOULD FEEL COMFORTABLE WITH NOW
------------------------------------------------------------

After this file, you should be comfortable with:
- reading *_test.go files
- understanding t.Errorf vs t.Fatalf
- testing success paths
- testing error paths
- running go test confidently
*/

