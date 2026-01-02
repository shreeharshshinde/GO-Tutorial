package main

import "fmt"

/*
This file explains TYPE ASSERTIONS and TYPE SWITCHES in Go.

Context:
- You already know how to put values into interfaces (boxing).
- Now you learn how to safely get them back out (unboxing).

Key problem:
- If a value is stored as `any` (or interface{}),
  the compiler does NOT know its concrete type.

Key solutions:
- Type assertions: x.(T)
- Type switches
*/

func main() {

	fmt.Println("--- 1. The Core Problem ---")
	// We store a string inside an empty interface.
	// At runtime, the value is still a string,
	// but at compile time, Go only knows it is 'any'.
	var data any = "Kubernetes"

	// This DOES NOT compile:
	// The compiler does not know that 'data' is a string.
	// fmt.Println(len(data))

	fmt.Println("Stored value:", data)


	fmt.Println("\n--- 2. Unsafe Type Assertion (May Panic) ---")
	// Syntax:
	//   value := interfaceValue.(ConcreteType)
	//
	// We are telling the compiler:
	// "Trust me, I KNOW this is a string."
	strValue := data.(string)
	fmt.Println("Recovered string:", strValue)

	// If the assertion is WRONG, the program PANICS.
	// Uncommenting below will crash the program.
	// wrong := data.(int)


	fmt.Println("\n--- 3. Safe Type Assertion (Comma OK Idiom) ---")
	// Safe form:
	//   value, ok := interfaceValue.(ConcreteType)
	//
	// If the assertion fails:
	// - ok == false
	// - value is the zero value of the type
	intValue, ok := data.(int)
	if !ok {
		fmt.Println(" -> Safety check failed: data is NOT an int.")
	} else {
		fmt.Println(" -> Found integer:", intValue)
	}

	// This is the recommended form when:
	// - Input is dynamic
	// - Source is untrusted
	// - Panic is unacceptable


	fmt.Println("\n--- 4. Multiple Safe Assertions (Not Ideal) ---")
	// You COULD check many types manually like this,
	// but it quickly becomes ugly and error-prone.

	if s, ok := data.(string); ok {
		fmt.Println("String length:", len(s))
	} else if b, ok := data.(bool); ok {
		fmt.Println("Boolean value:", b)
	} else {
		fmt.Println("Unknown type")
	}

	// This is why TYPE SWITCHES exist.


	fmt.Println("\n--- 5. Type Switch (Correct Tool) ---")
	// Type switches allow branching logic
	// based on the ACTUAL dynamic type.
	analyzeType(42)
	analyzeType("Hello")
	analyzeType(true)
	analyzeType(3.14)
	analyzeType([]int{1, 2, 3})
}

/*
analyzeType accepts ANY value.

At compile time:
- input is of type `any`

At runtime:
- Go inspects the actual concrete type
- Executes the matching case
*/
func analyzeType(input any) {

	// Syntax:
	// switch v := input.(type)
	//
	// Rules:
	// - Only valid inside a switch
	// - v has the concrete type in each case
	switch v := input.(type) {

	case int:
		fmt.Printf(" [Type Switch] Int: %d (square=%d)\n", v, v*v)

	case string:
		fmt.Printf(" [Type Switch] String: %q (len=%d)\n", v, len(v))

	case bool:
		fmt.Printf(" [Type Switch] Bool: %t\n", v)

	case float64:
		fmt.Printf(" [Type Switch] Float64: %.2f\n", v)

	case []int:
		fmt.Printf(" [Type Switch] Slice of int, length=%d\n", len(v))

	default:
		// Always include default
		fmt.Printf(" [Type Switch] Unknown type: %T\n", v)
	}
}

/*
Important notes:

1. Type assertions and type switches work ONLY on interfaces.
2. They inspect the DYNAMIC type, not the static type.
3. Unsafe assertions panic on mismatch.
4. Safe assertions return (value, false).
5. Type switches are preferred when handling many possible types.
6. This pattern is used heavily in:
   - Kubernetes API machinery
   - YAML / JSON decoding
   - Logging
   - Plugin systems
*/
