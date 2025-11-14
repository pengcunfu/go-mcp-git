package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func createTestRepo(t *testing.T) (string, *git.Repository) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize repository
	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add and commit the file
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}

	_, err = worktree.Add("test.txt")
	if err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	return tempDir, repo
}

func TestOperations_Status(t *testing.T) {
	tempDir, _ := createTestRepo(t)
	defer os.RemoveAll(tempDir)

	ops := NewOperations("Test User", "test@example.com")

	// Test clean status
	status, err := ops.Status(tempDir)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	if status != "working tree clean" {
		t.Errorf("Expected clean status, got: %s", status)
	}

	// Create a new file to test dirty status
	newFile := filepath.Join(tempDir, "new.txt")
	err = os.WriteFile(newFile, []byte("new content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	status, err = ops.Status(tempDir)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	if status == "working tree clean" {
		t.Errorf("Expected dirty status, got clean")
	}
}

func TestOperations_Add(t *testing.T) {
	tempDir, _ := createTestRepo(t)
	defer os.RemoveAll(tempDir)

	ops := NewOperations("Test User", "test@example.com")

	// Create a new file
	newFile := filepath.Join(tempDir, "new.txt")
	err := os.WriteFile(newFile, []byte("new content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Add the file
	result, err := ops.Add(tempDir, []string{"new.txt"})
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if result != "Files staged successfully" {
		t.Errorf("Expected success message, got: %s", result)
	}
}

func TestOperations_Commit(t *testing.T) {
	tempDir, _ := createTestRepo(t)
	defer os.RemoveAll(tempDir)

	ops := NewOperations("Test User", "test@example.com")

	// Create and add a new file
	newFile := filepath.Join(tempDir, "new.txt")
	err := os.WriteFile(newFile, []byte("new content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	_, err = ops.Add(tempDir, []string{"new.txt"})
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Commit the changes
	result, err := ops.Commit(tempDir, "Test commit")
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	if !contains(result, "Changes committed successfully with hash") {
		t.Errorf("Expected commit success message, got: %s", result)
	}
}

func TestOperations_CreateBranch(t *testing.T) {
	tempDir, _ := createTestRepo(t)
	defer os.RemoveAll(tempDir)

	ops := NewOperations("Test User", "test@example.com")

	// Create a new branch
	result, err := ops.CreateBranch(tempDir, "test-branch", "")
	if err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}

	expected := "Created branch 'test-branch' from 'HEAD'"
	if result != expected {
		t.Errorf("Expected '%s', got: %s", expected, result)
	}
}

func TestOperations_Checkout(t *testing.T) {
	tempDir, _ := createTestRepo(t)
	defer os.RemoveAll(tempDir)

	ops := NewOperations("Test User", "test@example.com")

	// Create a new branch first
	_, err := ops.CreateBranch(tempDir, "test-branch", "")
	if err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}

	// Checkout the branch
	result, err := ops.Checkout(tempDir, "test-branch")
	if err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}

	expected := "Switched to branch 'test-branch'"
	if result != expected {
		t.Errorf("Expected '%s', got: %s", expected, result)
	}
}

func TestOperations_Log(t *testing.T) {
	tempDir, _ := createTestRepo(t)
	defer os.RemoveAll(tempDir)

	ops := NewOperations("Test User", "test@example.com")

	// Get log
	commits, err := ops.Log(tempDir, 10, "", "")
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	if len(commits) == 0 {
		t.Error("Expected at least one commit")
	}

	// Check if the commit contains expected fields
	firstCommit := commits[0]
	if !contains(firstCommit, "Commit:") || !contains(firstCommit, "Author:") || !contains(firstCommit, "Date:") || !contains(firstCommit, "Message:") {
		t.Errorf("Commit format incorrect: %s", firstCommit)
	}
}

func TestOperations_Branch(t *testing.T) {
	tempDir, _ := createTestRepo(t)
	defer os.RemoveAll(tempDir)

	ops := NewOperations("Test User", "test@example.com")

	// Create a test branch
	_, err := ops.CreateBranch(tempDir, "test-branch", "")
	if err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}

	// List local branches
	result, err := ops.Branch(tempDir, "local", "", "")
	if err != nil {
		t.Fatalf("Branch failed: %v", err)
	}

	if !contains(result, "test-branch") {
		t.Errorf("Expected test-branch in result, got: %s", result)
	}
}

func TestOperations_Reset(t *testing.T) {
	tempDir, _ := createTestRepo(t)
	defer os.RemoveAll(tempDir)

	ops := NewOperations("Test User", "test@example.com")

	// Create and add a new file
	newFile := filepath.Join(tempDir, "new.txt")
	err := os.WriteFile(newFile, []byte("new content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	_, err = ops.Add(tempDir, []string{"new.txt"})
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Reset staged changes
	result, err := ops.Reset(tempDir)
	if err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	expected := "All staged changes reset"
	if result != expected {
		t.Errorf("Expected '%s', got: %s", expected, result)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
