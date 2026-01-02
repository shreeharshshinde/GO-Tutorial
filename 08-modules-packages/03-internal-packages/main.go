package main

import "fmt"

/*
This file demonstrates the IDEA of internal packages.

You cannot fully demonstrate internal/ violations
inside a single package without multiple modules,
so this file explains what would and would not compile.
*/

func main() {
	fmt.Println("=== 08.3 internal/ Packages ===")

	fmt.Println("Key rule:")
	fmt.Println("Packages inside internal/ can only be imported")
	fmt.Println("by code within the parent module tree.")

	fmt.Println()
	fmt.Println("This rule is enforced at COMPILE TIME.")
	fmt.Println("CNCF projects rely on this for architecture safety.")
}
