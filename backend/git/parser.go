package git

import (
	"fmt"
	"strconv"
	"strings"

	"git-gui/backend/types"
)

// ParseGitStatus parses the output of `git status --porcelain` into FileStatus structs.
func ParseGitStatus(output string) ([]types.FileStatus, error) {
	if strings.TrimSpace(output) == "" {
		return []types.FileStatus{}, nil
	}

	var files []types.FileStatus
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	for _, line := range lines {
		if len(line) < 4 {
			continue
		}

		indexStatus := line[0]
		workTreeStatus := line[1]
		path := strings.TrimSpace(line[3:])

		file := types.FileStatus{Path: path}

		switch {
		case indexStatus == '?' && workTreeStatus == '?':
			file.Status = types.StatusUntracked
			file.Staged = false
		case indexStatus == 'A':
			file.Status = types.StatusAdded
			file.Staged = true
		case indexStatus == 'D':
			file.Status = types.StatusDeleted
			file.Staged = true
		case indexStatus == 'R':
			file.Status = types.StatusRenamed
			file.Staged = true
		case indexStatus == 'M':
			file.Status = types.StatusModified
			file.Staged = true
		case workTreeStatus == 'M':
			file.Status = types.StatusModified
			file.Staged = false
		case workTreeStatus == 'D':
			file.Status = types.StatusDeleted
			file.Staged = false
		default:
			file.Status = types.StatusModified
			file.Staged = false
		}

		files = append(files, file)
	}

	return files, nil
}

// ParseBranches parses the output of `git branch` into Branch structs.
func ParseBranches(output string) ([]types.Branch, error) {
	if strings.TrimSpace(output) == "" {
		return []types.Branch{}, nil
	}

	var branches []types.Branch
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	for _, line := range lines {
		if len(line) < 3 {
			continue
		}

		isCurrent := line[0] == '*'
		name := strings.TrimSpace(line[2:])

		branches = append(branches, types.Branch{
			Name:      name,
			IsCurrent: isCurrent,
			IsRemote:  false,
		})
	}

	return branches, nil
}

// ParseDiff parses unified diff output into a DiffResult.
func ParseDiff(filePath, output string) (*types.DiffResult, error) {
	result := &types.DiffResult{
		FilePath: filePath,
		Diff:     output,
		Hunks:    []types.DiffHunk{},
	}

	if strings.TrimSpace(output) == "" {
		return result, nil
	}

	lines := strings.Split(output, "\n")
	var currentHunk *types.DiffHunk

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			if currentHunk != nil {
				result.Hunks = append(result.Hunks, *currentHunk)
			}
			hunk, err := ParseHunkHeader(line)
			if err != nil {
				currentHunk = &types.DiffHunk{Header: line}
			} else {
				currentHunk = hunk
			}
		} else if currentHunk != nil {
			currentHunk.Lines = append(currentHunk.Lines, line)
		}
	}

	if currentHunk != nil {
		result.Hunks = append(result.Hunks, *currentHunk)
	}

	return result, nil
}

// ParseHunkHeader parses a @@ hunk header line like "@@ -1,3 +1,4 @@".
func ParseHunkHeader(header string) (*types.DiffHunk, error) {
	hunk := &types.DiffHunk{Header: header}

	// Extract the range info between @@ markers
	parts := strings.SplitN(header, "@@", 3)
	if len(parts) < 2 {
		return hunk, fmt.Errorf("invalid hunk header: %s", header)
	}

	rangeInfo := strings.TrimSpace(parts[1])
	ranges := strings.Split(rangeInfo, " ")
	if len(ranges) < 2 {
		return hunk, fmt.Errorf("invalid range info: %s", rangeInfo)
	}

	// Parse old range (-X,Y)
	oldRange := strings.TrimPrefix(ranges[0], "-")
	oldParts := strings.Split(oldRange, ",")
	if len(oldParts) >= 1 {
		hunk.OldStart, _ = strconv.Atoi(oldParts[0])
	}
	if len(oldParts) >= 2 {
		hunk.OldLines, _ = strconv.Atoi(oldParts[1])
	}

	// Parse new range (+X,Y)
	newRange := strings.TrimPrefix(ranges[1], "+")
	newParts := strings.Split(newRange, ",")
	if len(newParts) >= 1 {
		hunk.NewStart, _ = strconv.Atoi(newParts[0])
	}
	if len(newParts) >= 2 {
		hunk.NewLines, _ = strconv.Atoi(newParts[1])
	}

	return hunk, nil
}

// ExtractCommitSHA extracts the short commit SHA from git commit output.
func ExtractCommitSHA(output string) string {
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
