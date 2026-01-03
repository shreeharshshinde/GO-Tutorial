# Module 10 — Go Project Structure (CNCF Style)
## 10.2 `cmd/` vs `pkg/` vs `internal/` — Dependency Boundaries

This document explains **where code should live** in large Go projects and **why**.

If you deeply understand this file, you will:
- stop placing code in the wrong layer
- understand reviewer comments like “wrong package” or “leaks abstraction”
- read Kubernetes codebases without confusion
- design Go systems that scale

This is **architecture**, not style.

---

## 1. The Core Question This Module Answers

> Where should this code live?

Specifically:
- Should it be in `cmd/`?
- Should it be in `pkg/`?
- Should it be in `internal/`?

Choosing wrong causes:
- tight coupling
- unstable APIs
- broken refactors
- PR rejections in CNCF projects

---

## 2. The Mental Model (Memorize This)

Think in terms of **dependency direction**:

cmd/
↓
pkg/ internal/
↓ ↓
(lower-level code)

yaml
Copy code

Rules:
- `cmd/` depends on everything
- `internal/` depends on lower layers
- `pkg/` must NOT depend on `cmd/`
- `internal/` must NOT be imported externally

Dependencies flow **inward**, never outward.

---

## 3. `cmd/` — Application Wiring Layer

Purpose:
- entry points
- binary configuration
- dependency wiring
- flag parsing
- startup orchestration

What belongs here:
- `main.go`
- CLI flags
- environment variable parsing
- selecting implementations

What does NOT belong here:
- business logic
- domain rules
- reusable libraries

Rule of thumb:
> If it’s in `cmd/`, it should be boring.

Example (Kubernetes):
- `cmd/kube-apiserver/main.go` wires the API server
- logic lives elsewhere

---

## 4. `pkg/` — Public API Layer

Purpose:
- reusable libraries
- stable APIs
- cross-project reuse

Putting code in `pkg/` means:
> “Other projects may depend on this.”

Implications:
- breaking changes are expensive
- versioning matters
- API design matters

Kubernetes uses `pkg/` cautiously because public APIs are hard to undo.

Rule:
> Only put code in `pkg/` if you are willing to support it long-term.

---

## 5. `internal/` — Private Implementation Layer

Purpose:
- implementation details
- unstable logic
- glue code
- helpers
- evolving designs

Key feature:
- **compiler enforced privacy**

Only code inside the module can import `internal/`.

This allows:
- aggressive refactoring
- fast iteration
- safe experimentation

Most Kubernetes logic lives here.

---

## 6. Concrete Example

Imagine you’re building a controller.

Where does each part go?

| Component | Location |
|----------|---------|
| main.go | `cmd/controller/main.go` |
| Flag parsing | `cmd/controller` |
| Controller logic | `internal/controller` |
| Queue handling | `internal/queue` |
| Reusable rate limiter | `pkg/ratelimit` (maybe) |

If you’re unsure → **default to `internal/`**.

---

## 7. Why Kubernetes Prefers `internal/`

Kubernetes:
- evolves rapidly
- has huge surface area
- must avoid breaking users

Using `internal/`:
- protects internals
- allows redesigns
- keeps public APIs small

This is intentional and strategic.

---

## 8. Common Beginner Mistakes (Very Important)

- Putting logic in `cmd/` because “it’s simpler”
- Putting everything in `pkg/`
- Avoiding `internal/` because it “feels restrictive”
- Treating folders as namespaces
- Copying layouts blindly

These mistakes scale poorly.

---

## 9. How Reviewers Think (CNCF Perspective)

When reviewers see code, they ask:

- Is this reusable?
- Is this stable?
- Who should be allowed to import this?
- Will this block refactoring later?

Folder placement answers these questions **implicitly**.

---

## 10. Decision Checklist (Use This)

Before placing code, ask:

1. Is this binary-specific?
   → `cmd/`

2. Is this reusable by other projects?
   → `pkg/`

3. Is this internal logic that may change?
   → `internal/`

If unsure → choose `internal/`.

---

## 11. Why This Matters for You

As a contributor:
- Correct placement builds reviewer trust
- Wrong placement raises red flags instantly

As a maintainer:
- Clear boundaries protect the codebase
- Refactors become safe and predictable

Architecture is a long-term investment.

---

## Final Takeaways (Memorize These)

- `cmd/` wires things together
- `pkg/` is a public promise
- `internal/` is your safety net
- Dependency direction matters
- Structure communicates intent

If you get this right, **everything else becomes easier**.