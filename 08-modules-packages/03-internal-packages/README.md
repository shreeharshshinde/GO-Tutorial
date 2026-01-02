# 08.3 — `internal/` Packages (Enforced Encapsulation in Go)

This section explains **how Go enforces architectural boundaries**
without classes, inheritance, or access modifiers like `protected`.

The `internal/` directory is a **compiler-enforced rule**, not a convention.

Understanding this is mandatory for reading Kubernetes, Kubeflow,
and most CNCF Go codebases.

---

## 1. The Problem Go Is Solving

In large systems:
- Not everything should be public
- Some code is implementation detail
- Packages should not be imported arbitrarily

Other languages use:
- `protected`
- `package-private`
- complex visibility rules

Go uses **`internal/`**.

---

## 2. What `internal/` Means (Exact Rule)

Any package inside a directory named `internal`:

<module-root>/internal/foo

pgsql
Copy code

❗ **Can only be imported by code inside:**

<module-root>/...

yaml
Copy code

If code outside that tree tries to import it → **compile-time error**.

This is enforced by the Go compiler.

---

## 3. Example Directory Layout

my-module/
├── go.mod
├── main.go
├── api/
│ └── server.go
├── internal/
│ ├── config/
│ │ └── loader.go
│ └── cache/
│ └── cache.go

yaml
Copy code

Allowed imports:
- `api` → `internal/config`
- `main` → `internal/cache`

❌ Not allowed:
- Any other module importing `internal/config`

---

## 4. Why This Is Better Than Conventions

Without `internal/`:
- Developers rely on comments (“do not import this”)
- Mistakes happen
- Refactoring becomes impossible

With `internal/`:
- Compiler enforces boundaries
- Architecture is protected
- Large teams can move fast safely

---

## 5. How Kubernetes Uses `internal/`

Kubernetes uses `internal/` for:
- internal helpers
- unstable APIs
- implementation details
- things NOT meant for users

Public APIs live in:
- `pkg/`
- `api/`
- `client-go`

Private logic lives in:
- `internal/`

This prevents accidental coupling.

---

## 6. `internal/` vs Unexported Identifiers

Important distinction:

- **Unexported names** (`lowercase`)  
  → visible within the *same package*

- **`internal/` packages**  
  → visible only within a *module subtree*

They solve **different problems**.

---

## 7. Common Beginner Mistakes

❌ Putting reusable code in `internal/`  
❌ Importing `internal/` from another module  
❌ Using `internal/` to “hide” bad APIs  
❌ Overusing `internal/` everywhere  

`internal/` is for **architecture**, not fear.

---

## 8. CNCF Design Philosophy

In CNCF projects:
- Public APIs are sacred
- Internal code can change freely
- `internal/` protects maintainers from breaking users

This is how Kubernetes evolves without chaos.

---

## Summary (Memorize This)

- `internal/` is enforced by the compiler
- It creates module-level encapsulation
- It protects architecture, not variables
- CNCF projects rely on it heavily
- If you see `internal/`, do not import it casually

If you understand this, Go codebases feel **much safer**.