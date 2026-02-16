package main

import (
	"fmt"
	"strconv"
	"strings"
)

// parseGitStatus parses the output of `git status --porcelain` into FileStatus structs.
func parseGitStatus(output string) ([]FileStatus, error) {
	if strings.TrimSpace(output) == "" {
		return []FileStatus{}, nil
	}

	var files []FileStatus
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	for _, line := range lines {
		if len(line) < 4 {
			continue
		}

		indexStatus := line[0]
		workTreeStatus := line[1]
		path := strings.TrimSpace(line[3:])

		file := FileStatus{Path: path}

		switch {
		case indexStatus == '?' && workTreeStatus == '?':
			file.Status = StatusUntracked
			file.Staged = false
		case indexStatus == 'A':
			file.Status = StatusAdded
			file.Staged = true
		case indexStatus == 'D':
			file.Status = StatusDeleted
			file.Staged = true
		case indexStatus == 'R':
			file.Status = StatusRenamed
			file.Staged = true
		case indexStatus == 'M':
			file.Status = StatusModified
			file.Staged = true
		case workTreeStatus == 'M':
			file.Status = StatusModified
			file.Staged = false
		case workTreeStatus == 'D':
			file.Status = StatusDeleted
			file.Staged = false
		default:
			file.Status = StatusModified
			file.Staged = false
		}

		files = append(files, file)
	}

	return files, nil
}

// parseBranches parses the output of `git branch` into Branch structs.
func parseBranches(output string) ([]Branch, error) {
	if strings.TrimSpace(output) == "" {
		return []Branch{}, nil
	}

	var branches []Branch
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	for _, line := range lines {
		if len(line) < 3 {
			continue
		}

		isCurrent := line[0] == '*'
		name := strings.TrimSpace(line[2:])

		branches = append(branches, Branch{
			Name:      name,
			IsCurrent: isCurrent,
			IsRemote:  false,
		})
	}

	return branches, nil
}

// parseDiff parses unified diff output into a DiffResult.
func parseDiff(filePath, output string) (*DiffResult, error) {
	result := &DiffResult{
		FilePath: filePath,
		Diff:     output,
		Hunks:    []DiffHunk{},
	}

	if strings.TrimSpace(output) == "" {
		return result, nil
	}

	lines := strings.Split(output, "\n")
	var currentHunk *DiffHunk

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			if currentHunk != nil {
				result.Hunks = append(result.Hunks, *currentHunk)
			}
			hunk, err := parseHunkHeader(line)
			if err != nil {
				currentHunk = &DiffHunk{Header: line}
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

// parseHunkHeader parses a @@ hunk header line like "@@ -1,3 +1,4 @@".
func parseHunkHeader(header string) (*DiffHunk, error) {
	hunk := &DiffHunk{Header: header}

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
