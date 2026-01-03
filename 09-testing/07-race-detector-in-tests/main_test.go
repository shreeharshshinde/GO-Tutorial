package racedetectortests

import (
	"sync"
	"testing"
)

/*
============================================================
MODULE 09 — TESTING IN GO
STEP 09.7 — RACE DETECTOR IN TESTS (COMPLETE DEEP DIVE)
============================================================

This file is EXECUTABLE DOCUMENTATION.

Goal:
Understand WHAT the Go race detector is,
HOW it works,
WHEN to use it,
and WHY CNCF projects treat it as mandatory.

If you understand THIS file,
you will prevent entire classes of production bugs.
*/

/*
------------------------------------------------------------
1. WHAT IS A DATA RACE?
------------------------------------------------------------

A data race occurs when:
- Two or more goroutines access the same memory
- At least one is a write
- There is NO synchronization

Data races cause:
- Silent corruption
- Heisenbugs
- Impossible-to-reproduce failures

Most dangerous bug category in systems code.
*/

/*
------------------------------------------------------------
2. THE RACE DETECTOR
------------------------------------------------------------

Go provides a built-in race detector.

You run it with:
	go test -race ./...

It instruments your program at runtime and detects:
- concurrent unsynchronized memory access

This is NOT a linter.
This is NOT static analysis.
This is RUNTIME detection.
*/

/*
------------------------------------------------------------
3. UNSAFE CODE (INTENTIONALLY RACY)
------------------------------------------------------------

This counter has NO synchronization.
It looks innocent.
It is BROKEN.
*/

type UnsafeCounter struct {
	value int
}

func (c *UnsafeCounter) Inc() {
	c.value++
}

func (c *UnsafeCounter) Value() int {
	return c.value
}

/*
------------------------------------------------------------
4. TEST THAT PASSES BUT IS WRONG
------------------------------------------------------------

IMPORTANT:
This test MAY PASS without -race.
That does NOT mean the code is correct.

This is why -race exists.
*/

func TestUnsafeCounter(t *testing.T) {
	counter := &UnsafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}

	wg.Wait()

	// This assertion might pass by luck
	_ = counter.Value()
}

/*
------------------------------------------------------------
5. RUN THIS TEST WITH RACE DETECTOR
------------------------------------------------------------

	go test -race ./09-testing/07-race-detector-in-tests

You WILL see output like:

	WARNING: DATA RACE
	Read at ...
	Write at ...

This is expected and GOOD.
*/

/*
------------------------------------------------------------
6. SAFE VERSION (PROPERLY SYNCHRONIZED)
------------------------------------------------------------
*/

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

/*
------------------------------------------------------------
7. TEST THAT PASSES WITH -race
------------------------------------------------------------
*/

func TestSafeCounter(t *testing.T) {
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}

	wg.Wait()

	if counter.Value() != 100 {
		t.Fatalf("counter = %d; want 100", counter.Value())
	}
}

/*
------------------------------------------------------------
8. WHAT THE RACE DETECTOR CAN AND CANNOT DO
------------------------------------------------------------

CAN detect:
- shared variable races
- map concurrent writes
- slice concurrent access
- pointer aliasing races

CANNOT detect:
- deadlocks
- logical ordering bugs
- missed signals
- algorithmic mistakes

It is necessary but not sufficient.
*/

/*
------------------------------------------------------------
9. PERFORMANCE COST
------------------------------------------------------------

-race:
- increases memory usage
- slows execution (2–5x)

This is why:
- used in CI
- not in production builds
*/

/*
------------------------------------------------------------
10. CNCF POLICY ON RACES
------------------------------------------------------------

In Kubernetes / CNCF projects:
- Race detector failures = BLOCKING
- PRs must pass -race
- Flaky race tests are fixed immediately
- Unsafe concurrency is rejected outright

Race-free code is a HARD requirement.
*/

/*
------------------------------------------------------------
11. WHEN YOU SHOULD RUN -race
------------------------------------------------------------

You should run -race:
- before opening a PR
- when touching concurrency
- when modifying shared state
- when tests behave strangely

Best practice:
	go test -race ./...
*/

/*
------------------------------------------------------------
12. COMMON BEGINNER MISTAKES
------------------------------------------------------------

❌ "It passed locally, so it's fine"
❌ Ignoring race warnings
❌ Using maps without locks
❌ Assuming atomicity of int++
❌ Disabling race detector to make CI green
*/

/*
------------------------------------------------------------
13. FINAL MENTAL MODEL
------------------------------------------------------------

Tests answer:
	"Does this work?"

Race detector answers:
	"Is this SAFE to run concurrently?"

You need BOTH.
*/

/*
------------------------------------------------------------
14. WHAT YOU SHOULD BE COMFORTABLE WITH NOW
------------------------------------------------------------

After this file, you should:
- Know what data races are
- Understand how -race works
- Recognize racy code immediately
- Fix races using locks or channels
- Pass CNCF CI race checks confidently
*/
