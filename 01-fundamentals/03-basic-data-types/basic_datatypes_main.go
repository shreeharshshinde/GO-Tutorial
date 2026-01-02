package main

import (
	"fmt"
)

// ==========================================
// 1. TYPE ALIASING (The "Enum" Pattern)
// ==========================================
// K8s does not have "enums". It uses type aliasing to create strict categories.
// This prevents you from passing a "NodeCondition" string into a "PodPhase" function.

type PodPhase string

const (
	PodPending   PodPhase = "Pending"
	PodRunning   PodPhase = "Running"
	PodSucceeded PodPhase = "Succeeded"
	PodFailed    PodPhase = "Failed"
)

func main() {
	fmt.Println("--- 1. Integers & Sizing ---")
	// K8s is specific about sizes.
	// 'int' depends on the OS (32 or 64 bit).
	// 'int32' and 'int64' are fixed size. K8s APIs use these strictly.
	var loopCounter int = 10         // Good for local loops
	var replicas int32 = 3           // K8s Spec: Replicas is ALWAYS int32
	var resourceVersion int64 = 1024 // K8s Spec: ResourceVersion is usually int64

	fmt.Printf("Loop Counter (int): %d\n", loopCounter)
	fmt.Printf("Replicas (int32): %d\n", replicas)
	fmt.Printf("ResourceVersion (int64): %d\n", resourceVersion)

	// CONVERSION REQUIRED: You cannot add int32 to int without casting!
	// total := loopCounter + replicas // This would ERROR
	total := int32(loopCounter) + replicas
	fmt.Printf("Total (converted): %d\n", total)

	fmt.Println("\n--- 2. Floats ---")
	// Used for CPU quotas (e.g., 0.5 CPU)
	var cpuLimit float64 = 0.5
	fmt.Printf("CPU Limit: %0.2f\n", cpuLimit)

	fmt.Println("\n--- 3. Strings (Raw vs Interpreted) ---")
	// Normal String: Recognizes escape characters like \n
	var podName string = "nginx-pod\n"

	// Raw String Literal (Backticks): Ignores escape characters.
	// K8s uses this HEAVILY for ConfigMaps, Scripts, and embedding YAML/JSON.
	var configMapData string = `
								apiVersion: v1
								kind: ConfigMap
								metadata:
								name: game-config
								data:
								  	  game.properties: |
									enemies=aliens
									lives=3
								`
	fmt.Printf("Pod Name: %s", podName)
	fmt.Println("ConfigMap YAML Snippet:")
	fmt.Println(configMapData)

	fmt.Println("\n--- 4. Booleans ---")
	// Used for flags like "Ready", "Started", "EnableServiceLinks"
	var isReady bool = true
	var isTerminating bool = false // Default is false
	fmt.Printf("Is Pod Ready? %t\n", isReady)
	fmt.Printf("Is Pod Terminating? %t\n", isTerminating)

	fmt.Println("\n--- 5. Type Safety & Aliasing (K8s Enums) ---")
	var currentPhase PodPhase = PodRunning

	// var wrongPhase PodPhase = "SomeRandomString" // This works locally but is logically wrong.
	// Strict functions demand the type:
	checkPodStatus(currentPhase)

	// checkPodStatus("Running") // ERROR: Cannot use string as PodPhase

	fmt.Println("\n--- 6. Zero Values (Crucial for K8s) ---")
	// In Go, variables declared without a value get a "Zero Value".
	// K8s has to distinguish between "User set replicas to 0" vs "User didn't set replicas".
	var emptyInt int       // 0
	var emptyString string // "" (empty string)
	var emptyBool bool     // false
	var emptyPtr *int      // nil (This is how K8s detects "missing" values!)

	fmt.Printf("Zero Int: %d\n", emptyInt)
	fmt.Printf("Zero String: '%s'\n", emptyString)
	fmt.Printf("Zero Bool: %t\n", emptyBool)
	fmt.Printf("Zero Pointer: %v\n", emptyPtr)

	if emptyPtr == nil {
		fmt.Println(" -> Pointer is nil. In K8s, this means 'field not specified in YAML'.")
	}
}

// This function ONLY accepts the custom type 'PodPhase', not just any string.
func checkPodStatus(phase PodPhase) {
	fmt.Printf("Checking logic for phase: %s\n", phase)
}