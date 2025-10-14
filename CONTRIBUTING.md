# Contributing to Saturn

Welcome to the Saturn Timer Daemon hacknight! This guide will help you get started with contributing to the project and earning points by solving issues.


## 🚀 Getting Started

### Prerequisites

- **Go 1.25+** installed on your system
- **Git** for version control
- A **GitHub account**
- A **code editor** (VS Code, GoLand, etc.)
- Basic understanding of **Go concurrency** and **HTTP APIs**

### First Time Setup

1. **Fork the repository** to your GitHub account

2. **Clone your fork**:
   ```bash
   git clone https://github.com/acmpesuecc/saturn.git
   cd saturn
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/acmpesuecc/saturn.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Verify everything works**:
   ```bash
   go build
   go test ./...
   ```

6. **Run the server**:
   ```bash
   go run . --webhook_url="http://localhost:3000/webhook"
   ```

7. **Test an endpoint** (in a new terminal):
   ```bash
   curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"event_id": "test1", "timeout_seconds": 10, "emit": "Hello!"}' \
     http://localhost:3000/register
   ```

## 💻 Development Setup

### Project Structure

```
saturn/
├── main.go              # Entry point, HTTP server setup
├── timer.go             # Core timer logic and handlers
├── types.go             # Type definitions
├── go.mod               # Go module dependencies
├── README.md            # Project documentation
└── CONTRIBUTING.md      # This file
```

### Understanding the Codebase

**Key Components:**

- **`Timer` struct**: Main timer management structure
- **`TimerMap`**: Thread-safe map storing active timers
- **`TimerMapValue`**: Stores timer handle, duration, and init time
- **Handlers**: `RegisterHandler`, `CancelHandler`, `RemainingHandler`, `ExtendHandler`

**API Endpoints:**

- `POST /register` - Create a new timer
- `POST /cancel` - Cancel an existing timer
- `POST /remaining` - Check remaining time for a timer
- `POST /extend` - Extend a timer's duration
- `GET /test` - Health check endpoint
- `POST /webhook` - Test webhook receiver


## 🔄 Contribution Workflow

### 1. Choose an Issue

Browse [issues](https://github.com/acmpesuecc/saturn/issues) and pick an issue that interests you.

### 2. Create a Branch

```bash
# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create a feature branch
git checkout -b issue-<number>-<short-description>

# Examples:
# git checkout -b issue-1-fix-race-condition
# git checkout -b issue-2-batch-operations
```

### 3. Implement Your Solution

- Read the issue description carefully
- Understand the requirements
- Write clean, idiomatic Go code
- Follow the existing code style
- Add comments for complex logic

### 4. Write Tests

Every contribution **must** include tests:

```go
// Example test structure
func TestYourFeature(t *testing.T) {
    // Setup
    timer := CreateTimer("http://localhost:3000/webhook", "./test-logs/")

    // Test logic
    // ...

    // Assertions
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

### 5. Run Tests Locally

```bash
# Run all tests
go test ./...

# Run with race detection (IMPORTANT for bug fixes!)
go test -race ./...

# Run specific test
go test -run TestYourFeature

# Run with verbose output
go test -v ./...
```

### 6. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "Fix: Resolve race condition in timer cancellation (Issue #1)"

# More commit message examples:
# "Feature: Add batch timer operations endpoint (Issue #2)"
# "Task: Implement timer persistence layer (Issue #3)"
# "Docs: Update API documentation for statistics endpoint"
```

### 7. Push and Create Pull Request

```bash
git push origin issue-<number>-<short-description>
```

Then create a PR on GitHub with:
- **Title**: `[Issue #X] Short description`
- **Description**: What you changed and why
- **Screenshots/Logs**: If applicable

## 📝 Code Standards

### Go Style Guidelines

1. **Follow Go conventions**:
   - Use `gofmt` to format your code
   - Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
   - Use meaningful variable names

2. **Error Handling**:
   ```go
   // Good
   if err != nil {
       t.Logger.Error().Err(err).Msg("Failed to process request")
       w.WriteHeader(http.StatusInternalServerError)
       return
   }

   // Avoid ignoring errors
   ```

3. **Concurrency Safety**:
   ```go
   // Always protect shared state with locks
   t.State.Lock()
   defer t.State.Unlock()

   // Access shared data
   ```

4. **Logging**:
   ```go
   // Use structured logging
   t.Logger.Info().
       Str("event_id", eventID).
       Int("timeout", timeout).
       Msg("Timer registered")
   ```

5. **JSON Tags**:
   ```go
   // Always use json tags for API types
   type YourType struct {
       FieldName string `json:"field_name"`
   }
   ```

### Type Definitions

- Add all new types to `types.go`
- Use clear, descriptive names
- Add comments for exported types
- Group related types together

### API Design

- Follow RESTful conventions
- Use appropriate HTTP status codes
- Return JSON responses with consistent structure
- Include error messages in responses

## 🧪 Testing Guidelines

### Test Requirements

Every PR must include:

1. **Unit Tests** for new functions
2. **Integration Tests** for API endpoints
3. **Race Tests** for concurrent code (run with `-race`)
4. **Edge Case Tests** for error conditions

### Writing Good Tests

```go
func TestRegisterHandler(t *testing.T) {
    // Setup
    timer := CreateTimer("http://localhost:3000/webhook", "./test-logs/")

    // Create test request
    requestBody := RegisterEvent{
        EventID:     "test-event-1",
        TimeoutSecs: 10,
        Emit:        "test message",
    }
    bodyBytes, _ := json.Marshal(requestBody)

    req := httptest.NewRequest("POST", "/register", bytes.NewReader(bodyBytes))
    w := httptest.NewRecorder()

    // Execute
    timer.RegisterHandler(w, req)

    // Assert
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    // Cleanup
    timer.State.Lock()
    delete(timer.State.TimerMap, "test-event-1")
    timer.State.Unlock()
}
```

## ❓ Getting Help

### Resources

- **Go Documentation**: https://go.dev/doc/
- **Go Concurrency Patterns**: https://go.dev/blog/pipelines
- **Effective Go**: https://go.dev/doc/effective_go
- **Saturn README**: [README.md](./README.md)
- **Issue Details**: [HACKNIGHT_ISSUES.md](./HACKNIGHT_ISSUES.md)

### Getting Support

1. **Read the issue description** thoroughly
2. **Check existing code** for similar patterns
3. **Ask questions** in the PR comments
4. **Review Go documentation** for language-specific questions

### Debugging Tips

1. **Use the logger**: Add debug logs to understand flow
   ```go
   t.Logger.Debug().Str("key", value).Msg("Debug message")
   ```

2. **Test in isolation**: Write small test cases for your logic

3. **Use Go's race detector**: Catches concurrency bugs
   ```bash
   go test -race ./...
   ```

4. **Print state**: Debug by printing the timer map contents

5. **Use Postman/curl**: Test API endpoints manually

## 🌟 Good Luck!
