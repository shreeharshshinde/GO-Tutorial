package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// --- 1. Basic Struct ---
// Structs are collections of fields. They are Value Types (copied when passed).
type Server struct {
	Name string
	IP   string
	Port int
}

// --- 2. Composition (Embedding) ---
// Go does not have Inheritance (no "extends" keyword).
// Instead, it uses Composition. We "embed" one struct into another.
type BaseConfig struct {
	Environment string // e.g., "Production", "Dev"
	DebugMode   bool
}

type DatabaseConfig struct {
	// Embedding BaseConfig allows us to access its fields directly
	BaseConfig

	DBName     string
	Connection string
}

// --- 3. JSON Tags & Serialization ---
// Capitalized fields are EXPORTED (Public). Lowercase fields are PRIVATE.
// The `json:"..."` tag tells the encoder how to name the field in JSON.
type APIResponse struct {
	Status  int    `json:"status_code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
	secret  string `json:"-"`
}

func main() {
	fmt.Println("--- 1. Struct Composition ---")
	db := DatabaseConfig{
		BaseConfig: BaseConfig{
			Environment: "Production",
			DebugMode:   false,
		},
		DBName:     "users_db",
		Connection: "postgres://localhost:5432",
	}

	fmt.Printf("Env: %s, DB: %s\n", db.Environment, db.DBName)

	fmt.Println("\n--- 2. JSON Marshaling (Go -> JSON) ---")
	response := APIResponse{
		Status:  200,
		Message: "Success",
		secret:  "this-will-not-show-up",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

	fmt.Println("\n--- 3. JSON Unmarshaling (JSON -> Go) ---")
	jsonInput := `{"status_code": 404, "message": "Not Found", "extra_field": "ignored"}`

	var incomingResp APIResponse
	if err := json.Unmarshal([]byte(jsonInput), &incomingResp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed Struct: %+v\n", incomingResp)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 4. Zero Value of Structs ---")
	// --------------------------------------------------------------------

	var s Server
	fmt.Printf("Zero-value struct: %+v\n", s)
	// All fields are zero-values; no constructor needed

	// --------------------------------------------------------------------
	fmt.Println("\n--- 5. Structs Are Value Types (Copy Semantics) ---")
	// --------------------------------------------------------------------

	s1 := Server{Name: "A", IP: "1.1.1.1", Port: 80}
	s2 := s1
	s2.Port = 443

	fmt.Println("s1:", s1)
	fmt.Println("s2:", s2)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 6. Struct Pointers (Mutation) ---")
	// --------------------------------------------------------------------

	updatePort(&s1)
	fmt.Println("After pointer update:", s1)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 7. Pointer vs Value Receiver (Conceptual) ---")
	// --------------------------------------------------------------------
	// Value receiver -> operates on copy
	// Pointer receiver -> mutates original
	// Rule: If method mutates or struct is large, use pointer receiver

	// --------------------------------------------------------------------
	fmt.Println("\n--- 8. Embedded Structs with Pointer ---")
	// --------------------------------------------------------------------

	type AppConfig struct {
		*BaseConfig
		Name string
	}

	app := AppConfig{
		BaseConfig: &BaseConfig{
			Environment: "Dev",
			DebugMode:   true,
		},
		Name: "API",
	}

	fmt.Println("App Env:", app.Environment)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 9. Anonymous Structs ---")
	// --------------------------------------------------------------------

	temp := struct {
		ID   int
		Name string
	}{
		ID:   1,
		Name: "temp",
	}
	fmt.Printf("Anonymous struct: %+v\n", temp)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 10. Struct Comparability ---")
	// --------------------------------------------------------------------

	type A struct {
		X int
		Y string
	}

	a1 := A{1, "x"}
	a2 := A{1, "x"}
	fmt.Println("a1 == a2:", a1 == a2)

	// Structs are NOT comparable if they contain:
	// slices, maps, funcs

	// --------------------------------------------------------------------
	fmt.Println("\n--- 11. Structs Containing Slices (Aliasing Trap) ---")
	// --------------------------------------------------------------------

	type Pod struct {
		Containers []string
	}

	p1 := Pod{Containers: []string{"c1", "c2"}}
	p2 := p1 // shallow copy
	p2.Containers[0] = "evil"

	fmt.Println("p1:", p1)
	fmt.Println("p2:", p2)

	// Correct deep copy
	p3 := Pod{
		Containers: append([]string(nil), p1.Containers...),
	}
	p3.Containers[0] = "safe"

	fmt.Println("p1 after deep copy:", p1)
	fmt.Println("p3:", p3)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 12. Structs with Maps (Same Aliasing Rule) ---")
	// --------------------------------------------------------------------

	type Cache struct {
		Data map[string]int
	}

	c1 := Cache{Data: map[string]int{"a": 1}}
	c2 := c1
	c2.Data["a"] = 999

	fmt.Println("c1:", c1)
	fmt.Println("c2:", c2)

	// Deep copy required for maps too

	// --------------------------------------------------------------------
	fmt.Println("\n--- 13. omitempty Zero-Value Trap ---")
	// --------------------------------------------------------------------

	resp := APIResponse{
		Status:  200,
		Message: "OK",
		Data:    "",
	}

	b, _ := json.Marshal(resp)
	fmt.Println(string(b))
	// Empty string omitted, not serialized

	// --------------------------------------------------------------------
	fmt.Println("\n--- 14. Unknown JSON Fields Are Ignored ---")
	// --------------------------------------------------------------------

	jsonExtra := `{"status_code":200,"message":"OK","unknown":"ignored"}`
	var r APIResponse
	_ = json.Unmarshal([]byte(jsonExtra), &r)
	fmt.Printf("After unknown field JSON: %+v\n", r)

	// --------------------------------------------------------------------
	fmt.Println("\n--- 15. JSON Pointer Requirement ---")
	// --------------------------------------------------------------------
	// json.Unmarshal REQUIRES pointer
	// json.Unmarshal(data, value) ❌
	// json.Unmarshal(data, &value) ✅

	// --------------------------------------------------------------------
	fmt.Println("\n--- 16. Struct Alignment & Memory (Advanced) ---")
	// --------------------------------------------------------------------
	// Field order affects memory size due to padding
	// Place larger fields first for memory efficiency

	// --------------------------------------------------------------------
	fmt.Println("\n--- 17. Export Rules (Package Boundary) ---")
	// --------------------------------------------------------------------
	// Uppercase fields & structs are visible across packages
	// Lowercase fields are package-private

	// --------------------------------------------------------------------
	fmt.Println("\n--- 18. Tags Are Compile-Time Metadata ---")
	// --------------------------------------------------------------------
	// Struct tags are strings; Go does NOT validate them
	// Typos silently break behavior

	// --------------------------------------------------------------------
	fmt.Println("\n--- 19. Struct Equality != Semantic Equality ---")
	// --------------------------------------------------------------------
	// Two structs can be equal in fields but not semantically equal
	// Example: timestamps, cached fields, internal state

	// --------------------------------------------------------------------
	fmt.Println("\n--- 20. When NOT to Use Structs ---")
	// --------------------------------------------------------------------
	// Use structs for:
	// - Data modeling
	// - API objects
	// - Configuration
	// Avoid structs for:
	// - Behavior-only abstractions (use interfaces)
}

// -------------------- Helper Functions --------------------

func updatePort(s *Server) {
	s.Port = 8080
}
