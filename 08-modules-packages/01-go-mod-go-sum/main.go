package main

import (
	"fmt"

	// External dependency to prove module resolution works
	"k8s.io/apimachinery/pkg/util/runtime"
)

/*
============================================================
MODULE 08 — MODULES & PACKAGES
STEP 08.1 — go.mod & go.sum (COMPLETE DEEP DIVE)
============================================================

This file is meant to be READ, not just run.

Goal:
Understand how Go modules enable large-scale projects
like Kubernetes, Kubeflow, and CNCF ecosystems.

If you understand THIS file + its comments,
you will not be afraid of go.mod / go.sum ever again.
*/

/*
------------------------------------------------------------
1. WHAT IS A GO MODULE (IMPORTANT)
------------------------------------------------------------

A Go module is:
- A versioned unit of source code
- With a single dependency graph
- Defined by a go.mod file

It is NOT:
- Just a folder
- Just a Git repo
- Just a package

Think:
	module = repo + version + dependency contract
*/

/*
------------------------------------------------------------
2. THE module DIRECTIVE (PUBLIC API ROOT)
------------------------------------------------------------

Example go.mod line:

	module k8s.io/client-go

Meaning:
- This string is the ROOT import path
- All packages live under this namespace
- It is part of the PUBLIC API

Changing it is a BREAKING CHANGE.

This is why CNCF projects use stable vanity domains:
	k8s.io
	knative.dev
	sigs.k8s.io
*/

/*
------------------------------------------------------------
3. THE go DIRECTIVE (MISUNDERSTOOD BUT CRITICAL)
------------------------------------------------------------

Example:

	go 1.22

This does NOT mean:
- "You must install Go 1.22"

It DOES mean:
- Language semantics baseline
- Standard library behavior
- Module resolution rules

Changing this line can subtly change program behavior.

Kubernetes treats this line VERY seriously.
*/

/*
------------------------------------------------------------
4. require DIRECTIVE (DEPENDENCY GRAPH)
------------------------------------------------------------

Example:

	require k8s.io/apimachinery v0.30.1

This is NOT just about what THIS file imports.

It defines:
- The full dependency DAG
- What versions are acceptable
- What gets downloaded and verified

This is why indirect dependencies exist.
*/

/*
------------------------------------------------------------
5. go.sum (SUPPLY CHAIN SECURITY)
------------------------------------------------------------

go.sum stores:
- Cryptographic hashes of modules
- Hash of module source
- Hash of module go.mod

This guarantees:
- Reproducible builds
- Tamper detection
- Safe proxy usage

You must NEVER edit go.sum manually.

CNCF CI will reject PRs if go.sum is incorrect.
*/

/*
------------------------------------------------------------
6. GLOBAL MODULE CACHE (WHY GO SCALES)
------------------------------------------------------------

Go stores dependencies in a global cache:

	$GOMODCACHE

Benefits:
- One copy per version
- Fast builds
- No node_modules explosion
- Deterministic resolution

This is why Kubernetes builds stay sane.
*/

/*
------------------------------------------------------------
7. go mod tidy (NON-OPTIONAL)
------------------------------------------------------------

Command:

	go mod tidy

What it guarantees:
- No unused dependencies
- No missing dependencies
- Clean go.mod
- Correct go.sum

In CNCF projects:
- Forgetting tidy = guaranteed review comment
*/

/*
------------------------------------------------------------
8. replace DIRECTIVE (DANGEROUS IF MISUSED)
------------------------------------------------------------

Example:

	replace k8s.io/client-go => ../client-go

Used ONLY for:
- Local development
- Testing multiple modules together

RULE:
- NEVER commit replace unless maintainers ask
- Very common beginner PR mistake
*/

/*
------------------------------------------------------------
9. WHY KUBERNETES USES MULTIPLE MODULES
------------------------------------------------------------

Kubernetes is NOT one giant module.

Instead:
- k8s.io/api
- k8s.io/apimachinery
- k8s.io/client-go

Why?
- Independent versioning
- Smaller dependency graphs
- Faster CI
- Clear API boundaries

This is deliberate architecture.
*/

/*
------------------------------------------------------------
10. THIS FILE'S ROLE
------------------------------------------------------------

This file:
- Proves module resolution works
- Imports a real CNCF dependency
- Exists mainly for documentation

Most module logic lives in go.mod, not here.
*/

func main() {
	fmt.Println("=== Go Modules Deep Dive (08.1) ===")

	// This function comes from k8s.io/apimachinery
	// If this builds and runs:
	// - go.mod is correct
	// - go.sum is verified
	// - module resolution works
	defer runtime.HandleCrash()

	fmt.Println("Go module system is active and verified.")
}

/*
============================================================
FINAL TAKEAWAYS (MEMORIZE THESE)
============================================================

1. go.mod defines your module's contract
2. go.sum guarantees integrity and reproducibility
3. module paths are public APIs
4. versions are promises
5. go mod tidy is mandatory
6. replace is dangerous if committed
7. CNCF projects rely heavily on modules

If you understand THIS file,
you can confidently navigate any CNCF Go repo.
*/
