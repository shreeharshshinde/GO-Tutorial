package main

import (
	// Standard library package
	"fmt"
	
	// Renaming imports (Aliasing)
	// PRO TIP: This is standard in K8s codebases. 
	// We often rename "math/rand" to "rand" or specific K8s APIs like "appsv1".
	str "strings" 
)

func main() {
	clusterName := "production-cluster-01"

	// Using the standard package
	fmt.Println("Cluster Name:", clusterName)

	// Using the aliased package (str instead of strings)
	// In K8s, this is used to handle versioning, e.g., appsv1 vs appsv2
	upperName := str.ToUpper(clusterName)
	
	fmt.Printf("Normalized Cluster Name: %s\n", upperName)
}