# go-mcp-git: Go 语言实现的 Git MCP 服务器

[English](README_EN.md) | 中文

## 概述

一个使用 Go 语言实现的模型上下文协议（Model Context Protocol）服务器，用于 Git 仓库交互和自动化。该服务器提供工具，让大型语言模型能够读取、搜索和操作 Git 仓库。

这是原始 Python mcp-server-git 的 Go 语言移植版本，提供了更好的性能和更简单的部署。

### 工具列表

#### 基础操作
1. `git_status` - 显示工作树状态
2. `git_init` - **新增** 初始化新的Git仓库
3. `git_add` - 将文件内容添加到暂存区
4. `git_commit` - 将更改记录到仓库
5. `git_reset` - 取消暂存所有已暂存的更改

#### 分支管理
6. `git_branch` - 列出 Git 分支
7. `git_create_branch` - 创建新分支
8. `git_checkout` - 切换分支

#### 差异和日志
9. `git_diff_unstaged` - 显示工作目录中尚未暂存的更改
10. `git_diff_staged` - 显示已暂存待提交的更改
11. `git_diff` - 显示分支或提交之间的差异
12. `git_log` - 显示提交日志，支持可选的日期过滤
13. `git_show` - 显示提交的内容

#### 远程操作
14. `git_push` - **新增** 推送更改到远程仓库
15. `git_list_repositories` - **新增** 列出目录中的Git仓库

#### 标签管理
16. `git_create_tag` - **新增** 创建Git标签（支持轻量级和注释标签）
17. `git_delete_tag` - **新增** 删除Git标签
18. `git_list_tags` - **新增** 列出Git标签（支持模式过滤）
19. `git_push_tags` - **新增** 推送标签到远程仓库

#### 高级功能
20. `git_raw_command` - **新增** 直接执行原始Git命令（绕过shell包装问题）

## 安装

### 使用 Go 安装

```bash
go install github.com/pengcunfu/go-mcp-git@latest
```

### 从源码构建

```bash
git clone https://github.com/pengcunfu/go-mcp-git.git
cd go-mcp-git
go build -o go-mcp-git ./cmd/server
```

## 使用方法

```bash
go-mcp-git --repository /path/to/git/repo
```

### git_raw_command 工具特别说明

`git_raw_command` 工具是专门为解决在某些环境（如Windsurf IDE）中Git命令被shell包装导致的引号转义问题而设计的。

**问题场景：**
当执行包含引号的Git命令时，例如：
```bash
git tag -a v0.0.1 -m "发布v0.0.1版本 - 初始MCP Git服务器实现"
```

在PowerShell环境中可能被错误包装为：
```powershell
Invoke-Expression "git tag -a v0.0.1 -m "发布v0.0.1版本 - 初始MCP Git服务器实现""
```

这会导致引号转义错误和命令执行失败。

**解决方案：**
使用 `git_raw_command` 工具可以直接执行原始Git命令，绕过shell包装：

```json
{
  "repo_path": "/path/to/repository",
  "command": "git tag -a v0.0.1 -m \"发布v0.0.1版本 - 初始MCP Git服务器实现\""
}
```

**支持的命令示例：**
- `git tag -a v1.0.0 -m "Release version 1.0.0"`
- `git commit --amend -m "Updated commit message"`
- `git push origin --tags`
- `git config user.name "Your Name"`

### 新增工具使用示例

#### 标签管理
```json
// 创建注释标签
{
  "repo_path": "/path/to/repository",
  "tag_name": "v1.0.0",
  "message": "Release version 1.0.0",
  "annotated": true
}

// 列出所有标签
{
  "repo_path": "/path/to/repository"
}

// 推送特定标签
{
  "repo_path": "/path/to/repository",
  "remote": "origin",
  "tag_name": "v1.0.0"
}

// 推送所有标签
{
  "repo_path": "/path/to/repository",
  "remote": "origin"
}
```

#### 仓库初始化和推送
```json
// 初始化新仓库
{
  "repo_path": "/path/to/new/repository",
  "bare": false
}

// 推送到远程仓库
{
  "repo_path": "/path/to/repository",
  "remote": "origin",
  "tags": true
}
```

#### 仓库发现
```json
// 递归搜索Git仓库
{
  "search_path": "/path/to/search",
  "recursive": true
}
```

## 配置

### 与 Claude Desktop 一起使用

在您的 `claude_desktop_config.json` 中添加以下配置：

```json
{
  "mcpServers": {
    "go-mcp-git": {
      "command": "D:\\Tools\\MCP\\go-mcp-git\\go-mcp-git.exe",
      "args": [
        "--repository",
        "https://github.com/pengcunfu/go-mcp-git.git"
      ]
    }
  }
}
```

## 许可证

本 MCP 服务器采用 Apache 2.0 许可证。
