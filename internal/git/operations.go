package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Operations provides Git operations
type Operations struct{}

// NewOperations creates a new Git operations instance
func NewOperations() *Operations {
	return &Operations{}
}

// Status returns the working tree status
func (g *Operations) Status(repoPath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	if status.IsClean() {
		return "working tree clean", nil
	}

	var result strings.Builder
	for file, fileStatus := range status {
		result.WriteString(fmt.Sprintf("%s %s\n", string(fileStatus.Staging)+string(fileStatus.Worktree), file))
	}

	return strings.TrimSpace(result.String()), nil
}

// DiffUnstaged returns unstaged changes
func (g *Operations) DiffUnstaged(repoPath string, contextLines int) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	// Get HEAD commit
	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", fmt.Errorf("failed to get commit: %w", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return "", fmt.Errorf("failed to get tree: %w", err)
	}

	// For simplicity, we'll return a placeholder for unstaged changes
	// A full implementation would compare the working tree with HEAD
	_ = tree // avoid unused variable error

	// Get working tree status to check for unstaged changes
	status, err := worktree.Status()
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	var unstagedFiles []string
	for file, fileStatus := range status {
		if fileStatus.Worktree != git.Unmodified {
			unstagedFiles = append(unstagedFiles, file)
		}
	}

	if len(unstagedFiles) == 0 {
		return "no unstaged changes", nil
	}

	var result strings.Builder
	for _, file := range unstagedFiles {
		result.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", file, file))
		result.WriteString(fmt.Sprintf("--- a/%s\n", file))
		result.WriteString(fmt.Sprintf("+++ b/%s\n", file))
		// Note: For simplicity, we're showing a basic diff format
		// A full implementation would show the actual line-by-line differences
		result.WriteString("@@ unstaged changes @@\n")
	}

	return strings.TrimSpace(result.String()), nil
}

// DiffStaged returns staged changes
func (g *Operations) DiffStaged(repoPath string, contextLines int) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	// Get HEAD commit
	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", fmt.Errorf("failed to get commit: %w", err)
	}

	_, err = commit.Tree()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD tree: %w", err)
	}

	// Get index (staged changes)
	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	var stagedFiles []string
	for file, fileStatus := range status {
		if fileStatus.Staging != git.Unmodified {
			stagedFiles = append(stagedFiles, file)
		}
	}

	if len(stagedFiles) == 0 {
		return "no staged changes", nil
	}

	var result strings.Builder
	for _, file := range stagedFiles {
		result.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", file, file))
		result.WriteString(fmt.Sprintf("--- a/%s\n", file))
		result.WriteString(fmt.Sprintf("+++ b/%s\n", file))
		result.WriteString("@@ staged changes @@\n")
	}

	return strings.TrimSpace(result.String()), nil
}

// Diff returns differences between current state and target
func (g *Operations) Diff(repoPath, target string, contextLines int) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	// Resolve target reference
	_, err = repo.Reference(plumbing.ReferenceName("refs/heads/"+target), true)
	if err != nil {
		// Try as a commit hash
		targetHash := plumbing.NewHash(target)
		_, err = repo.CommitObject(targetHash)
		if err != nil {
			return "", fmt.Errorf("failed to resolve target '%s': %w", target, err)
		}
	}

	// Get current HEAD
	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("diff between HEAD (%s) and %s\n", head.Hash().String()[:7], target))
	result.WriteString("(detailed diff implementation would go here)\n")

	return result.String(), nil
}

// Commit creates a new commit with the given message
func (g *Operations) Commit(repoPath, message string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	// Create commit
	hash, err := worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "MCP Git Server",
			Email: "mcp-git@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}

	return fmt.Sprintf("Changes committed successfully with hash %s", hash.String()), nil
}

// Add stages files for commit
func (g *Operations) Add(repoPath string, files []string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	for _, file := range files {
		if file == "." {
			// Add all files
			_, err = worktree.Add(".")
			if err != nil {
				return "", fmt.Errorf("failed to add all files: %w", err)
			}
		} else {
			_, err = worktree.Add(file)
			if err != nil {
				return "", fmt.Errorf("failed to add file %s: %w", file, err)
			}
		}
	}

	return "Files staged successfully", nil
}

// Reset unstages all staged changes
func (g *Operations) Reset(repoPath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	// Get HEAD commit
	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	err = worktree.Reset(&git.ResetOptions{
		Commit: head.Hash(),
		Mode:   git.MixedReset,
	})
	if err != nil {
		return "", fmt.Errorf("failed to reset: %w", err)
	}

	return "All staged changes reset", nil
}

// Log returns commit history
func (g *Operations) Log(repoPath string, maxCount int, startTimestamp, endTimestamp string) ([]string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get commit iterator
	commitIter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}
	defer commitIter.Close()

	var commits []string
	count := 0

	// Parse timestamps if provided
	var startTime, endTime *time.Time
	if startTimestamp != "" {
		t, err := parseTimestamp(startTimestamp)
		if err != nil {
			return nil, fmt.Errorf("invalid start timestamp: %w", err)
		}
		startTime = &t
	}
	if endTimestamp != "" {
		t, err := parseTimestamp(endTimestamp)
		if err != nil {
			return nil, fmt.Errorf("invalid end timestamp: %w", err)
		}
		endTime = &t
	}

	err = commitIter.ForEach(func(commit *object.Commit) error {
		if count >= maxCount {
			return fmt.Errorf("max count reached")
		}

		// Filter by timestamp if provided
		if startTime != nil && commit.Author.When.Before(*startTime) {
			return nil
		}
		if endTime != nil && commit.Author.When.After(*endTime) {
			return nil
		}

		commitStr := fmt.Sprintf("Commit: %s\nAuthor: %s\nDate: %s\nMessage: %s\n",
			commit.Hash.String(),
			commit.Author.Name,
			commit.Author.When.Format(time.RFC3339),
			strings.TrimSpace(commit.Message))

		commits = append(commits, commitStr)
		count++
		return nil
	})

	if err != nil && err.Error() != "max count reached" {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	return commits, nil
}

// CreateBranch creates a new branch
func (g *Operations) CreateBranch(repoPath, branchName, baseBranch string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	var baseRef *plumbing.Reference
	if baseBranch != "" {
		baseRef, err = repo.Reference(plumbing.ReferenceName("refs/heads/"+baseBranch), true)
		if err != nil {
			return "", fmt.Errorf("failed to find base branch %s: %w", baseBranch, err)
		}
	} else {
		baseRef, err = repo.Head()
		if err != nil {
			return "", fmt.Errorf("failed to get HEAD: %w", err)
		}
	}

	// Create new branch
	branchRef := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+branchName), baseRef.Hash())
	err = repo.Storer.SetReference(branchRef)
	if err != nil {
		return "", fmt.Errorf("failed to create branch: %w", err)
	}

	baseName := "HEAD"
	if baseBranch != "" {
		baseName = baseBranch
	}

	return fmt.Sprintf("Created branch '%s' from '%s'", branchName, baseName), nil
}

// Checkout switches to a branch
func (g *Operations) Checkout(repoPath, branchName string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branchName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to checkout branch: %w", err)
	}

	return fmt.Sprintf("Switched to branch '%s'", branchName), nil
}

// Show displays the contents of a commit
func (g *Operations) Show(repoPath, revision string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	// Parse revision
	hash := plumbing.NewHash(revision)
	commit, err := repo.CommitObject(hash)
	if err != nil {
		return "", fmt.Errorf("failed to get commit %s: %w", revision, err)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Commit: %s\n", commit.Hash.String()))
	result.WriteString(fmt.Sprintf("Author: %s\n", commit.Author.Name))
	result.WriteString(fmt.Sprintf("Date: %s\n", commit.Author.When.Format(time.RFC3339)))
	result.WriteString(fmt.Sprintf("Message: %s\n\n", strings.TrimSpace(commit.Message)))

	// Show diff (simplified)
	if len(commit.ParentHashes) > 0 {
		parent, err := repo.CommitObject(commit.ParentHashes[0])
		if err == nil {
			parentTree, _ := parent.Tree()
			commitTree, _ := commit.Tree()
			if parentTree != nil && commitTree != nil {
				changes, err := parentTree.Diff(commitTree)
				if err == nil {
					for _, change := range changes {
						result.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", change.From.Name, change.To.Name))
					}
				}
			}
		}
	}

	return result.String(), nil
}

// Branch lists branches
func (g *Operations) Branch(repoPath, branchType, contains, notContains string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	var refs []*plumbing.Reference
	var result strings.Builder

	switch branchType {
	case "local":
		branchRefs, err := repo.Branches()
		if err != nil {
			return "", fmt.Errorf("failed to get local branches: %w", err)
		}
		err = branchRefs.ForEach(func(ref *plumbing.Reference) error {
			refs = append(refs, ref)
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("failed to iterate branches: %w", err)
		}

	case "remote":
		remoteRefs, err := repo.References()
		if err != nil {
			return "", fmt.Errorf("failed to get references: %w", err)
		}
		err = remoteRefs.ForEach(func(ref *plumbing.Reference) error {
			if ref.Name().IsRemote() {
				refs = append(refs, ref)
			}
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("failed to iterate remote references: %w", err)
		}

	case "all":
		allRefs, err := repo.References()
		if err != nil {
			return "", fmt.Errorf("failed to get references: %w", err)
		}
		err = allRefs.ForEach(func(ref *plumbing.Reference) error {
			if ref.Name().IsBranch() || ref.Name().IsRemote() {
				refs = append(refs, ref)
			}
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("failed to iterate references: %w", err)
		}

	default:
		return "", fmt.Errorf("invalid branch type: %s", branchType)
	}

	// Get current branch
	head, err := repo.Head()
	var currentBranch string
	if err == nil {
		currentBranch = head.Name().Short()
	}

	for _, ref := range refs {
		branchName := ref.Name().Short()
		if ref.Name().IsRemote() {
			branchName = strings.TrimPrefix(string(ref.Name()), "refs/remotes/")
		}

		// Mark current branch
		prefix := "  "
		if branchName == currentBranch {
			prefix = "* "
		}

		result.WriteString(fmt.Sprintf("%s%s\n", prefix, branchName))
	}

	return strings.TrimSpace(result.String()), nil
}

// parseTimestamp parses various timestamp formats
func parseTimestamp(timestamp string) (time.Time, error) {
	// Try different formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
		"Jan 2 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestamp); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", timestamp)
}

// RawCommand executes a raw Git command directly
func (g *Operations) RawCommand(repoPath, command string) (string, error) {
	// Parse the command to extract git subcommand and arguments
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	// Ensure the first part is "git"
	if parts[0] != "git" {
		return "", fmt.Errorf("command must start with 'git'")
	}

	// Remove "git" from the beginning
	args := parts[1:]
	
	// Create the command
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	
	// Execute the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git command failed: %s\nOutput: %s", err.Error(), string(output))
	}
	
	return string(output), nil
}

// Init initializes a new Git repository
func (g *Operations) Init(repoPath string, bare bool) (string, error) {
	if repoPath == "" {
		return "", fmt.Errorf("repository path cannot be empty")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	var repo *git.Repository
	var err error

	if bare {
		repo, err = git.PlainInitWithOptions(repoPath, &git.PlainInitOptions{
			Bare: true,
		})
	} else {
		repo, err = git.PlainInit(repoPath, false)
	}

	if err != nil {
		return "", fmt.Errorf("failed to initialize repository: %w", err)
	}

	repoType := "regular"
	if bare {
		repoType = "bare"
	}

	_ = repo // avoid unused variable warning
	return fmt.Sprintf("Initialized empty Git repository (%s) in %s", repoType, repoPath), nil
}

// Push pushes changes to remote repository
func (g *Operations) Push(repoPath, remote, refspec string, tags bool) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	// Get remote
	if remote == "" {
		remote = "origin"
	}

	remoteObj, err := repo.Remote(remote)
	if err != nil {
		return "", fmt.Errorf("failed to get remote '%s': %w", remote, err)
	}

	// Prepare push options
	pushOptions := &git.PushOptions{}

	// If refspec is provided, use it
	if refspec != "" {
		pushOptions.RefSpecs = []config.RefSpec{config.RefSpec(refspec)}
	}

	// If tags flag is set, push tags
	if tags {
		pushOptions.RefSpecs = append(pushOptions.RefSpecs, config.RefSpec("refs/tags/*:refs/tags/*"))
	}

	err = remoteObj.Push(pushOptions)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return "Everything up-to-date", nil
		}
		return "", fmt.Errorf("failed to push: %w", err)
	}

	result := fmt.Sprintf("Successfully pushed to %s", remote)
	if tags {
		result += " (including tags)"
	}
	if refspec != "" {
		result += fmt.Sprintf(" with refspec: %s", refspec)
	}

	return result, nil
}

// ListRepositories lists Git repositories in a directory
func (g *Operations) ListRepositories(searchPath string, recursive bool) ([]string, error) {
	if searchPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		searchPath = cwd
	}

	var repositories []string

	if recursive {
		err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Continue walking even if there's an error
			}

			if info.IsDir() && info.Name() == ".git" {
				repoPath := filepath.Dir(path)
				repositories = append(repositories, repoPath)
				return filepath.SkipDir // Don't walk into .git directory
			}

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory: %w", err)
		}
	} else {
		// Check if the current directory is a Git repository
		gitDir := filepath.Join(searchPath, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			repositories = append(repositories, searchPath)
		}
	}

	return repositories, nil
}

// CreateTag creates a new Git tag
func (g *Operations) CreateTag(repoPath, tagName, message string, annotated bool) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	// Get HEAD commit
	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	if annotated {
		// Create annotated tag
		_, err = repo.CreateTag(tagName, head.Hash(), &git.CreateTagOptions{
			Tagger: &object.Signature{
				Name:  "MCP Git Server",
				Email: "mcp-git@example.com",
				When:  time.Now(),
			},
			Message: message,
		})
	} else {
		// Create lightweight tag
		tagRef := plumbing.NewHashReference(plumbing.ReferenceName("refs/tags/"+tagName), head.Hash())
		err = repo.Storer.SetReference(tagRef)
	}

	if err != nil {
		return "", fmt.Errorf("failed to create tag: %w", err)
	}

	tagType := "lightweight"
	if annotated {
		tagType = "annotated"
	}

	result := fmt.Sprintf("Created %s tag '%s' at %s", tagType, tagName, head.Hash().String()[:7])
	if message != "" {
		result += fmt.Sprintf(" with message: %s", message)
	}

	return result, nil
}

// DeleteTag deletes a Git tag
func (g *Operations) DeleteTag(repoPath, tagName string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	// Check if tag exists
	tagRef := plumbing.ReferenceName("refs/tags/" + tagName)
	_, err = repo.Reference(tagRef, true)
	if err != nil {
		return "", fmt.Errorf("tag '%s' not found: %w", tagName, err)
	}

	// Delete the tag
	err = repo.Storer.RemoveReference(tagRef)
	if err != nil {
		return "", fmt.Errorf("failed to delete tag: %w", err)
	}

	return fmt.Sprintf("Deleted tag '%s'", tagName), nil
}

// ListTags lists all Git tags
func (g *Operations) ListTags(repoPath string, pattern string) ([]string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	tagRefs, err := repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	var tags []string
	err = tagRefs.ForEach(func(ref *plumbing.Reference) error {
		tagName := strings.TrimPrefix(string(ref.Name()), "refs/tags/")
		
		// Apply pattern filter if provided
		if pattern != "" {
			matched, err := filepath.Match(pattern, tagName)
			if err != nil {
				return err
			}
			if !matched {
				return nil
			}
		}
		
		tags = append(tags, tagName)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate tags: %w", err)
	}

	return tags, nil
}

// PushTags pushes tags to remote repository
func (g *Operations) PushTags(repoPath, remote string, tagName string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	if remote == "" {
		remote = "origin"
	}

	remoteObj, err := repo.Remote(remote)
	if err != nil {
		return "", fmt.Errorf("failed to get remote '%s': %w", remote, err)
	}

	var refSpecs []config.RefSpec
	var message string

	if tagName != "" {
		// Push specific tag
		refSpecs = []config.RefSpec{config.RefSpec("refs/tags/" + tagName + ":refs/tags/" + tagName)}
		message = fmt.Sprintf("Pushed tag '%s' to %s", tagName, remote)
	} else {
		// Push all tags
		refSpecs = []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")}
		message = fmt.Sprintf("Pushed all tags to %s", remote)
	}

	err = remoteObj.Push(&git.PushOptions{
		RefSpecs: refSpecs,
	})

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return "Everything up-to-date", nil
		}
		return "", fmt.Errorf("failed to push tags: %w", err)
	}

	return message, nil
}
