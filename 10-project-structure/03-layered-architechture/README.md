# Module 10 — Go Project Structure (CNCF Style)
## 10.3 Layered Architecture — How Large Go Systems Are Organized

This document explains **how large Go systems are layered internally** and why
Kubernetes, Kubeflow, and other CNCF projects look “indirect” or “over-abstracted”
to newcomers.

This is not accidental complexity.
This is **controlled complexity**.

---

## 1. The Problem Layered Architecture Solves

In small projects:
- everything can talk to everything
- shortcuts are tempting
- changes are cheap

In large projects:
- tight coupling destroys velocity
- changes ripple unpredictably
- testing becomes painful
- refactoring becomes impossible

Layered architecture exists to **control dependency flow**.

---

## 2. The Core Rule (Memorize This)

> **Dependencies flow inward, never outward.**

Higher-level layers:
- depend on lower-level layers

Lower-level layers:
- must NOT depend on higher-level layers

This rule is enforced socially (code review) and structurally (packages).

---

## 3. A Typical CNCF Layer Stack

Most large Go systems follow a variation of this:

+---------------------------+
| cmd / main |
+---------------------------+
| orchestration / wiring |
+---------------------------+
| application logic |
+---------------------------+
| domain logic |
+---------------------------+
| infrastructure / adapters |
+---------------------------+

yaml
Copy code

Each layer has **one job**.

---

## 4. `cmd/` — The Entry Layer

Responsibilities:
- program startup
- flag parsing
- environment configuration
- dependency wiring
- lifecycle management

Characteristics:
- very thin
- no business logic
- hard to unit test (and that’s fine)

In Kubernetes:
- `cmd/kube-apiserver`
- `cmd/kube-controller-manager`

---

## 5. Orchestration Layer (Glue Code)

This layer:
- connects components
- wires interfaces to implementations
- manages lifecycles

Often lives in:
- `internal/app`
- `internal/server`
- `internal/controller`

This is where:
- controllers are assembled
- informers are attached
- queues are initialized

---

## 6. Application Logic Layer

This layer answers:
> “What should happen?”

Examples:
- reconcile loops
- workflows
- use-case orchestration
- business rules at a high level

This code:
- depends on interfaces
- is highly testable
- contains little infrastructure detail

In Kubernetes:
- controller reconcile logic
- admission control logic

---

## 7. Domain Logic Layer (Core Rules)

This is the **heart of the system**.

Characteristics:
- pure logic
- minimal dependencies
- no I/O
- no networking
- no Kubernetes clients

Examples:
- validation rules
- state transitions
- invariants
- scheduling decisions

This layer should be:
- easy to test
- stable over time
- independent of frameworks

---

## 8. Infrastructure / Adapter Layer

This layer deals with the outside world:
- Kubernetes API
- databases
- filesystems
- network calls

Characteristics:
- messy
- failure-prone
- slow
- highly mocked or faked in tests

Lives in:
- `internal/client`
- `internal/store`
- `internal/kube`
- `internal/io`

---

## 9. Why Kubernetes Looks “Indirect”

Newcomers often ask:
> “Why doesn’t Kubernetes just do the thing directly?”

Answer:
- because direct calls create tight coupling
- because logic must be testable
- because APIs evolve
- because infrastructure fails

Indirection is **intentional insulation**.

---

## 10. Example: Kubernetes Controller Flow

Simplified flow:

cmd/
↓
controller setup
↓
reconcile loop (application logic)
↓
domain decision
↓
kube client call (infrastructure)

yaml
Copy code

Each layer:
- knows only what it needs
- does one job well

---

## 11. Testing Benefits of Layers

Layered design enables:
- unit tests for domain logic
- fake infrastructure for controllers
- fast feedback loops
- fewer flaky tests

This is why CNCF projects invest heavily in structure.

---

## 12. Common Beginner Mistakes

- Putting infrastructure logic in domain code
- Letting low-level packages import high-level ones
- Skipping layers “for simplicity”
- Mixing concerns in one package
- Treating layering as optional

Layer violations scale poorly.

---

## 13. How Reviewers Detect Layer Violations

Reviewers look for:
- suspicious imports
- infrastructure leaking into logic
- business rules inside clients
- circular dependencies

These are architectural red flags.

---

## 14. Practical Guidance

When adding code, ask:
- What layer does this belong to?
- What should this code depend on?
- Who should be allowed to import it?

If unsure:
- push logic inward
- keep infrastructure outward

---

## Final Takeaways (Memorize These)

- Layers control dependency direction
- Inward dependencies only
- Domain logic stays pure
- Infrastructure stays isolated
- Indirection enables scale
- Kubernetes is layered by necessity, not preference

If you understand layering, large Go codebases stop feeling chaotic.