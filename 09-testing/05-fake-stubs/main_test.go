package fakesandstubs

import (
	"errors"
	"testing"
)

/*
============================================================
MODULE 09 — TESTING IN GO
STEP 09.5 — FAKES vs STUBS (KUBERNETES STYLE)
============================================================

This file is EXECUTABLE DOCUMENTATION.

Goal:
Understand the DIFFERENCE between:
- Stubs
- Fakes
- (and why Kubernetes prefers fakes)

If you understand THIS file,
client-go fake tests will finally make sense.
*/

/*
------------------------------------------------------------
1. TEST DOUBLES — TERMINOLOGY (IMPORTANT)
------------------------------------------------------------

There are multiple kinds of "mock-like" objects:

- Stub   → returns fixed answers
- Fake   → has working logic (but simplified)
- Mock   → verifies interactions (09.4)
- Spy    → records calls

In Go + CNCF projects:
- Stubs and Fakes are preferred
- Heavy mocking is discouraged
*/

/*
------------------------------------------------------------
2. THE INTERFACE (DEPENDENCY CONTRACT)
------------------------------------------------------------
*/

type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

/*
------------------------------------------------------------
3. STUB IMPLEMENTATION
------------------------------------------------------------

A STUB:
- returns hard-coded data
- has NO internal state
- is used to force code paths
*/

type StubStore struct {
	value string
	err   error
}

func (s *StubStore) Get(key string) (string, error) {
	return s.value, s.err
}

func (s *StubStore) Set(key, value string) error {
	return s.err
}

/*
Use a STUB when:
- You only care about ONE response
- You are testing error handling
- You don't care about realistic behavior
*/

// ==========================================================
// 4. FAKE IMPLEMENTATION
// ==========================================================

/*
A FAKE:
- has real (but simplified) logic
- maintains internal state
- behaves like a real system

Kubernetes uses FAKES heavily.
*/

type FakeStore struct {
	data map[string]string
}

func NewFakeStore() *FakeStore {
	return &FakeStore{
		data: make(map[string]string),
	}
}

func (f *FakeStore) Get(key string) (string, error) {
	val, ok := f.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return val, nil
}

func (f *FakeStore) Set(key, value string) error {
	f.data[key] = value
	return nil
}

/*
------------------------------------------------------------
5. CODE UNDER TEST
------------------------------------------------------------
*/

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) SaveUser(id, name string) error {
	if id == "" {
		return errors.New("empty id")
	}
	return s.store.Set(id, name)
}

func (s *Service) LoadUser(id string) (string, error) {
	return s.store.Get(id)
}

/*
------------------------------------------------------------
6. TESTING WITH A STUB
------------------------------------------------------------

We use a STUB when:
- We want to simulate failure
- We don't care about state
*/

func TestService_WithStub(t *testing.T) {
	stub := &StubStore{
		err: errors.New("database down"),
	}

	service := NewService(stub)

	err := service.SaveUser("1", "alice")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

/*
------------------------------------------------------------
7. TESTING WITH A FAKE
------------------------------------------------------------

We use a FAKE when:
- We want realistic behavior
- We want to test multiple operations
- We want state transitions
*/

func TestService_WithFake(t *testing.T) {
	fake := NewFakeStore()
	service := NewService(fake)

	// Save data
	err := service.SaveUser("1", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Load data
	name, err := service.LoadUser("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if name != "alice" {
		t.Fatalf("LoadUser returned %q; want %q", name, "alice")
	}
}

/*
------------------------------------------------------------
8. WHY KUBERNETES PREFERS FAKES
------------------------------------------------------------

Kubernetes uses FAKES because:
- Controllers have complex flows
- State matters
- Multiple interactions happen
- Tests should resemble real behavior

Examples:
- client-go fake clients
- fake informers
- fake API servers

These are NOT mocks.
They are FAKES.
*/

/*
------------------------------------------------------------
9. STUBS vs FAKES — WHEN TO USE WHAT
------------------------------------------------------------

Use STUB when:
- Testing error paths
- Forcing edge cases
- Keeping tests minimal

Use FAKE when:
- Testing workflows
- Testing state changes
- Testing sequences of calls
*/

/*
------------------------------------------------------------
10. COMMON BEGINNER MISTAKES
------------------------------------------------------------

❌ Using mocks everywhere
❌ Writing fake logic too complex
❌ Testing fake behavior instead of service behavior
❌ Sharing fake state across tests
❌ Using globals instead of constructors
*/

/*
------------------------------------------------------------
11. CNCF TESTING PHILOSOPHY
------------------------------------------------------------

CNCF projects prefer:
- Deterministic tests
- Readable tests
- Realistic behavior
- Minimal magic

Fakes enable this.
*/

/*
------------------------------------------------------------
12. WHAT YOU SHOULD BE COMFORTABLE WITH NOW
------------------------------------------------------------

After this file, you should:
- Know the difference between stubs and fakes
- Choose the right test double
- Understand client-go fake usage
- Read Kubernetes controller tests
- Avoid over-mocking
*/
