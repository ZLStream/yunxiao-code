# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

阿里云云效 CLI 工具（yx-code），用于在终端中完成云效代码仓库的日常开发工作流。Go 语言编写，使用 cobra 框架构建 CLI，通过 npm 分发跨平台二进制文件。

## 常用命令

```bash
# 本地构建
make build-local        # 构建当前平台二进制文件到 bin/
make build              # 构建所有平台（darwin-amd64/arm64, linux-amd64）

# 安装到系统
sudo make install       # 安装到 /usr/local/bin/yx-code

# 发布到 npm（需先构建）
make publish

# 直接运行（开发调试）
go run ./cmd/yx-code <command>

# 清理
make clean
```

## 代码架构

```
cmd/yx-code/main.go     # CLI 入口，所有子命令定义（commit/push/mr/diff/clone/init）
internal/
  api/client.go         # 云效 OpenAPI 客户端：GetProjectId、CreateMergeRequest、GetCompare
  config/config.go      # 配置加载：文件(.yunxiao.yaml) → 环境变量 → CLI参数，优先级递增
  git/git.go            # Git 操作封装：GetCurrentBranch、GetRepoName、AddAndCommit、Clone、Push
  review/reviewer.go    # 代码差异格式化输出
claude-plugin/skills/   # Agent Skills 插件（yx-commands 系列斜杠命令）
```

## 核心逻辑

### 配置优先级（从低到高）
1. 默认值 `openapi-rdc.aliyuncs.com`
2. `.yunxiao.yaml` 配置文件（从当前目录向上查找）
3. 环境变量：`YUNXIAO_DOMAIN`、`YUNXIAO_ORGANIZATION_ID`、`YUNXIAO_TOKEN`
4. CLI 参数：`--domain`、`--org`、`--token`

### API 认证
请求头使用 `x-yunxiao-token` 传递个人访问令牌，API 格式：
```
https://{domain}/oapi/v1/codeup/organizations/{org_id}/repositories/{repo_id}/changeRequests
```

### 受保护分支
`main`、`master`、`develop` 分支禁止直接推送，需通过 Merge Request 提交代码。

## Go 依赖

- `github.com/spf13/cobra` - CLI 框架
- `gopkg.in/yaml.v3` - YAML 配置解析

## Skill 插件开发

Skills 位于 `claude-plugin/skills/` 目录，每个子目录是一个独立 skill，包含 `SKILL.md` 文件。npm 安装时通过 `scripts/postinstall.js` 自动将各 skill 链接到 `~/.claude/skills/`。

新增 skill 时：
1. 在 `claude-plugin/skills/` 下创建目录和 `SKILL.md`
2. 更新 `package.json` 的 `files` 字段确保包含该目录
3. 更新 `postinstall.js` 的可用命令提示