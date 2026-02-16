package types

// StatusType represents the git status of a file.
type StatusType string

const (
	StatusModified  StatusType = "modified"
	StatusAdded     StatusType = "added"
	StatusDeleted   StatusType = "deleted"
	StatusUntracked StatusType = "untracked"
	StatusRenamed   StatusType = "renamed"
)

// GitRepo represents the git repository state.
type GitRepo struct {
	Path          string `json:"Path"`
	CurrentBranch string `json:"CurrentBranch"`
}

// FileStatus represents a single file's git status.
type FileStatus struct {
	Path   string     `json:"Path"`
	Status StatusType `json:"Status"`
	Staged bool       `json:"Staged"`
}

// Branch represents a git branch.
type Branch struct {
	Name      string `json:"Name"`
	IsCurrent bool   `json:"IsCurrent"`
	IsRemote  bool   `json:"IsRemote"`
}

// DiffResult represents the diff output for a file.
type DiffResult struct {
	FilePath string     `json:"FilePath"`
	Diff     string     `json:"Diff"`
	Hunks    []DiffHunk `json:"Hunks"`
}

// DiffHunk represents an individual change block in a diff.
type DiffHunk struct {
	Header   string   `json:"Header"`
	OldStart int      `json:"OldStart"`
	OldLines int      `json:"OldLines"`
	NewStart int      `json:"NewStart"`
	NewLines int      `json:"NewLines"`
	Lines    []string `json:"Lines"`
}

// CommitResult represents the result of a commit operation.
type CommitResult struct {
	Success   bool   `json:"Success"`
	CommitSHA string `json:"CommitSHA"`
	Message   string `json:"Message"`
}
