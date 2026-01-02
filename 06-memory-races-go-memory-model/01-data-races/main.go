package main

import (
	"fmt"
	"sync"
	"time"
)

/*
This file explains DATA RACES.

Definition (precise):
A data race occurs when:
1. Two or more goroutines access the SAME memory
2. At least one access is a WRITE
3. There is NO synchronization between them

Result:
- Undefined behavior
- Incorrect values
- Silent corruption
- Impossible-to-debug bugs
*/

// ==========================================================
// 1. A RACY PROGRAM (DO NOT COPY THIS PATTERN)
// ==========================================================

var counter int // shared memory, NOT protected

func main() {

	fmt.Println("--- 1. Data Race Example ---")

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // RACE: read-modify-write
		}()
	}

	wg.Wait()
	fmt.Println("Final counter value:", counter)

	/*
	EXPECTED: 1000
	REALITY:  often less (e.g., 937, 982, etc.)

	Why?
	counter++ is NOT atomic.
	It expands to:
	1. read counter
	2. add 1
	3. write counter
	*/
	time.Sleep(100 * time.Millisecond)

	// ==========================================================
	// 2. FIXING THE RACE WITH A MUTEX
	// ==========================================================
	fmt.Println("\n--- 2. Fixing with Mutex ---")

	counter = 0
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("Final counter value:", counter)

	// ==========================================================
	// 3. WHY RACES ARE DANGEROUS
	// ==========================================================
	fmt.Println("\n--- 3. Why Data Races Are Dangerous ---")

	/*
	Data races are NOT just “wrong numbers”.

	They can cause:
	- Corrupted maps
	- Crashes (fatal error: concurrent map writes)
	- Memory reordering bugs
	- Security vulnerabilities

	Worst of all:
	- They may NOT fail every time
	*/

	fmt.Println("Races are timing-dependent bugs.")
}

/*
============================================================
KEY TAKEAWAYS
============================================================

1. A race is a correctness bug, not a performance issue
2. If one goroutine writes, ALL access must be synchronized
3. “It seems to work” means NOTHING
4. Mutexes establish ORDER and VISIBILITY

============================================================
KUBERNETES CONTEXT
============================================================

- Controller caches
- Informer state
- Shared metrics
- Leader election

ALL of these rely on race-free memory access.

A single race can:
- Corrupt controller state
- Cause infinite reconcile loops
- Trigger cascading failures
*/
