# Git GUI - Testing Strategy & Setup

## Project Setup & Testing Strategy

### Initial Project Setup

**Create Wails Project:**
```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Create new project with Svelte template
wails init -n git-gui -t svelte

# Navigate to project
cd git-gui
```

**Install Additional Dependencies:**

**Backend (Go):**
```bash
# Testing libraries
go get -u github.com/stretchr/testify/assert
go get -u github.com/stretchr/testify/mock
```

**Frontend (TypeScript/Svelte):**
```bash
# Testing libraries
npm install --save-dev @testing-library/svelte
npm install --save-dev @testing-library/jest-dom
npm install --save-dev vitest
npm install --save-dev @vitest/ui
npm install --save-dev jsdom
```

**Project Configuration:**

**Vite Config (vite.config.ts):**
```typescript
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/setupTests.ts']
  }
})
```

**TypeScript Config (tsconfig.json):**
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "lib": ["ES2020", "DOM"],
    "moduleResolution": "bundler",
    "strict": true,
    "skipLibCheck": true,
    "types": ["vitest/globals", "@testing-library/jest-dom"]
  }
}
```

### Testing Strategy

#### Core Principles

1. **One test, one concern** - Each test function should verify exactly one behavior or outcome
2. **Mock external dependencies** - Don't rely on actual git commands, file systems, or network in tests
3. **Precise assertions** - Test what matters, not implementation details
4. **Avoid test proliferation** - Don't write 10 tests when 1 well-designed test covers the cases

#### Backend Testing (Go)

**What to Test:**
- Git command parsing logic
- Data transformations (git output → structs)
- Error handling paths
- Business logic in App methods

**What NOT to Test:**
- Actual git command execution (mock it)
- External file system operations
- Network calls

**Test Organization:**
```
app_test.go           - Tests for App struct methods
git_parser_test.go    - Tests for parsing git output
git_executor_test.go  - Tests for git command execution (mocked)
```

**Example Test Patterns:**

**Good - Single concern, comprehensive:**
```go
// Tests the complete parsing logic with realistic input
func TestParseGitStatus_AllStatusTypes(t *testing.T) {
    input := ` M modified.txt
M  staged.txt
A  added.txt
D  deleted.txt
?? untracked.txt`

    files, err := parseGitStatus(input)
    
    assert.NoError(t, err)
    assert.Len(t, files, 5)
    
    // Verify each type is parsed correctly
    assert.Equal(t, "modified.txt", files[0].Path)
    assert.Equal(t, StatusModified, files[0].Status)
    assert.False(t, files[0].Staged)
    
    assert.Equal(t, "staged.txt", files[1].Path)
    assert.True(t, files[1].Staged)
    
    assert.Equal(t, StatusAdded, files[2].Status)
    assert.Equal(t, StatusDeleted, files[3].Status)
    assert.Equal(t, StatusUntracked, files[4].Status)
}
```

**Bad - Multiple concerns, excessive tests:**
```go
// DON'T DO THIS - 5 tests when 1 would suffice
func TestParseGitStatus_Modified(t *testing.T) { /* ... */ }
func TestParseGitStatus_Staged(t *testing.T) { /* ... */ }
func TestParseGitStatus_Added(t *testing.T) { /* ... */ }
func TestParseGitStatus_Deleted(t *testing.T) { /* ... */ }
func TestParseGitStatus_Untracked(t *testing.T) { /* ... */ }
```

**Mocking Git Commands:**
```go
// Create a mockable interface for git execution
type GitExecutor interface {
    Execute(args ...string) (string, error)
}

// Real implementation
type RealGitExecutor struct{}

func (e *RealGitExecutor) Execute(args ...string) (string, error) {
    cmd := exec.Command("git", args...)
    output, err := cmd.Output()
    return string(output), err
}

// Mock for testing
type MockGitExecutor struct {
    mock.Mock
}

func (m *MockGitExecutor) Execute(args ...string) (string, error) {
    ret := m.Called(args)
    return ret.String(0), ret.Error(1)
}

// Test using mock
func TestGetGitStatus_Success(t *testing.T) {
    mockExec := new(MockGitExecutor)
    mockExec.On("Execute", []string{"status", "--porcelain"}).
        Return(" M file.txt\n", nil)
    
    app := &App{executor: mockExec}
    files, err := app.GetGitStatus()
    
    assert.NoError(t, err)
    assert.Len(t, files, 1)
    assert.Equal(t, "file.txt", files[0].Path)
    mockExec.AssertExpectations(t)
}
```

**Table-Driven Tests (When Appropriate):**
Use table-driven tests only when testing the same logic with different inputs, not for different behaviors:

```go
// Good use of table-driven test - same function, different inputs
func TestParseBranchName(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"* main\n", "main"},
        {"* feature/test\n", "feature/test"},
        {"  develop\n", "develop"},
    }
    
    for _, tt := range tests {
        result := parseBranchName(tt.input)
        assert.Equal(t, tt.expected, result)
    }
}
```

**Error Testing:**
```go
// Test both success and error paths in separate focused tests
func TestCommitFiles_EmptyMessage(t *testing.T) {
    app := &App{}
    _, err := app.CommitFiles([]string{"file.txt"}, "")
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "commit message required")
}

func TestCommitFiles_NoFiles(t *testing.T) {
    app := &App{}
    _, err := app.CommitFiles([]string{}, "message")
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "no files to commit")
}
```

#### Frontend Testing (Svelte + Vitest)

**What to Test:**
- Component rendering with different props/state
- User interactions (clicks, typing, selections)
- State management (store updates)
- API call integration (mocked)

**What NOT to Test:**
- CSS styling specifics
- Third-party library internals
- Svelte framework behavior

**Test Organization:**
```
FileList.test.ts       - FileList component tests
DiffViewer.test.ts     - DiffViewer component tests
BranchSelector.test.ts - BranchSelector component tests
CommitPanel.test.ts    - CommitPanel component tests
api.test.ts           - API wrapper tests (mocked backend calls)
stores.test.ts        - Store behavior tests
```

**Example Test Patterns:**

**Good - Component renders and handles interaction:**
```typescript
import { render, fireEvent } from '@testing-library/svelte'
import FileList from './FileList.svelte'
import { files } from './stores'

// Mock the backend API
vi.mock('./wailsjs/go/main/App', () => ({
  GetGitStatus: vi.fn(() => Promise.resolve([
    { Path: 'file1.txt', Status: 'modified', Staged: false },
    { Path: 'file2.txt', Status: 'added', Staged: true }
  ]))
}))

test('FileList displays files and handles checkbox toggle', async () => {
  const { getByText, getByRole } = render(FileList)
  
  // Files should be displayed
  expect(getByText('file1.txt')).toBeInTheDocument()
  expect(getByText('file2.txt')).toBeInTheDocument()
  
  // Find checkbox for file1
  const checkbox = getByRole('checkbox', { name: /file1.txt/ })
  expect(checkbox).not.toBeChecked()
  
  // Click checkbox
  await fireEvent.click(checkbox)
  
  // Checkbox should now be checked
  expect(checkbox).toBeChecked()
})
```

**Bad - Testing too many things at once:**
```typescript
// DON'T DO THIS - too many concerns in one test
test('FileList everything', async () => {
  // Tests rendering, clicking, API calls, state updates, 
  // CSS classes, event handlers, etc. all in one test
})
```

**Store Testing:**
```typescript
import { get } from 'svelte/store'
import { files, stagedFiles } from './stores'

test('stagedFiles derived store filters correctly', () => {
  // Set test data
  files.set([
    { Path: 'file1.txt', Status: 'modified', Staged: true },
    { Path: 'file2.txt', Status: 'modified', Staged: false },
    { Path: 'file3.txt', Status: 'added', Staged: true }
  ])
  
  // Verify derived store
  const staged = get(stagedFiles)
  expect(staged).toHaveLength(2)
  expect(staged[0].Path).toBe('file1.txt')
  expect(staged[1].Path).toBe('file3.txt')
})
```

**API Wrapper Testing (Mocked Backend):**
```typescript
import { vi } from 'vitest'
import { commitFiles } from './api'
import * as App from './wailsjs/go/main/App'

// Mock the Wails backend
vi.mock('./wailsjs/go/main/App')

test('commitFiles calls backend and updates state on success', async () => {
  const mockCommitResult = {
    Success: true,
    CommitSHA: 'abc123',
    Message: 'Committed'
  }
  
  App.CommitFiles.mockResolvedValue(mockCommitResult)
  
  const result = await commitFiles(['file1.txt'], 'test commit')
  
  expect(App.CommitFiles).toHaveBeenCalledWith(['file1.txt'], 'test commit')
  expect(result.Success).toBe(true)
})

test('commitFiles handles errors gracefully', async () => {
  App.CommitFiles.mockRejectedValue(new Error('Git error'))
  
  await expect(commitFiles(['file1.txt'], 'test')).rejects.toThrow('Git error')
})
```

**Component Interaction Testing:**
```typescript
test('CommitPanel disables buttons when conditions not met', async () => {
  const { getByText } = render(CommitPanel)
  
  const commitBtn = getByText('Commit')
  const commitPushBtn = getByText('Commit & Push')
  
  // No files staged, no message - buttons disabled
  expect(commitBtn).toBeDisabled()
  expect(commitPushBtn).toBeDisabled()
  
  // Add message but still no staged files - still disabled
  const textarea = getByRole('textbox')
  await fireEvent.input(textarea, { target: { value: 'test message' } })
  
  expect(commitBtn).toBeDisabled()
  expect(commitPushBtn).toBeDisabled()
})
```

### Test Coverage Guidelines

**Aim for meaningful coverage, not 100%:**
- Backend: 70-80% line coverage on business logic
- Frontend: 60-70% coverage on components
- Focus on critical paths, not trivial getters/setters

**Don't test:**
- Auto-generated Wails bindings
- Third-party library code
- Simple pass-through functions
- CSS/styling

**Do test:**
- Parsing logic
- State transformations
- User interactions
- Error handling
- Edge cases (empty lists, null values, etc.)

### Running Tests

**Backend Tests:**
```bash
# Run all Go tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific test file
go test ./app_test.go -v

# Run specific test function
go test -run TestGetGitStatus -v
```

**Frontend Tests:**
```bash
# Run all frontend tests
npm test

# Run with UI
npm run test:ui

# Run with coverage
npm run test:coverage

# Run in watch mode
npm test -- --watch

# Run specific test file
npm test FileList.test.ts
```

**Pre-commit Testing:**
```bash
# Run all tests before committing
go test ./... && npm test
```

### CI/CD Considerations (Optional)

If setting up continuous integration:
```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.21'
      - uses: actions/setup-node@v2
        with:
          node-version: '18'
      - run: go test ./... -cover
      - run: npm ci && npm test
```

### Test File Naming Convention

**Backend:**
- `filename.go` → `filename_test.go`
- Place test file in same directory as source file

**Frontend:**
- `Component.svelte` → `Component.test.ts`
- `api.ts` → `api.test.ts`
- Place test file adjacent to source file

### Mock Data Fixtures
