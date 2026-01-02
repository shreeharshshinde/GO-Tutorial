package main

import (
	"fmt"
	"time"
)

/*
This file explains REAL-WORLD INTERFACE PATTERNS in Go.

By now you already know:
- What interfaces are
- How duck typing works
- any / type assertions
- Embedding and nil interface pitfalls

This file answers the final questions:
- WHEN should interfaces be used?
- WHERE should interfaces live?
- HOW are interfaces used in large systems (like Kubernetes)?
- WHEN should you NOT use interfaces?
*/

//
// 1. ACCEPT INTERFACES, RETURN CONCRETE TYPES
//

/*
Rule:
- Functions should ACCEPT interfaces
- Functions should RETURN concrete types

Why:
- Callers depend on behavior
- Creators should decide implementation details
*/

// Logger defines behavior.
type Logger interface {
	Log(msg string)
}

// Concrete implementation.
type ConsoleLogger struct{}

func (c *ConsoleLogger) Log(msg string) {
	fmt.Println("[LOG]", msg)
}

// BAD: Returning an interface hides the concrete type.
/*
func NewLoggerBad() Logger {
	return &ConsoleLogger{}
}
*/

// GOOD: Return the concrete type.
func NewLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

//
// 2. INTERFACES ARE ABOUT CALLERS, NOT IMPLEMENTERS
//

/*
Key idea:
- The CONSUMER defines the interface
- Not the PROVIDER

This avoids forcing unnecessary methods.
*/

type Service struct {
	log Logger
}

func NewService(l Logger) *Service {
	return &Service{log: l}
}

func (s *Service) Run() {
	s.log.Log("service running")
}

//
// 3. SMALL INTERFACES (INTERFACE SEGREGATION)
//

/*
Go favors SMALL interfaces.

From standard library:
- io.Reader
- io.Writer
- io.Closer

Large interfaces are hard to implement and test.
*/

type Reader interface {
	Read() string
}

type Writer interface {
	Write(data string)
}

// Composed interface.
type ReadWriter interface {
	Reader
	Writer
}

//
// 4. INTERFACE COMPOSITION IN PRACTICE
//

type MemoryBuffer struct {
	data string
}

func (m *MemoryBuffer) Read() string {
	return m.data
}

func (m *MemoryBuffer) Write(d string) {
	m.data = d
}

//
// 5. INTERFACES ENABLE TESTING (PRIMARY USE CASE)
//

/*
Interfaces are NOT for abstraction first.
They are for:
- Swapping implementations
- Testing
*/

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (r *RealClock) Now() time.Time {
	return time.Now()
}

// Fake implementation for tests.
type FakeClock struct {
	t time.Time
}

func (f *FakeClock) Now() time.Time {
	return f.t
}

type Scheduler struct {
	clock Clock
}

func NewScheduler(c Clock) *Scheduler {
	return &Scheduler{clock: c}
}

func (s *Scheduler) PrintTime() {
	fmt.Println("Time:", s.clock.Now())
}

//
// 6. INTERFACE POLLUTION (ANTI-PATTERN)
//

/*
BAD PRACTICE:
Creating interfaces too early or too large.

Do NOT do this:
*/

type UserService interface {
	CreateUser()
	DeleteUser()
	UpdateUser()
	FindUser()
}

/*
Why this is bad:
- Hard to mock
- Hard to evolve
- Forces unnecessary methods

Correct approach:
- Start with concrete types
- Introduce interfaces ONLY when needed
*/

//
// 7. WHEN NOT TO USE INTERFACES
//

/*
Do NOT use interfaces when:
- There is only ONE implementation
- You control both caller and callee
- There is no need for mocking
- You are still exploring the design
*/

type SimpleCounter struct {
	count int
}

func (s *SimpleCounter) Increment() {
	s.count++
}

//
// 8. INTERFACES AT PACKAGE BOUNDARIES
//

/*
Interfaces belong at PACKAGE BOUNDARIES.

Good:
- database package defines DB interface
- consumer depends on DB interface

Bad:
- structs embedding interfaces internally without need
*/

//
// 9. KUBERNETES MENTAL MODEL
//

/*
Kubernetes uses interfaces for:
- Clients
- Informers
- Listers
- Controllers

Concrete implementations are hidden.
Behavior is exposed.
*/

//
// 10. FINAL DESIGN RULES (READ CAREFULLY)
//

/*
1. Prefer concrete types first
2. Introduce interfaces when behavior varies
3. Keep interfaces small
4. Accept interfaces, return concrete types
5. Interfaces belong to consumers
6. Interfaces enable testing, not inheritance
7. Avoid interface pollution
*/

func main() {
	fmt.Println("--- Accept Interfaces, Return Concrete ---")
	logger := NewLogger()
	logger.Log("hello")

	fmt.Println("\n--- Service Using Interface ---")
	service := NewService(logger)
	service.Run()

	fmt.Println("\n--- Interface Composition ---")
	buf := &MemoryBuffer{}
	buf.Write("hello world")
	fmt.Println("Read:", buf.Read())

	fmt.Println("\n--- Interfaces for Testing ---")
	fakeClock := &FakeClock{t: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	scheduler := NewScheduler(fakeClock)
	scheduler.PrintTime()

	fmt.Println("\n--- Concrete Type Without Interface ---")
	counter := &SimpleCounter{}
	counter.Increment()
	counter.Increment()
	fmt.Println("Counter:", counter.count)
}
