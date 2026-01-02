package main

import "fmt"

func main() {
	fmt.Println("--- 1. Arrays (Fixed Size, Value Type) ---")
	// Arrays have a fixed size. The size is PART of the type.
	// [3]int is distinct from [4]int.
	var grades [3]int = [3]int{90, 85, 95}

	// Arrays are "Value Types". When you assign one to another,
	// Go performs a full COPY of the data.
	gradesCopy := grades
	gradesCopy[0] = 100 // Modifying the copy...

	fmt.Println("Original:", grades)     // [90 85 95] - Unchanged
	fmt.Println("Copy:    ", gradesCopy) // [100 85 95]

	fmt.Println("\n--- 2. Slices (Dynamic Window on Arrays) ---")
	// A slice is a lightweight structure with three fields:
	// 1. Pointer (to the underlying array)
	// 2. Length (elements you can see)
	// 3. Capacity (total space allocated)

	// Create a slice using a literal (Go creates the array hidden in background)
	nums := []int{10, 20, 30, 40, 50}

	// Slicing an existing slice
	// This creates a new "Window" looking at the SAME underlying array.
	subSlice := nums[1:3] // Indices 1 up to (but not including) 3 -> {20, 30}

	fmt.Printf("Original: %v\n", nums)
	fmt.Printf("SubSlice: %v\n", subSlice)

	// DANGER ZONE: Modifying a sub-slice affects the original!
	subSlice[0] = 999
	fmt.Println("After modifying subSlice[0] = 999:")
	fmt.Printf("Original is CHANGED: %v\n", nums)

	fmt.Println("\n--- 3. Length vs Capacity & 'append' ---")
	// 'make' allocates an array and creates a slice view of it.
	// make([]type, len, cap)
	users := make([]string, 0, 5)

	fmt.Printf("Start: Len=%d, Cap=%d, Ptr=%p\n", len(users), cap(users), users)

	users = append(users, "Alice")
	users = append(users, "Bob")
	fmt.Printf("After 2 appends: Len=%d, Cap=%d, Ptr=%p\n", len(users), cap(users), users)

	// What happens when we exceed capacity?
	// Go creates a NEW larger array, COPIES data over, and updates the pointer.
	users = append(users, "Charlie", "Dave", "Eve", "Frank")

	fmt.Printf("After overflow:  Len=%d, Cap=%d, Ptr=%p\n", len(users), cap(users), users)
	fmt.Println("Notice the pointer address changed! A new array was allocated.")

	// --------------------------------------------------------------------
	fmt.Println("\n--- 4. nil Slice vs Empty Slice (IMPORTANT) ---")
	// --------------------------------------------------------------------

	var nilSlice []int      // nil slice
	emptySlice := []int{}   // empty slice
	makeSlice := make([]int, 0)

	fmt.Println("nilSlice == nil:", nilSlice == nil)   // true
	fmt.Println("emptySlice == nil:", emptySlice == nil) // false
	fmt.Println("makeSlice == nil:", makeSlice == nil) // false

	fmt.Println("len(nilSlice):", len(nilSlice))
	fmt.Println("len(emptySlice):", len(emptySlice))

	// len/cap are same, but semantics differ (JSON, APIs, Kubernetes)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 5. Slice Assignment vs copy() ---")
	// --------------------------------------------------------------------

	src := []int{1, 2, 3}

	// This does NOT copy data (aliasing)
	alias := src
	alias[0] = 999
	fmt.Println("src after alias modification:", src)

	// Correct deep copy
	dst := make([]int, len(src))
	copy(dst, src)
	dst[0] = 111

	fmt.Println("src after deep copy:", src)
	fmt.Println("dst:", dst)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 6. Full Slice Expression (Capacity Control) ---")
	// --------------------------------------------------------------------

	base := []int{1, 2, 3, 4, 5}

	// Normal slicing (capacity allows overwrite)
	unsafe := base[1:3]

	// Full slice expression: low : high : max
	safe := base[1:3:3]

	unsafe = append(unsafe, 999)
	fmt.Println("After unsafe append, base is CHANGED:", base)

	base = []int{1, 2, 3, 4, 5}
	safe = append(safe, 999)
	fmt.Println("After safe append, base is SAFE:", base)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 7. Slices Passed to Functions ---")
	// --------------------------------------------------------------------

	list := []string{"A", "B"}

	addWrong(list)
	fmt.Println("After addWrong:", list)

	list = addCorrect(list)
	fmt.Println("After addCorrect:", list)

	addWithPointer(&list)
	fmt.Println("After addWithPointer:", list)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 8. Slices Inside Structs (Aliasing Bug) ---")
	// --------------------------------------------------------------------

	type Pod struct {
		Containers []string
	}

	p1 := Pod{Containers: []string{"c1", "c2"}}
	p2 := p1 // shallow copy
	p2.Containers[0] = "hacked"

	fmt.Println("p1:", p1) // corrupted
	fmt.Println("p2:", p2)

	// Correct deep copy
	p3 := Pod{
		Containers: append([]string(nil), p1.Containers...),
	}
	p3.Containers[0] = "safe"

	fmt.Println("p1 after deep copy:", p1)
	fmt.Println("p3:", p3)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 9. Slices + Goroutines (Race Condition Pattern) ---")
	// --------------------------------------------------------------------

	for _, v := range []int{1, 2, 3} {
		v := v // capture value
		go func() {
			fmt.Println("goroutine value:", v)
		}()
	}

	// NOTE: In real code, use sync.WaitGroup
}

// -------------------- Helper Functions --------------------

func addWrong(s []string) {
	// Modifies local copy only
	s = append(s, "WRONG")
}

func addCorrect(s []string) []string {
	return append(s, "CORRECT")
}

func addWithPointer(s *[]string) {
	*s = append(*s, "POINTER")
}
