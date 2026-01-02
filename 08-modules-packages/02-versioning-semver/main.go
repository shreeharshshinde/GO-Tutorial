package main

import "fmt"

/*
============================================================
MODULE 08 — MODULES & PACKAGES
STEP 08.2 — VERSIONING & SEMANTIC IMPORT VERSIONING
============================================================

This file explains:
- How Go versions modules
- Why v2+ appears in import paths
- How Semantic Versioning actually works in Go
- Why CNCF projects (Kubernetes) are extremely strict here

If you understand THIS file,
you will not break public APIs accidentally.
*/

/*
------------------------------------------------------------
1. SEMANTIC VERSIONING (SemVer) — RECAP
------------------------------------------------------------

Format:
	vMAJOR.MINOR.PATCH

Rules:
- MAJOR: breaking API changes
- MINOR: backward-compatible features
- PATCH: backward-compatible bug fixes

Example:
	v1.2.3
*/

/*
------------------------------------------------------------
2. HOW GO INTERPRETS VERSIONS
------------------------------------------------------------

Go takes Semantic Versioning VERY seriously.

Critical rule:
- MAJOR versions v2+ are considered BREAKING
- Therefore, they MUST be explicit in import paths

This is called:
SEMANTIC IMPORT VERSIONING
*/

/*
------------------------------------------------------------
3. THE MOST IMPORTANT RULE (MEMORIZE)
------------------------------------------------------------

For Go modules:

v0 or v1:
	import "example.com/mylib"

v2+:
	import "example.com/mylib/v2"

The MAJOR version becomes PART OF THE IMPORT PATH.
*/

/*
------------------------------------------------------------
4. WHY THIS RULE EXISTS
------------------------------------------------------------

Problem Go is solving:

Imagine:
- Library v1 exists
- Library v2 breaks APIs
- Two dependencies want different versions

Without versioned imports:
- Dependency hell
- Diamond dependency problem
- Runtime crashes

Go fixes this at COMPILE TIME.
*/

/*
------------------------------------------------------------
5. WHAT HAPPENS IF YOU TAG v2 WITHOUT /v2
------------------------------------------------------------

If a module author does this:

	git tag v2.0.0

BUT keeps:

	module example.com/mylib

Go will:
- REJECT the module
- Refuse to download it

This is a HARD ERROR.

Go enforces correctness at the ecosystem level.
*/

/*
------------------------------------------------------------
6. CORRECT v2+ MODULE LAYOUT
------------------------------------------------------------

Repo structure:

example.com/mylib
├── go.mod          (module example.com/mylib)
├── feature.go      (v1 API)
├── v2/
│   ├── go.mod      (module example.com/mylib/v2)
│   └── feature.go  (v2 API)

Each MAJOR version:
- Has its OWN go.mod
- Is a DIFFERENT module
*/

/*
------------------------------------------------------------
7. HOW CONSUMERS USE v2 MODULES
------------------------------------------------------------

In go.mod:

	require example.com/mylib/v2 v2.1.0

In code:

	import "example.com/mylib/v2"

This allows:
- v1 and v2 to coexist
- Zero ambiguity
- Explicit upgrades
*/

/*
------------------------------------------------------------
8. WHY KUBERNETES ALMOST NEVER USES v2+
------------------------------------------------------------

Kubernetes strategy:
- Keep APIs backward compatible
- Avoid breaking changes
- Release new features via MINOR versions

That’s why you see:
	k8s.io/apimachinery v0.30.1
	k8s.io/client-go v0.30.1

Instead of:
	v2, v3, v4 chaos

This is deliberate ecosystem stability.
*/

/*
------------------------------------------------------------
9. v0 VERSIONS (PRE-STABLE)
------------------------------------------------------------

v0.x.y means:
- API is NOT stable
- Breaking changes allowed at ANY time

Go treats v0 specially:
- NO /v0 in import path
- Still considered unstable

Kubernetes lived in v0 for a LONG time.
*/

/*
------------------------------------------------------------
10. WHY CNCF PROJECTS PIN EXACT VERSIONS
------------------------------------------------------------

Example:

	require k8s.io/client-go v0.30.1

Why not @latest?
- Determinism
- Reproducibility
- Safe rollouts
- Auditable supply chain

Floating versions are NOT allowed in CNCF.
*/

/*
------------------------------------------------------------
11. TRANSITIVE VERSION SELECTION (MVS)
------------------------------------------------------------

Go uses:
MINIMAL VERSION SELECTION (MVS)

Meaning:
- Chooses the MINIMUM version that satisfies all requirements
- No surprise upgrades
- No version conflicts at runtime

This is why go.mod stays stable.
*/

/*
------------------------------------------------------------
12. COMMON MISTAKES (DO NOT DO THESE)
------------------------------------------------------------

❌ Tag v2 without /v2 module path
❌ Import v2 module without /v2
❌ Assume Go resolves version conflicts magically
❌ Change public APIs in v1 without bumping major
❌ Use floating versions in go.mod
*/

/*
------------------------------------------------------------
13. HOW THIS AFFECTS YOU AS A CONTRIBUTOR
------------------------------------------------------------

When contributing to CNCF projects:
- You almost NEVER introduce v2
- You preserve backward compatibility
- You treat APIs as contracts
- You avoid breaking users

Breaking changes require:
- KEPs / design docs
- Multi-release deprecation
- Extreme care
*/

/*
------------------------------------------------------------
14. THIS FILE'S ROLE
------------------------------------------------------------

This file exists to:
- Be READ
- Explain rules
- Prevent ecosystem-breaking mistakes

Most versioning logic lives in go.mod + tags.
*/

func main() {
	fmt.Println("=== Module Versioning Deep Dive (08.2) ===")
	fmt.Println("Semantic Import Versioning enforced by Go.")
	fmt.Println("Major versions are explicit and intentional.")
}

/*
============================================================
FINAL TAKEAWAYS (CRITICAL)
============================================================

1. Go enforces Semantic Versioning at compile time
2. v2+ MUST appear in import paths
3. Each major version is a separate module
4. This prevents dependency hell
5. CNCF projects avoid breaking changes
6. Versioning is an API promise

If you understand THIS file,
you will not break downstream users.
*/
