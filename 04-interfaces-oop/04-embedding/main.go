package main

import "fmt"

/*
This file explains EMBEDDING in Go.

Embedding is how Go achieves:
- Code reuse
- Composition over inheritance
- Shared behavior without class hierarchies

This is CRUCIAL for Kubernetes:
Every Kubernetes resource embeds ObjectMeta, which is why
Pod, Service, Deployment all have Name, Namespace, Labels, etc.
*/

// ==========================================
// 1. STRUCT EMBEDDING (COMPOSITION)
// ==========================================

// Logger is a reusable component that provides logging behavior.
type Logger struct {
	Level string
}

// Log is a method on Logger.
// When Logger is embedded, this method can be promoted.
func (l *Logger) Log(msg string) {
	fmt.Printf("[%s] %s\n", l.Level, msg)
}

// Identity is a reusable component that provides identity fields.
type Identity struct {
	ID   int
	Name string
}

/*
User is a composite type.

Instead of "extending" Identity or Logger,
we EMBED them.

Embedding rules:
- The embedded type name becomes a field name
- Methods and fields may be PROMOTED
*/
type User struct {
	// Anonymous fields (no explicit field name)
	Identity
	Logger

	Email string
}

// ==========================================
// 2. FIELD & METHOD PROMOTION
// ==========================================
// Fields and methods of embedded structs can be accessed
// directly on the outer struct.

// ==========================================
// 3. NAME SHADOWING (FIELD CONFLICTS)
// ==========================================

// Admin embeds Identity but also defines its own Name field.
type Admin struct {
	Identity // has Name
	Name     string // shadows Identity.Name
	Role     string
}

/*
Shadowing rules:
- The outer field wins
- The embedded field is still accessible via qualification
*/

func main() {
	fmt.Println("--- 1. Field & Method Promotion ---")

	u := User{
		Identity: Identity{ID: 101, Name: "Alice"},
		Logger:   Logger{Level: "INFO"},
		Email:    "alice@example.com",
	}

	// Field promotion:
	// ID and Name come from Identity
	fmt.Printf("User: %s (ID: %d)\n", u.Name, u.ID)

	// Method promotion:
	// Log() comes from Logger
	u.Log("User logged in successfully.")

	/*
	What actually happens:
	u.Log(...) is rewritten by Go as:
	u.Logger.Log(...)
	*/

	fmt.Println("\n--- 2. Accessing Embedded Struct Explicitly ---")
	// You can always access the embedded struct directly
	fmt.Println("Logger Level:", u.Logger.Level)
	fmt.Println("Identity Name:", u.Identity.Name)

	fmt.Println("\n--- 3. Name Shadowing (Conflict Resolution) ---")

	a := Admin{
		Identity: Identity{Name: "Inner Identity Name"},
		Name:     "Outer Admin Name",
		Role:     "SuperUser",
	}

	// Outer field takes priority
	fmt.Println("Admin Name (Outer):", a.Name)

	// Inner field must be accessed explicitly
	fmt.Println("Identity Name (Inner):", a.Identity.Name)

	// Both fields coexist safely
	fmt.Println("Admin Role:", a.Role)

	fmt.Println("\n--- 4. Embedding Enables Interface Satisfaction ---")
	/*
	If an embedded struct implements an interface,
	the outer struct ALSO implements it.

	This is widely used in Kubernetes:
	- ObjectMeta implements interfaces
	- Pod embeds ObjectMeta
	- Pod automatically satisfies those interfaces
	*/

	// User implements logging behavior without redefining Log()
	var logger interface {
		Log(string)
	}

	logger = &u
	logger.Log("Logging via interface")

	fmt.Println("\n--- 5. Embedding vs Inheritance (Key Differences) ---")
	/*
	Go embedding is NOT inheritance:
	- No 'is-a' relationship
	- No method overriding
	- No super / base keyword

	It is:
	- has-a relationship
	- explicit composition
	- predictable behavior
	*/

	fmt.Println("\n--- 6. Pointer vs Value Embedding ---")
	/*
	You can embed:
	- Identity
	- *Identity

	Embedding a pointer allows shared mutation.
	Kubernetes commonly embeds pointers for metadata.
	*/

	type Resource struct {
		*Identity
	}

	r := Resource{
		Identity: &Identity{ID: 500, Name: "resource"},
	}

	fmt.Println("Resource Name:", r.Name)
	r.Name = "updated resource"
	fmt.Println("Updated Name:", r.Identity.Name)
}
