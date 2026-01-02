package main

import (
	"fmt"
	"runtime"
	"time"
)

/*
This file explains GOROUTINES and the GO SCHEDULER.

Why this matters:
- Go does NOT use threads like Java or C++
- Goroutines are lightweight, cheap, and multiplexed
- Go has its OWN scheduler inside your binary

Key ideas covered:
- Goroutine lifecycle
- Why main() controls program lifetime
- The Go scheduler vs OS scheduler
- CPU cores and GOMAXPROCS
- Closures and the loop variable trap
- Why this matters in Kubernetes controllers
*/

func main() {

	// ==========================================================
	// 1. PROGRAM LIFECYCLE: main() IS THE BOSS
	// ==========================================================
	fmt.Println("--- 1. main() vs Goroutines ---")

	// Launch a goroutine.
	go func() {
		fmt.Println(" [Background] Goroutine is running")
	}()

	/*
	Important rule:
	- The program exits when main() exits
	- ALL goroutines are terminated immediately
	- There is NO implicit waiting

	This is why sleep is used here ONLY for demo.
	In real code, use sync.WaitGroup or channels.
	*/
	time.Sleep(100 * time.Millisecond)

	// ==========================================================
	// 2. GO SCHEDULER & CPU CORES
	// ==========================================================
	fmt.Println("\n--- 2. Go Scheduler & CPUs ---")

	cores := runtime.NumCPU()
	fmt.Printf("Machine has %d CPU cores\n", cores)

	/*
	GOMAXPROCS controls how many OS threads (P) can run Go code.
	By default:
	- Go uses all available CPU cores
	- This changed in Go 1.5+ (automatic)
	*/
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

	// ==========================================================
	// 3. WHAT IS A GOROUTINE REALLY?
	// ==========================================================
	fmt.Println("\n--- 3. Goroutines vs Threads ---")

	/*
	Java / C++:
	- 1 thread ≈ 1 OS thread
	- Large stack (~1MB)
	- Context switching is slow (kernel mode)

	Go:
	- Goroutines are NOT OS threads
	- Initial stack ~2KB (grows automatically)
	- Context switching happens in user space
	*/

	// Launch many goroutines cheaply
	for i := 0; i < 3; i++ {
		go func(id int) {
			fmt.Println(" Goroutine:", id)
		}(i)
	}

	time.Sleep(100 * time.Millisecond)

	// ==========================================================
	// 4. THE M:N SCHEDULER (G-P-M MODEL)
	// ==========================================================
	fmt.Println("\n--- 4. The M:N Scheduler (Conceptual) ---")

	/*
	G (Goroutine):
	- Your function + stack
	- Cheap, thousands or millions possible

	M (Machine):
	- Actual OS thread
	- Expensive, limited by OS

	P (Processor):
	- Scheduling context
	- Owns a run queue of goroutines

	The scheduler maps:
		M OS threads
		N goroutines
	using P as the bridge
	*/

	/*
	CRITICAL BEHAVIOR:
	- If a goroutine blocks (I/O, syscall, sleep)
	- The scheduler DETACHES it from the OS thread
	- Another goroutine runs on that thread

	This is why:
	- One OS thread can handle thousands of goroutines
	*/

	// ==========================================================
	// 5. CLOSURES & LOOP VARIABLE TRAP
	// ==========================================================
	fmt.Println("\n--- 5. Closures & Loop Variables ---")

	numbers := []int{10, 20, 30, 40, 50}

	/*
	What a closure does:
	- It CAPTURES variables from the outer scope
	- It does NOT capture values, it captures variables

	This is subtle and extremely important.
	*/

	// Historically INCORRECT pattern (buggy in older Go)
	for _, v := range numbers {
		_ = v // mark v as intentionally unused (example-only)
		go func() {
			// In Go < 1.22, this often prints:
			// 50, 50, 50, 50, 50
			// because 'v' is reused
			// fmt.Println("Unsafe:", v)
		}()
	}

	/*
	Go 1.22+ FIX:
	- Loop variables are now copied per iteration
	- But DO NOT rely on this for clarity
	*/

	// CORRECT, EXPLICIT, SAFE pattern
	for _, v := range numbers {
		go func(val int) {
			fmt.Println(" Safe:", val)
		}(v)
	}

	time.Sleep(300 * time.Millisecond)

	// ==========================================================
	// 6. WHY THIS MATTERS IN SYSTEMS (KUBERNETES)
	// ==========================================================
	fmt.Println("\n--- 6. Why This Matters in Systems Code ---")

	/*
	Kubernetes example:

	for _, pod := range podList {
		go processPod(pod) // MUST pass pod
	}

	WRONG:
	go func() {
		processPod(pod) // pod may change!
	}()

	Result:
	- Every goroutine processes the LAST pod
	- Causes incorrect reconciliation
	- Real production bugs
	*/

	// ==========================================================
	// 7. GOROUTINES HAVE NO RETURN VALUES
	// ==========================================================
	fmt.Println("\n--- 7. No Return Values ---")

	/*
	go func() int {
		return 42
	}()

	This is INVALID.

	If you need data back:
	- Use channels
	- Or shared memory with synchronization

	Channels are the NEXT topic.
	*/

	// ==========================================================
	// 8. SCHEDULER FAIRNESS (runtime.Gosched)
	// ==========================================================
	fmt.Println("\n--- 8. Scheduler Yielding ---")

	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(" A:", i)
			runtime.Gosched() // Yield execution
		}
	}()

	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(" B:", i)
			runtime.Gosched()
		}
	}()

	time.Sleep(300 * time.Millisecond)

	/*
	runtime.Gosched():
	- Yields the processor
	- Allows other goroutines to run
	- Rarely needed in real code
	*/
}
