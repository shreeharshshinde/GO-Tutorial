package main

import "fmt"

/*
This file explains INTERFACES in Go.

Key ideas:
- Interfaces define behavior, not data
- Implementation is implicit (duck typing)
- Interfaces enable polymorphism
- The empty interface (any) can hold any value
- Interface values internally store (type, value)

This file does NOT cover nil-interface pitfalls.
That is the next topic.
*/

// ==========================================
// 1. THE CONTRACT (Interface)
// ==========================================
// An interface defines WHAT can be done.
// It does NOT define HOW it is done.
type CloudProvider interface {
	LaunchInstance(name string) string
}

/*
Important:
- Interfaces contain ONLY method signatures
- No fields
- No implementations
*/

// ==========================================
// 2. IMPLICIT IMPLEMENTATION (DUCK TYPING)
// ==========================================

// Struct 1: AWS
type AWS struct {
	Region string
}

// AWS implements CloudProvider because it has the required method.
// No "implements" keyword is needed.
func (a AWS) LaunchInstance(name string) string {
	return fmt.Sprintf("AWS (%s): Launching EC2 instance '%s'", a.Region, name)
}

// Struct 2: Azure
type Azure struct {
	SubscriptionID string
}

// Azure ALSO implements CloudProvider.
// Go only checks method names and signatures.
func (az Azure) LaunchInstance(name string) string {
	return fmt.Sprintf(
		"Azure (Sub-%s): Starting VM '%s'",
		az.SubscriptionID,
		name,
	)
}

/*
This is duck typing:
- If it has the right methods, it satisfies the interface
- Go does not care about the concrete type
*/

// ==========================================
// 3. POLYMORPHISM
// ==========================================
// This function works with ANY CloudProvider.
// It does not know or care whether it is AWS or Azure.
func provisionServer(cloud CloudProvider, serverName string) {
	fmt.Println(" -> Provisioning request received...")
	result := cloud.LaunchInstance(serverName)
	fmt.Println("    Result:", result)
}

/*
Why this matters:
- Business logic depends on interfaces
- Concrete implementations can change
- This is heavily used in Kubernetes and CNCF projects
*/

// ==========================================
// 4. THE EMPTY INTERFACE (any)
// ==========================================
// interface{} (alias: any) has ZERO methods.
// Since every type has at least zero methods,
// EVERY type satisfies it.
func printAnything(data any) {
	fmt.Printf(
		"[Generic Log] Value: %v (Type: %T)\n",
		data,
		data,
	)
}

/*
Use cases:
- Logging
- Serialization
- Generic containers
- Configuration values

Caution:
- You lose compile-time type safety
- You must use type assertions or type switches
*/

// ==========================================
// 5. INTERFACE INTERNALS (TYPE, VALUE)
// ==========================================
// An interface value stores TWO things internally:
// 1. Concrete type
// 2. Concrete value
//
// It is NOT just a pointer.

func main() {
	fmt.Println("--- 1. Polymorphism (Duck Typing) ---")

	// Concrete implementations
	myAWS := AWS{Region: "us-east-1"}
	myAzure := Azure{SubscriptionID: "99-xyz"}

	// Same function, different behavior
	provisionServer(myAWS, "web-server-01")
	provisionServer(myAzure, "db-server-01")

	fmt.Println("\n--- 2. The 'any' Type (Universal Container) ---")

	printAnything(42)            // int
	printAnything("Hello World") // string
	printAnything(myAWS)         // struct

	fmt.Println("\n--- 3. Interface Internals (Type, Value Pair) ---")

	var c CloudProvider

	// At this point:
	// c = (type: nil, value: nil)
	if c == nil {
		fmt.Println("Interface is currently nil.")
	}

	// Assign a concrete value
	c = myAWS

	// Now:
	// c = (type: AWS, value: {Region: "us-east-1"})
	fmt.Printf("Assigned AWS. Interface holds type: %T\n", c)

	// Calling method uses dynamic dispatch
	fmt.Println(c.LaunchInstance("cache-server"))
}
