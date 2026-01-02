package main

import "fmt"

/*
This file explains the MOST COMMON and MOST DANGEROUS
interface bug in Go: NIL INTERFACES.

Many Go developers (even experienced ones) get this wrong.

Key idea:
- An interface is NOT just a pointer
- It is a pair: (Type, Value)

This file explains:
- What a nil interface really is
- Interface holding a nil concrete value
- Why == nil lies sometimes
- How this causes real Kubernetes bugs
*/

// ==========================================
// 1. Interface Internals (Recap)
// ==========================================

// An interface value internally stores:
//   1. Concrete Type
//   2. Concrete Value
//
// interfaceValue == nil
// ONLY when BOTH type AND value are nil.

// ==========================================
// 2. A Simple Interface
// ==========================================

type Worker interface {
	Work()
}

type Engineer struct {
	Name string
}

func (e *Engineer) Work() {
	fmt.Println("Engineer working:", e.Name)
}

// ==========================================
// 3. The First Trap: Interface Holding nil
// ==========================================

func main() {

	fmt.Println("--- 1. Truly Nil Interface ---")

	var w Worker
	// w is (type=nil, value=nil)

	if w == nil {
		fmt.Println("w is nil (no type, no value)")
	}

	fmt.Println("\n--- 2. Interface Holding nil Concrete Value ---")

	var eng *Engineer = nil
	// eng is a nil *Engineer pointer

	w = eng
	// w is now:
	// (type=*Engineer, value=nil)

	// This is the trap
	if w == nil {
		fmt.Println("w is nil")
	} else {
		fmt.Println("w is NOT nil (type is known)")
	}

	// Even though the underlying pointer is nil,
	// the interface itself is NOT nil.

	fmt.Println("\n--- 3. Calling Method on Interface with nil Value ---")
	// This will PANIC because Work() uses e.Name
	// Uncomment to see the crash.
	// w.Work()

	fmt.Println("\n--- 4. Why This Happens ---")
	/*
	Comparing w == nil checks:
	- Is the TYPE nil?
	- Is the VALUE nil?

	Here:
	- TYPE  = *Engineer (not nil)
	- VALUE = nil

	So w != nil
	*/

	fmt.Println("\n--- 5. Common Real-World Bug Pattern ---")

	err := doWork(false)

	if err != nil {
		fmt.Println("Error occurred:", err)
	} else {
		fmt.Println("No error")
	}

	/*
	EXPECTED:
	No error

	ACTUAL:
	err != nil (because interface holds typed nil)

	This exact bug has caused:
	- Kubernetes controller crashes
	- Memory leaks
	- Infinite retries
	*/

	fmt.Println("\n--- 6. The Broken Function ---")
	fmt.Println("doWork returns (*CustomError)(nil) as error")

	fmt.Println("\n--- 7. The Correct Fix ---")
	fixedErr := doWorkFixed(false)

	if fixedErr != nil {
		fmt.Println("Error occurred:", fixedErr)
	} else {
		fmt.Println("No error (correct)")
	}

	fmt.Println("\n--- 8. How to Safely Check for Nil ---")
	checkNilInterface()
}

// ==========================================
// 4. Error Interface Trap (VERY IMPORTANT)
// ==========================================

// error is just an interface:
// type error interface {
//     Error() string
// }

type CustomError struct {
	Msg string
}

func (e *CustomError) Error() string {
	return e.Msg
}

// Broken version
func doWork(fail bool) error {
	if !fail {
		var err *CustomError = nil
		return err // returns (type=*CustomError, value=nil)
	}
	return &CustomError{Msg: "work failed"}
}

// Correct version
func doWorkFixed(fail bool) error {
	if !fail {
		return nil // returns (type=nil, value=nil)
	}
	return &CustomError{Msg: "work failed"}
}

// ==========================================
// 5. How to Detect This Bug
// ==========================================

func checkNilInterface() {
	var err error

	var ce *CustomError = nil
	err = ce

	fmt.Println("err == nil:", err == nil)

	// Safe way: type assertion
	if err != nil {
		if ce, ok := err.(*CustomError); ok && ce == nil {
			fmt.Println("Interface holds nil *CustomError")
		}
	}
}

/*
RULES TO REMEMBER (READ THIS CAREFULLY):

1. interface == nil ONLY when (type=nil, value=nil)
2. interface holding nil concrete value != nil
3. Returning typed nil as interface is a BUG
4. error interface is the most common victim
5. Always return literal nil when no error
6. Be suspicious of err != nil checks
7. Kubernetes code is full of defensive nil checks because of this
*/
