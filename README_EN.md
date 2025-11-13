# go-mcp-git: A Git MCP Server in Go

English | [中文](README.md)

## Overview

A Model Context Protocol server for Git repository interaction and automation, implemented in Go. This server provides tools to read, search, and manipulate Git repositories via Large Language Models.

This is a Go port of the original Python mcp-server-git, offering improved performance and easier deployment.

### Tools

1. `git_status` - Shows the working tree status
2. `git_diff_unstaged` - Shows changes in working directory not yet staged
3. `git_diff_staged` - Shows changes that are staged for commit
4. `git_diff` - Shows differences between branches or commits
5. `git_commit` - Records changes to the repository
6. `git_add` - Adds file contents to the staging area
7. `git_reset` - Unstages all staged changes
8. `git_log` - Shows the commit logs with optional date filtering
9. `git_create_branch` - Creates a new branch
10. `git_checkout` - Switches branches
11. `git_show` - Shows the contents of a commit
12. `git_branch` - List Git branches

## Installation

### Using Go

```bash
go install github.com/pengcunfu/go-mcp-git@latest
```

### From Source

```bash
git clone https://github.com/pengcunfu/go-mcp-git.git
cd go-mcp-git
go build -o go-mcp-git ./cmd/server
```

## Usage

```bash
go-mcp-git --repository /path/to/git/repo
```

## Configuration

### Usage with Claude Desktop

Add this to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "git": {
      "command": "go-mcp-git",
      "args": ["--repository", "path/to/git/repo"]
    }
  }
}
```

## License

This MCP server is licensed under the Apache 2.0 License.