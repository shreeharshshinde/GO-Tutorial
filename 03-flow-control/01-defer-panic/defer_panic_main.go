package main

import (
	"fmt"
	"os"
)

func main() {

	// ============================================================
	// --- 1. Defer: The Cleanup Crew ---
	// ============================================================
	// 'defer' schedules a function call to run AFTER the surrounding
	// function returns.
	//
	// IMPORTANT:
	// - defer runs on normal return
	// - defer runs on panic
	// - defer does NOT run on os.Exit()

	fmt.Println("1. Opening a file...")
	f, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}

	// Best Practice:
	// Defer cleanup immediately after acquiring a resource.
	defer closeFile(f)

	fmt.Println("2. Writing data...")
	fmt.Fprintln(f, "Hello, Go Systems Engineering!")

	// ============================================================
	// --- 2. Stacked Defers (LIFO Order) ---
	// ============================================================
	// Defers execute in Last-In-First-Out order (like a stack)

	fmt.Println("\n--- Stacked Defers ---")
	defer fmt.Println("Cleanup Step 1 (Runs Last)")
	defer fmt.Println("Cleanup Step 2 (Runs Second)")
	defer fmt.Println("Cleanup Step 3 (Runs First)")

	// ============================================================
	// --- 3. Panic & Recover (Crash Protection) ---
	// ============================================================
	// panic:
	// - Stops normal execution
	// - Unwinds the stack
	// - Executes deferred calls
	//
	// recover:
	// - Stops the panic
	// - Only works inside a deferred function
	// - Only works in the same goroutine

	fmt.Println("\n--- Panic Protection ---")
	safeExecute()

	fmt.Println("4. Main function continues... Program did not crash!")

	// ============================================================
	// --- 4. Defer Arguments Are Evaluated Immediately ---
	// ============================================================

	fmt.Println("\n--- Defer Argument Evaluation ---")
	x := 10
	defer fmt.Println("Deferred x =", x)
	x = 20
	fmt.Println("Current x =", x)
	// Deferred prints 10, not 20

	// ============================================================
	// --- 5. Defer Inside Loops (VERY IMPORTANT) ---
	// ============================================================

	fmt.Println("\n--- Defer Inside Loop ---")
	for i := 0; i < 3; i++ {
		defer fmt.Println("Deferred loop value:", i)
	}
	// Output order (when main exits):
	// 2
	// 1
	// 0

	// ============================================================
	// --- 6. Defer Is NOT Free (But Worth It) ---
	// ============================================================
	// defer has small overhead
	// Use defer for correctness and safety
	// Not for tight performance-critical loops

	// ============================================================
	// --- 7. Panic Is NOT Error Handling ---
	// ============================================================
	// panic is for:
	// - Programmer bugs
	// - Impossible states
	// - Corrupted invariants
	//
	// Errors are for:
	// - I/O failures
	// - User input
	// - Network issues

	// ============================================================
	// --- 8. recover() ONLY Works in Deferred Functions ---
	// ============================================================

	fmt.Println("\n--- recover() Scope ---")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered panic:", r)
			}
		}()
		panic("panic inside anonymous function")
	}()

	// ============================================================
	// --- 9. recover() Does NOT Work Across Goroutines ---
	// ============================================================

	fmt.Println("\n--- recover() and Goroutines ---")
	fmt.Println("Panics must be recovered in the SAME goroutine")

	/*
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("This will NOT run")
			}
		}()
		panic("goroutine panic")
	}()
	*/

	// ============================================================
	// --- 10. Order of Defers During Panic ---
	// ============================================================

	fmt.Println("\n--- Defers During Panic ---")
	func() {
		defer fmt.Println("Deferred 1")
		defer fmt.Println("Deferred 2")
		panic("boom")
	}()
	// Output:
	// Deferred 2
	// Deferred 1

	// ============================================================
	// --- 11. Named Return Values + Defer ---
	// ============================================================

	fmt.Println("\n--- Named Return + Defer ---")
	fmt.Println("Result:", namedReturn())

	// ============================================================
	// --- 12. defer DOES NOT Run on os.Exit ---
	// ============================================================
	// This is CRITICAL knowledge.
	// os.Exit terminates the program immediately.
	// No deferred calls run.

	/*
	defer fmt.Println("This will NEVER run")
	os.Exit(1)
	*/

	// ============================================================
	// --- 13. Panic Value Can Be ANY Type ---
	// ============================================================

	fmt.Println("\n--- Panic Value Types ---")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered type: %T, value: %v\n", r, r)
			}
		}()
		panic(404)
	}()

	// ============================================================
	// --- 14. Re-panicking After Recover ---
	// ============================================================

	fmt.Println("\n--- Re-Panic Pattern ---")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Logging panic:", r)
				panic(r) // rethrow
			}
		}()
		panic("fatal error")
	}()
}

// ============================================================
// Helper Functions
// ============================================================

func closeFile(f *os.File) {
	fmt.Println("3. Defer Triggered: Closing file now.")
	_ = f.Close()
	_ = os.Remove("test.txt")
}

func safeExecute() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("   [RECOVERED] Caught a panic:", r)
		}
	}()

	fmt.Println("   -> About to panic...")
	panic("Something went terribly wrong!")
}

func namedReturn() (result int) {
	defer func() {
		result = result + 10
	}()
	return 5
}
