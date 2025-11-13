package git

// GitStatus represents the parameters for git status
type GitStatus struct {
	RepoPath string `json:"repo_path"`
}

// GitDiffUnstaged represents the parameters for git diff (unstaged)
type GitDiffUnstaged struct {
	RepoPath     string `json:"repo_path"`
	ContextLines int    `json:"context_lines,omitempty"`
}

// GitDiffStaged represents the parameters for git diff --cached
type GitDiffStaged struct {
	RepoPath     string `json:"repo_path"`
	ContextLines int    `json:"context_lines,omitempty"`
}

// GitDiff represents the parameters for git diff with target
type GitDiff struct {
	RepoPath     string `json:"repo_path"`
	Target       string `json:"target"`
	ContextLines int    `json:"context_lines,omitempty"`
}

// GitCommit represents the parameters for git commit
type GitCommit struct {
	RepoPath string `json:"repo_path"`
	Message  string `json:"message"`
}

// GitAdd represents the parameters for git add
type GitAdd struct {
	RepoPath string   `json:"repo_path"`
	Files    []string `json:"files"`
}

// GitReset represents the parameters for git reset
type GitReset struct {
	RepoPath string `json:"repo_path"`
}

// GitLog represents the parameters for git log
type GitLog struct {
	RepoPath       string `json:"repo_path"`
	MaxCount       int    `json:"max_count,omitempty"`
	StartTimestamp string `json:"start_timestamp,omitempty"`
	EndTimestamp   string `json:"end_timestamp,omitempty"`
}

// GitCreateBranch represents the parameters for creating a branch
type GitCreateBranch struct {
	RepoPath   string `json:"repo_path"`
	BranchName string `json:"branch_name"`
	BaseBranch string `json:"base_branch,omitempty"`
}

// GitCheckout represents the parameters for git checkout
type GitCheckout struct {
	RepoPath   string `json:"repo_path"`
	BranchName string `json:"branch_name"`
}

// GitShow represents the parameters for git show
type GitShow struct {
	RepoPath string `json:"repo_path"`
	Revision string `json:"revision"`
}

// GitBranch represents the parameters for git branch
type GitBranch struct {
	RepoPath    string `json:"repo_path"`
	BranchType  string `json:"branch_type"`
	Contains    string `json:"contains,omitempty"`
	NotContains string `json:"not_contains,omitempty"`
}

// Default number of context lines for diff operations
const DefaultContextLines = 3
