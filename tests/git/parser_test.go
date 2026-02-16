package git_test

import (
	"testing"

	"git-gui/backend/git"
	"git-gui/backend/types"

	"github.com/stretchr/testify/assert"
)

func TestParseGitStatus_AllStatusTypes(t *testing.T) {
	input := ` M modified.txt
M  staged.txt
A  added.txt
D  deleted.txt
?? untracked.txt
R  renamed.txt`

	files, err := git.ParseGitStatus(input)

	assert.NoError(t, err)
	assert.Len(t, files, 6)

	assert.Equal(t, "modified.txt", files[0].Path)
	assert.Equal(t, types.StatusModified, files[0].Status)
	assert.False(t, files[0].Staged)

	assert.Equal(t, "staged.txt", files[1].Path)
	assert.Equal(t, types.StatusModified, files[1].Status)
	assert.True(t, files[1].Staged)

	assert.Equal(t, "added.txt", files[2].Path)
	assert.Equal(t, types.StatusAdded, files[2].Status)
	assert.True(t, files[2].Staged)

	assert.Equal(t, "deleted.txt", files[3].Path)
	assert.Equal(t, types.StatusDeleted, files[3].Status)
	assert.True(t, files[3].Staged)

	assert.Equal(t, "untracked.txt", files[4].Path)
	assert.Equal(t, types.StatusUntracked, files[4].Status)
	assert.False(t, files[4].Staged)

	assert.Equal(t, "renamed.txt", files[5].Path)
	assert.Equal(t, types.StatusRenamed, files[5].Status)
	assert.True(t, files[5].Staged)
}

func TestParseGitStatus_EmptyOutput(t *testing.T) {
	files, err := git.ParseGitStatus("")

	assert.NoError(t, err)
	assert.Empty(t, files)
}

func TestParseGitStatus_WhitespaceOnly(t *testing.T) {
	files, err := git.ParseGitStatus("   \n  \n")

	assert.NoError(t, err)
	assert.Empty(t, files)
}

func TestParseGitStatus_WorktreeDeleted(t *testing.T) {
	input := " D deleted-from-worktree.txt"

	files, err := git.ParseGitStatus(input)

	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, types.StatusDeleted, files[0].Status)
	assert.False(t, files[0].Staged)
}

func TestParseGitStatus_SkipsShortLines(t *testing.T) {
	input := "ab\n M valid.txt\nxy"

	files, err := git.ParseGitStatus(input)

	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, "valid.txt", files[0].Path)
}

func TestParseBranches_MultipleBranches(t *testing.T) {
	input := `* main
  develop
  feature/login`

	branches, err := git.ParseBranches(input)

	assert.NoError(t, err)
	assert.Len(t, branches, 3)

	assert.Equal(t, "main", branches[0].Name)
	assert.True(t, branches[0].IsCurrent)
	assert.False(t, branches[0].IsRemote)

	assert.Equal(t, "develop", branches[1].Name)
	assert.False(t, branches[1].IsCurrent)

	assert.Equal(t, "feature/login", branches[2].Name)
	assert.False(t, branches[2].IsCurrent)
}

func TestParseBranches_EmptyOutput(t *testing.T) {
	branches, err := git.ParseBranches("")

	assert.NoError(t, err)
	assert.Empty(t, branches)
}

func TestParseBranches_SingleBranch(t *testing.T) {
	input := "* main\n"

	branches, err := git.ParseBranches(input)

	assert.NoError(t, err)
	assert.Len(t, branches, 1)
	assert.Equal(t, "main", branches[0].Name)
	assert.True(t, branches[0].IsCurrent)
}

func TestParseDiff_WithHunks(t *testing.T) {
	input := `diff --git a/file.txt b/file.txt
index abc123..def456 100644
--- a/file.txt
+++ b/file.txt
@@ -1,3 +1,4 @@
 line1
+added line
 line2
 line3`

	result, err := git.ParseDiff("file.txt", input)

	assert.NoError(t, err)
	assert.Equal(t, "file.txt", result.FilePath)
	assert.Equal(t, input, result.Diff)
	assert.Len(t, result.Hunks, 1)

	hunk := result.Hunks[0]
	assert.Equal(t, 1, hunk.OldStart)
	assert.Equal(t, 3, hunk.OldLines)
	assert.Equal(t, 1, hunk.NewStart)
	assert.Equal(t, 4, hunk.NewLines)
	assert.Len(t, hunk.Lines, 4)
	assert.Equal(t, " line1", hunk.Lines[0])
	assert.Equal(t, "+added line", hunk.Lines[1])
}

func TestParseDiff_EmptyOutput(t *testing.T) {
	result, err := git.ParseDiff("file.txt", "")

	assert.NoError(t, err)
	assert.Equal(t, "file.txt", result.FilePath)
	assert.Empty(t, result.Hunks)
}

func TestParseDiff_MultipleHunks(t *testing.T) {
	input := `@@ -1,3 +1,4 @@
 line1
+new line
 line2
@@ -10,3 +11,2 @@
 line10
-removed line
 line11`

	result, err := git.ParseDiff("multi.txt", input)

	assert.NoError(t, err)
	assert.Len(t, result.Hunks, 2)

	assert.Equal(t, 1, result.Hunks[0].OldStart)
	assert.Equal(t, 3, result.Hunks[0].OldLines)
	assert.Len(t, result.Hunks[0].Lines, 3)

	assert.Equal(t, 10, result.Hunks[1].OldStart)
	assert.Equal(t, 3, result.Hunks[1].OldLines)
	assert.Equal(t, 11, result.Hunks[1].NewStart)
	assert.Equal(t, 2, result.Hunks[1].NewLines)
	assert.Len(t, result.Hunks[1].Lines, 3)
}

func TestParseHunkHeader_Standard(t *testing.T) {
	hunk, err := git.ParseHunkHeader("@@ -1,3 +1,4 @@")

	assert.NoError(t, err)
	assert.Equal(t, 1, hunk.OldStart)
	assert.Equal(t, 3, hunk.OldLines)
	assert.Equal(t, 1, hunk.NewStart)
	assert.Equal(t, 4, hunk.NewLines)
}

func TestParseHunkHeader_WithContext(t *testing.T) {
	hunk, err := git.ParseHunkHeader("@@ -10,6 +12,8 @@ func main()")

	assert.NoError(t, err)
	assert.Equal(t, 10, hunk.OldStart)
	assert.Equal(t, 6, hunk.OldLines)
	assert.Equal(t, 12, hunk.NewStart)
	assert.Equal(t, 8, hunk.NewLines)
}

func TestParseHunkHeader_InvalidFormat(t *testing.T) {
	_, err := git.ParseHunkHeader("not a hunk header")

	assert.Error(t, err)
}

func TestExtractCommitSHA(t *testing.T) {
	output := "[main abc1234] Add new feature\n 1 file changed, 2 insertions(+)\n"

	sha := git.ExtractCommitSHA(output)

	assert.Equal(t, "abc1234", sha)
}

func TestExtractCommitSHA_NoMatch(t *testing.T) {
	sha := git.ExtractCommitSHA("no commit info here")

	assert.Equal(t, "", sha)
}
