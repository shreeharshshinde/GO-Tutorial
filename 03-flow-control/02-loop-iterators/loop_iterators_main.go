package main

import "fmt"

func main() {

	fmt.Println("--- 1. Basic for loop ---")
	// Go has only one loop keyword: for
	for i := 0; i < 3; i++ {
		fmt.Println("i =", i)
	}

	fmt.Println("\n--- 2. for as while loop ---")
	// for can behave like while
	j := 0
	for j < 3 {
		fmt.Println("j =", j)
		j++
	}

	fmt.Println("\n--- 3. Infinite loop ---")
	// for without condition runs forever
	k := 0
	for {
		if k == 3 {
			break
		}
		fmt.Println("k =", k)
		k++
	}

	fmt.Println("\n--- 4. break and continue ---")
	for i := 0; i < 5; i++ {
		if i == 2 {
			continue // skip this iteration
		}
		if i == 4 {
			break // exit loop completely
		}
		fmt.Println("value:", i)
	}

	fmt.Println("\n--- 5. for with multiple variables ---")
	for a, b := 0, 10; a < b; a, b = a+1, b-1 {
		fmt.Println("a:", a, "b:", b)
	}

	fmt.Println("\n--- 6. range over slice ---")
	nums := []int{10, 20, 30}
	for index, value := range nums {
		fmt.Println("index:", index, "value:", value)
	}
	// value is a COPY, not a reference

	fmt.Println("\n--- 7. Modifying slice during range (common mistake) ---")
	for _, v := range nums {
		v = v * 10 // does NOT modify nums
	}
	fmt.Println("nums after wrong update:", nums)

	// correct way: use index
	for i := range nums {
		nums[i] = nums[i] * 10
	}
	fmt.Println("nums after correct update:", nums)

	fmt.Println("\n--- 8. range over array ---")
	arr := [3]int{1, 2, 3}
	for i, v := range arr {
		fmt.Println("i:", i, "v:", v)
	}

	fmt.Println("\n--- 9. range over map ---")
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range m {
		fmt.Println("key:", k, "value:", v)
	}
	// map iteration order is random

	fmt.Println("\n--- 10. range over string ---")
	// range over string iterates over runes (UTF-8 safe)
	s := "Go✓"
	for i, r := range s {
		fmt.Printf("index: %d rune: %c\n", i, r)
	}

	fmt.Println("\n--- 11. byte-wise string iteration ---")
	for i := 0; i < len(s); i++ {
		fmt.Printf("byte index: %d value: %x\n", i, s[i])
	}

	fmt.Println("\n--- 12. ignoring index or value ---")
	for _, v := range nums {
		fmt.Println("value only:", v)
	}
	for i := range nums {
		fmt.Println("index only:", i)
	}

	fmt.Println("\n--- 13. Loop variable capture (very important) ---")
	for _, v := range []int{1, 2, 3} {
		go func() {
			fmt.Println("wrong capture:", v)
		}()
	}

	// correct way
	for _, v := range []int{1, 2, 3} {
		v := v
		go func() {
			fmt.Println("correct capture:", v)
		}()
	}

	fmt.Println("\n--- 14. labeled break ---")
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outer
			}
			fmt.Println("i:", i, "j:", j)
		}
	}

	fmt.Println("\n--- 15. labeled continue ---")
outer2:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				continue outer2
			}
			fmt.Println("i:", i, "j:", j)
		}
	}

	fmt.Println("\n--- 16. range over channel ---")
	ch := make(chan int)

	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
		}
		close(ch)
	}()

	for v := range ch {
		fmt.Println("received:", v)
	}
	// loop ends when channel is closed

	fmt.Println("\n--- 17. range vs index loop (performance note) ---")
	type Big struct {
		Data [1024]int
	}

	bigSlice := []Big{{}, {}}

	// range copies value
	for _, v := range bigSlice {
		_ = v
	}

	// index avoids copy
	for i := range bigSlice {
		_ = bigSlice[i]
	}
}
