package backend

import (
	"errors"
	"testing"

	"git-gui/backend/git"
	"git-gui/backend/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGitExecutor implements git.GitExecutor for testing.
type MockGitExecutor struct {
	mock.Mock
}

func (m *MockGitExecutor) Execute(args ...string) (string, error) {
	callArgs := m.Called(args)
	return callArgs.String(0), callArgs.Error(1)
}

func newTestApp(executor git.GitExecutor) *App {
	return &App{
		executor: executor,
		repo:     &types.GitRepo{Path: "/test/repo", CurrentBranch: "main"},
	}
}

func TestGetGitStatus_Success(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"status", "--porcelain"}).
		Return(" M file.txt\n?? new.txt\n", nil)

	app := newTestApp(mockExec)
	files, err := app.GetGitStatus()

	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Equal(t, "file.txt", files[0].Path)
	assert.Equal(t, types.StatusModified, files[0].Status)
	assert.Equal(t, "new.txt", files[1].Path)
	assert.Equal(t, types.StatusUntracked, files[1].Status)
	mockExec.AssertExpectations(t)
}

func TestGetGitStatus_NoRepo(t *testing.T) {
	app := &App{}
	_, err := app.GetGitStatus()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no repository initialized")
}

func TestGetGitStatus_GitError(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"status", "--porcelain"}).
		Return("", errors.New("git failed"))

	app := newTestApp(mockExec)
	_, err := app.GetGitStatus()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get git status")
}

func TestGetGitDiff_UnstagedDiff(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"diff", "file.txt"}).
		Return("@@ -1,3 +1,4 @@\n line1\n+new\n line2\n", nil)

	app := newTestApp(mockExec)
	result, err := app.GetGitDiff("file.txt")

	assert.NoError(t, err)
	assert.Equal(t, "file.txt", result.FilePath)
	assert.Len(t, result.Hunks, 1)
	mockExec.AssertExpectations(t)
}

func TestGetGitDiff_FallsBackToStagedDiff(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"diff", "file.txt"}).
		Return("", nil)
	mockExec.On("Execute", []string{"diff", "--cached", "file.txt"}).
		Return("@@ -1,2 +1,3 @@\n line1\n+staged\n", nil)

	app := newTestApp(mockExec)
	result, err := app.GetGitDiff("file.txt")

	assert.NoError(t, err)
	assert.Len(t, result.Hunks, 1)
	mockExec.AssertExpectations(t)
}

func TestGetGitDiff_NoRepo(t *testing.T) {
	app := &App{}
	_, err := app.GetGitDiff("file.txt")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no repository initialized")
}

func TestGetBranches_Success(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"branch"}).
		Return("* main\n  develop\n", nil)

	app := newTestApp(mockExec)
	branches, err := app.GetBranches()

	assert.NoError(t, err)
	assert.Len(t, branches, 2)
	assert.Equal(t, "main", branches[0].Name)
	assert.True(t, branches[0].IsCurrent)
	assert.Equal(t, "develop", branches[1].Name)
	assert.False(t, branches[1].IsCurrent)
	mockExec.AssertExpectations(t)
}

func TestGetCurrentBranch_Success(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"branch", "--show-current"}).
		Return("main\n", nil)

	app := newTestApp(mockExec)
	branch, err := app.GetCurrentBranch()

	assert.NoError(t, err)
	assert.Equal(t, "main", branch)
}

func TestSwitchBranch_Success(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"checkout", "develop"}).
		Return("Switched to branch 'develop'\n", nil)

	app := newTestApp(mockExec)
	err := app.SwitchBranch("develop")

	assert.NoError(t, err)
	assert.Equal(t, "develop", app.repo.CurrentBranch)
	mockExec.AssertExpectations(t)
}

func TestSwitchBranch_Error(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"checkout", "nonexistent"}).
		Return("", errors.New("branch not found"))

	app := newTestApp(mockExec)
	err := app.SwitchBranch("nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to switch to branch nonexistent")
}

func TestCreateBranch_Success(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"checkout", "-b", "feature/new"}).
		Return("Switched to a new branch 'feature/new'\n", nil)

	app := newTestApp(mockExec)
	err := app.CreateBranch("feature/new")

	assert.NoError(t, err)
	assert.Equal(t, "feature/new", app.repo.CurrentBranch)
	mockExec.AssertExpectations(t)
}

func TestCommitFiles_Success(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"add", "file1.txt", "file2.txt"}).
		Return("", nil)
	mockExec.On("Execute", []string{"commit", "-m", "test commit"}).
		Return("[main abc1234] test commit\n 2 files changed\n", nil)

	app := newTestApp(mockExec)
	result, err := app.CommitFiles([]string{"file1.txt", "file2.txt"}, "test commit")

	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "abc1234", result.CommitSHA)
	assert.Equal(t, "test commit", result.Message)
	mockExec.AssertExpectations(t)
}

func TestCommitFiles_NoFiles(t *testing.T) {
	app := newTestApp(new(MockGitExecutor))
	_, err := app.CommitFiles([]string{}, "message")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no files to commit")
}

func TestCommitFiles_EmptyMessage(t *testing.T) {
	app := newTestApp(new(MockGitExecutor))
	_, err := app.CommitFiles([]string{"file.txt"}, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit message required")
}

func TestCommitFiles_StageError(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"add", "bad.txt"}).
		Return("", errors.New("path not found"))

	app := newTestApp(mockExec)
	_, err := app.CommitFiles([]string{"bad.txt"}, "message")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stage files")
}

func TestPushChanges_Success(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"push"}).
		Return("Everything up-to-date\n", nil)

	app := newTestApp(mockExec)
	err := app.PushChanges()

	assert.NoError(t, err)
	mockExec.AssertExpectations(t)
}

func TestPushChanges_Error(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"push"}).
		Return("", errors.New("no remote configured"))

	app := newTestApp(mockExec)
	err := app.PushChanges()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to push")
}

func TestCommitAndPush_Success(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"add", "file.txt"}).
		Return("", nil)
	mockExec.On("Execute", []string{"commit", "-m", "push me"}).
		Return("[main def5678] push me\n", nil)
	mockExec.On("Execute", []string{"push"}).
		Return("", nil)

	app := newTestApp(mockExec)
	result, err := app.CommitAndPush([]string{"file.txt"}, "push me")

	assert.NoError(t, err)
	assert.True(t, result.Success)
	mockExec.AssertExpectations(t)
}

func TestCommitAndPush_PushFails(t *testing.T) {
	mockExec := new(MockGitExecutor)
	mockExec.On("Execute", []string{"add", "file.txt"}).
		Return("", nil)
	mockExec.On("Execute", []string{"commit", "-m", "msg"}).
		Return("[main abc1234] msg\n", nil)
	mockExec.On("Execute", []string{"push"}).
		Return("", errors.New("remote rejected"))

	app := newTestApp(mockExec)
	_, err := app.CommitAndPush([]string{"file.txt"}, "msg")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit succeeded but push failed")
}

func TestGetCurrentRepo_NoRepo(t *testing.T) {
	app := &App{}
	_, err := app.GetCurrentRepo()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no repository initialized")
}

func TestGetCurrentRepo_Success(t *testing.T) {
	app := newTestApp(nil)
	repo, err := app.GetCurrentRepo()

	assert.NoError(t, err)
	assert.Equal(t, "/test/repo", repo.Path)
	assert.Equal(t, "main", repo.CurrentBranch)
}

func TestGetRepoRoot_Success(t *testing.T) {
	app := newTestApp(nil)
	root, err := app.GetRepoRoot()

	assert.NoError(t, err)
	assert.Equal(t, "/test/repo", root)
}

func TestGetRepoRoot_NoRepo(t *testing.T) {
	app := &App{}
	_, err := app.GetRepoRoot()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no repository initialized")
}
