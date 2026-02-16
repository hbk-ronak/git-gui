# Git GUI - Product & MVP Specification

## Overview

A lightweight desktop application for managing git workflows with a clean, minimal interface. This tool provides a thin wrapper around core git commands with a focus on the most common developer workflows: staging files, committing, and pushing changes.

## Tech Stack

### Framework
- **Wails v2** - Desktop application framework
- **Go** - Backend (git command wrapper)
- **Svelte** - Frontend UI framework

### Why Wails?
- Lightweight bundle size (~10MB vs Electron's ~120MB)
- Uses native OS webview (no Chromium bundle)
- Single binary distribution
- Clean separation between backend (Go) and frontend (web tech)
- Type-safe bindings auto-generated between Go and TypeScript

## User Stories

### Story 1: View changed files and their diffs
**As a developer**, I want to see all my changed files and view their diffs, so that I can review what I've modified before committing.

**Acceptance Criteria:**
- I see a list of all modified/untracked files from `git status --porcelain`
- When I click on a file, the right pane shows `git diff <file>`
- The diff shows added lines (green/+) and removed lines (red/-)
- I can click between files to view different diffs

### Story 2: Stage and commit files
**As a developer**, I want to select which files to commit and write a commit message, so that I can create focused commits.

**Acceptance Criteria:**
- I can check/uncheck files to stage them for commit
- I can enter a commit message in the text area
- "Commit" button runs `git add <checked-files> && git commit -m "message"`
- "Commit & Push" button runs the above plus `git push`
- I see success/error feedback after git operations
- Checked files and commit message clear after successful commit

### Story 3: Switch and create branches
**As a developer**, I want to switch between branches and create new ones, so that I can organize my work.

**Acceptance Criteria:**
- I see the current branch name in a dropdown
- Clicking the dropdown shows all local branches from `git branch`
- Selecting a branch runs `git checkout <branch>`
- Clicking "New Branch" opens a dialog with a name input
- Entering a name and clicking "Create" runs `git checkout -b <name>`
- The UI updates to show the new current branch

## Git Command Mapping

### Core Commands Used

| Operation | Git Command | Notes |
|-----------|-------------|-------|
| Get status | `git status --porcelain` | Machine-readable format |
| Get diff | `git diff <file>` | Unstaged changes |
| Get diff (staged) | `git diff --cached <file>` | Staged changes |
| List branches | `git branch` | Local branches |
| Current branch | `git branch --show-current` | Active branch name |
| Switch branch | `git checkout <branch>` | Change branch |
| Create branch | `git checkout -b <name>` | Create and switch |
| Stage files | `git add <files...>` | Stage for commit |
| Commit | `git commit -m "message"` | Create commit |
| Push | `git push` | Push to remote |

### Parsing Requirements

**git status --porcelain format:**
```
 M file1.txt          # Modified, not staged
M  file2.txt          # Modified, staged
A  file3.txt          # Added, staged
?? file4.txt          # Untracked
```

**git branch format:**
```
* main                # Current branch (asterisk)
  feature/test
  bugfix/issue-123
```

**git diff format:**
Standard unified diff format with hunks starting with `@@`

## Out of Scope for MVP

The following features are explicitly **not** included in v1:

- Line-by-line or hunk-by-hunk staging
- Merge conflict resolution
- Git history/log viewing
- Stash management (beyond basic stash)
- Remote repository management
- Submodule support
- Git hooks configuration
- Rebase operations
- Cherry-picking
- Multiple repository windows
- Diff syntax highlighting (basic +/- is sufficient)

## Success Criteria

The MVP is successful when a developer can:

1. Open the app and see their changed files
2. Click on a file and see what changed
3. Select multiple files to commit
4. Write a commit message
5. Commit the changes
6. Push to remote with one click
7. Switch between branches
8. Create a new branch

All operations should complete in under 2 seconds for typical repositories (<1000 files changed).

## Additional Notes

- The app should auto-detect the git repository from the current working directory on startup
- If no git repository is found, show a friendly message to navigate to a git repo
- The UI should be responsive and not block during git operations
- Consider adding keyboard shortcuts for common operations (Ctrl+Enter to commit, etc.)
- Use native system notifications for success/error feedback when appropriate
