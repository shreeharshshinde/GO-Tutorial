package main

import (
	"context"
	"fmt"
	"time"
)

/*
This file explains context.Context — the MOST IMPORTANT
library for Cloud Native Go.

In Kubernetes:
- Almost every function starts with: func(ctx context.Context, ...)
- Context controls TIME, CANCELLATION, and REQUEST-SCOPED DATA

Key ideas:
- Context is a cancellation signal
- Context propagates down the call stack
- Context prevents goroutine leaks
*/

// ==========================================================
// 1. CONTEXT KEYS (BEST PRACTICE)
// ==========================================================

/*
NEVER use built-in types (string, int) as context keys.

Why:
- Context is shared across packages
- String keys can collide

Correct pattern:
- Define an unexported key type
*/

type traceIDKeyType struct{}

var traceIDKey = traceIDKeyType{}

func main() {

	// ==========================================================
	// 1. TIMEOUTS (DEADLINE PATTERN)
	// ==========================================================
	fmt.Println("--- 1. Context with Timeout ---")

	/*
	context.WithTimeout:
	- Creates a child context
	- Automatically cancels after duration
	- ALWAYS call cancel() to release resources early
	*/

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// This job needs 5 seconds, but timeout is 2 seconds
	makeRequest(ctx, "Slow-API-Call", 5*time.Second)

	// ==========================================================
	// 2. MANUAL CANCELLATION
	// ==========================================================
	fmt.Println("\n--- 2. Manual Cancellation ---")

	ctx2, stop := context.WithCancel(context.Background())

	// Simulate user pressing "Cancel"
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println(" [Main] User cancelled the operation!")
		stop()
	}()

	makeRequest(ctx2, "Long-Running-Task", 10*time.Second)

	// ==========================================================
	// 3. CONTEXT VALUES (REQUEST SCOPING)
	// ==========================================================
	fmt.Println("\n--- 3. Context Values (Trace IDs) ---")

	/*
	Context values are for:
	- Request IDs
	- Trace IDs
	- Authentication metadata

	NOT for:
	- Optional function parameters
	- Business logic data
	*/

	valCtx := context.WithValue(
		context.Background(),
		traceIDKey,
		"abcd-1234",
	)

	processData(valCtx)
}

// ==========================================================
// INTERRUPTIBLE WORK FUNCTION
// ==========================================================

func makeRequest(ctx context.Context, name string, duration time.Duration) {
	fmt.Printf(" [%s] Starting (needs %v)\n", name, duration)

	select {
	case <-time.After(duration):
		// Job finished before cancellation
		fmt.Printf(" [%s] Completed successfully\n", name)

	case <-ctx.Done():
		// Context cancelled or timed out
		fmt.Printf(" [%s] ABORTED! Reason: %v\n", name, ctx.Err())
	}
}

// ==========================================================
// CONTEXT VALUE CONSUMPTION
// ==========================================================

func processData(ctx context.Context) {
	// Safe retrieval using typed key
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		traceID = "unknown"
	}

	fmt.Printf(" [Worker] Processing request with TraceID: %s\n", traceID)
}

/*
============================================================
DEEP CONCEPT: THE CONTEXT TREE
============================================================

Contexts form a TREE:

Root:
- context.Background()  (lives forever)

Child:
- context.WithTimeout(parent, 5s)

Grandchild:
- context.WithCancel(child)

RULE:
If a parent context is cancelled,
ALL children are cancelled IMMEDIATELY.

This is how Kubernetes cleans up:
- HTTP request dies
- Database query stops
- Goroutines exit
- Metrics stop collecting
- Logs stop emitting

============================================================
COMMON MISTAKES (VERY IMPORTANT)
============================================================

❌ Forgetting to call cancel()
   → memory + goroutine leaks

❌ Ignoring ctx.Done()
   → cancellation signal is useless

❌ Using context.Background() inside libraries
   → breaks cancellation chain

❌ Storing large data in context
   → context is not a data bag

============================================================
GOOD vs BAD INTERRUPTIBILITY
============================================================

BAD (cannot be cancelled):
	time.Sleep(10 * time.Second)

GOOD (interruptible):
	select {
	case <-time.After(10 * time.Second):
		// work done
	case <-ctx.Done():
		return
	}

============================================================
KUBERNETES CONTEXT
============================================================

- Every controller reconcile loop receives ctx
- API server cancels ctx on client disconnect
- Informers stop when ctx is cancelled
- Controllers shut down gracefully using ctx

context.Context is NOT optional.
It is the backbone of Cloud Native Go.
*/
