# Module 10 — Go Project Structure (CNCF Style)
## 10.6 Logging, Metrics & Tracing — Observability in CNCF Go

This document explains **how observability is designed and implemented** in large
Go systems like Kubernetes, Kubeflow, and other CNCF projects.

If you understand this file, you will:
- understand why logging code is everywhere
- know why `fmt.Println` is forbidden in production code
- understand `klog`, `logr`, metrics, and tracing at a conceptual level
- read Kubernetes code without feeling overwhelmed by observability concerns

---

## 1. What Observability Really Means

Observability answers three fundamental questions:

1. **Logs** — What happened?
2. **Metrics** — How often / how much?
3. **Tracing** — Where did time go?

In distributed systems:
- bugs are not reproducible locally
- failures happen under load
- behavior emerges across services

Observability is not optional. It is **core functionality**.

---

## 2. Logging in CNCF Projects

### Logging is for humans

Logs answer:
- What happened?
- Why did it happen?
- What context was involved?

In CNCF Go projects:
- logs are structured
- logs are contextual
- logs are machine-parseable

---

## 3. Why `fmt.Println` Is Forbidden

Problems with `fmt.Println`:
- no log levels
- no structure
- no context
- no control over output
- impossible to aggregate

Using `fmt.Println` in Kubernetes code will get your PR rejected.

---

## 4. Structured Logging

Structured logging means:
- logs are key/value pairs
- machines can index them
- humans can still read them

Example (conceptual):

"msg": "pod scheduled"
"pod": "nginx-123"
"node": "node-a"

yaml
Copy code

This allows:
- filtering
- aggregation
- correlation

---

## 5. Logging Libraries in CNCF Go

### `klog`
- historical Kubernetes logger
- supports verbosity levels
- widely used in core Kubernetes

### `logr`
- logging interface (abstraction)
- decouples logging API from implementation
- used heavily in controller-runtime

Pattern:
> Code depends on `logr.Logger`, not on a concrete logger.

This matches dependency inversion principles.

---

## 6. Passing Loggers Explicitly

In CNCF projects:
- loggers are passed explicitly
- stored on structs
- enriched with context

Example mental model:

controller := Controller{
log: baseLogger.WithValues("controller", "foo"),
}

yaml
Copy code

Avoid:
- global loggers
- package-level logging state

---

## 7. Metrics — Quantifying Behavior

Metrics answer:
- How many?
- How fast?
- How often?
- How full?

Metrics are numeric and aggregated.

Examples:
- request count
- error rate
- queue depth
- reconcile duration

---

## 8. Prometheus as the CNCF Standard

CNCF projects almost universally use:
- Prometheus-compatible metrics
- pull-based scraping
- numeric time series

Design principles:
- metrics are cheap
- labels are bounded
- cardinality is controlled

High-cardinality metrics can kill systems.

---

## 9. Metrics Placement Rules

Good metrics:
- measure outcomes, not internals
- exist at boundaries
- are stable over time

Bad metrics:
- per-user labels
- per-request IDs
- unbounded dimensions

Metrics are APIs — breaking them hurts operators.

---

## 10. Tracing — Understanding Latency

Tracing answers:
> Where did time go?

In distributed systems:
- a single request spans many services
- latency accumulates across boundaries

Tracing provides:
- end-to-end visibility
- causal relationships
- latency breakdowns

---

## 11. Context Is the Backbone of Tracing

In Go:
- tracing context flows via `context.Context`
- spans are attached to contexts
- cancellation propagates automatically

This is why:
- every Kubernetes function takes `ctx context.Context`

Context is not optional. It is infrastructure.

---

## 12. Why Observability Code Is Everywhere

Newcomers often ask:
> “Why is there so much logging and metrics code?”

Answer:
- systems fail in production
- debugging after the fact requires data
- observability is cheaper than outages

Most observability code exists for **future failures**, not current ones.

---

## 13. Common Beginner Mistakes

- Logging without context
- Logging too much at high verbosity
- Logging too little at error paths
- Adding high-cardinality metric labels
- Ignoring tracing context

These mistakes cause real production incidents.

---

## 14. How Reviewers Evaluate Observability Code

Reviewers look for:
- structured logs
- meaningful log levels
- stable metric names
- bounded labels
- proper context usage

Good observability code builds reviewer confidence.

---

## Final Takeaways (Memorize These)

- Observability is core, not optional
- Logging is structured and contextual
- Metrics are numeric and bounded
- Tracing explains latency
- Context carries cancellation and tracing
- Explicit is better than implicit

If you understand observability, you understand how systems are operated in reality.