package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pengcunfu/go-mcp-git/internal/git"
	"github.com/pengcunfu/go-mcp-git/internal/mcp"
)

// Server represents the MCP Git server
type Server struct {
	mcpServer  *mcp.Server
	gitOps     *git.Operations
	repository string
	verbose    int
}

// New creates a new MCP Git server
func New(repository string, verbose int) *Server {
	mcpServer := mcp.NewServer("go-mcp-git", "0.0.1")
	gitOps := git.NewOperations()

	server := &Server{
		mcpServer:  mcpServer,
		gitOps:     gitOps,
		repository: repository,
		verbose:    verbose,
	}

	server.registerTools()
	return server
}

// Serve starts the MCP server
func (s *Server) Serve(ctx context.Context) error {
	if s.verbose > 0 {
		log.Printf("Starting MCP Git server")
		if s.repository != "" {
			log.Printf("Using repository: %s", s.repository)
		}
	}

	return s.mcpServer.Serve(ctx)
}

// registerTools registers all Git tools with the MCP server
func (s *Server) registerTools() {
	// Git Status
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_status",
		Description: "Shows the working tree status",
		InputSchema: s.createSchema("GitStatus", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
			},
			"required": []string{"repo_path"},
		}),
	}, s.handleGitStatus)

	// Git Diff Unstaged
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_diff_unstaged",
		Description: "Shows changes in working directory not yet staged",
		InputSchema: s.createSchema("GitDiffUnstaged", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"context_lines": map[string]interface{}{
					"type":        "integer",
					"description": "Number of context lines to show",
					"default":     git.DefaultContextLines,
				},
			},
			"required": []string{"repo_path"},
		}),
	}, s.handleGitDiffUnstaged)

	// Git Diff Staged
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_diff_staged",
		Description: "Shows changes that are staged for commit",
		InputSchema: s.createSchema("GitDiffStaged", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"context_lines": map[string]interface{}{
					"type":        "integer",
					"description": "Number of context lines to show",
					"default":     git.DefaultContextLines,
				},
			},
			"required": []string{"repo_path"},
		}),
	}, s.handleGitDiffStaged)

	// Git Diff
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_diff",
		Description: "Shows differences between branches or commits",
		InputSchema: s.createSchema("GitDiff", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"target": map[string]interface{}{
					"type":        "string",
					"description": "Target branch or commit to compare with",
				},
				"context_lines": map[string]interface{}{
					"type":        "integer",
					"description": "Number of context lines to show",
					"default":     git.DefaultContextLines,
				},
			},
			"required": []string{"repo_path", "target"},
		}),
	}, s.handleGitDiff)

	// Git Commit
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_commit",
		Description: "Records changes to the repository",
		InputSchema: s.createSchema("GitCommit", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"message": map[string]interface{}{
					"type":        "string",
					"description": "Commit message",
				},
			},
			"required": []string{"repo_path", "message"},
		}),
	}, s.handleGitCommit)

	// Git Add
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_add",
		Description: "Adds file contents to the staging area",
		InputSchema: s.createSchema("GitAdd", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"files": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
					"description": "Array of file paths to stage",
				},
			},
			"required": []string{"repo_path", "files"},
		}),
	}, s.handleGitAdd)

	// Git Reset
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_reset",
		Description: "Unstages all staged changes",
		InputSchema: s.createSchema("GitReset", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
			},
			"required": []string{"repo_path"},
		}),
	}, s.handleGitReset)

	// Git Log
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_log",
		Description: "Shows the commit logs with optional date filtering",
		InputSchema: s.createSchema("GitLog", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"max_count": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of commits to show",
					"default":     10,
				},
				"start_timestamp": map[string]interface{}{
					"type":        "string",
					"description": "Start timestamp for filtering commits",
				},
				"end_timestamp": map[string]interface{}{
					"type":        "string",
					"description": "End timestamp for filtering commits",
				},
			},
			"required": []string{"repo_path"},
		}),
	}, s.handleGitLog)

	// Git Create Branch
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_create_branch",
		Description: "Creates a new branch",
		InputSchema: s.createSchema("GitCreateBranch", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"branch_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the new branch",
				},
				"base_branch": map[string]interface{}{
					"type":        "string",
					"description": "Base branch to create from (defaults to current branch)",
				},
			},
			"required": []string{"repo_path", "branch_name"},
		}),
	}, s.handleGitCreateBranch)

	// Git Checkout
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_checkout",
		Description: "Switches branches",
		InputSchema: s.createSchema("GitCheckout", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"branch_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of branch to checkout",
				},
			},
			"required": []string{"repo_path", "branch_name"},
		}),
	}, s.handleGitCheckout)

	// Git Show
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_show",
		Description: "Shows the contents of a commit",
		InputSchema: s.createSchema("GitShow", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"revision": map[string]interface{}{
					"type":        "string",
					"description": "The revision (commit hash, branch name, tag) to show",
				},
			},
			"required": []string{"repo_path", "revision"},
		}),
	}, s.handleGitShow)

	// Git Branch
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_branch",
		Description: "List Git branches",
		InputSchema: s.createSchema("GitBranch", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"branch_type": map[string]interface{}{
					"type":        "string",
					"description": "Whether to list local branches ('local'), remote branches ('remote') or all branches('all')",
					"enum":        []string{"local", "remote", "all"},
					"default":     "local",
				},
				"contains": map[string]interface{}{
					"type":        "string",
					"description": "The commit sha that branch should contain",
				},
				"not_contains": map[string]interface{}{
					"type":        "string",
					"description": "The commit sha that branch should NOT contain",
				},
			},
			"required": []string{"repo_path"},
		}),
	}, s.handleGitBranch)

	// Git Raw Command
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_raw_command",
		Description: "Execute a raw Git command directly (bypasses shell wrapping issues)",
		InputSchema: s.createSchema("GitRawCommand", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"command": map[string]interface{}{
					"type":        "string",
					"description": "Raw Git command to execute (e.g., 'git tag -a v0.0.1 -m \"Release v0.0.1\"')",
				},
			},
			"required": []string{"repo_path", "command"},
		}),
	}, s.handleGitRawCommand)

	// Git Init
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_init",
		Description: "Initialize a new Git repository",
		InputSchema: s.createSchema("GitInit", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path where to initialize the repository",
				},
				"bare": map[string]interface{}{
					"type":        "boolean",
					"description": "Initialize as bare repository",
					"default":     false,
				},
			},
			"required": []string{"repo_path"},
		}),
	}, s.handleGitInit)

	// Git Push
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_push",
		Description: "Push changes to remote repository",
		InputSchema: s.createSchema("GitPush", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"remote": map[string]interface{}{
					"type":        "string",
					"description": "Remote name (default: origin)",
					"default":     "origin",
				},
				"refspec": map[string]interface{}{
					"type":        "string",
					"description": "Refspec to push (e.g., 'refs/heads/main:refs/heads/main')",
				},
				"tags": map[string]interface{}{
					"type":        "boolean",
					"description": "Push tags along with commits",
					"default":     false,
				},
			},
			"required": []string{"repo_path"},
		}),
	}, s.handleGitPush)

	// Git List Repositories
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_list_repositories",
		Description: "List Git repositories in a directory",
		InputSchema: s.createSchema("GitListRepositories", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"search_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to search for repositories (default: current directory)",
				},
				"recursive": map[string]interface{}{
					"type":        "boolean",
					"description": "Search recursively in subdirectories",
					"default":     false,
				},
			},
		}),
	}, s.handleGitListRepositories)

	// Git Create Tag
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_create_tag",
		Description: "Create a new Git tag",
		InputSchema: s.createSchema("GitCreateTag", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"tag_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the tag to create",
				},
				"message": map[string]interface{}{
					"type":        "string",
					"description": "Tag message (for annotated tags)",
				},
				"annotated": map[string]interface{}{
					"type":        "boolean",
					"description": "Create annotated tag (default: true)",
					"default":     true,
				},
			},
			"required": []string{"repo_path", "tag_name"},
		}),
	}, s.handleGitCreateTag)

	// Git Delete Tag
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_delete_tag",
		Description: "Delete a Git tag",
		InputSchema: s.createSchema("GitDeleteTag", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"tag_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the tag to delete",
				},
			},
			"required": []string{"repo_path", "tag_name"},
		}),
	}, s.handleGitDeleteTag)

	// Git List Tags
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_list_tags",
		Description: "List Git tags",
		InputSchema: s.createSchema("GitListTags", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"pattern": map[string]interface{}{
					"type":        "string",
					"description": "Pattern to filter tags (glob pattern)",
				},
			},
			"required": []string{"repo_path"},
		}),
	}, s.handleGitListTags)

	// Git Push Tags
	s.mcpServer.RegisterTool(mcp.Tool{
		Name:        "git_push_tags",
		Description: "Push tags to remote repository",
		InputSchema: s.createSchema("GitPushTags", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"repo_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to Git repository",
				},
				"remote": map[string]interface{}{
					"type":        "string",
					"description": "Remote name (default: origin)",
					"default":     "origin",
				},
				"tag_name": map[string]interface{}{
					"type":        "string",
					"description": "Specific tag name to push (leave empty to push all tags)",
				},
			},
			"required": []string{"repo_path"},
		}),
	}, s.handleGitPushTags)
}

// createSchema creates a JSON schema for tool input
func (s *Server) createSchema(title string, schemaData map[string]interface{}) interface{} {
	schema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   title,
	}
	
	// Copy all fields from schemaData to schema
	for key, value := range schemaData {
		schema[key] = value
	}
	
	return schema
}

// getRepoPath returns the repository path, using the provided path or the configured default
func (s *Server) getRepoPath(providedPath string) string {
	if providedPath != "" {
		return providedPath
	}
	if s.repository != "" {
		return s.repository
	}
	// Default to current directory
	cwd, _ := os.Getwd()
	return cwd
}

// Tool handlers

func (s *Server) handleGitStatus(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	
	result, err := s.gitOps.Status(repoPath)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: fmt.Sprintf("Repository status:\n%s", result),
	}}, nil
}

func (s *Server) handleGitDiffUnstaged(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	contextLines := getInt(arguments, "context_lines", git.DefaultContextLines)
	
	result, err := s.gitOps.DiffUnstaged(repoPath, contextLines)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: fmt.Sprintf("Unstaged changes:\n%s", result),
	}}, nil
}

func (s *Server) handleGitDiffStaged(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	contextLines := getInt(arguments, "context_lines", git.DefaultContextLines)
	
	result, err := s.gitOps.DiffStaged(repoPath, contextLines)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: fmt.Sprintf("Staged changes:\n%s", result),
	}}, nil
}

func (s *Server) handleGitDiff(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	target := getString(arguments, "target")
	contextLines := getInt(arguments, "context_lines", git.DefaultContextLines)
	
	result, err := s.gitOps.Diff(repoPath, target, contextLines)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: fmt.Sprintf("Diff with %s:\n%s", target, result),
	}}, nil
}

func (s *Server) handleGitCommit(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	message := getString(arguments, "message")
	
	result, err := s.gitOps.Commit(repoPath, message)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitAdd(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	files := getStringSlice(arguments, "files")
	
	result, err := s.gitOps.Add(repoPath, files)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitReset(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	
	result, err := s.gitOps.Reset(repoPath)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitLog(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	maxCount := getInt(arguments, "max_count", 10)
	startTimestamp := getString(arguments, "start_timestamp")
	endTimestamp := getString(arguments, "end_timestamp")
	
	commits, err := s.gitOps.Log(repoPath, maxCount, startTimestamp, endTimestamp)
	if err != nil {
		return nil, err
	}

	result := "Commit history:\n"
	for _, commit := range commits {
		result += commit + "\n"
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitCreateBranch(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	branchName := getString(arguments, "branch_name")
	baseBranch := getString(arguments, "base_branch")
	
	result, err := s.gitOps.CreateBranch(repoPath, branchName, baseBranch)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitCheckout(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	branchName := getString(arguments, "branch_name")
	
	result, err := s.gitOps.Checkout(repoPath, branchName)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitShow(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	revision := getString(arguments, "revision")
	
	result, err := s.gitOps.Show(repoPath, revision)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitBranch(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	branchType := getString(arguments, "branch_type")
	if branchType == "" {
		branchType = "local"
	}
	contains := getString(arguments, "contains")
	notContains := getString(arguments, "not_contains")
	
	result, err := s.gitOps.Branch(repoPath, branchType, contains, notContains)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

// Helper functions for extracting values from arguments

func getString(args map[string]interface{}, key string) string {
	if val, ok := args[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getInt(args map[string]interface{}, key string, defaultVal int) int {
	if val, ok := args[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case json.Number:
			if i, err := v.Int64(); err == nil {
				return int(i)
			}
		}
	}
	return defaultVal
}

func getStringSlice(args map[string]interface{}, key string) []string {
	if val, ok := args[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(slice))
			for _, item := range slice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return []string{}
}

func getBool(args map[string]interface{}, key string, defaultVal bool) bool {
	if val, ok := args[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultVal
}

func (s *Server) handleGitRawCommand(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	command := getString(arguments, "command")
	
	result, err := s.gitOps.RawCommand(repoPath, command)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitInit(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := getString(arguments, "repo_path")
	bare := getBool(arguments, "bare", false)
	
	result, err := s.gitOps.Init(repoPath, bare)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitPush(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	remote := getString(arguments, "remote")
	refspec := getString(arguments, "refspec")
	tags := getBool(arguments, "tags", false)
	
	result, err := s.gitOps.Push(repoPath, remote, refspec, tags)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitListRepositories(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	searchPath := getString(arguments, "search_path")
	recursive := getBool(arguments, "recursive", false)
	
	repositories, err := s.gitOps.ListRepositories(searchPath, recursive)
	if err != nil {
		return nil, err
	}

	if len(repositories) == 0 {
		return []mcp.TextContent{{
			Type: "text",
			Text: "No Git repositories found",
		}}, nil
	}

	result := "Found Git repositories:\n"
	for _, repo := range repositories {
		result += fmt.Sprintf("- %s\n", repo)
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: strings.TrimSpace(result),
	}}, nil
}

func (s *Server) handleGitCreateTag(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	tagName := getString(arguments, "tag_name")
	message := getString(arguments, "message")
	annotated := getBool(arguments, "annotated", true)
	
	result, err := s.gitOps.CreateTag(repoPath, tagName, message, annotated)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitDeleteTag(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	tagName := getString(arguments, "tag_name")
	
	result, err := s.gitOps.DeleteTag(repoPath, tagName)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}

func (s *Server) handleGitListTags(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	pattern := getString(arguments, "pattern")
	
	tags, err := s.gitOps.ListTags(repoPath, pattern)
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return []mcp.TextContent{{
			Type: "text",
			Text: "No tags found",
		}}, nil
	}

	result := "Tags:\n"
	for _, tag := range tags {
		result += fmt.Sprintf("- %s\n", tag)
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: strings.TrimSpace(result),
	}}, nil
}

func (s *Server) handleGitPushTags(ctx context.Context, arguments map[string]interface{}) ([]mcp.TextContent, error) {
	repoPath := s.getRepoPath(getString(arguments, "repo_path"))
	remote := getString(arguments, "remote")
	tagName := getString(arguments, "tag_name")
	
	result, err := s.gitOps.PushTags(repoPath, remote, tagName)
	if err != nil {
		return nil, err
	}

	return []mcp.TextContent{{
		Type: "text",
		Text: result,
	}}, nil
}
