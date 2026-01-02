package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

/*
============================================================
MODULE 07 — ERROR HANDLING
STEP 07.5 — RETRY & BACKOFF PATTERNS (DEEP DIVE)
============================================================

Core question:
"When an operation fails, should we retry?"

In distributed systems:
- Networks fail
- APIs timeout
- Leaders change
- Pods restart

Retrying blindly is DANGEROUS.
Retrying correctly is ESSENTIAL.
*/

// ==========================================================
// 1. ERROR CLASSIFICATION (RETRYABLE vs FATAL)
// ==========================================================

var (
	ErrTemporaryFailure = errors.New("temporary failure")
	ErrPermanentFailure = errors.New("permanent failure")
)

/*
Rule:
- NOT all errors are retryable
- Retrying permanent errors makes systems worse
*/

func flakyOperation() error {
	n := rand.Intn(10)

	switch {
	case n < 5:
		return ErrTemporaryFailure // retryable
	case n < 7:
		return ErrPermanentFailure // fatal
	default:
		return nil // success
	}
}

// ==========================================================
// 2. NAIVE RETRY (BAD PRACTICE)
// ==========================================================

func naiveRetry() error {
	for {
		err := flakyOperation()
		if err == nil {
			return nil
		}
		fmt.Println("retrying after error:", err)
	}
}

/*
Why this is BAD:
- Infinite retries
- No delay
- Hot loops
- Can DOS your own system
*/

// ==========================================================
// 3. BOUNDED RETRIES (MINIMUM SAFETY)
// ==========================================================

func boundedRetry(maxAttempts int) error {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := flakyOperation()
		if err == nil {
			return nil
		}

		if errors.Is(err, ErrPermanentFailure) {
			// Fatal error: retrying is pointless
			return err
		}

		fmt.Printf("attempt %d failed: %v\n", attempt, err)
	}

	return fmt.Errorf("retry limit exceeded")
}

// ==========================================================
// 4. EXPONENTIAL BACKOFF (REQUIRED)
// ==========================================================

/*
Backoff strategy:
- First retry quickly
- Each subsequent retry waits longer
- Prevents overload
*/

func retryWithBackoff(ctx context.Context, maxAttempts int) error {
	baseDelay := 100 * time.Millisecond

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := flakyOperation()
		if err == nil {
			return nil
		}

		if errors.Is(err, ErrPermanentFailure) {
			return err
		}

		// Calculate exponential backoff
		backoff := time.Duration(
			math.Pow(2, float64(attempt-1)),
		) * baseDelay

		fmt.Printf(
			"attempt %d failed (%v), backing off for %v\n",
			attempt, err, backoff,
		)

		select {
		case <-time.After(backoff):
			// continue retrying
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("retry limit exceeded")
}

// ==========================================================
// 5. JITTER (CRITICAL IN REAL SYSTEMS)
// ==========================================================

/*
Without jitter:
- All clients retry at the same time
- Thundering herd problem

With jitter:
- Retry times are randomized
*/

func retryWithJitter(ctx context.Context, maxAttempts int) error {
	baseDelay := 100 * time.Millisecond

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := flakyOperation()
		if err == nil {
			return nil
		}

		if errors.Is(err, ErrPermanentFailure) {
			return err
		}

		exp := math.Pow(2, float64(attempt-1))
		jitter := rand.Float64() + 0.5 // 0.5x – 1.5x

		backoff := time.Duration(exp * jitter * float64(baseDelay))

		fmt.Printf(
			"attempt %d failed (%v), jitter backoff %v\n",
			attempt, err, backoff,
		)

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("retry limit exceeded")
}

// ==========================================================
// 6. IDEMPOTENCY (VERY IMPORTANT)
// ==========================================================

/*
Retrying is SAFE only if the operation is IDEMPOTENT.

Idempotent:
- Doing it once or 10 times has the same effect

Examples:
- PUT /resource
- Ensure state exists
- Update cache

Non-idempotent:
- Charge credit card
- Send email
- Increment counter
*/

func idempotentOperation(state *int) error {
	// safe: sets state, not increments
	*state = 42
	return ErrTemporaryFailure
}

// ==========================================================
// 7. KUBERNETES PATTERN
// ==========================================================

/*
Kubernetes controllers:

- Reconcile() is idempotent
- Errors trigger retries
- Backoff prevents API overload
- Context cancels retries on shutdown

This file models EXACTLY that.
*/

// ==========================================================
// 8. MAIN — DEMONSTRATIONS
// ==========================================================

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== Retry & Backoff Patterns ===")

	fmt.Println("\n-- Bounded retry --")
	err := boundedRetry(5)
	fmt.Println("result:", err)

	fmt.Println("\n-- Retry with backoff + context --")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = retryWithBackoff(ctx, 10)
	fmt.Println("result:", err)

	fmt.Println("\n-- Retry with jitter --")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()

	err = retryWithJitter(ctx2, 10)
	fmt.Println("result:", err)

	fmt.Println("\n-- Idempotency demo --")
	var state int
	err = idempotentOperation(&state)
	fmt.Println("state:", state, "error:", err)
}

/*
============================================================
KEY TAKEAWAYS (CRITICAL)
============================================================

1. Never retry blindly
2. Classify errors (retryable vs fatal)
3. Always bound retries
4. Use exponential backoff
5. Add jitter in distributed systems
6. Respect context cancellation
7. Ensure idempotency before retrying

============================================================
KUBERNETES CONTEXT
============================================================

- Controller-runtime uses backoff
- API retries are bounded
- Context stops retries on shutdown
- Non-idempotent operations are avoided

If you understand THIS file,
you understand real-world Go error recovery.
*/
