# go-mcp-git: Go 语言实现的 Git MCP 服务器

[English](README_EN.md) | 中文

## 概述

一个使用 Go 语言实现的模型上下文协议（Model Context Protocol）服务器，用于 Git 仓库交互和自动化。该服务器提供工具，让大型语言模型能够读取、搜索和操作 Git 仓库。

这是原始 Python mcp-server-git 的 Go 语言移植版本，提供了更好的性能和更简单的部署。

### 工具列表

1. `git_status` - 显示工作树状态
2. `git_diff_unstaged` - 显示工作目录中尚未暂存的更改
3. `git_diff_staged` - 显示已暂存待提交的更改
4. `git_diff` - 显示分支或提交之间的差异
5. `git_commit` - 将更改记录到仓库
6. `git_add` - 将文件内容添加到暂存区
7. `git_reset` - 取消暂存所有已暂存的更改
8. `git_log` - 显示提交日志，支持可选的日期过滤
9. `git_create_branch` - 创建新分支
10. `git_checkout` - 切换分支
11. `git_show` - 显示提交的内容
12. `git_branch` - 列出 Git 分支

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

## 配置

### 与 Claude Desktop 一起使用

在您的 `claude_desktop_config.json` 中添加以下配置：

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

## 许可证

本 MCP 服务器采用 Apache 2.0 许可证。
