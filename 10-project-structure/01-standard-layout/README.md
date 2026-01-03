# Module 10 — Go Project Structure (CNCF Style)
## 10.1 Standard Go Project Layout (Reality vs Myth)

This document explains **how real-world Go projects are structured**, especially
large **CNCF projects like Kubernetes and Kubeflow**.

This is not about aesthetics or conventions.
This is about **scaling teams, codebases, and maintenance over years**.

If you understand this document, you will:
- stop feeling lost in large Go repos
- understand why code lives where it does
- avoid common structural PR rejections

---

## 1. The Big Myth: “Go Has a Standard Layout”

Go does **not** have an officially mandated project structure.

There is:
- no required folder layout
- no enforced architecture
- no framework-level opinion

What *does* exist is:
- strong community convergence
- patterns that survived large-scale use
- CNCF-driven best practices

If you see the same layout everywhere, it’s because it **worked at scale**, not
because Go requires it.

---

## 2. The De Facto CNCF Project Layout

Most large CNCF Go projects look roughly like this:

repo/
├── cmd/
├── pkg/
├── internal/
├── api/
├── configs/
├── test/
├── go.mod
├── go.sum

yaml
Copy code

Not all directories are required.
Each directory exists for a **specific architectural reason**.

---

## 3. `cmd/` — Entry Points (Binaries)

The `cmd/` directory contains **entry points**.

Each subdirectory usually builds **one binary**.

Example from Kubernetes:

cmd/
├── kube-apiserver/
│ └── main.go
├── kube-controller-manager/
│ └── main.go

yaml
Copy code

Rules for `cmd/`:
- contains `main` packages
- minimal logic
- mostly wiring and configuration
- no business logic

If logic grows large inside `cmd/`, it’s usually a design smell.

---

## 4. `pkg/` — Public, Reusable Libraries

`pkg/` contains **code intended for reuse by other projects**.

Putting code in `pkg/` is a promise:
> “This API is safe for external users to depend on.”

Characteristics:
- stable APIs
- careful versioning
- backward compatibility matters

Kubernetes uses `pkg/` sparingly because public APIs are expensive to maintain.

---

## 5. `internal/` — Private Implementation

`internal/` contains **implementation details**.

This is not just a convention — it is **compiler enforced**.

Rules:
- only code within the parent module can import `internal/`
- external users are blocked at compile time

Purpose:
- safe refactoring
- hiding unstable logic
- protecting architecture

This is one of the most important tools used by CNCF projects.

---

## 6. `api/` — API Types and Schemas

The `api/` directory usually contains:
- versioned API types
- request/response schemas
- CRD definitions

Example (Kubernetes):

api/
├── core/
│ └── v1/
├── apps/
│ └── v1/

yaml
Copy code

Code in `api/` is treated as a **contract**.
Changes here are heavily scrutinized.

---

## 7. `test/` — Integration and End-to-End Tests

The `test/` directory is **not** for unit tests.

Unit tests:
- live next to code
- use `*_test.go`

`test/` usually contains:
- integration tests
- end-to-end tests
- environment setup
- cluster bootstrapping logic

Kubernetes uses this directory extensively.

---

## 8. `configs/` — Configuration and Manifests

`configs/` typically stores:
- sample configuration files
- YAML manifests
- deployment examples

These are often consumed by:
- users
- CI pipelines
- documentation

---

## 9. What Go the Language Does *Not* Care About

Go does **not** care about:
- folder names (except `internal/`)
- where `main.go` lives
- repository layout

The compiler only cares about:
- packages
- imports
- module boundaries

Humans care about structure.
At scale, structure matters more than syntax.

---

## 10. Why CNCF Projects Care Deeply About Structure

At CNCF scale, projects have:
- thousands of files
- hundreds of contributors
- multi-year lifetimes

Structure enables:
- safe parallel development
- predictable ownership
- fast code reviews
- reduced cognitive load

Without structure, large projects collapse under their own weight.

---

## 11. Common Beginner Mistakes

- Putting business logic in `cmd/`
- Exposing everything via `pkg/`
- Avoiding `internal/`
- Treating folders as namespaces
- Copying layouts without understanding intent

Structure without reasoning is harmful.

---

## 12. Your Learning Repository

Your current learning repository:
- single module
- topic-based folders
- progressive structure

This is **perfect for learning**.

Production repositories evolve from this foundation, not the other way around.

---

## Final Takeaways (Memorize These)

- Go has no official layout
- CNCF provides de facto standards
- `cmd/` = binaries
- `pkg/` = public APIs
- `internal/` = private implementation
- `api/` = contracts
- Structure exists to support scale, not aesthetics

Understanding structure is the **first step to fearlessly reading large Go projects**.