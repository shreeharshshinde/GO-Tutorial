package main

import "fmt"

const x int64 = 10
const (
	idKey   = "id"
	nameKey = "name"
)
const z = 20 * 10

func constant() {
	const y = "hello"
	fmt.Println(x)
	fmt.Println(y)
	// x = x + 1 // cannot assign to x (untyped int constant 10)
	/* y = "bye" // cannot assign to y (const "hello") */
	fmt.Println(x)
	fmt.Println(y)
}
