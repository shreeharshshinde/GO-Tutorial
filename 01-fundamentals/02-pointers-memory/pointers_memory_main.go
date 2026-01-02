package main

import "fmt"

type Pod struct {
	Name   string
	Status string
}

func main() {
	fmt.Println("--- Pointers in Kubernetes Context ---")

	// 1. POINTER DECLARATION
	// Create a pod. Memory is allocated for this struct.
	myPod := Pod{Name: "nginx", Status: "Pending"}

	// Create a pointer to that memory location.
	// podPointer holds the Hex address (e.g., 0xc000010200), not the data.
	var podPointer *Pod = &myPod

	fmt.Printf("Pod Address in Memory: %p\n", podPointer)

	// 2. PASS BY REFERENCE (The "Controller" Pattern)
	// We pass the address. The function jumps to that address and edits the data.
	updateStatus(podPointer)

	// The original 'myPod' is changed because we followed the map (pointer) to it.
	fmt.Printf("Current Pod Status: %s\n", myPod.Status)
}

// updateStatus requires a pointer (*Pod). 
// If we removed the *, this function would edit a useless copy.
func updateStatus(p *Pod) {
	fmt.Println("Controller: Updating pod status to Running...")
	p.Status = "Running"
}