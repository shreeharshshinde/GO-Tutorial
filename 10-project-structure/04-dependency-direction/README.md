# Module 10 — Go Project Structure (CNCF Style)
## 10.4 Dependency Direction & Inversion — Why Interfaces Exist Everywhere

This document explains **dependency direction**, **dependency inversion**, and why
large Go projects (especially Kubernetes) rely heavily on interfaces.

If you understand this file, you will:
- stop being confused by “too many interfaces”
- understand why controllers look indirect
- design testable, maintainable Go systems
- write PRs that align with CNCF expectations

---

## 1. The Core Problem: Tight Coupling

In tightly coupled systems:
- high-level logic depends on low-level details
- changes ripple unpredictably
- testing becomes hard
- refactoring becomes risky

Example of tight coupling:
- controller directly creates a Kubernetes client
- business logic directly touches the API
- code cannot be tested without real infrastructure

This does not scale.

---

## 2. The Fundamental Rule (Repeat This)

> **High-level code must not depend on low-level code.  
> Both should depend on abstractions.**

This is called **Dependency Inversion**.

This principle is language-agnostic.
Go simply makes it *visible*.

---

## 3. What “Dependency Direction” Means in Practice

Dependency direction answers:
> Who is allowed to import whom?

Correct direction:

Application Logic
↓
Interfaces
↓
Infrastructure

css
Copy code

Incorrect direction:

Application Logic → Infrastructure (BAD)

yaml
Copy code

When direction is wrong:
- logic becomes untestable
- infra leaks everywhere
- changes become expensive

---

## 4. Interfaces Are NOT for Polymorphism First

In Go, interfaces are primarily for:
- dependency inversion
- decoupling
- testability

They are NOT:
- class hierarchies
- OOP inheritance replacements

Interfaces exist to **invert ownership of dependencies**.

---

## 5. Who Owns the Interface? (Very Important)

Rule:
> **The consumer owns the interface, not the implementation.**

This is a defining Go idiom.

Example:
- controller defines `Client` interface
- Kubernetes client implements it
- fake client implements it for tests

NOT the other way around.

---

## 6. Kubernetes Example (Mental Model)

In Kubernetes controllers:

- Controller logic depends on an interface
- Real client implements that interface
- Fake client implements the same interface in tests

This allows:
- real cluster usage
- unit testing without a cluster
- clean separation of concerns

---

## 7. Dependency Injection in Go (No Frameworks)

Go uses **explicit dependency injection**.

Instead of:
- global variables
- service locators
- magic frameworks

Go prefers:
- constructors
- function parameters
- explicit wiring in `cmd/`

This is why `main.go` looks verbose.

That verbosity is **intentional clarity**.

---

## 8. Example Dependency Flow (Controller)

Typical flow:

cmd/main.go
├── create real kube client
├── create informers
├── create controller
└── inject dependencies

yaml
Copy code

Controller code:
- does NOT know how clients are created
- does NOT care if it’s real or fake
- only depends on behavior

---

## 9. Why Kubernetes Avoids Global State

Global state:
- hides dependencies
- breaks test isolation
- makes concurrency dangerous

Dependency injection:
- makes dependencies explicit
- enables deterministic tests
- simplifies reasoning

This is why globals are discouraged.

---

## 10. How This Affects Testing

Good dependency direction enables:
- mock implementations
- fake clients
- fast unit tests
- clear failure modes

Bad dependency direction leads to:
- integration-only tests
- slow feedback
- flaky CI
- untestable code

---

## 11. Common Beginner Misunderstandings

- “Why so many interfaces?”
- “Why not just call the client directly?”
- “This seems over-engineered”

At scale:
- direct calls don’t survive
- interfaces reduce long-term complexity
- architecture pays off over years

---

## 12. How Reviewers Spot Problems Instantly

CNCF reviewers look for:
- imports from infra into logic
- controllers creating clients internally
- interfaces defined in implementation packages
- globals instead of injected dependencies

These are architectural red flags.

---

## 13. Practical Checklist

Before writing code, ask:
- What layer am I in?
- What should I depend on?
- Can this be tested without infrastructure?
- Who owns this interface?

If unsure:
- depend on an interface
- inject dependencies from above

---

## Final Takeaways (Memorize These)

- Dependency direction matters more than code style
- High-level logic must not depend on low-level details
- Interfaces exist to invert dependencies
- Consumers own interfaces
- Explicit wiring beats magic
- Kubernetes architecture is interface-driven by necessity

Once this clicks, Go architecture becomes predictable instead of confusing.