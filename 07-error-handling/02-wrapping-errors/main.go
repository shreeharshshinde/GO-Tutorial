package main

import (
	"errors"
	"fmt"
)

/*
============================================================
MODULE 07 — ERROR HANDLING
STEP 07.2 — WRAPPING ERRORS (DEEP DIVE)
============================================================

Problem:
We want to add context to errors as they move up the stack.

But we must NOT lose:
- the original error
- the ability to compare errors
- the ability to classify errors

Naive string-based errors break systems.
Go solves this with ERROR WRAPPING.
*/

// ==========================================================
// 1. THE PROBLEM: STRING-BASED ERROR CHAINS (BROKEN)
// ==========================================================

var ErrDiskFull = errors.New("disk full")

func writeToDiskBroken() error {
	return ErrDiskFull
}

func saveFileBroken() error {
	err := writeToDiskBroken()
	if err != nil {
		// ❌ BROKEN: original error is now just a string
		return fmt.Errorf("save failed: %v", err)
	}
	return nil
}

func brokenComparison() {
	err := saveFileBroken()

	// This FAILS because ErrDiskFull is lost
	if err == ErrDiskFull {
		fmt.Println("disk is full (this will NEVER run)")
	}
}

/*
Why this is broken:
- err.Error() is just a string
- identity is lost
- retries / classification become impossible
*/

// ==========================================================
// 2. THE SOLUTION: ERROR WRAPPING WITH %w
// ==========================================================

func writeToDisk() error {
	return ErrDiskFull
}

func saveFile() error {
	err := writeToDisk()
	if err != nil {
		// ✅ CORRECT: wraps the error
		return fmt.Errorf("save failed: %w", err)
	}
	return nil
}

/*
Key rule:
- Use %w exactly ONCE per error
- %w preserves the original error
*/

// ==========================================================
// 3. WHAT WRAPPING ACTUALLY DOES
// ==========================================================

func inspectWrappedError() {
	err := saveFile()

	fmt.Println("Full error message:")
	fmt.Println(" ", err)

	/*
	The message becomes:
	"save failed: disk full"

	But internally, err now contains:
	- outer message
	- pointer to inner error (ErrDiskFull)
	*/
}

// ==========================================================
// 4. ERROR CHAINS (MULTI-LAYER WRAPPING)
// ==========================================================

func apiHandler() error {
	if err := saveFile(); err != nil {
		return fmt.Errorf("api handler failed: %w", err)
	}
	return nil
}

/*
Error chain looks like:

api handler failed
 └── save failed
     └── disk full
*/

// ==========================================================
// 5. WHY WRAPPING IS REQUIRED FOR LARGE SYSTEMS ?
// ==========================================================

/*
In distributed systems:
- Low-level errors know WHAT failed
- High-level errors know WHERE it failed

Wrapping preserves BOTH.
*/

// ==========================================================
// 6. DO NOT DOUBLE-WRAP THE SAME ERROR
// ==========================================================

func badDoubleWrap() error {
	err := writeToDisk()
	if err != nil {
		// ❌ WRONG: pointless extra layer
		return fmt.Errorf("%w", err)
	}
	return nil
}

/*
Rule:
- Add MEANING, not noise
- Each wrap should add context
*/

// ==========================================================
// 7. SENTINEL ERRORS + WRAPPING (PREVIEW)
// ==========================================================

func handleError() {
	err := apiHandler()

	/*
	We CANNOT compare directly:
		err == ErrDiskFull  ❌
	But we CAN unwrap (next step).
	*/

	fmt.Println("Error received:", err)
}

// ==========================================================
// 8. PANIC WARNING
// ==========================================================

/*
Wrapping errors does NOT stop panics.
Panics should be:
- rare
- programmer-only
- unrecoverable

Never wrap panics as errors.
*/

// ==========================================================
// 9. MAIN — DEMO
// ==========================================================

func main() {
	fmt.Println("=== Error Wrapping Demo ===")

	fmt.Println("\n-- Broken wrapping --")
	brokenComparison()

	fmt.Println("\n-- Correct wrapping --")
	err := apiHandler()
	fmt.Println("Returned error:", err)

	fmt.Println("\n-- Inspect wrapped error --")
	inspectWrappedError()

	fmt.Println("\nProgram continues normally.")
}

/*
============================================================
KEY TAKEAWAYS
============================================================

1. Never lose the original error
2. Use fmt.Errorf(... %w ...)
3. Wrapping builds an error chain
4. One wrap = one layer of meaning
5. Strings alone are insufficient
6. Wrapping enables classification and retries

============================================================
KUBERNETES CONTEXT
============================================================

Kubernetes relies on wrapping for:
- Retry decisions
- Error classification
- Observability
- Root cause analysis

If you don’t wrap errors properly:
- controllers retry incorrectly
- failures become opaque
- bugs become impossible to diagnose
*/