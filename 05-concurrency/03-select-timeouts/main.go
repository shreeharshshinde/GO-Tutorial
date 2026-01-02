package main

import (
	"fmt"
	"time"
)

/*
This file explains the `select` statement in Go.

Problem:
- A normal receive (msg := <-ch) blocks forever
- Real systems must:
  - Wait on multiple channels
  - Abort slow operations
  - Stay responsive without blocking

Solution:
- select acts like a switch, but for CHANNEL OPERATIONS

This file covers:
- Racing multiple channels
- Timeouts using time.After
- Non-blocking reads and writes
- Why this is critical in Kubernetes
*/

func main() {

	// ==========================================================
	// 1. MULTIPLEXING (RACING CHANNELS)
	// ==========================================================
	fmt.Println("--- 1. Racing Channels (select) ---")

	fastWorker := make(chan string)
	slowWorker := make(chan string)

	go func() {
		time.Sleep(500 * time.Millisecond)
		fastWorker <- "Fast Result"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		slowWorker <- "Slow Result"
	}()

	/*
	select blocks until ONE case is ready.
	If multiple are ready, one is chosen randomly.
	*/

	select {
	case msg := <-fastWorker:
		fmt.Println("Winner:", msg)
	case msg := <-slowWorker:
		fmt.Println("Winner:", msg)
	}

	/*
	This pattern is used when:
	- Multiple workers race
	- First successful result wins
	- Others can be ignored or cancelled
	*/

	// ==========================================================
	// 2. TIMEOUT PATTERN (SYSTEMS CRITICAL)
	// ==========================================================
	fmt.Println("\n--- 2. Timeout Pattern ---")

	dbQuery := make(chan string)

	go func() {
		// Simulate slow operation
		time.Sleep(3 * time.Second)
		dbQuery <- "Query Results"
	}()

	/*
	time.After(d):
	- Returns a channel
	- Sends current time after duration d
	- Used to enforce deadlines
	*/

	select {
	case res := <-dbQuery:
		fmt.Println("Success:", res)
	case <-time.After(1 * time.Second):
		fmt.Println("Error: Operation timed out!")
	}

	/*
	IMPORTANT:
	- time.After allocates a timer
	- In tight loops, prefer time.NewTimer
	*/

	// ==========================================================
	// 3. NON-BLOCKING READS
	// ==========================================================
	fmt.Println("\n--- 3. Non-Blocking Reads ---")

	messages := make(chan string, 1)
	messages <- "Hello"

	/*
	select + default:
	- If no channel is ready, default executes
	- This makes the operation NON-BLOCKING
	*/

	select {
	case msg := <-messages:
		fmt.Println("Received:", msg)
	default:
		fmt.Println("No messages available")
	}

	// Channel is now empty
	select {
	case msg := <-messages:
		fmt.Println("Received:", msg)
	default:
		fmt.Println("No messages available")
	}

	// ==========================================================
	// 4. NON-BLOCKING WRITES
	// ==========================================================
	fmt.Println("\n--- 4. Non-Blocking Writes ---")

	queue := make(chan int, 1)

	select {
	case queue <- 1:
		fmt.Println("Enqueued job")
	default:
		fmt.Println("Queue full, dropping job")
	}

	select {
	case queue <- 2:
		fmt.Println("Enqueued job")
	default:
		fmt.Println("Queue full, dropping job")
	}

	// ==========================================================
	// 5. SELECT WITH STOP SIGNAL (CONTROLLERS)
	// ==========================================================
	fmt.Println("\n--- 5. Stop Signal Pattern ---")

	stopCh := make(chan struct{})

	go func() {
		time.Sleep(1 * time.Second)
		close(stopCh) // broadcast stop signal
	}()

	work := 0
	for {
		select {
		case <-stopCh:
			fmt.Println("Stop signal received. Shutting down.")
			return
		default:
			// Keep working
			work++
			fmt.Println("Working...", work)
			time.Sleep(300 * time.Millisecond)
		}
	}
}

/*
============================================================
DEEP CONCEPTS (READ CAREFULLY)
============================================================

1. select WITHOUT default
- Blocks until a case is ready

2. select WITH default
- Never blocks
- Enables polling / responsiveness

3. time.After
- Channel-based timeout
- Ideal for one-off deadlines

4. Multiple ready cases
- Go picks one randomly
- Prevents starvation

============================================================
KUBERNETES CONTEXT
============================================================

1. API Timeouts
- Client requests use select + time.After
- Prevents hanging forever when API server is slow

2. Controller Shutdown
- stopCh is checked in select
- Allows graceful termination

3. Event Loops
- Controllers must remain responsive
- Non-blocking select is essential

============================================================
COMMON MISTAKES
============================================================

- Forgetting default → unintended blocking
- Forgetting timeout → goroutine leaks
- Using time.After in tight loops
- Ignoring stop signals
*/
