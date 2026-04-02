# yx-commands

阿里云云效代码仓库 yx-commands Skills，基于 [Agent Skills](https://agentskills.io/) 开放标准，让 Claude Code、Cursor 等 AI 编程工具可以自动完成代码提交、推送、创建合并请求等开发工作流。

## 功能特性

- 🤖 **智能提交** - AI 分析代码变更，自动生成 commit message
- 📤 **自动推送** - 推送当前分支到远程仓库
- 🔀 **创建 MR** - 自动创建合并请求，默认目标分支 develop
- 📊 **代码审查** - 查看分支间的代码差异统计
- ⚡ **一键工作流** - 提交 → 推送 → 创建 MR 一条龙完成

## 安装

```bash
npm install -g yunxiao-code
```

安装后自动部署 skill 到 `~/.claude/skills/`，Claude Code 和 Cursor 会自动加载。

支持平台：macOS (Intel/Apple Silicon)、Linux (x64)

## 快速开始

开发完成后，输入：

```
/yx-commands-commit-push-mr
```

AI 自动完成：提交 → 推送 → 创建 MR

## Skill 命令

| 命令 | 功能 | 使用场景 |
|------|------|----------|
| `/yx-commands-commit` | 智能提交 | AI 分析变更生成 commit message |
| `/yx-commands-push` | 推送代码 | 推送当前分支到远程 |
| `/yx-commands-mr` | 创建 MR | 创建合并请求（默认目标 develop） |
| `/yx-commands-review` | 代码审查 | 查看分支差异统计 |
| `/yx-commands-commit-push-mr` | 完整工作流 | 一键完成提交、推送、创建 MR |

## 典型工作流

### 一键完成（推荐）

开发完成后，直接输入：

```
/yx-commands-commit-push-mr
```

AI 会自动：
1. 分析代码变更，生成 commit message
2. 执行 `git add .` 和 `git commit`
3. 推送到远程仓库
4. 创建合并请求

### 分步执行

```
/yx-commands-commit     → AI 分析变更并提交
/yx-commands-push       → 推送到远程
/yx-commands-mr         → 创建合并请求
```

### 代码审查

查看当前分支与目标分支的差异：

```
/yx-commands-review
```

## 配置

### 初始化

在云效仓库目录下创建 `.yunxiao.yaml`：

```bash
yx-code init
```

填写 `organization_id` 和 `token`（云效控制台 → 个人设置 → 个人访问令牌）。

### 环境变量

```bash
export YUNXIAO_ORGANIZATION_ID="your-org-id"
export YUNXIAO_TOKEN="your-token"
```

### 配置优先级

1. 默认值 `openapi-rdc.aliyuncs.com`
2. `.yunxiao.yaml` 配置文件（从当前目录向上查找）
3. 环境变量：`YUNXIAO_ORGANIZATION_ID`、`YUNXIAO_TOKEN`
4. CLI 参数：`--org`、`--token`

## CLI 命令（底层实现）

Skill 底层通过 `yx-code` CLI 与云效 API 交互：

```bash
yx-code init              # 初始化配置
yx-code clone <url> -b <branch>  # 克隆仓库并创建分支
yx-code commit -m "msg"   # 提交代码
yx-code push              # 推送代码
yx-code mr -m "title"     # 创建合并请求
yx-code diff              # 查看代码差异
```

## 兼容性

| 工具 | 支持 | 说明 |
|------|------|------|
| [Claude Code](https://claude.ai/code) | ✅ | 原生支持 Agent Skills |
| [Cursor](https://cursor.com) | ✅ | 兼容 `~/.claude/skills/` 目录 |

## 相关链接

- [Agent Skills 标准](https://agentskills.io/)
- [阿里云云效官方文档](https://help.aliyun.com/zh/yunxiao/)
- [云效 OpenAPI 文档](https://help.aliyun.com/zh/yunxiao/developer-reference/)

## 从源码安装

```bash
git clone https://github.com/ZLStream/yunxiao-code.git
cd yunxiao-code
make build-local
sudo make install

# 手动链接 skills
for dir in claude-plugin/skills/*/; do
  ln -s "$(pwd)/$dir" ~/.claude/skills/"$(basename $dir)"
done
```

## License

MIT
