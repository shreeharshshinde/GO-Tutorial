package main

import (
	"fmt"
	"time"
)

/*
This file explains CHANNELS in Go.

Mental model:
- Goroutines are workers
- Channels are conveyor belts
- Data moves, ownership transfers

Golden Rule of Go:
"Do not communicate by sharing memory;
 share memory by communicating."

This file covers:
- Unbuffered channels (synchronous)
- Buffered channels (asynchronous)
- Blocking behavior
- Deadlocks
- Directional channels
- Closing channels
*/

func main() {

	// ==========================================================
	// 1. UNBUFFERED CHANNELS (SYNCHRONOUS)
	// ==========================================================
	fmt.Println("--- 1. Unbuffered Channel (Handshake) ---")

	/*
	Unbuffered channel:
	- Capacity = 0
	- Sender blocks until receiver is ready
	- Receiver blocks until sender is ready
	- Forces synchronization
	*/

	handshake := make(chan string)

	go func() {
		fmt.Println(" [Sender] Sending packet...")
		handshake <- "Packet A" // BLOCKS until main receives
		fmt.Println(" [Sender] Packet received by main!")
	}()

	time.Sleep(1 * time.Second) // simulate delay

	fmt.Println(" [Main] Waiting to receive...")
	msg := <-handshake // Receiver unblocks sender
	fmt.Println(" [Main] Got:", msg)

	/*
	Think of unbuffered channels like a phone call:
	Both sides must be present.
	*/

	// ==========================================================
	// 2. BUFFERED CHANNELS (ASYNCHRONOUS)
	// ==========================================================
	fmt.Println("\n--- 2. Buffered Channel (Queue) ---")

	/*
	Buffered channel:
	- Has capacity > 0
	- Sender blocks ONLY when buffer is full
	- Receiver blocks ONLY when buffer is empty
	- Enables burst handling
	*/

	queue := make(chan string, 2)

	queue <- "Job 1"
	queue <- "Job 2"
	fmt.Println("Added 2 jobs without blocking.")

	/*
	At this point:
	- Buffer is full
	- Another send would BLOCK forever (deadlock)
	*/

	// queue <- "Job 3" // Uncommenting this causes DEADLOCK

	fmt.Println("Reading:", <-queue)
	fmt.Println("Reading:", <-queue)

	/*
	Think of buffered channels like email:
	You can send messages even if the receiver is busy,
	until the mailbox fills up.
	*/

	// ==========================================================
	// 3. BLOCKING & DEADLOCKS
	// ==========================================================
	fmt.Println("\n--- 3. Blocking & Deadlocks ---")

	/*
	Deadlock rule:
	- All goroutines are blocked
	- Go runtime detects this and PANICS

	Common causes:
	- Sending with no receiver
	- Receiving with no sender
	- range on channel that is never closed
	*/

	// ==========================================================
	// 4. DIRECTIONAL CHANNELS
	// ==========================================================
	fmt.Println("\n--- 4. Directional Channels & Closing ---")

	workStream := make(chan int, 5)

	/*
	Directional channels restrict behavior at COMPILE TIME.
	This prevents whole classes of bugs.
	*/

	go producer(workStream)
	consumer(workStream)
}

// ==========================================================
// PRODUCER: WRITE-ONLY CHANNEL
// ==========================================================

// chan<- int means SEND-ONLY
// Reading from ch here would be a compile-time error
func producer(ch chan<- int) {
	/*
	Closing a channel:
	- Signals "no more values"
	- Only the sender should close
	- Closing twice PANICS
	*/
	defer close(ch)

	for i := 1; i <= 3; i++ {
		fmt.Printf(" [Producer] Sending %d\n", i)
		ch <- i
	}

	fmt.Println(" [Producer] Closed channel.")
}

// ==========================================================
// CONSUMER: READ-ONLY CHANNEL
// ==========================================================

// <-chan int means RECEIVE-ONLY
// Sending here would be a compile-time error
func consumer(ch <-chan int) {
	/*
	range on channel:
	- Receives values until channel is CLOSED
	- If channel is never closed → DEADLOCK
	*/

	for num := range ch {
		fmt.Printf(" [Consumer] Processed %d\n", num)
	}

	fmt.Println(" [Consumer] Channel empty and closed.")
}

/*
============================================================
DEEP CONCEPTS (READ CAREFULLY)
============================================================

1. UNBUFFERED CHANNELS
- Enforce synchronization
- Useful when order and handoff matter

2. BUFFERED CHANNELS
- Absorb bursts
- Used in worker pools and queues

3. CLOSING CHANNELS
- Required when receiver uses range
- Sender owns the responsibility to close

4. DIRECTIONAL CHANNELS
- Encode intent in types
- Prevent misuse
- Common in Kubernetes APIs

============================================================
KUBERNETES CONTEXT
============================================================

Kubernetes controllers use WORKQUEUES:
- Backed by buffered channels
- Handle bursts of events (Pod updates, Node changes)
- Prevent controller crashes under load

Directional channels ensure:
- Producers cannot consume
- Consumers cannot produce
- System stays correct by construction
*/
