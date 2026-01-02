package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
============================================================
STEP 6.3 — THE GO RACE DETECTOR (DEEP DIVE)
============================================================

What the race detector is:
- A dynamic analysis tool
- Instruments your program at runtime
- Detects concurrent UNSYNCHRONIZED memory access

What it is NOT:
- A compiler check
- A static analyzer
- A proof of correctness

Golden Rule:
If the race detector reports a race,
YOUR PROGRAM IS WRONG.
*/

/*
============================================================
1. A PROGRAM WITH A DATA RACE
============================================================
*/

var counter int

func raceExample() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // DATA RACE
		}()
	}

	wg.Wait()
}

/*
Run this program normally:
	go run main.go

Likely output:
	(no error)
	counter value might look correct

Now run WITH race detector:
	go run -race main.go

You WILL see:
	WARNING: DATA RACE
*/

/*
============================================================
2. FIXING THE RACE (DETECTOR GOES SILENT)
============================================================
*/

func fixedWithMutex() {
	counter = 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
}

/*
Race detector result:
- No warnings
- Program is now race-free
*/

/*
============================================================
3. RACE DETECTOR vs LOGIC BUGS
============================================================

IMPORTANT:
The race detector ONLY detects data races.
It does NOT detect logic bugs.
*/

func logicBugButNoRace() {
	ch := make(chan int)

	go func() {
		ch <- 1
		ch <- 2 // DEADLOCK, but no data race
	}()

	fmt.Println(<-ch)
}

/*
This program:
- Has NO data race
- But will DEADLOCK

Race detector:
- Says nothing
- Because this is a LOGIC bug, not a race
*/

/*
============================================================
4. WHAT THE RACE DETECTOR CANNOT SEE
============================================================
*/

func partialSyncBug() {
	var a int
	var b int32

	go func() {
		a = 42
		// Atomic write only synchronizes b
	}()

	// Read b atomically
	if atomic.LoadInt32(&b) == 1 {
		// a is read WITHOUT synchronization
		// This read is NOT guaranteed to see 42
		fmt.Println("a =", a)
	}
}

/*
This is a REAL bug (memory visibility issue),
but the race detector MAY NOT catch it.

Why?
- No concurrent conflicting access detected
- Bug is ordering/visibility related

Race detector is powerful, but NOT omniscient.
*/

/*
============================================================
5. HOW THE RACE DETECTOR WORKS (CONCEPTUAL)
============================================================

- Instruments every read/write
- Tracks goroutine access history
- Builds happens-before graph
- Reports conflicting accesses

Cost:
- 5–10x slower
- Much higher memory usage

This is why:
- -race is for testing
- NEVER for production binaries
*/

/*
============================================================
6. WHEN TO USE THE RACE DETECTOR
============================================================

ALWAYS use -race when:
- Writing concurrent code
- Modifying shared state
- Writing controllers, caches, queues
- Running CI for Go services

In Kubernetes:
- CI pipelines run with -race
- Many bugs are caught ONLY this way
*/

/*
============================================================
7. ABSOLUTE RULES (NON-NEGOTIABLE)
============================================================

RULE 1:
If -race reports a race → fix it. No exceptions.

RULE 2:
"Seems harmless" races are STILL races.

RULE 3:
Protect ALL shared data, not just writes.

RULE 4:
Never silence race detector warnings.

RULE 5:
Race-free ≠ correct, but racy = incorrect.
*/

func main() {
	fmt.Println("=== Race Detector Demo ===")

	fmt.Println("\n-- Running racy example --")
	raceExample()
	fmt.Println("Counter:", counter)

	fmt.Println("\n-- Running fixed example --")
	fixedWithMutex()
	fmt.Println("Counter:", counter)

	fmt.Println("\n-- Logic bug example (no race) --")
	fmt.Println("This will deadlock if uncommented.")
	// logicBugButNoRace()
}
