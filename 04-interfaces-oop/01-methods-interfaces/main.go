package main

import "fmt"

/*
This file explains METHODS and RECEIVERS in Go.

Key ideas covered here:
- What a method is
- Value receivers vs pointer receivers
- Copy vs mutation behavior
- Automatic dereferencing by Go
- Method sets (T vs *T)
- Methods on non-struct types
- Nil receiver pattern
- Performance implications
- Rules used in real-world Go (CNCF / Kubernetes)

No prior knowledge is assumed.
*/

// Server is a simple struct used to demonstrate receiver behavior.
type Server struct {
	Host     string
	Port     int
	Hits     int
}

/*
VALUE RECEIVER: func (s Server)

What happens:
- Go creates a COPY of the Server struct
- The method operates on the copy
- Any changes are lost after the method returns

Properties:
- Safe (original data cannot be modified)
- Can be expensive for large structs
- Good for read-only behavior
*/
func (s Server) ShowInfo() {
	fmt.Printf("[View] %s:%d (Hits=%d)\n", s.Host, s.Port, s.Hits)
}

/*
POINTER RECEIVER: func (s *Server)

What happens:
- Go passes the MEMORY ADDRESS of the struct
- The method operates on the original value
- Changes persist after the call

Properties:
- No copying (better performance)
- Allows mutation
- Must be used carefully

Rule of thumb:
In systems programming, use pointer receivers by default.
*/
func (s *Server) RecordHit() {
	s.Hits++
	fmt.Println("-> Hit recorded")
}

/*
IMPORTANT: Automatic dereferencing

Even though RecordHit expects *Server,
this call is valid:

    srv.RecordHit()

Go automatically rewrites it as:

    (&srv).RecordHit()

Similarly, value-receiver methods can be called on pointers.
This is done at COMPILE TIME.
*/

////////////////////////////////////////////////////////////////////////////////
// METHODS ON NON-STRUCT TYPES
////////////////////////////////////////////////////////////////////////////////

/*
Go allows methods on ANY named type,
not just structs.

This enables domain-specific behavior
without inheritance.
*/
type IntList []int

/*
This method belongs to IntList.
It cannot be attached to []int directly,
because []int is a built-in type.
*/
func (l IntList) Sum() int {
	total := 0
	for _, v := range l {
		total += v
	}
	return total
}

/*
Rules for method definitions:
- The receiver type must be defined in the SAME package
- You cannot add methods to types from other packages
- You cannot add methods to built-in types
*/

////////////////////////////////////////////////////////////////////////////////
// NIL RECEIVER PATTERN (UNIQUE TO GO)
////////////////////////////////////////////////////////////////////////////////

/*
TreeNode demonstrates the "nil receiver" pattern.

In Go:
- Calling a method on a nil pointer is allowed
- The method itself must handle nil safely
*/
type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

/*
This method is SAFE even if the receiver is nil.

This is extremely common in:
- trees
- linked lists
- recursive data structures
- Kubernetes API types
*/
func (t *TreeNode) SumValues() int {
	if t == nil {
		return 0
	}
	return t.Value + t.Left.SumValues() + t.Right.SumValues()
}

/*
Important distinction:
- Nil RECEIVER is allowed
- Nil INTERFACE is different (covered later)
*/

////////////////////////////////////////////////////////////////////////////////
// METHOD SETS (CRITICAL CONCEPT)
////////////////////////////////////////////////////////////////////////////////

/*
Method set determines which methods belong to a type.

For this file:

Method set of Server:
- ShowInfo()

Method set of *Server:
- ShowInfo()
- RecordHit()

Implications:
- Server does NOT have RecordHit()
- *Server has BOTH methods

This directly affects interface satisfaction.
*/

////////////////////////////////////////////////////////////////////////////////
// METHODS ARE JUST FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

/*
These two are equivalent:

    srv.ShowInfo()
    ShowInfo(srv)

The receiver is just syntactic sugar.
There is no "this" or "self" keyword in Go.
*/

////////////////////////////////////////////////////////////////////////////////
// PERFORMANCE NOTE
////////////////////////////////////////////////////////////////////////////////

/*
Value receivers COPY the entire struct.
If the struct is large, this is expensive.

Pointer receivers:
- avoid copying
- are faster
- are preferred in most real systems
*/

////////////////////////////////////////////////////////////////////////////////
// CONSISTENCY RULE (VERY IMPORTANT)
////////////////////////////////////////////////////////////////////////////////

/*
If ANY method uses a pointer receiver,
use pointer receivers for ALL methods.

Why:
- Predictable method sets
- Fewer interface bugs
- Matches standard library & Kubernetes style
*/

////////////////////////////////////////////////////////////////////////////////
// MAIN: Examples tying everything together
////////////////////////////////////////////////////////////////////////////////

func main() {
	fmt.Println("--- Pointer vs Value Receivers ---")

	srv := Server{Host: "127.0.0.1", Port: 8080}

	// Value receiver (copy)
	srv.ShowInfo()

	// Pointer receiver (mutates original)
	srv.RecordHit()
	srv.ShowInfo()

	fmt.Println("\n--- Methods on Non-Struct Types ---")

	nums := IntList{1, 2, 3, 4}
	fmt.Println("Sum:", nums.Sum())

	fmt.Println("\n--- Nil Receiver Pattern ---")

	root := &TreeNode{
		Value: 10,
		Left:  &TreeNode{Value: 5},
		Right: nil,
	}

	// Safe even though Right is nil
	fmt.Println("Tree sum:", root.SumValues())

	var empty *TreeNode
	fmt.Println("Nil tree sum:", empty.SumValues())

	fmt.Println("\n--- Copy vs Pointer ---")

	copySrv := srv
	ptrSrv := &srv

	copySrv.RecordHit() // operates on COPY
	ptrSrv.RecordHit()  // operates on ORIGINAL

	fmt.Println("Original Hits:", srv.Hits)
	fmt.Println("Copy Hits:", copySrv.Hits)
}
