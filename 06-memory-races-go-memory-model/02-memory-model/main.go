package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
============================================================
STEP 6.2 — THE GO MEMORY MODEL (DEEP DIVE)
============================================================

Key Question:
"Why can code that looks correct still behave incorrectly?"

Answer:
Because CPUs, compilers, and Go itself are allowed to
REORDER memory operations unless you explicitly prevent it.

The Go Memory Model defines:
- WHEN a write by one goroutine becomes visible to another
- WHAT guarantees synchronization primitives provide
- WHAT is *not* guaranteed without synchronization

------------------------------------------------------------
CORE IDEA:
Visibility is NOT guaranteed without synchronization.
------------------------------------------------------------
*/

/*
============================================================
1. THE ILLUSION OF SEQUENTIAL THINKING
============================================================

Humans think code runs top-to-bottom.
CPUs and compilers do NOT.

This code looks obvious.
It is NOT safe.
*/

var ready bool
var data int

func exampleWithoutSync() {
	go func() {
		data = 42      // write 1
		ready = true  // write 2
	}()

	// Busy-wait (DO NOT DO THIS)
	for !ready {
	}

	// ❌ This print is NOT guaranteed to print 42
	fmt.Println("Data:", data)
}

/*
Why?

Because the compiler / CPU may reorder writes:
- ready = true may become visible BEFORE data = 42

Result:
- ready == true
- data == 0 (old value)

This is NOT a race detector issue.
This is a MEMORY ORDERING issue.
*/

/*
============================================================
2. HAPPENS-BEFORE (THE MOST IMPORTANT CONCEPT)
============================================================

"happens-before" is a GUARANTEE relationship.

If A happens-before B:
- All writes in A are visible in B

If NOT:
- Visibility is undefined
*/

/*
============================================================
3. MUTEX CREATES HAPPENS-BEFORE
============================================================
*/

var mu sync.Mutex
var safeData int
var safeReady bool

func exampleWithMutex() {
	go func() {
		mu.Lock()
		safeData = 42
		safeReady = true
		mu.Unlock()
	}()

	for {
		mu.Lock()
		if safeReady {
			fmt.Println("Safe Data:", safeData)
			mu.Unlock()
			return
		}
		mu.Unlock()
	}
}

/*
Why this works:

Unlock() happens-before Lock()
→ All writes before Unlock() are visible after Lock()

Mutex provides:
- Mutual exclusion
- Memory visibility
- Ordering guarantees
*/

/*
============================================================
4. CHANNELS ALSO CREATE HAPPENS-BEFORE
============================================================
*/

func exampleWithChannel() {
	ch := make(chan int)

	go func() {
		ch <- 42 // send
	}()

	val := <-ch // receive
	fmt.Println("Channel Data:", val)
}

/*
Rule:
A send on a channel happens-before the corresponding receive.

Channels are synchronization primitives,
NOT just communication tools.
*/

/*
============================================================
5. ATOMICS: VISIBILITY WITHOUT LOCKS
============================================================
*/

var atomicData int64

func exampleWithAtomic() {
	go func() {
		atomic.StoreInt64(&atomicData, 42)
	}()

	val := atomic.LoadInt64(&atomicData)
	fmt.Println("Atomic Data:", val)
}

/*
Atomic operations guarantee:
- No tearing
- Correct visibility
- Correct ordering (for that variable)

BUT:
- Only for that variable
- No protection for invariants across multiple variables
*/

/*
============================================================
6. THE MOST DANGEROUS BUG: PARTIAL SYNCHRONIZATION
============================================================
*/

var a int
var b int32

func brokenPartialSync() {
	go func() {
		a = 10
		atomic.StoreInt32(&b, 1)
	}()

	if atomic.LoadInt32(&b) == 1 {
		// ❌ NOT guaranteed that a == 10
		fmt.Println("a =", a)
	}
}

/*
Why this is broken:

Atomic guarantees apply ONLY to the atomic variable (b).
There is NO happens-before relationship for 'a'.

This bug passes code review.
This bug passes tests.
This bug FAILS in production.
*/

/*
============================================================
7. WHAT GO GUARANTEES (AND WHAT IT DOES NOT)
============================================================

WITHOUT synchronization:
- Writes may not be visible
- Reads may see stale values
- Instructions may be reordered

WITH synchronization:
- Mutex Lock/Unlock → happens-before
- Channel send/receive → happens-before
- Atomic operations → happens-before (per variable)

GO DOES NOT GUARANTEE:
- Sequential consistency by default
- Visibility without synchronization
- That "sleep" fixes races
*/

/*
============================================================
8. PRACTICAL RULES (SYSTEMS GO)
============================================================

RULE 1:
If one goroutine writes, ALL access must synchronize.

RULE 2:
Synchronization is about VISIBILITY, not just locking.

RULE 3:
Choose ONE synchronization mechanism per data:
- Mutex OR
- Channel OR
- Atomics

RULE 4:
Never mix atomics and non-atomics on related state.

RULE 5:
time.Sleep is NOT synchronization.
*/

/*
============================================================
9. KUBERNETES CONTEXT
============================================================

Kubernetes relies heavily on:
- happens-before guarantees
- informer caches
- shared controller state

Examples:
- Informer update → controller read
- Leader election flags
- Shared metrics

A missing happens-before edge can:
- Break reconciliation
- Cause infinite loops
- Corrupt cluster state

This is why Kubernetes code is VERY strict
about synchronization.
*/

/*
============================================================
MAIN (INTENTIONALLY EMPTY)
============================================================

We do NOT call these functions automatically.

This file is for:
- Reading
- Reasoning
- Understanding guarantees

Not for running blindly.
*/

func main() {
	fmt.Println("Read the code. Reason about the memory model.")
}
