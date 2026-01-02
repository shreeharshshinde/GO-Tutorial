package main

import "fmt"

/*
============================================================
MODULE 08 — MODULES & PACKAGES
STEP 08.4 — PUBLIC VS PRIVATE APIs (COMPLETE DEEP DIVE)
============================================================

This file is EXECUTABLE DOCUMENTATION.

Goal:
Understand how Go defines PUBLIC vs PRIVATE APIs,
why exported names are contracts,
and how CNCF projects like Kubernetes avoid breaking users.

If you understand THIS file,
you will stop accidentally breaking APIs.
*/

/*
------------------------------------------------------------
1. THE ONLY VISIBILITY RULE IN GO
------------------------------------------------------------

Go has exactly ONE rule for visibility:

- Names starting with CAPITAL letters are EXPORTED (public)
- Names starting with lowercase letters are UNEXPORTED (private)

There are:
- no keywords like public / private / protected
- no package-private modifiers

Capitalization IS the API boundary.
*/

// ==========================================================
// 2. PUBLIC TYPES ARE API CONTRACTS
// ==========================================================

// Controller is EXPORTED.
// Once exported, this type becomes part of the public API.
type Controller struct {
	// retries is UNEXPORTED.
	// This field is NOT part of the public contract.
	// We can change it freely.
	retries int
}

/*
IMPORTANT:
Exporting a TYPE means:
- users can construct it
- users can depend on its behavior
- you must support it long-term

CNCF projects are EXTREMELY careful here.
*/

// ==========================================================
// 3. CONSTRUCTORS ARE GATEKEEPERS
// ==========================================================

// NewController is an EXPORTED constructor.
// This is the RECOMMENDED way to create public types.
func NewController() *Controller {
	return &Controller{
		retries: 3,
	}
}

/*
Why constructors matter:
- validation
- default values
- backward compatibility
- controlled evolution

Kubernetes almost NEVER exposes structs
without constructors.
*/

// ==========================================================
// 4. EXPORTED METHODS ARE ALSO CONTRACTS
// ==========================================================

// Run is an EXPORTED method.
// Its behavior is now a public promise.
func (c *Controller) Run() {
	fmt.Println("Controller running with retries =", c.retries)
}

/*
Changing exported method behavior can BREAK users,
even if the signature stays the same.
*/

// ==========================================================
// 5. PRIVATE METHODS ARE SAFE TO CHANGE
// ==========================================================

// resetRetries is UNEXPORTED.
// This can change freely without breaking users.
func (c *Controller) resetRetries() {
	c.retries = 0
}

/*
CNCF codebases have LOTS of private helpers.
This gives maintainers freedom to refactor.
*/

// ==========================================================
// 6. STRUCT FIELDS VS METHODS (CRITICAL DESIGN RULE)
// ==========================================================

/*
BAD DESIGN (do NOT do this in public APIs):

	type Config struct {
		Timeout int
	}

Why this is bad:
- fields cannot be validated
- fields cannot evolve
- fields cannot be deprecated cleanly
*/

// GOOD DESIGN:

type Config struct {
	timeout int // private field
}

// Timeout exposes behavior, not structure.
func (c *Config) Timeout() int {
	return c.timeout
}

/*
Kubernetes prefers:
- private fields
- public methods

This allows API evolution without breakage.
*/

// ==========================================================
// 7. PACKAGES ARE APIs, NOT FOLDERS
// ==========================================================

/*
In Go:
- the PACKAGE is the API unit
- not the file
- not the struct

Exporting a package means:
- everything exported inside it matters
- folder structure becomes part of design

This is why Kubernetes is very careful
with package names and locations.
*/

// ==========================================================
// 8. DEFAULT TO PRIVATE (CNCF RULE)
// ==========================================================

/*
CNCF rule of thumb:

- Start unexported
- Export ONLY when required
- Assume exported code is forever

Common review comment in Kubernetes:
"Does this really need to be exported?"
*/

// ==========================================================
// 9. COMMON BEGINNER MISTAKES (DO NOT DO THESE)
// ==========================================================

/*
❌ Exporting everything "just in case"
❌ Exposing struct fields directly
❌ Changing exported behavior casually
❌ Treating packages like internal folders
❌ Forgetting that capitalization is a contract
*/

// ==========================================================
// 10. HOW KUBERNETES APPLIES THIS
// ==========================================================

/*
Kubernetes design patterns:

- Very small public APIs
- Heavy use of unexported helpers
- Constructors for public types
- Minimal exported fields
- Long deprecation cycles

This is how Kubernetes stays stable
while evolving rapidly.
*/

// ==========================================================
// 11. MAIN — DEMONSTRATION
// ==========================================================

func main() {
	fmt.Println("=== 08.4 Public vs Private APIs ===")

	ctrl := NewController()
	ctrl.Run()

	fmt.Println()
	fmt.Println("Key insight:")
	fmt.Println("Exported names are promises.")
	fmt.Println("Private names give freedom.")
}

/*
============================================================
FINAL TAKEAWAYS (MEMORIZE THESE)
============================================================

1. Capital letter = public API
2. Public APIs are long-term contracts
3. Private APIs are safe to change
4. Exported fields are dangerous
5. Prefer constructors over struct literals
6. Packages define API boundaries
7. CNCF projects default to private

If you internalize THIS file,
you will design safe, evolvable Go APIs.
*/
