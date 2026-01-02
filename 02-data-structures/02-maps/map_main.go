package main

import "fmt"

func main() {
	fmt.Println("--- 1. Initialization ---")
	// Method A: Literal (Good for static data)
	// map[KeyType]ValueType
	servers := map[string]string{
		"srv-01": "192.168.1.10",
		"srv-02": "192.168.1.11",
	}

	// Method B: make() (Good if you need to add data later)
	// We can hint capacity (10) to avoid resizing overhead, just like slices.
	users := make(map[string]int, 10)
	users["admin"] = 1
	users["guest"] = 1001

	fmt.Printf("Servers: %v\n", servers)
	fmt.Printf("Users:   %v\n", users)

	fmt.Println("\n--- 2. Access & The 'Comma OK' Idiom ---")
	// Reading a value
	ip := servers["srv-01"]
	fmt.Println("IP for srv-01:", ip)

	// DANGER: What if the key doesn't exist?
	// Go returns the "Zero Value" for the type (empty string for string, 0 for int).
	// It does NOT panic or throw an error.
	unknownIP := servers["srv-99"]
	fmt.Printf("IP for srv-99 (Missing): '%s'\n", unknownIP)

	// The "Comma OK" idiom checks if the key ACTUALLY exists.
	// val, ok := map[key]
	val, exists := servers["srv-99"]
	if !exists {
		fmt.Println(" -> Key 'srv-99' does not exist in map.")
	} else {
		fmt.Println(" -> Found:", val)
	}

	fmt.Println("\n--- 3. Deleting ---")
	delete(servers, "srv-02")
	fmt.Println("After deleting srv-02:", servers)

	fmt.Println("\n--- 4. Iteration (Random Order) ---")
	// The order of keys in a range loop is randomized intentionally by Go.
	// Do not write code that expects "alpha" to come before "beta".
	for key, value := range users {
		fmt.Printf("Key: %s, Value: %d\n", key, value)
	}

	fmt.Println("\n--- 5. Reference Semantics ---")
	// Maps are reference types. Passing them to functions shares the data.
	resetUser(users)
	fmt.Println("After function call (admin should be 0):", users["admin"])

	// --------------------------------------------------------------------
	fmt.Println("\n--- 6. nil Map vs Empty Map (IMPORTANT) ---")
	// --------------------------------------------------------------------

	var nilMap map[string]int        // nil map
	emptyMap := make(map[string]int) // empty but initialized

	fmt.Println("nilMap == nil:", nilMap == nil)
	fmt.Println("emptyMap == nil:", emptyMap == nil)

	// Reading from nil map is SAFE
	fmt.Println("Read from nilMap:", nilMap["ghost"]) // zero value

	// Writing to nil map PANICS
	// nilMap["x"] = 1 // ❌ panic: assignment to entry in nil map

	emptyMap["ok"] = 1 // ✅ works
	fmt.Println("emptyMap after write:", emptyMap)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 7. Map Assignment Does NOT Copy Data ---")
	// --------------------------------------------------------------------

	original := map[string]int{"a": 1, "b": 2}
	alias := original // NO deep copy

	alias["a"] = 999
	fmt.Println("original after alias modification:", original)

	// Correct way to deep copy a map
	copied := make(map[string]int, len(original))
	for k, v := range original {
		copied[k] = v
	}
	copied["a"] = 111

	fmt.Println("original after deep copy:", original)
	fmt.Println("copied:", copied)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 8. Maps Are NOT Comparable ---")
	// --------------------------------------------------------------------

	// You CANNOT do this:
	// if original == copied {} // ❌ compile-time error

	// Only comparison allowed is against nil
	fmt.Println("original == nil:", original == nil)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 9. Maps Are NOT Safe for Concurrent Writes ---")
	// --------------------------------------------------------------------

	shared := make(map[int]int)
	_ = shared // mark as intentionally unused for demonstration

	// This pattern WILL cause runtime panic if done concurrently:
	/*
		go func() { shared[1] = 1 }()
		go func() { shared[2] = 2 }()
	*/

	fmt.Println("Maps need sync.Mutex or sync.Map for concurrency")

	// --------------------------------------------------------------------
	fmt.Println("\n--- 10. Deleting While Ranging Is SAFE ---")
	// --------------------------------------------------------------------

	temp := map[string]int{"a": 1, "b": 2, "c": 3}
	for k := range temp {
		delete(temp, k) // safe in Go
	}
	fmt.Println("Map after deleting during iteration:", temp)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 11. Zero Value Trap (Comma OK Needed) ---")
	// --------------------------------------------------------------------

	counts := map[string]int{"apple": 0}

	v1 := counts["apple"]
	v2 := counts["banana"]

	fmt.Println("apple:", v1)  // 0 (exists)
	fmt.Println("banana:", v2) // 0 (does NOT exist)

	_, ok := counts["banana"]
	fmt.Println("banana exists?", ok)
}

func resetUser(m map[string]int) {
	// This modifies the ORIGINAL map
	m["admin"] = 0
}
