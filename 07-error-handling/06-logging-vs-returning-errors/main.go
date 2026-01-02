package main

import (
	"context"
	"errors"
	"fmt"
	"log"
)

/*
============================================================
MODULE 07 — ERROR HANDLING
STEP 07.6 — LOGGING vs RETURNING ERRORS (DEEP DIVE)
============================================================

This file answers ONE critical question:

"Where should errors be logged?"

Short answer:
- LOW-LEVEL code: return errors
- HIGH-LEVEL code: log errors ONCE

Bad logging causes:
- duplicate logs
- noisy systems
- impossible debugging

Kubernetes has very strict rules here.
*/

// ==========================================================
// 1. THE CARDINAL RULE
// ==========================================================

/*
An error should be LOGGED EXACTLY ONCE.

Either:
- you HANDLE it → log it
OR
- you CANNOT handle it → return it

NEVER BOTH.
*/

// ==========================================================
// 2. BAD PATTERN: LOG AND RETURN (DO NOT DO THIS)
// ==========================================================

func readConfigBad() error {
	err := errors.New("failed to read config")

	// ❌ WRONG: logging here
	log.Println("[readConfigBad]", err)

	return err
}

func startAppBad() {
	err := readConfigBad()
	if err != nil {
		// ❌ WRONG: logging AGAIN
		log.Println("[startAppBad]", err)
	}
}

/*
Output:
[readConfigBad] failed to read config
[startAppBad] failed to read config

Same error logged twice → noisy → useless
*/

// ==========================================================
// 3. CORRECT PATTERN: RETURN ONLY
// ==========================================================

func readConfigGood() error {
	return errors.New("failed to read config")
}

func startAppGood() {
	err := readConfigGood()
	if err != nil {
		// ✅ LOG ONCE, at boundary
		log.Println("[startAppGood] startup failed:", err)
	}
}

/*
Rule:
- Low-level functions NEVER log
- They return errors upward
*/

// ==========================================================
// 4. ADD CONTEXT, NOT LOGS
// ==========================================================

func connectDB() error {
	return errors.New("connection refused")
}

func initStorage() error {
	if err := connectDB(); err != nil {
		// ✅ Add context
		return fmt.Errorf("storage init failed: %w", err)
	}
	return nil
}

func mainInit() {
	if err := initStorage(); err != nil {
		// ✅ Log once, with full context
		log.Println("[mainInit]", err)
	}
}

// ==========================================================
// 5. WHEN LOGGING IS CORRECT IN LOW-LEVEL CODE
// ==========================================================

/*
Low-level code MAY log if:
- it is intentionally swallowing the error
- it is doing best-effort cleanup
*/

func cleanupTempFiles() {
	err := errors.New("permission denied")

	// ✅ OK: error is handled & swallowed
	log.Println("[cleanup] warning:", err)
}

// ==========================================================
// 6. CONTEXT-AWARE LOGGING (CRITICAL)
// ==========================================================

type ctxKey string

const requestIDKey ctxKey = "requestID"

func logWithContext(ctx context.Context, msg string, err error) {
	reqID := ctx.Value(requestIDKey)
	log.Printf("[req=%v] %s: %v\n", reqID, msg, err)
}

func processRequest(ctx context.Context) error {
	return errors.New("timeout talking to database")
}

func handleRequest(ctx context.Context) {
	err := processRequest(ctx)
	if err != nil {
		// ✅ log WITH context
		logWithContext(ctx, "request failed", err)
	}
}

// ==========================================================
// 7. DO NOT LOG AND PANIC
// ==========================================================

func badPanic() {
	err := errors.New("corrupted state")

	// ❌ WRONG
	log.Println("fatal error:", err)
	panic(err)
}

/*
Why wrong?
- Panic already prints stack trace
- Logging causes duplication
*/

// ==========================================================
// 8. PANIC ONLY FOR PROGRAMMER ERRORS
// ==========================================================

func mustHaveConfig(cfg *string) {
	if cfg == nil {
		panic("BUG: config must not be nil")
	}
}

// ==========================================================
// 9. KUBERNETES CONTROLLER PATTERN
// ==========================================================

/*
Typical controller loop:

for {
	err := reconcile()
	if err != nil {
		log.Error(err, "reconcile failed")
		return err // triggers retry
	}
}

- reconcile() NEVER logs
- controller logs ONCE
- retry is automatic
*/

func reconcile() error {
	return errors.New("API server unavailable")
}

func controllerLoop() {
	err := reconcile()
	if err != nil {
		log.Println("[controller] reconcile failed:", err)
		// return err → requeue
	}
}

// ==========================================================
// 10. COMMON REVIEW COMMENTS (REAL WORLD)
// ==========================================================

/*
❌ "This function logs and returns error"
❌ "Please remove this log, caller already logs"
❌ "Do not log in libraries"
❌ "Add context, not logs"
❌ "This error is logged multiple times"
*/

// ==========================================================
// 11. MAIN — DEMONSTRATION
// ==========================================================

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	fmt.Println("=== Logging vs Returning Errors ===")

	fmt.Println("\n-- Bad logging --")
	startAppBad()

	fmt.Println("\n-- Correct logging --")
	startAppGood()

	fmt.Println("\n-- Context logging --")
	ctx := context.WithValue(context.Background(), requestIDKey, "req-123")
	handleRequest(ctx)

	fmt.Println("\n-- Controller pattern --")
	controllerLoop()

	fmt.Println("\nProgram finished.")
}

/*
============================================================
KEY TAKEAWAYS (MEMORIZE THESE)
============================================================

1. Log errors EXACTLY ONCE
2. Low-level code returns errors
3. High-level code logs errors
4. Never log and return the same error
5. Add context via wrapping, not logs
6. Context-aware logs are mandatory
7. Panics are NOT error handling

============================================================
KUBERNETES CONTEXT
============================================================

- Libraries never log
- Controllers log once
- Errors drive retries
- Duplicate logs are considered bugs
- Clean logs = debuggable systems

If you follow THIS file,
your PRs will not get logging-related review comments.
*/
