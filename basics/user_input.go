// All Go programs start with a package declaration.
// 'main' is the special package for executable programs.
package main

// Import statements bring in other packages for us to use.
import (
	"bufio"  // For buffered I/O, like reading text
	"fmt"    // For formatted I/O, like printing to the console
	"os"     // Provides an interface to operating system functionality
	"strconv"// For string conversions (e.g., string to number)
	"strings"// For string manipulation
)

// The main function is the entry point of our application.
func userInput() {
	// --- Reading a String Input ---
	fmt.Print("Please enter your name: ")

	// Create a new reader that reads from standard input (the keyboard).
	reader := bufio.NewReader(os.Stdin)

	// Read all text until the user presses the Enter key ('\n').
	// This returns the text and a potential error.
	userName, _ := reader.ReadString('\n')

	// The input includes the newline character, so we remove it for clean output.
	userName = strings.TrimSpace(userName)

	fmt.Println("Hello,", userName)

	// --- Reading a Numeric Input ---
	fmt.Print("Please enter your age: ")

	// We can reuse the same reader.
	ageInput, _ := reader.ReadString('\n')

	// Convert the string input to a number.
	// We must first trim whitespace (like the newline character).
	// strconv.ParseFloat expects a string and the bit size (64 for float64).
	age, err := strconv.ParseFloat(strings.TrimSpace(ageInput), 64)

	// It's crucial to handle potential errors.
	// What if the user typed "hello" instead of a number?
	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
	} else {
		fmt.Println("Next year, you will be", age+1, "years old.")
	}
}