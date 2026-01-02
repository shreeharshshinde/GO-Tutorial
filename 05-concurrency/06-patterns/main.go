package main

import (
	"fmt"
	"time"
)

/*
This file explains the WORKER POOL pattern.

Problem: Unbounded Concurrency
- Launching a goroutine per job does NOT scale
- 10,000 jobs → 10,000 goroutines
- This exhausts:
  - CPU
  - Memory
  - File descriptors
  - Databases

Solution: Worker Pool
- Fixed number of workers
- Shared queue of jobs
- Backpressure prevents overload

This pattern is foundational in:
- Kubernetes controllers
- Job schedulers
- Data pipelines
*/

// ==========================================================
// JOB & RESULT DEFINITIONS
// ==========================================================

// Job represents a unit of work.
type Job struct {
	ID       int
	Filename string
}

// Result represents completed work.
type Result struct {
	JobID    int
	Status   string
	Duration time.Duration
}

func main() {

	// ==========================================================
	// 1. CONFIGURATION
	// ==========================================================
	const numWorkers = 3
	const numJobs = 10

	fmt.Println("--- Starting Worker Pool ---")
	fmt.Printf("Workers: %d, Jobs: %d\n", numWorkers, numJobs)

	/*
	Buffered channels are essential here.

	jobs:
	- Acts as the work queue
	- Buffers pending jobs

	results:
	- Collects finished work
	- Buffered so workers don’t block unnecessarily
	*/
	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	// ==========================================================
	// 2. START WORKERS (CONSUMERS)
	// ==========================================================
	/*
	We launch a FIXED number of goroutines.

	Important:
	- Workers block waiting for jobs
	- No CPU is wasted while waiting
	- This caps concurrency
	*/
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// ==========================================================
	// 3. PRODUCE JOBS
	// ==========================================================
	/*
	The producer (main goroutine) sends jobs.

	If jobs channel becomes full:
	- main() BLOCKS
	- This is BACKPRESSURE
	- System self-regulates instead of crashing
	*/
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{
			ID:       j,
			Filename: fmt.Sprintf("data_%d.csv", j),
		}
	}

	/*
	CRITICAL RULE:
	- Only the SENDER closes the channel
	- Closing signals: "no more work is coming"
	*/
	close(jobs)

	// ==========================================================
	// 4. COLLECT RESULTS
	// ==========================================================
	/*
	We expect EXACTLY numJobs results.

	This avoids:
	- Waiting forever
	- Needing to close results channel
	*/
	for i := 1; i <= numJobs; i++ {
		res := <-results
		fmt.Printf(
			" [Main] Result for Job %d: %s (took %v)\n",
			res.JobID,
			res.Status,
			res.Duration,
		)
	}

	fmt.Println("--- All jobs completed ---")
}

// ==========================================================
// WORKER FUNCTION
// ==========================================================

/*
worker:
- Runs forever until jobs channel is CLOSED
- Processes ONE job at a time
- Sends result back to results channel

Key idea:
- Workers SHARE the same jobs channel
- This naturally load-balances work
*/
func worker(id int, jobs <-chan Job, results chan<- Result) {

	/*
	range on channel:
	- Receives until channel is closed AND empty
	- Clean shutdown without extra signaling
	*/
	for job := range jobs {
		fmt.Printf(
			"   [Worker %d] Started Job %d (%s)\n",
			id,
			job.ID,
			job.Filename,
		)

		start := time.Now()

		// Simulate expensive work
		time.Sleep(500 * time.Millisecond)

		results <- Result{
			JobID:    job.ID,
			Status:   "Success",
			Duration: time.Since(start),
		}
	}

	fmt.Printf("   [Worker %d] Stopping (no more jobs)\n", id)
}

/*
============================================================
DEEP CONCEPT: BACKPRESSURE
============================================================

Backpressure means:
- Producers slow down when consumers are overloaded

How this code enforces it:
- jobs channel has finite capacity
- When full, producer blocks
- Workers must free slots before more jobs enter

This prevents:
- Memory blow-ups
- DB overload
- File descriptor exhaustion

============================================================
WHY THIS IS BETTER THAN UNBOUNDED GOROUTINES
============================================================

BAD:
	for _, job := range jobs {
		go process(job)
	}

RESULT:
- Unlimited goroutines
- System collapse under load

GOOD:
- Fixed workers
- Predictable resource usage
- Stable throughput

============================================================
KUBERNETES CONTEXT
============================================================

In Kubernetes controllers you often see:

	NewController(..., threadiness = 2)

threadiness == number of workers

Meaning:
- Only 2 pods are reconciled concurrently
- Queue absorbs bursts
- API server is protected
- Controller stays responsive

This file models EXACTLY that design.
*/
