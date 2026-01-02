package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
This file explains SYNC PRIMITIVES in Go.

Why they exist:
- Channels move data between goroutines
- Mutexes protect shared data in place

Problem:
- Concurrent writes to shared memory cause RACE CONDITIONS
- Races lead to crashes OR silent corruption (worst case)

Solutions covered here:
- sync.Mutex
- sync.RWMutex
- atomic operations

Golden rule:
Use channels to communicate.
Use mutexes to protect shared state.
*/

// ==========================================================
// 1. MUTEX: PROTECTING A SHARED RESOURCE
// ==========================================================

// SafeCounter protects an integer from concurrent access.
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

/*
Inc increments the counter safely.

Lock rules:
- Lock BEFORE touching shared data
- Unlock AFTER finishing
- Always use defer for Unlock
*/
func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Critical section
	c.value++
}

/*
Value reads the counter safely.
Reads also need protection.
*/
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// ==========================================================
// 2. RWMutex: MANY READERS, ONE WRITER
// ==========================================================

/*
RWMutex allows:
- Multiple readers at the same time
- Only one writer
- Writers block readers and other writers

Used heavily in:
- Kubernetes informers
- Shared caches
*/

type SafeCache struct {
	mu    sync.RWMutex
	items map[string]string
}

// Get acquires a READ lock.
func (c *SafeCache) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.items[key]
}

// Set acquires a WRITE lock.
func (c *SafeCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

// ==========================================================
// 3. ATOMIC OPERATIONS (LOCK-FREE)
// ==========================================================

/*
Atomics:
- Operate at CPU instruction level
- Faster than mutexes
- Limited to simple operations

Use atomics when:
- You need counters, flags, or stats
- No complex invariants
*/

type AtomicCounter struct {
	value int64
}

func (a *AtomicCounter) Inc() {
	atomic.AddInt64(&a.value, 1)
}

func (a *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&a.value)
}

func main() {

	// ==========================================================
	// 1. Mutex Example (Preventing Data Races)
	// ==========================================================
	fmt.Println("--- 1. sync.Mutex (Preventing Races) ---")

	counter := SafeCounter{}
	var wg sync.WaitGroup

	// Launch many goroutines touching the same memory
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}

	wg.Wait()

	/*
	Without mutex:
	- Result would be non-deterministic
	- Could be < 1000
	- Could silently corrupt memory
	*/
	fmt.Printf("Final Count: %d (Expected: 1000)\n", counter.Value())

	// ==========================================================
	// 2. RWMutex Example (Read-Heavy Workload)
	// ==========================================================
	fmt.Println("\n--- 2. sync.RWMutex (Read-Heavy Cache) ---")

	cache := SafeCache{
		items: make(map[string]string),
	}

	// Single write
	cache.Set("k8s-version", "v1.30")

	// Many concurrent reads
	for i := 0; i < 5; i++ {
		go func(id int) {
			val := cache.Get("k8s-version")
			fmt.Printf(" [Reader %d] Found: %s\n", id, val)
		}(i)
	}

	time.Sleep(300 * time.Millisecond)

	// ==========================================================
	// 3. Atomic Counter Example
	// ==========================================================
	fmt.Println("\n--- 3. Atomic Operations (Fast Counters) ---")

	var atomicCounter AtomicCounter
	wg = sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomicCounter.Inc()
		}()
	}

	wg.Wait()
	fmt.Printf("Atomic Count: %d (Expected: 1000)\n", atomicCounter.Value())
}

/*
============================================================
DEEP CONCEPTS (READ CAREFULLY)
============================================================

1. WHY MUTEXES EXIST
- CPU instructions are not atomic by default
- Read-modify-write is NOT safe concurrently
- Mutex enforces mutual exclusion

2. DEADLOCKS
- Forgetting Unlock() freezes your program
- defer Unlock() is mandatory
- Never lock twice in the same goroutine

3. RWMutex TRADEOFF
- Faster for read-heavy workloads
- Slower if writes are frequent
- Writers starve readers temporarily

4. ATOMICS VS MUTEX
- Atomics are faster
- Atomics are limited
- Mutexes handle complex state

============================================================
KUBERNETES CONTEXT
============================================================

- Informer caches use RWMutex
- Controllers read state constantly (RLock)
- API events update state occasionally (Lock)
- Atomic counters track metrics efficiently

============================================================
COMMON MISTAKES
============================================================

- Using mutex when channels are better
- Forgetting to lock on reads
- Returning pointers to protected data
- Overusing RWMutex when writes are frequent
- Using atomics for complex logic
*/
