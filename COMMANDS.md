# Go Command Cheat Sheet for Systems Engineers

A reference guide for the Go toolchain, ranging from basic execution to advanced profiling and module management.

## 1. The Essentials (Daily Workflow)
These are the commands you will use 90% of the time.

| Command | Description | Example |
| :--- | :--- | :--- |
| `go run` | Compiles and runs the file immediately without saving an executable. Perfect for quick scripts. | `go run main.go` |
| `go build` | Compiles source code into a binary executable. | `go build -o myapp main.go` |
| `go install` | Compiles and moves the binary to `$GOPATH/bin`. Useful for installing tools. | `go install github.com/user/tool@latest` |
| `go fmt` | **Auto-formats** your code to standard Go style. (K8s CI fails if you don't run this). | `go fmt ./...` |
| `go vet` | Examines code for suspicious constructs (common bugs) that the compiler might miss. | `go vet ./...` |

---

## 2. Module Management (`go.mod`)
Kubernetes relies heavily on modules to manage massive dependency trees.

| Command | Description | Example |
| :--- | :--- | :--- |
| `go mod init` | Initializes a new module in the current directory. Creates `go.mod`. | `go mod init github.com/user/project` |
| `go mod tidy` | Adds missing modules and removes unused ones. **Run this often.** | `go mod tidy` |
| `go get` | Downloads and installs a specific version of a dependency. | `go get k8s.io/client-go@v0.28.0` |
| `go list -m all` | Lists all dependencies currently in use. | `go list -m all` |
| `go mod vendor` | Copies all dependencies into a local `vendor/` folder. (K8s uses this for stability). | `go mod vendor` |

---

## 3. Testing & Verification
Kubernetes has strict testing requirements. You must master these flags.

| Command | Description |
| :--- | :--- |
| `go test ./...` | Runs all tests in the current directory and subdirectories. |
| `go test -v` | Runs tests in **Verbose** mode (shows every test name and output). |
| `go test -cover` | Shows the percentage of code covered by tests. |
| `go test -race` | **CRITICAL:** Detects race conditions in concurrent code. Always run this for K8s controllers. |
| `go test -bench=.` | Runs performance benchmarks. |

---

## 4. Advanced & Tooling (The "Expert" Zone)
Tools for understanding what is happening under the hood (Profiling, Tracing, Assembly).

| Command | Description | Usage Context |
| :--- | :--- | :--- |
| `go tool pprof` | Interactive tool to analyze CPU/Memory profiles. | Debugging high CPU usage in a Pod. |
| `go tool trace` | Visualizes the execution trace of your program. | Debugging Goroutine scheduling latency. |
| `go doc` | Prints documentation for a package or symbol without opening a browser. | `go doc json.Unmarshal` |
| `go env` | Prints Go environment variables. | Checking `GOOS`/`GOARCH` (Linux vs Mac). |
| `go work` | Manages "Workspaces" (working on multiple modules simultaneously). | Developing a K8s fork alongside a custom controller. |

## 5. Cross-Compilation (Build for Linux on Mac/Windows)
Kubernetes runs on Linux. If you are on a Mac/Windows, you **must** cross-compile before pushing a binary to a Docker container.

```bash
# Build a Linux binary from a Mac/Windows machine
GOOS=linux GOARCH=amd64 go build -o my-controller main.go
```

