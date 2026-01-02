package main

import "fmt"

/*
============================================================
MODULE 08 — MODULES & PACKAGES
STEP 08.5 — MULTI-MODULE REPOSITORIES (KUBERNETES STYLE)
============================================================

This file is EXECUTABLE DOCUMENTATION.

Goal:
Understand why large Go projects (Kubernetes, Kubeflow, etc.)
use MULTIPLE go.mod files in a single repository,
and when YOU should (and should NOT) do this.

This is advanced Go architecture knowledge.
*/

/*
------------------------------------------------------------
1. THE SIMPLE WORLD (SINGLE MODULE)
------------------------------------------------------------

Most Go projects start like this:

repo/
├── go.mod
├── main.go
├── pkg/
└── internal/

This is PERFECT for:
- small to medium projects
- applications
- learning
- most services

You should default to THIS.
*/

/*
------------------------------------------------------------
2. THE PROBLEM AT CNCF SCALE
------------------------------------------------------------

Kubernetes is NOT:
- one binary
- one library
- one consumer

It is:
- libraries
- CLIs
- controllers
- APIs
- reused by thousands of projects

Problems with a single module:
- versioning everything together
- huge dependency graph
- slow CI
- forced upgrades for all consumers
*/

/*
------------------------------------------------------------
3. THE SOLUTION: MULTI-MODULE REPOS
------------------------------------------------------------

A multi-module repo contains MULTIPLE go.mod files.

Each go.mod defines:
- a version boundary
- an API surface
- an independent lifecycle
*/

/*
------------------------------------------------------------
4. REAL KUBERNETES EXAMPLE (SIMPLIFIED)
------------------------------------------------------------

kubernetes/
├── go.mod                  (root tooling / internal builds)
├── staging/
│   ├── src/
│   │   ├── k8s.io/api/
│   │   │   └── go.mod
│   │   ├── k8s.io/apimachinery/
│   │   │   └── go.mod
│   │   └── k8s.io/client-go/
│   │       └── go.mod

Each of these is a SEPARATE MODULE.
*/

/*
------------------------------------------------------------
5. WHY THIS MATTERS
------------------------------------------------------------

Example:
- client-go users ONLY need client-go
- they do NOT need kube-apiserver internals
- they do NOT want massive dependency graphs

Multi-module repos allow:
- independent versioning
- smaller dependency trees
- faster builds
- safer upgrades
*/

/*
------------------------------------------------------------
6. VERSIONING BENEFITS
------------------------------------------------------------

With multi-module repos:
- client-go can be v0.30.1
- api can be v0.30.1
- internal tooling can evolve freely

Breaking client-go does NOT require:
- breaking Kubernetes internals
*/

/*
------------------------------------------------------------
7. HOW MODULES DEPEND ON EACH OTHER
------------------------------------------------------------

Inside the repo, modules often use replace directives:

replace k8s.io/apimachinery => ../apimachinery

This allows:
- local development
- single-repo workflows

IMPORTANT:
Replace directives are usually NOT published.
They are cleaned before release.
*/

/*
------------------------------------------------------------
8. WHEN YOU SHOULD USE MULTI-MODULE
------------------------------------------------------------

GOOD reasons:
- multiple public libraries
- independent consumers
- different stability guarantees
- different release cycles

BAD reasons:
- “it feels cleaner”
- premature optimization
- avoiding refactors
*/

/*
------------------------------------------------------------
9. COMMON BEGINNER MISTAKES
------------------------------------------------------------

❌ Creating many go.mod files too early
❌ Splitting modules without versioning needs
❌ Forgetting replace cleanup
❌ Circular dependencies between modules
❌ Treating modules like folders

Multi-module is an ADVANCED tool.
*/

/*
------------------------------------------------------------
10. CNCF CONTRIBUTOR MENTAL MODEL
------------------------------------------------------------

When reading CNCF repos:
- Identify module boundaries FIRST
- Check which go.mod you are in
- Understand public vs internal modules
- Do NOT casually move code across modules

Many PR mistakes come from ignoring this.
*/

/*
------------------------------------------------------------
11. YOU AS A CONTRIBUTOR
------------------------------------------------------------

In Kubeflow / Kubernetes PRs:
- you usually touch ONE module
- you run go test within that module
- you do NOT restructure modules casually
- module changes require maintainer discussion
*/

/*
------------------------------------------------------------
12. HOW THIS AFFECTS YOUR LEARNING REPO
------------------------------------------------------------

For YOU right now:
- ONE module is correct
- You are doing the right thing
- Multi-module knowledge is for READING repos
- Not for premature structuring
*/

/*
------------------------------------------------------------
13. RELATION TO PREVIOUS TOPICS
------------------------------------------------------------

Multi-module repos build on:
- Semantic Versioning (08.2)
- internal/ packages (08.3)
- Public vs Private APIs (08.4)

These topics ONLY make sense together.
*/

/*
------------------------------------------------------------
14. MAIN — EXECUTION
------------------------------------------------------------
*/

func main() {
	fmt.Println("=== 08.5 Multi-Module Repositories ===")

	fmt.Println("Large CNCF projects use multiple go.mod files")
	fmt.Println("to control versioning, dependencies, and stability.")

	fmt.Println()
	fmt.Println("Rule of thumb:")
	fmt.Println("Start single-module.")
	fmt.Println("Split only when consumers demand it.")
}

/*
============================================================
FINAL TAKEAWAYS (CRITICAL)
============================================================

1. A Go module is a version boundary
2. Large repos often contain many modules
3. Kubernetes uses multi-module architecture heavily
4. Each module has its own go.mod
5. Multi-module enables independent evolution
6. This is an ADVANCED tool — use carefully
7. Understanding it makes CNCF repos readable

With this, Module 08 is COMPLETE.
*/
