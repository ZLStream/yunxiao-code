# yx-code

阿里云云效 CLI 工具，支持 Git 提交、推送、克隆及创建合并请求，让开发者可以在终端中完成日常开发工作流，无需打开浏览器。

## 功能特性

- 🚀 **克隆仓库** - 支持从主分支自动创建新分支
- 📝 **提交代码** - 一键 add + commit
- 📤 **推送代码** - 自动推送当前分支到远程
- 🔀 **合并请求** - 快速创建 MR，支持指定目标分支
- 📊 **代码差异** - 查看分支间的代码差异统计
- 🤖 **Claude Code 集成** - 支持 Claude Code Skills/Commands 自动化工作流

## 安装

### npm 安装（推荐）

```bash
npm install -g yunxiao-code
```

支持平台：
- macOS (Intel/Apple Silicon)
- Linux (x64)

### 从源码编译

```bash
git clone https://github.com/ZLStream/yunxiao-code.git
cd yunxiao-code
make build-local
sudo make install
```

或手动编译：

```bash
go build -o yx-code ./cmd/yx-code
sudo mv yx-code /usr/local/bin/
```

## 系统要求

- Go 1.26.1+（从源码编译时）
- Git CLI（运行时依赖）
- Node.js 14+（npm 安装时）

## 快速开始

### 1. 初始化配置

```bash
yx-code init
```

在当前目录生成 `.yunxiao.yaml` 配置文件：

```yaml
# 云效合并请求 CLI 配置
domain: openapi-rdc.aliyuncs.com
organization_id: ""
token: ""
```

填写 `organization_id` 和 `token`（个人访问令牌）。

> **获取 Token**: 云效控制台 → 个人设置 → 个人访问令牌 → 创建令牌

### 2. 配置项说明

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `domain` | 云效 API 域名 | `openapi-rdc.aliyuncs.com` |
| `organization_id` | 组织 ID（必填） | 无 |
| `token` | 个人访问令牌（必填） | 无 |

### 3. 环境变量配置

也可以通过环境变量配置：

```bash
export YUNXIAO_DOMAIN="openapi-rdc.aliyuncs.com"
export YUNXIAO_ORGANIZATION_ID="your-org-id"
export YUNXIAO_TOKEN="your-token"
```

### 4. 配置优先级（从低到高）

1. 默认值
2. `.yunxiao.yaml` 配置文件（从当前目录向上逐级查找）
3. 环境变量：`YUNXIAO_DOMAIN`、`YUNXIAO_ORGANIZATION_ID`、`YUNXIAO_TOKEN`
4. CLI 参数：`--domain`、`--org`、`--token`

## 命令参考

### yx-code init

初始化云效配置文件。

```bash
yx-code init
```

### yx-code clone

克隆云效代码仓库。

```bash
yx-code clone <git-url> [flags]
```

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `<git-url>` | | 是 | Git 仓库地址 |
| `--branch` | `-b` | 否 | 新分支名称，指定后会从主分支创建新分支 |
| `--path` | `-p` | 否 | 克隆目标路径，默认从 URL 提取仓库名 |

```bash
# 基本克隆
yx-code clone https://codeup.aliyun.com/your-org/your-repo.git

# 克隆并创建新分支
yx-code clone https://codeup.aliyun.com/your-org/your-repo.git -b feature/new-feature
```

### yx-code commit

提交代码变更。

```bash
yx-code commit -m "<message>"
```

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--message` | `-m` | 是 | 提交信息 |

```bash
yx-code commit -m "feat: 添加用户登录功能"
```

### yx-code push

推送当前分支到远程仓库。

```bash
yx-code push
```

### yx-code mr

创建合并请求。

```bash
yx-code mr -m "<title>" [flags]
```

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--message` | `-m` | 是 | MR 标题 |
| `--description` | `-d` | 否 | MR 描述 |
| `--target` | `-t` | 否 | 目标分支，默认: develop |

```bash
# 基本用法（目标分支默认 develop）
yx-code mr -m "添加用户登录功能"

# 指定目标分支
yx-code mr -m "添加用户登录功能" -t develop

# 带描述
yx-code mr -m "添加用户登录功能" -d "实现用户登录、注销功能"
```

### yx-code diff

查看分支间的代码差异统计。

```bash
yx-code diff [flags]
```

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--target` | `-t` | 否 | 目标分支，默认: develop |
| `--source` | `-s` | 否 | 源分支，默认: 当前分支 |

```bash
# 查看当前分支与 develop 的差异
yx-code diff

# 指定分支
yx-code diff -t main -s feature/new-feature
```

### 全局参数

所有子命令支持以下全局参数，用于覆盖配置文件：

```bash
yx-code <command> --domain <域名> --org <组织ID> --token <令牌>
```

| 参数 | 说明 |
|------|------|
| `--domain` | 云效 API 域名 |
| `--org` | 组织 ID |
| `--token` | 个人访问令牌 |

## 工作流示例

```bash
# 1. 初始化配置
yx-code init
# 编辑 .yunxiao.yaml 填写凭证

# 2. 克隆仓库并创建特性分支
yx-code clone https://codeup.aliyun.com/org/repo.git -b feature/login

# 3. 进入仓库目录
cd repo

# 4. 开发完成后提交
yx-code commit -m "feat: 实现登录功能"

# 5. 推送到远程
yx-code push

# 6. 创建合并请求
yx-code mr -m "添加用户登录功能"

# 7. 查看代码差异
yx-code diff
```

## Claude Code 集成

本工具提供 Claude Code Skills/Commands 支持，可在 Claude Code 中通过斜杠命令自动化云效工作流。

### 可用命令

| 命令 | 功能 |
|------|------|
| `/yx-commands:commit` | 分析代码变更并提交 |
| `/yx-commands:push` | 推送代码到远程 |
| `/yx-commands:mr` | 创建合并请求 |
| `/yx-commands:review` | 查看分支代码差异 |
| `/yx-commands:commit-push-mr` | 一键完成提交、推送、创建MR |

### 使用示例

在 Claude Code 中直接输入命令：

```
/yx-commands:commit-push-mr
```

Claude 会自动：
1. 分析代码变更
2. 生成 commit message
3. 执行 commit、push
4. 创建合并请求

### 安装

npm 安装后自动部署到 `~/.claude/commands/yx-commands/` 和 `~/.claude/skills/yx-commands/`。

## 相关链接

- [阿里云云效官方文档](https://help.aliyun.com/zh/yunxiao/)
- [云效 OpenAPI 文档](https://help.aliyun.com/zh/yunxiao/developer-reference/)

## License

MIT