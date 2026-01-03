package mockingwithinterfaces

import (
	"errors"
	"testing"
)

/*
============================================================
MODULE 09 — TESTING IN GO
STEP 09.4 — MOCKING WITH INTERFACES (COMPLETE DEEP DIVE)
============================================================

This file is EXECUTABLE DOCUMENTATION.

Goal:
Understand how Go replaces traditional mocking frameworks
(Mockito, Jest, etc.) using INTERFACES and DEPENDENCY INJECTION.

If you understand THIS file,
you will understand how Kubernetes tests behavior.
*/

/*
------------------------------------------------------------
1. THE CORE IDEA (VERY IMPORTANT)
------------------------------------------------------------

In Go:
- You DO NOT mock concrete implementations
- You mock BEHAVIOR via interfaces

Rule:
"Depend on interfaces, not implementations"

This is NOT a testing trick.
This is a DESIGN PRINCIPLE.
*/

/*
------------------------------------------------------------
2. THE PROBLEM WE ARE SOLVING
------------------------------------------------------------

We have a Service that:
- depends on an external system (database)
- we do NOT want to hit the real DB in tests

We want:
- fast tests
- deterministic behavior
- full control over failures
*/

// ==========================================================
// 3. DEFINE THE INTERFACE (THE CONTRACT)
// ==========================================================

// Database defines the behavior we depend on.
// This is the KEY abstraction.
type Database interface {
	Save(key, value string) error
}

/*
IMPORTANT:
- Interfaces belong to the CONSUMER
- NOT the implementation

This is a major Go design principle.
*/

// ==========================================================
// 4. CODE UNDER TEST (DEPENDS ON INTERFACE)
// ==========================================================

type Service struct {
	db Database
}

// NewService injects the dependency.
func NewService(db Database) *Service {
	return &Service{db: db}
}

// StoreData is the behavior we want to test.
func (s *Service) StoreData(key, value string) error {
	if key == "" {
		return errors.New("empty key")
	}
	return s.db.Save(key, value)
}

/*
------------------------------------------------------------
5. REAL IMPLEMENTATION (NOT USED IN TEST)
------------------------------------------------------------

In production, this could be:
- SQL DB
- etcd
- Kubernetes API server

We DO NOT use it in tests.
*/

type RealDatabase struct{}

func (r *RealDatabase) Save(key, value string) error {
	// Imagine real I/O here
	return nil
}

// ==========================================================
// 6. MOCK IMPLEMENTATION (FOR TESTING)
// ==========================================================

type MockDatabase struct {
	// Observability
	savedKey   string
	savedValue string

	// Control behavior
	forceError bool
}

func (m *MockDatabase) Save(key, value string) error {
	if m.forceError {
		return errors.New("mock failure")
	}
	m.savedKey = key
	m.savedValue = value
	return nil
}

/*
------------------------------------------------------------
7. TESTING WITH THE MOCK
------------------------------------------------------------

Notice:
- No mocking framework
- No reflection
- No magic
- Just Go code
*/

func TestService_StoreData(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		value         string
		forceDBError  bool
		expectError   bool
		expectedKey   string
		expectedValue string
	}{
		{
			name:          "successful save",
			key:           "user",
			value:         "alice",
			expectedKey:   "user",
			expectedValue: "alice",
			expectError:   false,
		},
		{
			name:         "empty key validation",
			key:          "",
			value:        "data",
			expectError:  true,
		},
		{
			name:         "database failure",
			key:          "user",
			value:        "bob",
			forceDBError: true,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDatabase{
				forceError: tt.forceDBError,
			}

			service := NewService(mockDB)

			err := service.StoreData(tt.key, tt.value)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify behavior (interaction testing)
			if mockDB.savedKey != tt.expectedKey {
				t.Fatalf(
					"savedKey = %q; want %q",
					mockDB.savedKey,
					tt.expectedKey,
				)
			}

			if mockDB.savedValue != tt.expectedValue {
				t.Fatalf(
					"savedValue = %q; want %q",
					mockDB.savedValue,
					tt.expectedValue,
				)
			}
		})
	}
}

/*
------------------------------------------------------------
8. WHY THIS IS BETTER THAN MOCKING FRAMEWORKS
------------------------------------------------------------

Traditional mocks:
- Use reflection
- Hide behavior
- Break silently
- Hard to refactor

Go mocks:
- Are real code
- Are type-safe
- Break at compile time
- Encourage good design
*/

/*
------------------------------------------------------------
9. HOW KUBERNETES USES THIS PATTERN
------------------------------------------------------------

Kubernetes tests:
- mock clients
- fake informers
- stub API servers
- control error paths

Example:
- client-go fake clients
- controller-runtime envtest
- API server fakes

ALL rely on interfaces.
*/

/*
------------------------------------------------------------
10. COMMON BEGINNER MISTAKES
------------------------------------------------------------

❌ Depending on concrete types
❌ Defining interfaces in implementation packages
❌ Over-mocking
❌ Mocking things you don't own
❌ Using global variables instead of injection
*/

/*
------------------------------------------------------------
11. DESIGN INSIGHT (IMPORTANT)
------------------------------------------------------------

Good Go code:
- Is easy to test
- Because it was designed well

If code is hard to test:
- The design is usually wrong

Testing reveals design flaws.
*/

/*
------------------------------------------------------------
12. WHAT YOU SHOULD BE COMFORTABLE WITH NOW
------------------------------------------------------------

After this file, you should:
- Design code around interfaces
- Inject dependencies
- Write your own mocks
- Read Kubernetes mock-based tests
- Avoid mocking frameworks confidently
*/
