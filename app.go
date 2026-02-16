package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// App is the main application struct.
type App struct {
	ctx      context.Context
	executor GitExecutor
	repo     *GitRepo
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Auto-detect git repo from current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	_ = a.InitRepo(cwd)
}

// InitRepo initializes the app with a git repository at the given path.
func (a *App) InitRepo(path string) error {
	valid, err := a.ValidateRepo(path)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("not a git repository: %s", path)
	}

	executor := NewGitExecutor(path)
	root, err := executor.Execute("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("failed to find repo root: %w", err)
	}

	repoPath := strings.TrimSpace(root)
	a.executor = NewGitExecutor(repoPath)
	a.repo = &GitRepo{Path: repoPath}

	branch, err := a.GetCurrentBranch()
	if err == nil {
		a.repo.CurrentBranch = branch
	}

	return nil
}

// GetCurrentRepo returns the current git repository info.
func (a *App) GetCurrentRepo() (*GitRepo, error) {
	if a.repo == nil {
		return nil, errors.New("no repository initialized")
	}
	return a.repo, nil
}

// ValidateRepo checks whether the given path is inside a git repository.
func (a *App) ValidateRepo(path string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	executor := NewGitExecutor(absPath)
	_, err = executor.Execute("rev-parse", "--git-dir")
	return err == nil, nil
}

// GetRepoRoot returns the root directory of the git repository.
func (a *App) GetRepoRoot() (string, error) {
	if a.repo == nil {
		return "", errors.New("no repository initialized")
	}
	return a.repo.Path, nil
}

// GetGitStatus returns the list of changed files in the repository.
func (a *App) GetGitStatus() ([]FileStatus, error) {
	if a.executor == nil {
		return nil, errors.New("no repository initialized")
	}

	output, err := a.executor.Execute("status", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	files, err := parseGitStatus(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse git status: %w", err)
	}

	return files, nil
}

// GetGitDiff returns the diff for a specific file.
func (a *App) GetGitDiff(filePath string) (*DiffResult, error) {
	if a.executor == nil {
		return nil, errors.New("no repository initialized")
	}

	// Try unstaged diff first
	output, err := a.executor.Execute("diff", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get diff for %s: %w", filePath, err)
	}

	// If no unstaged diff, try staged diff
	if strings.TrimSpace(output) == "" {
		output, err = a.executor.Execute("diff", "--cached", filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get staged diff for %s: %w", filePath, err)
		}
	}

	result, err := parseDiff(filePath, output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse diff for %s: %w", filePath, err)
	}

	return result, nil
}

// GetBranches returns all local branches.
func (a *App) GetBranches() ([]Branch, error) {
	if a.executor == nil {
		return nil, errors.New("no repository initialized")
	}

	output, err := a.executor.Execute("branch")
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	branches, err := parseBranches(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse branches: %w", err)
	}

	return branches, nil
}

// GetCurrentBranch returns the name of the current branch.
func (a *App) GetCurrentBranch() (string, error) {
	if a.executor == nil {
		return "", errors.New("no repository initialized")
	}

	output, err := a.executor.Execute("branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(output), nil
}

// SwitchBranch switches to the specified branch.
func (a *App) SwitchBranch(name string) error {
	if a.executor == nil {
		return errors.New("no repository initialized")
	}

	_, err := a.executor.Execute("checkout", name)
	if err != nil {
		return fmt.Errorf("failed to switch to branch %s: %w", name, err)
	}

	a.repo.CurrentBranch = name
	return nil
}

// CreateBranch creates a new branch and switches to it.
func (a *App) CreateBranch(name string) error {
	if a.executor == nil {
		return errors.New("no repository initialized")
	}

	_, err := a.executor.Execute("checkout", "-b", name)
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %w", name, err)
	}

	a.repo.CurrentBranch = name
	return nil
}

// CommitFiles stages the specified files and creates a commit.
func (a *App) CommitFiles(files []string, message string) (*CommitResult, error) {
	if a.executor == nil {
		return nil, errors.New("no repository initialized")
	}
	if len(files) == 0 {
		return nil, errors.New("no files to commit")
	}
	if message == "" {
		return nil, errors.New("commit message required")
	}

	// Stage files
	args := append([]string{"add"}, files...)
	_, err := a.executor.Execute(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to stage files: %w", err)
	}

	// Commit
	output, err := a.executor.Execute("commit", "-m", message)
	if err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	// Extract short SHA from commit output
	sha := extractCommitSHA(output)

	return &CommitResult{
		Success:   true,
		CommitSHA: sha,
		Message:   message,
	}, nil
}

// PushChanges pushes committed changes to the remote.
func (a *App) PushChanges() error {
	if a.executor == nil {
		return errors.New("no repository initialized")
	}

	_, err := a.executor.Execute("push")
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

// CommitAndPush stages files, commits, and pushes in one operation.
func (a *App) CommitAndPush(files []string, message string) (*CommitResult, error) {
	result, err := a.CommitFiles(files, message)
	if err != nil {
		return nil, err
	}

	if err := a.PushChanges(); err != nil {
		return nil, fmt.Errorf("commit succeeded but push failed: %w", err)
	}

	return result, nil
}

// extractCommitSHA extracts the short commit SHA from git commit output.
func extractCommitSHA(output string) string {
	// git commit output typically contains "[branch SHA] message"
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "]") {
			start := strings.Index(line, " ")
			end := strings.Index(line, "]")
			if start != -1 && end != -1 && start < end {
				return strings.TrimSpace(line[start+1 : end])
			}
		}
	}
	return ""
}
