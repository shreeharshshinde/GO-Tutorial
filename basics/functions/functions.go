
package main

import "fmt"

// This function simulates making tea.
// Explaining the concept of defer
func main() {
    fmt.Println("Function starts")

    defer fmt.Println("First defer (runs third)")
    defer fmt.Println("Second defer (runs second)")
    defer fmt.Println("Third defer (runs first)")



    fmt.Println("1. Taking out a clean mug.")

    // We "defer" the cleanup task right after we acquire the resource (the mug).
    // This function call is now scheduled to run right before makeTea() exits.
    defer fmt.Println("5. Washing the mug.")

    fmt.Println("2. Boiling water.")
    fmt.Println("3. Pouring water and adding tea bag.")
    fmt.Println("4. Enjoying the tea.")

    // The 'makeTea' function is about to end here.
    fmt.Println("Function ends")

}
