# 变更日志

本文档记录了go-mcp-git项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

## [v0.0.2] - 2025-11-13

### 新增功能
- **仓库管理**
  - `git_init` - 初始化新的Git仓库（支持bare仓库选项）
  - `git_list_repositories` - 列出目录中的Git仓库（支持递归搜索）

- **远程操作**
  - `git_push` - 推送更改到远程仓库
    - 支持指定remote名称（默认origin）
    - 支持自定义refspec
    - 支持同时推送标签

- **完整标签管理系统**
  - `git_create_tag` - 创建Git标签
    - 支持轻量级标签和注释标签
    - 可自定义标签消息
  - `git_delete_tag` - 删除本地Git标签
  - `git_list_tags` - 列出Git标签
    - 支持glob模式过滤（如 `v*`, `release-*`）
  - `git_push_tags` - 推送标签到远程仓库
    - 支持推送特定标签
    - 支持推送所有标签

- **高级功能**
  - `git_raw_command` - 直接执行原始Git命令
    - 解决PowerShell环境中的引号转义问题
    - 绕过shell包装，确保命令原样执行
    - 支持所有Git命令和参数

### 改进
- **文档优化**
  - 重新组织README工具列表，按功能分类
  - 添加详细的使用示例和JSON参数说明
  - 新增git_raw_command特别说明和问题解决方案
  - 添加标签管理、仓库初始化等新工具的使用示例

- **代码结构**
  - 添加`getBool`辅助函数支持布尔参数解析
  - 改进错误处理和用户反馈
  - 增强参数验证和默认值处理

### 修复
- **CI/CD修复**
  - 移除CI工作流中的Docker构建步骤（项目无Dockerfile）
  - 修复构建配置中的main.go路径问题（从cmd/server更新为根目录）
  - 更新Makefile和release.yml中的构建路径

### 技术改进
- 新增必要的Go包导入（os、path/filepath、config等）
- 完善Git操作的错误处理和状态反馈
- 优化仓库路径解析和验证逻辑

## [v0.0.1] - 2025-11-13

### 新增功能
- **基础Git操作**
  - `git_status` - 显示工作树状态
  - `git_add` - 将文件内容添加到暂存区
  - `git_commit` - 将更改记录到仓库
  - `git_reset` - 取消暂存所有已暂存的更改

- **分支管理**
  - `git_branch` - 列出Git分支（本地、远程、全部）
  - `git_create_branch` - 创建新分支
  - `git_checkout` - 切换分支

- **差异和历史**
  - `git_diff_unstaged` - 显示工作目录中尚未暂存的更改
  - `git_diff_staged` - 显示已暂存待提交的更改  
  - `git_diff` - 显示分支或提交之间的差异
  - `git_log` - 显示提交日志（支持日期过滤和数量限制）
  - `git_show` - 显示提交的内容

### 技术实现
- **MCP服务器架构**
  - 基于Go语言实现的模型上下文协议服务器
  - 使用go-git库进行Git操作
  - JSON Schema验证输入参数
  - 结构化的工具注册和处理系统

- **项目结构**
  - `internal/server` - MCP服务器实现
  - `internal/git` - Git操作封装
  - `internal/mcp` - MCP协议实现
  - `main.go` - 服务器入口点

- **构建和部署**
  - Makefile支持多平台构建
  - GitHub Actions CI/CD流水线
  - 自动化测试和代码检查
  - 多架构二进制文件发布

### 文档
- 完整的README.md文档
- 工具使用说明和示例
- Claude Desktop集成配置指南
- 多语言文档支持（中文/英文）

---

## 版本说明

### 版本号规则
- **主版本号**: 不兼容的API修改
- **次版本号**: 向下兼容的功能性新增
- **修订号**: 向下兼容的问题修正

### 标签格式
- 发布版本: `v0.0.1`, `v0.0.2`
- 预发布版本: `v0.0.2-alpha.1`, `v0.0.2-beta.1`
- 开发版本: `v0.0.2-dev`

[未发布]: https://github.com/pengcunfu/go-mcp-git/compare/v0.0.2...HEAD
[v0.0.2]: https://github.com/pengcunfu/go-mcp-git/compare/v0.0.1...v0.0.2
[v0.0.1]: https://github.com/pengcunfu/go-mcp-git/releases/tag/v0.0.1