package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
============================================================
STEP 6.4 — FALSE SHARING & CACHE LINES (DEEP DIVE)
============================================================

Problem:
"My code is race-free and correct, but it is SLOW."

Cause:
FALSE SHARING — multiple goroutines modifying
different variables that live on the SAME CPU cache line.

This is NOT a data race.
This is a HARDWARE performance problem.
*/

/*
============================================================
1. WHAT IS A CACHE LINE?
============================================================

Modern CPUs:
- Do NOT load individual variables
- Load memory in chunks called CACHE LINES

Typical cache line size:
- 64 bytes

If two variables are within the same cache line:
- Writing to one INVALIDATES the entire line
- Other cores must reload it
- Performance collapses
*/

/*
============================================================
2. FALSE SHARING EXAMPLE (SLOW)
============================================================
*/

// Two counters placed next to each other in memory
type Counters struct {
	a int64
	b int64
}

func falseSharing() {
	var wg sync.WaitGroup
	var c Counters

	start := time.Now()

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 100_000_000; i++ {
			atomic.AddInt64(&c.a, 1)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100_000_000; i++ {
			atomic.AddInt64(&c.b, 1)
		}
	}()

	wg.Wait()
	fmt.Println("False sharing duration:", time.Since(start))
}

/*
Even though:
- a and b are DIFFERENT variables
- atomic operations are used
- there is NO data race

They likely share the SAME cache line.
Each write invalidates the other's cache.
*/

/*
============================================================
3. FIXING FALSE SHARING WITH PADDING
============================================================
*/

// Padding forces variables onto different cache lines
type PaddedCounters struct {
	a int64
	_ [56]byte // padding to fill cache line (64 - 8)
	b int64
}

func noFalseSharing() {
	var wg sync.WaitGroup
	var c PaddedCounters

	start := time.Now()

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 100_000_000; i++ {
			atomic.AddInt64(&c.a, 1)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100_000_000; i++ {
			atomic.AddInt64(&c.b, 1)
		}
	}()

	wg.Wait()
	fmt.Println("No false sharing duration:", time.Since(start))
}

/*
============================================================
4. WHY THIS WORKS
============================================================

Padding ensures:
- a is on one cache line
- b is on another cache line

Now:
- CPU cores do not fight over cache ownership
- Performance scales linearly
*/

/*
============================================================
5. WHY ATOMICS DO NOT SAVE YOU
============================================================

Atomics guarantee:
- Correctness
- Visibility
- Ordering (per variable)

Atomics DO NOT guarantee:
- Cache line isolation
- Performance

False sharing can happen even with perfect atomic code.
*/

/*
============================================================
6. WHEN THIS MATTERS
============================================================

False sharing matters when:
- High-frequency counters
- Metrics
- Hot paths
- Tight loops
- Multiple cores

It does NOT matter when:
- Low contention
- I/O bound code
- Infrequent writes
*/

/*
============================================================
7. REAL-WORLD EXAMPLES
============================================================

Kubernetes:
- Per-CPU stats
- Metrics counters
- Scheduler queues

Databases:
- Connection counters
- Lock-free structures

Go runtime itself:
- Uses padding extensively
*/

/*
============================================================
8. GOLDEN RULES
============================================================

RULE 1:
Correctness comes before performance.

RULE 2:
If performance collapses under load, suspect cache lines.

RULE 3:
Do NOT pad everything blindly.

RULE 4:
Measure before and after (benchmarks).

RULE 5:
Race-free does NOT mean fast.
*/

func main() {
	fmt.Println("Running false sharing demo (this may take time)...")

	falseSharing()
	noFalseSharing()

	fmt.Println("Done. Compare durations.")
}
