package main

import (
	"errors"
	"fmt"
)

/*
============================================================
MODULE 07 — ERROR HANDLING IN GO
STEP 07.1 — ERROR BASICS (DEEP DIVE)
============================================================

This file answers ONE core question:

"Why does Go handle errors this way, and why does it scale?"

Key idea:
Errors in Go are VALUES, not control flow.

No exceptions.
No try/catch.
No hidden jumps.

Everything is explicit.
*/

// ==========================================================
// 1. WHAT IS AN ERROR IN GO?
// ==========================================================

/*
In Go, error is just an interface:

	type error interface {
		Error() string
	}

That’s it.
No magic.
*/

func basicErrorExample() error {
	return errors.New("something went wrong")
}

// ==========================================================
// 2. THE MOST IMPORTANT GO IDIOM
// ==========================================================

/*
This pattern appears EVERYWHERE in Go codebases:

	result, err := doSomething()
	if err != nil {
		return err
	}

This is not verbosity.
This is CONTROL.
*/

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// ==========================================================
// 3. ERROR == CONTROL FLOW (EXPLICIT)
// ==========================================================

func errorAsControlFlow() {
	result, err := divide(10, 0)
	if err != nil {
		// Caller decides what to do
		fmt.Println("Error occurred:", err)
		return
	}

	fmt.Println("Result:", result)
}

/*
IMPORTANT:
- The function does NOT decide how to recover
- The CALLER decides
- This scales in distributed systems
*/

// ==========================================================
// 4. WHY GO DOES NOT USE EXCEPTIONS
// ==========================================================

/*
Exceptions break down in large systems:

- Hidden control flow
- Unclear ownership
- Hard to reason about in distributed systems
- Catastrophic in long-running servers

Go trades convenience for correctness.
*/

// ==========================================================
// 5. ERROR CREATION PATTERNS
// ==========================================================

// Pattern 1: errors.New (static error)
var ErrNotFound = errors.New("resource not found")

// Pattern 2: fmt.Errorf (dynamic message)
func openFile(name string) error {
	return fmt.Errorf("failed to open file %s", name)
}

// ==========================================================
// 6. RETURN ERRORS EARLY (FAIL FAST)
// ==========================================================

func processPipeline() error {
	if err := step1(); err != nil {
		return err
	}

	if err := step2(); err != nil {
		return err
	}

	if err := step3(); err != nil {
		return err
	}

	return nil
}

func step1() error { return nil }
func step2() error { return errors.New("step2 failed") }
func step3() error { return nil }

/*
This looks repetitive.
It is intentional.

Why?
- Clear execution path
- Easy debugging
- Easy tracing
*/

// ==========================================================
// 7. NEVER IGNORE ERRORS
// ==========================================================

func ignoringErrorsIsABug() {
	_, err := divide(10, 0)
	_ = err // ❌ swallowing the error

	/*
	This is one of the WORST patterns.

	Kubernetes code review rule:
	"If you ignore an error, explain WHY."
	*/
}

// ==========================================================
// 8. ERROR INFORMATION SHOULD BE LOCAL
// ==========================================================

func lowLevelFunction() error {
	return errors.New("connection refused")
}

func midLevelFunction() error {
	err := lowLevelFunction()
	if err != nil {
		// Add CONTEXT, do not lose original meaning
		return fmt.Errorf("mid-level operation failed: %v", err)
	}
	return nil
}

/*
This is NOT error wrapping yet.
We will fix this later using %w.
*/

// ==========================================================
// 9. PANIC VS ERROR (VERY IMPORTANT)
// ==========================================================

/*
Use ERROR when:
- Something can go wrong at runtime
- Caller might want to handle it

Use PANIC when:
- Programmer error
- Impossible state
- Broken invariants
*/

func mustDivide(a, b int) int {
	if b == 0 {
		panic("programmer error: division by zero")
	}
	return a / b
}

/*
Kubernetes rule:
- Panics should NEVER cross API boundaries
- Panics crash processes
*/

// ==========================================================
// 10. MAIN — SEE EVERYTHING TOGETHER
// ==========================================================

func main() {
	fmt.Println("=== Error Handling Basics ===")

	// Basic error
	err := basicErrorExample()
	fmt.Println("Basic error:", err)

	// Division example
	errorAsControlFlow()

	// Pipeline example
	if err := processPipeline(); err != nil {
		fmt.Println("Pipeline failed:", err)
	}

	// Sentinel error usage
	if err := openFile("config.yaml"); err != nil {
		fmt.Println("Open file error:", err)
	}

	fmt.Println("Program continues normally.")
}

/*
============================================================
KEY TAKEAWAYS (READ THIS TWICE)
============================================================

1. Errors are VALUES
2. Errors are returned, not thrown
3. Callers decide how to handle failures
4. Fail fast, return early
5. Never ignore errors silently
6. Panic is NOT error handling
7. Explicit > convenient

============================================================
KUBERNETES CONTEXT
============================================================

In Kubernetes:
- Every API call returns error
- Controllers propagate errors upward
- Retry logic depends on error types
- Panics are reserved for bugs, not runtime issues

If you understand THIS file,
you understand 80% of Go error handling.
*/