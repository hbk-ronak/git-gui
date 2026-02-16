# Git GUI - Command Runner

# Default recipe
default: dev

# Run the app in development mode
dev:
    /Users/chai/go/bin/wails dev

# Build the application
build:
    /Users/chai/go/bin/wails build

# Run all tests
test: test-backend test-frontend

# Run Go backend tests
test-backend:
    go test ./backend/... -v

# Run frontend tests
test-frontend:
    cd frontend && npx vitest run

# Lint all code
lint: lint-backend lint-frontend

# Lint Go code
lint-backend:
    go vet ./backend/...

# Lint frontend code
lint-frontend:
    cd frontend && npx svelte-check 2>/dev/null || true

# Clean build artifacts
clean:
    rm -rf build/bin/
    rm -rf frontend/dist/
    rm -rf frontend/node_modules/

# Install all dependencies
install:
    cd frontend && npm install
    go mod tidy

# Check environment
doctor:
    /Users/chai/go/bin/wails doctor
