package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitExecutor defines the interface for executing git commands.
type GitExecutor interface {
	Execute(args ...string) (string, error)
}

// RealGitExecutor executes git commands as subprocesses.
type RealGitExecutor struct {
	repoPath string
}

func NewGitExecutor(repoPath string) *RealGitExecutor {
	return &RealGitExecutor{repoPath: repoPath}
}

func (e *RealGitExecutor) Execute(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = e.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %s", strings.Join(args, " "), strings.TrimSpace(string(output)))
	}

	return string(output), nil
}
