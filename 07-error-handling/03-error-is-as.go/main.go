package main

import (
	"errors"
	"fmt"
	"os"
)

/*
============================================================
MODULE 07 — ERROR HANDLING
STEP 07.3 — errors.Is & errors.As (DEEP DIVE)
============================================================

This file explains HOW to:
- Detect specific errors inside wrapped error chains
- Branch logic safely (retry / abort / ignore)
- Extract structured error types

errors.Is  → answers: "Did this error occur?"
errors.As  → answers: "What TYPE of error is this?"

These functions are used heavily in Kubernetes.
*/

// ==========================================================
// 1. SENTINEL ERRORS (RECAP)
// ==========================================================

// Sentinel error: a shared, comparable value
var ErrPermissionDenied = errors.New("permission denied")

func readConfig() error {
	return ErrPermissionDenied
}

func loadApplication() error {
	if err := readConfig(); err != nil {
		// Wrap but DO NOT lose identity
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}

// ==========================================================
// 2. errors.Is — CHECKING ERROR IDENTITY
// ==========================================================

func checkWithErrorsIs() {
	err := loadApplication()

	/*
	errors.Is:
	- Walks the error chain
	- Compares against target
	- Works EVEN THROUGH WRAPPING
	*/

	if errors.Is(err, ErrPermissionDenied) {
		fmt.Println("Access denied: ask for correct permissions")
		return
	}

	fmt.Println("Unhandled error:", err)
}

/*
IMPORTANT:
NEVER do this:

	if err == ErrPermissionDenied ❌

This FAILS once wrapping is involved.
*/

// ==========================================================
// 3. REAL-WORLD EXAMPLE: OS ERRORS
// ==========================================================

func openImportantFile() error {
	_, err := os.Open("/path/that/does/not/exist")
	if err != nil {
		return fmt.Errorf("cannot open important file: %w", err)
	}
	return nil
}

func handleFileError() {
	err := openImportantFile()

	/*
	os.ErrNotExist is a sentinel error provided by Go.
	errors.Is works across wrapped OS errors.
	*/

	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("File does not exist. Creating default config.")
		return
	}

	fmt.Println("Unexpected error:", err)
}

// ==========================================================
// 4. errors.As — EXTRACTING ERROR TYPES
// ==========================================================

/*
Sometimes identity is NOT enough.
We need structured information.

errors.As allows type-based extraction.
*/

// Custom typed error
type ValidationError struct {
	Field string
	Issue string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s (%s)", e.Field, e.Issue)
}

func validateInput(input string) error {
	if input == "" {
		return &ValidationError{
			Field: "username",
			Issue: "cannot be empty",
		}
	}
	return nil
}

func processRequest() error {
	if err := validateInput(""); err != nil {
		return fmt.Errorf("request rejected: %w", err)
	}
	return nil
}

func handleValidationError() {
	err := processRequest()

	var vErr *ValidationError

	/*
	errors.As:
	- Walks the error chain
	- Finds FIRST error assignable to target type
	- Assigns it to vErr
	*/

	if errors.As(err, &vErr) {
		fmt.Printf("Invalid field '%s': %s\n", vErr.Field, vErr.Issue)
		return
	}

	fmt.Println("Other error:", err)
}

// ==========================================================
// 5. Is vs As — WHEN TO USE WHICH
// ==========================================================

/*
errors.Is:
- Use for SENTINEL ERRORS
- Retry / abort decisions
- Feature gating

errors.As:
- Use for TYPED ERRORS
- Extract metadata
- Detailed handling
*/

// ==========================================================
// 6. COMMON ANTI-PATTERNS (DO NOT DO THIS)
// ==========================================================

func badErrorHandling(err error) {

	// ❌ String comparison (fragile)
	if err != nil && err.Error() == "permission denied" {
		fmt.Println("Do not do this")
	}

	// ❌ Type assertion without As
	if _, ok := err.(*ValidationError); ok {
		fmt.Println("This breaks with wrapping")
	}
}

/*
Both fail when:
- errors are wrapped
- error messages change
*/

// ==========================================================
// 7. MULTI-LAYER ERROR CHAINS
// ==========================================================

func deepErrorChain() error {
	if err := loadApplication(); err != nil {
		return fmt.Errorf("startup failed: %w", err)
	}
	return nil
}

func inspectDeepChain() {
	err := deepErrorChain()

	if errors.Is(err, ErrPermissionDenied) {
		fmt.Println("Root cause detected deep in chain")
	}
}

// ==========================================================
// 8. KUBERNETES-STYLE ERROR HANDLING
// ==========================================================

/*
Kubernetes patterns:

if errors.Is(err, context.DeadlineExceeded) {
	// retry
}

if apierrors.IsNotFound(err) {
	// resource deleted, ignore
}

if errors.As(err, &net.OpError{}) {
	// network issue
}
*/

// ==========================================================
// 9. MAIN — DEMOS
// ==========================================================

func main() {
	fmt.Println("=== errors.Is & errors.As Demo ===")

	fmt.Println("\n-- errors.Is (sentinel) --")
	checkWithErrorsIs()

	fmt.Println("\n-- OS error handling --")
	handleFileError()

	fmt.Println("\n-- errors.As (typed error) --")
	handleValidationError()

	fmt.Println("\n-- Deep chain inspection --")
	inspectDeepChain()
}

/*
============================================================
KEY TAKEAWAYS (READ CAREFULLY)
============================================================

1. Always wrap errors with %w
2. Use errors.Is for sentinel errors
3. Use errors.As for typed errors
4. NEVER compare error strings
5. NEVER rely on direct equality with wrapped errors
6. Error handling must survive refactoring

============================================================
KUBERNETES CONTEXT
============================================================

Kubernetes controllers rely on errors.Is / errors.As to:
- Decide retries
- Ignore expected states
- Escalate fatal failures
- Preserve root causes

If you master THIS file,
you can read Kubernetes error handling confidently.
*/
