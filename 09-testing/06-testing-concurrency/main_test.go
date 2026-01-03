package testingconcurrency

import (
	"sync"
	"testing"
	"time"
)

/*
============================================================
MODULE 09 — TESTING IN GO
STEP 09.6 — TESTING CONCURRENCY (COMPLETE DEEP DIVE)
============================================================

This file is EXECUTABLE DOCUMENTATION.

Goal:
Understand how to SAFELY test concurrent Go code,
avoid flaky tests,
and write tests that are accepted in CNCF projects.

If you understand THIS file,
you will stop fearing goroutines in tests.
*/

/*
------------------------------------------------------------
1. WHY CONCURRENCY TESTING IS HARD
------------------------------------------------------------

Concurrency bugs are:
- timing-dependent
- non-deterministic
- often invisible locally
- catastrophic in production

Bad concurrency tests:
- rely on time.Sleep
- assume scheduling order
- pass locally, fail in CI

CNCF projects are VERY strict here.
*/

// ==========================================================
// 2. CODE UNDER TEST (THREAD-SAFE COUNTER)
// ==========================================================

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

/*
------------------------------------------------------------
3. BASIC CONCURRENT TEST (WITH WAITGROUP)
------------------------------------------------------------

Key rule:
- NEVER let goroutines outlive the test
- ALWAYS wait for completion
*/

func TestCounter_ConcurrentIncrement(t *testing.T) {
	counter := &Counter{}
	var wg sync.WaitGroup

	numGoroutines := 100
	incrementsPerG := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerG; j++ {
				counter.Inc()
			}
		}()
	}

	// CRITICAL: wait for all goroutines
	wg.Wait()

	expected := numGoroutines * incrementsPerG
	if counter.Value() != expected {
		t.Fatalf(
			"counter = %d; want %d",
			counter.Value(),
			expected,
		)
	}
}

/*
------------------------------------------------------------
4. WHY time.Sleep IS A CODE SMELL IN TESTS
------------------------------------------------------------

BAD TEST:

	go doWork()
	time.Sleep(1 * time.Second)
	assertSomething()

Why this is bad:
- depends on timing
- flaky under load
- slow CI
- hides race conditions

GOOD TEST:
- use synchronization primitives
- WaitGroups
- channels
*/

// ==========================================================
// 5. TESTING CHANNEL-BASED CONCURRENCY
// ==========================================================

func worker(input <-chan int, output chan<- int) {
	for v := range input {
		output <- v * 2
	}
}

func TestWorker_ChannelPipeline(t *testing.T) {
	input := make(chan int)
	output := make(chan int)

	go worker(input, output)

	go func() {
		defer close(input)
		for i := 1; i <= 5; i++ {
			input <- i
		}
	}()

	var results []int
	for i := 0; i < 5; i++ {
		results = append(results, <-output)
	}

	expected := []int{2, 4, 6, 8, 10}
	for i := range expected {
		if results[i] != expected[i] {
			t.Fatalf(
				"results[%d] = %d; want %d",
				i,
				results[i],
				expected[i],
			)
		}
	}
}

/*
------------------------------------------------------------
6. TESTING DEADLOCKS (IMPORTANT)
------------------------------------------------------------

Deadlocks are SILENT failures.
Tests may hang forever.

CNCF tests MUST have time bounds.
*/

func TestDeadlockProtection(t *testing.T) {
	done := make(chan struct{})

	go func() {
		time.Sleep(100 * time.Millisecond)
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("test timed out (possible deadlock)")
	}
}

/*
------------------------------------------------------------
7. USING CHANNELS FOR SYNCHRONIZATION
------------------------------------------------------------

Channels are not just for data.
They are synchronization signals.
*/

func TestGoroutineCompletion(t *testing.T) {
	done := make(chan struct{})

	go func() {
		// simulate work
		time.Sleep(50 * time.Millisecond)
		close(done)
	}()

	<-done // blocks until goroutine finishes
}

/*
------------------------------------------------------------
8. SUBTESTS + CONCURRENCY (SAFE PATTERN)
------------------------------------------------------------

When combining:
- table-driven tests
- subtests
- concurrency

ALWAYS capture loop variables.
*/

func TestConcurrentSubtests(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"double 1", 1, 2},
		{"double 2", 2, 4},
		{"double 3", 3, 6},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.input * 2
			if result != tt.want {
				t.Fatalf("got %d; want %d", result, tt.want)
			}
		})
	}
}

/*
------------------------------------------------------------
9. COMMON CONCURRENCY TEST SMELLS
------------------------------------------------------------

❌ time.Sleep instead of sync
❌ goroutines without waiting
❌ shared mutable state without locks
❌ assuming execution order
❌ flaky tests accepted as "normal"
*/

/*
------------------------------------------------------------
10. HOW KUBERNETES TESTS CONCURRENCY
------------------------------------------------------------

Kubernetes tests:
- use WaitGroups
- use channels as signals
- use contexts with timeouts
- fail fast on deadlocks
- avoid sleeps whenever possible

Flaky tests are aggressively fixed or deleted.
*/

/*
------------------------------------------------------------
11. WHAT YOU SHOULD BE COMFORTABLE WITH NOW
------------------------------------------------------------

After this file, you should:
- Test goroutines deterministically
- Avoid flaky timing-based tests
- Use WaitGroups correctly
- Use channels for synchronization
- Combine concurrency + subtests safely
*/
