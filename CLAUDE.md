# Development Standards for Claude Code

This document outlines the coding standards, principles, and workflows for this project. All code should follow these guidelines.

## Core Development Principles

### DRY (Don't Repeat Yourself)

**Definition:** Every piece of knowledge should have a single, unambiguous representation in the system.

**Application:**
- Extract common logic into reusable functions
- Create utility/helper modules for repeated operations
- Use shared constants instead of magic numbers/strings
- Centralize configuration and error messages

**Examples:**

**Good:**
```go
// Centralized git execution
func executeGitCommand(args ...string) (string, error) {
    cmd := exec.Command("git", args...)
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("git command failed: %w", err)
    }
    return string(output), nil
}

// All operations use the same helper
func getStatus() (string, error) {
    return executeGitCommand("status", "--porcelain")
}

func getCurrentBranch() (string, error) {
    return executeGitCommand("branch", "--show-current")
}
```

**Bad:**
```go
// Duplicated git execution logic
func getStatus() (string, error) {
    cmd := exec.Command("git", "status", "--porcelain")
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("git command failed: %w", err)
    }
    return string(output), nil
}

func getCurrentBranch() (string, error) {
    cmd := exec.Command("git", "branch", "--show-current")
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("git command failed: %w", err)
    }
    return string(output), nil
}
```

### YAGNI (You Aren't Gonna Need It)

**Definition:** Only implement features when you actually need them, not when you anticipate needing them.

**Application:**
- Don't add abstraction layers "just in case"
- No premature optimization
- Don't implement features not in the spec
- Keep dependencies minimal
- Use standard library when possible

**What NOT to build:**
- Configuration systems (use sensible defaults)
- Plugin architectures
- Caching layers (until proven necessary)
- Complex state machines
- Custom implementations of standard tools

**Examples:**

**Good:**
```go
// Simple, direct implementation
func CommitFiles(files []string, message string) error {
    if len(files) == 0 {
        return errors.New("no files to commit")
    }
    
    for _, file := range files {
        if err := executeGitCommand("add", file); err != nil {
            return err
        }
    }
    
    return executeGitCommand("commit", "-m", message)
}
```

**Bad:**
```go
// Over-engineered with unnecessary abstractions
type CommitStrategy interface {
    Execute(files []string, message string) error
}

type StandardCommitStrategy struct {
    executor GitExecutor
    validator MessageValidator
    logger Logger
}

func (s *StandardCommitStrategy) Execute(files []string, message string) error {
    // Complex abstraction for simple git commit
}
```

### KISS (Keep It Simple, Stupid)

**Definition:** Favor simple solutions over complex ones.

**Application:**
- Write straightforward, linear code
- Avoid deep nesting and complex control flow
- Prefer composition over inheritance
- Use clear, obvious implementations
- Don't be clever - be clear

**Examples:**

**Good:**
```typescript
async function commitFiles() {
    try {
        const result = await App.CommitFiles(selectedFiles, message)
        showSuccess("Committed successfully")
        await refreshFiles()
    } catch (error) {
        showError(error.message)
    }
}
```

**Bad:**
```typescript
async function commitFiles() {
    return commitQueue
        .enqueue(() => retryWithBackoff(() => App.CommitFiles(selectedFiles, message)))
        .then(handleCommitSuccess)
        .catch(handleCommitError)
        .finally(updateMetrics)
}
```

### Clean Code Principles

#### Naming Conventions

**Be Descriptive:**
- `getGitStatus()` not `getStatus()`
- `commitMessage` not `msg`
- `filepath` not `fp`

**Follow Language Conventions:**
- **Go:** PascalCase for exports, camelCase for private
- **TypeScript:** camelCase for variables/functions, PascalCase for components/types
- **Constants:** UPPER_SNAKE_CASE

**Examples:**
```go
// Good
func ParseGitStatus(output string) ([]FileStatus, error)
const MaxRetries = 3

// Bad
func parse_git_status(output string) ([]FileStatus, error)
const max_retries = 3
```

#### Function Design

**Single Responsibility:**
- Each function does one thing well
- If you can't describe it in one sentence, it's too complex

**Small Functions:**
- Aim for <50 lines per function
- Extract complex logic into helper functions

**Few Parameters:**
- Maximum 3-4 parameters
- Use structs/objects for more

**Early Returns:**
- Use guard clauses to reduce nesting
- Fail fast

**Examples:**

**Good:**
```go
func CommitFiles(files []string, message string) error {
    if len(files) == 0 {
        return errors.New("no files to commit")
    }
    
    if message == "" {
        return errors.New("commit message required")
    }
    
    if err := stageFiles(files); err != nil {
        return err
    }
    
    return createCommit(message)
}
```

**Bad:**
```go
func CommitFiles(files []string, message string) error {
    if len(files) > 0 {
        if message != "" {
            err := stageFiles(files)
            if err == nil {
                return createCommit(message)
            } else {
                return err
            }
        } else {
            return errors.New("commit message required")
        }
    } else {
        return errors.New("no files to commit")
    }
}
```

#### Code Organization

**Logical Grouping:**
- Group related functions together
- Separate concerns (UI logic vs business logic)
- Keep files under 300 lines

**No Magic Numbers:**
- Use named constants
- Make intent clear

**Examples:**
```go
// Good
const (
    MaxFilePathLength = 200
    DefaultWindowWidth = 1200
    DefaultWindowHeight = 800
)

// Bad
if len(path) > 200 { // What is 200?
    // ...
}
```

### Self-Documenting Code

**Definition:** Code should be clear enough to understand without extensive comments.

**Application:**
- Use descriptive names for variables, functions, types
- Write clear, linear logic
- Break complex operations into well-named functions
- Comments explain *why*, not *what*

**When to Comment:**
- Public API documentation
- Complex algorithms (explain the approach)
- Non-obvious workarounds or edge cases
- Important decisions and trade-offs

**When NOT to Comment:**
- Obvious operations
- Repeating what the code says
- Commented-out code (delete it)

**Examples:**

**Good:**
```go
// Good - self-explanatory
func isValidBranchName(name string) bool {
    if name == "" {
        return false
    }
    
    // Git doesn't allow spaces or special characters in branch names
    invalidChars := []string{" ", "~", "^", ":", "?", "*", "["}
    for _, char := range invalidChars {
        if strings.Contains(name, char) {
            return false
        }
    }
    
    return true
}
```

**Bad:**
```go
// Bad - unnecessary comments
func isValidBranchName(name string) bool {
    // Check if name is empty
    if name == "" {
        return false  // Return false if empty
    }
    
    // Create array of invalid characters
    invalidChars := []string{" ", "~", "^", ":", "?", "*", "["}
    // Loop through each invalid character
    for _, char := range invalidChars {
        // Check if name contains this character
        if strings.Contains(name, char) {
            return false  // Return false if found
        }
    }
    
    return true  // Return true if valid
}
```

## Error Handling

**Always Handle Errors:**
- Never silently ignore errors
- Provide context when wrapping errors
- Return errors, don't panic (except for truly unrecoverable situations)

**Examples:**

**Go:**
```go
// Good
func SwitchBranch(name string) error {
    if err := validateBranchExists(name); err != nil {
        return fmt.Errorf("cannot switch to branch %s: %w", name, err)
    }
    
    if err := executeGitCommand("checkout", name); err != nil {
        return fmt.Errorf("failed to checkout branch %s: %w", name, err)
    }
    
    return nil
}
```

**TypeScript:**
```typescript
// Good
try {
    const result = await App.CommitFiles(files, message)
    return result
} catch (error) {
    console.error("Commit failed:", error)
    throw new Error(`Failed to commit: ${error.message}`)
}
```

## Standard Development Tools

### Just Commands

This project uses `just` as a command runner. All projects should implement these standard commands:

#### Required Commands

**Testing:**
```just
# Run all tests
test:
    go test ./... -v
    npm test
```

**Linting:**
```just
# Lint all code
lint:
    go vet ./...
    cd frontend && npm run lint
```

**Building:**
```just
# Build the application
build:
    wails build
```

**Cleaning:**
```just
# Clean workspace
clean:
    rm -rf build/
    rm -rf frontend/dist/
```

## Tool Usage Rules

- ALWAYS use `just` commands for testing, linting, and building. NEVER run `go test`, `go vet`, `npm test`, or similar commands directly.
  - `just test` — run all tests
  - `just lint` — run all linters  
  - `just build` — build the application


## Progressive Development Approach

**CRITICAL: Do NOT attempt to build the entire project in one shot.**

### Development Workflow

1. **Start Small:** Implement one feature/component at a time
2. **Test Early:** Write tests as you go, not at the end
3. **Iterate:** Get one thing working before moving to the next
4. **Integrate:** Ensure each piece works with existing code before continuing

### Feature Implementation Order

**Phase 1 - Foundation:**
1. Set up project structure
2. Implement basic Git executor (mocked)
3. Write tests for Git executor
4. Implement one parser (e.g., status parser)
5. Write tests for parser

**Phase 2 - Core Backend:**
6. Implement remaining parsers
7. Implement App struct methods one by one
8. Test each method as implemented

**Phase 3 - Frontend Foundation:**
9. Set up Svelte project
10. Create basic layout structure
11. Implement stores
12. Test stores

**Phase 4 - Components:**
13. Build FileList component
14. Test FileList
15. Build DiffViewer component
16. Test DiffViewer
17. Build remaining components one at a time

**Phase 5 - Integration:**
18. Wire components together
19. Connect frontend to backend
20. Integration testing
21. Polish and bug fixes

### Why Progressive Development?

- **Catch issues early** - Problems are easier to debug in small pieces
- **Maintain quality** - Each piece is tested before moving on
- **Reduce rework** - Don't build on broken foundations
- **Stay focused** - One concern at a time
- **Easier debugging** - Smaller surface area to investigate

### Anti-Pattern: One-Shot Development

**DON'T DO THIS:**
```
❌ Write all backend code
❌ Write all frontend code
❌ Write all tests at the end
❌ Try to integrate everything at once
❌ Debug a giant ball of untested code
```

**DO THIS:**
```
✅ Implement one function
✅ Write tests for that function
✅ Verify it works
✅ Move to next function
✅ Keep everything working as you go
```

## Code Quality Checks

### Before Committing

Run these checks before every commit:

```bash
just test    # All tests pass
just lint    # No linting errors
just build   # Code compiles
```

### Code Review Checklist

- [ ] Follows DRY principle
- [ ] No premature optimization (YAGNI)
- [ ] Simple, clear code (KISS)
- [ ] Self-documenting with good names
- [ ] Proper error handling
- [ ] Tests written and passing
- [ ] No commented-out code
- [ ] No magic numbers
- [ ] Functions are small and focused

## Project Structure

```
.
├── CLAUDE.md           # This file - development standards
├── specs/              # Product specifications
│   ├── 1-product.md    # MVP and product requirements
│   ├── 2-ui.md         # UI specifications
│   ├── 3-backend.md    # Backend specifications
│   └── 4-tests.md      # Testing strategy
├── src/
│   ├── backend/        # Go backend code
│   └── frontend/       # Svelte frontend code
├── tests/              # Test files
└── justfile            # Command runner
```

## Summary

**Remember:**
1. **DRY** - Don't repeat yourself
2. **YAGNI** - You aren't gonna need it
3. **KISS** - Keep it simple
4. **Self-documenting** - Clear names, obvious code
5. **Use Just** - Standard commands for all tasks
6. **Progressive** - Build incrementally, test constantly
7. **No one-shotting** - Small steps, verify as you go

Following these principles will result in clean, maintainable, well-tested code.


You are PROHIBITED from creating any random markdown artifiacts without explicit instruction from the users.

Go version = 1.23.12
npm version = 11.9.0
node version = 25.6.1