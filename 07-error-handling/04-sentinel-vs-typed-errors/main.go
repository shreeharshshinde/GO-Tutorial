package main

import (
	"errors"
	"fmt"
)

/*
============================================================
MODULE 07 — ERROR HANDLING
STEP 07.4 — SENTINEL ERRORS vs TYPED ERRORS (DEEP DIVE)
============================================================

This file answers ONE design question:

"When should I use a sentinel error, and when should I define
a custom error type?"

This decision matters a LOT in large Go projects.
Bad choices lead to:
- brittle APIs
- unreadable error handling
- impossible refactors
*/

// ==========================================================
// 1. SENTINEL ERRORS — WHAT & WHY
// ==========================================================

/*
Sentinel error:
- A package-level variable
- Compared by identity (errors.Is)
- Represents a SINGLE condition

Example:
	var ErrNotFound = errors.New("not found")
*/

var ErrNotFound = errors.New("resource not found")

func findUser(id int) error {
	if id != 42 {
		return ErrNotFound
	}
	return nil
}

/*
Sentinel errors are GOOD when:
- There are only a FEW well-known failure modes
- Caller only cares "did this happen?"
- No extra metadata is required
*/

// ==========================================================
// 2. HANDLING SENTINEL ERRORS
// ==========================================================

func handleSentinel() {
	err := findUser(1)

	if errors.Is(err, ErrNotFound) {
		fmt.Println("User does not exist. Nothing to do.")
		return
	}

	if err != nil {
		fmt.Println("Unexpected error:", err)
	}
}

// ==========================================================
// 3. LIMITATIONS OF SENTINEL ERRORS
// ==========================================================

/*
Problems with sentinel errors:
- Cannot carry extra information
- Encourage too many global variables
- Hard to extend later without breaking APIs
*/

// BAD IDEA: explosion of sentinel errors
var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

/*
This does NOT scale.
*/

// ==========================================================
// 4. TYPED ERRORS — WHAT & WHY
// ==========================================================

/*
Typed error:
- A struct implementing Error()
- Carries structured data
- Enables richer handling via errors.As
*/

type ValidationError struct {
	Field string
	Rule  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s (%s)", e.Field, e.Rule)
}

func validateUser(username string) error {
	if username == "" {
		return &ValidationError{
			Field: "username",
			Rule:  "cannot be empty",
		}
	}
	return nil
}

/*
Typed errors are GOOD when:
- You need metadata
- You expect more variants in future
- Callers need structured access
*/

// ==========================================================
// 5. HANDLING TYPED ERRORS
// ==========================================================

func handleTyped() {
	err := validateUser("")

	var vErr *ValidationError
	if errors.As(err, &vErr) {
		fmt.Printf("Invalid input: field=%s rule=%s\n", vErr.Field, vErr.Rule)
		return
	}

	if err != nil {
		fmt.Println("Other error:", err)
	}
}

// ==========================================================
// 6. API DESIGN: STABILITY OVER TIME
// ==========================================================

/*
Sentinel errors are PART OF YOUR PUBLIC API.

If you remove or change them:
- You BREAK users
- You BREAK retries
- You BREAK controllers

Typed errors are MORE FLEXIBLE:
- You can add fields
- You can add variants
- Callers remain compatible
*/

// ==========================================================
// 7. HYBRID PATTERN (USED IN KUBERNETES)
// ==========================================================

/*
Best of both worlds:
- One sentinel for category
- Typed error for details
*/

var ErrInvalidInput = errors.New("invalid input")

type DetailedValidationError struct {
	Field string
	Rule  string
}

func (e *DetailedValidationError) Error() string {
	return fmt.Sprintf("invalid input: %s (%s)", e.Field, e.Rule)
}

func createUser(username string) error {
	if username == "" {
		// Wrap typed error with sentinel
		return fmt.Errorf("%w: %w",
			ErrInvalidInput,
			&DetailedValidationError{
				Field: "username",
				Rule:  "cannot be empty",
			},
		)
	}
	return nil
}

func handleHybrid() {
	err := createUser("")

	if errors.Is(err, ErrInvalidInput) {
		fmt.Println("Request rejected due to invalid input")

		var dErr *DetailedValidationError
		if errors.As(err, &dErr) {
			fmt.Printf("Details: field=%s rule=%s\n",
				dErr.Field, dErr.Rule)
		}
		return
	}
}

// ==========================================================
// 8. COMMON MISTAKES (DO NOT DO THESE)
// ==========================================================

func badDesignExamples() {

	// ❌ Sentinel error with dynamic message
	var ErrDynamic = fmt.Errorf("failed at %v", 123)
	_ = ErrDynamic

	// ❌ Using strings instead of types
	err := errors.New("user john failed")
	_ = err

	// ❌ Exposing internal error types unintentionally
}

// ==========================================================
// 9. DECISION TABLE (MEMORIZE THIS)
// ==========================================================

/*
USE SENTINEL ERRORS WHEN:
- Small fixed set
- Boolean-style decisions
- Public API contract

USE TYPED ERRORS WHEN:
- You need metadata
- You expect growth
- You want extensibility

USE HYBRID WHEN:
- Public category + internal detail
- Kubernetes-style APIs
*/

// ==========================================================
// 10. MAIN — DEMOS
// ==========================================================

func main() {
	fmt.Println("=== Sentinel vs Typed Errors ===")

	fmt.Println("\n-- Sentinel error --")
	handleSentinel()

	fmt.Println("\n-- Typed error --")
	handleTyped()

	fmt.Println("\n-- Hybrid pattern --")
	handleHybrid()
}

/*
============================================================
KEY TAKEAWAYS (CRITICAL)
============================================================

1. Sentinel errors define API contracts
2. Typed errors define rich failure information
3. Hybrid pattern is often best
4. Bad error design is VERY hard to fix later
5. Design errors for evolution, not convenience

============================================================
KUBERNETES CONTEXT
============================================================

Kubernetes uses:
- Sentinel errors for categories (IsNotFound, IsConflict)
- Typed errors for details (StatusError, APIStatus)

If you understand THIS file,
you can design stable Go APIs.
*/
