# yx-code CLI 文档

## 概述

`yx-code` 是阿里云云效 DevOps 平台的命令行工具，支持 Git 提交、推送、克隆及创建合并请求等操作，让开发者可以在终端中完成日常开发工作流，无需打开浏览器。

## 安装

### npm 安装（推荐）

```bash
npm install -g yx-code
```

### 从源码编译

```bash
git clone <repo-url>
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

## 配置

### 初始化配置文件

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

### 配置项说明

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `domain` | 云效 API 域名 | `openapi-rdc.aliyuncs.com` |
| `organization_id` | 组织 ID（必填） | 无 |
| `token` | 个人访问令牌（必填） | 无 |

### 环境变量配置

```bash
export YUNXIAO_DOMAIN="openapi-rdc.aliyuncs.com"
export YUNXIAO_ORGANIZATION_ID="your-org-id"
export YUNXIAO_TOKEN="your-token"
```

### 配置优先级（从低到高）

1. 默认值
2. `.yunxiao.yaml` 配置文件（从当前目录向上逐级查找）
3. 环境变量：`YUNXIAO_DOMAIN`、`YUNXIAO_ORGANIZATION_ID`、`YUNXIAO_TOKEN`
4. CLI 参数：`--domain`、`--org`、`--token`

---

## 命令参考

### yx-code init

初始化云效配置文件。

```bash
yx-code init
```

**命令参数**

无

**全局参数**

无

**示例**

```bash
yx-code init
```

**行为说明**

在当前目录创建 `.yunxiao.yaml` 模板文件，需手动填写 `organization_id` 和 `token`。

---

### yx-code clone

克隆云效代码仓库。

```bash
yx-code clone <git-url> [flags]
```

**命令参数**

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `<git-url>` | | 是 | Git 仓库地址 |
| `--branch` | `-b` | 否 | 新分支名称，指定后会从主分支创建新分支 |
| `--path` | `-p` | 否 | 克隆目标路径，默认从 URL 提取仓库名 |

**全局参数**

无

**示例**

```bash
# 基本克隆
yx-code clone https://codeup.aliyun.com/your-org/your-repo.git

# 克隆并创建新分支
yx-code clone https://codeup.aliyun.com/your-org/your-repo.git -b feature/new-feature

# 指定目标路径
yx-code clone https://codeup.aliyun.com/your-org/your-repo.git -p /path/to/destination
```

**行为说明**

- 不指定 `-b` 时，仅执行 `git clone`
- 指定 `-b` 时，自动检测主分支（main/master），从主分支创建新分支

---

### yx-code commit

提交代码变更。

```bash
yx-code commit -m "<message>"
```

**命令参数**

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--message` | `-m` | 是 | 提交信息 |

**全局参数**

无

**示例**

```bash
yx-code commit -m "feat: 添加用户登录功能"
```

**行为说明**

执行 `git add .` 后提交所有变更。

---

### yx-code push

推送当前分支到远程仓库。

```bash
yx-code push
```

**命令参数**

无

**全局参数**

无

**示例**

```bash
yx-code push
```

**行为说明**

- 自动获取当前分支名
- 推送到 `origin` 远程

---

### yx-code mr

创建合并请求。

```bash
yx-code mr -m "<title>" [flags]
```

**命令参数**

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--message` | `-m` | 是 | MR 标题 |
| `--description` | `-d` | 否 | MR 描述 |
| `--target` | `-t` | 否 | 目标分支，默认: `develop` |

**全局参数**

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--domain` | | 否 | 云效 API 域名，默认从配置文件读取 |
| `--org` | | 否 | 组织 ID，默认从配置文件读取 |
| `--token` | | 否 | 个人访问令牌，默认从配置文件读取 |

**示例**

```bash
# 基本用法（目标分支默认 develop）
yx-code mr -m "添加用户登录功能"

# 指定目标分支
yx-code mr -m "添加用户登录功能" -t main

# 带描述
yx-code mr -m "添加用户登录功能" -d "实现用户登录、注销功能"

# 使用全局参数覆盖配置
yx-code mr -m "标题" --org your-org-id --token your-token
```

**行为说明**

1. 获取当前分支作为源分支
2. 目标分支默认为 `develop`，可通过 `-t` 参数指定
3. 调用云效 API 创建合并请求
4. 创建成功后询问是否查看代码差异

---

### yx-code diff

查看分支间的代码差异统计。

```bash
yx-code diff [flags]
```

**命令参数**

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--target` | `-t` | 否 | 目标分支，默认: `develop` |
| `--source` | `-s` | 否 | 源分支，默认: 当前分支 |

**全局参数**

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--domain` | | 否 | 云效 API 域名，默认从配置文件读取 |
| `--org` | | 否 | 组织 ID，默认从配置文件读取 |
| `--token` | | 否 | 个人访问令牌，默认从配置文件读取 |

**示例**

```bash
# 查看当前分支与 develop 的差异
yx-code diff

# 指定分支
yx-code diff -t main -s feature/new-feature

# 使用全局参数覆盖配置
yx-code diff --org your-org-id --token your-token
```

**输出说明**

- 变更统计：文件数、新增行数、删除行数
- 每个文件的 diff 详情

---

### yx-code workitems

列出当前用户负责的工作项（需求、任务、缺陷）。

```bash
yx-code workitems [flags]
```

**命令参数**

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--status` | `-s` | 否 | 状态过滤：`undone`（排除已完成）/ `all` / 中文状态名（如：开发中），默认 `undone` |
| `--project` | `-p` | 否 | 指定项目标识符或名称，默认查询所有项目 |
| `--json` | | 否 | 以 JSON 格式输出 |

**全局参数**

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--domain` | | 否 | 云效 API 域名，默认从配置文件读取 |
| `--org` | | 否 | 组织 ID，默认从配置文件读取 |
| `--token` | | 否 | 个人访问令牌，默认从配置文件读取 |

**示例**

```bash
# 列出我负责的未完成工作项（默认）
yx-code workitems

# 列出所有状态的工作项
yx-code workitems --status all

# 列出指定状态的工作项
yx-code workitems --status 开发中

# 指定项目
yx-code workitems --project my-project

# JSON 格式输出（供 AI 整理展示）
yx-code workitems --json
```

**行为说明**

- 自动获取当前登录用户 ID，查询其负责的工作项
- 按需求、任务、缺陷三类分组展示
- 用户 ID 首次运行时通过 token 自动获取并缓存

---

### yx-code summary

生成指定时间范围内的工作项总结（需求、任务、缺陷），适合撰写周报、月报或绩效总结。

```bash
yx-code summary [flags]
```

**命令参数**

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--period` | `-p` | 否 | 时间周期：`week` / `month` / `quarter` / `year`，与 `--start`/`--end` 二选一 |
| `--start` | | 否 | 开始日期，格式：`YYYY-MM-DD` |
| `--end` | | 否 | 结束日期，格式：`YYYY-MM-DD`，默认今天 |
| `--status` | `-s` | 否 | 状态过滤：`all`（默认）/ `undone`（排除已完成）/ 中文状态名 |
| `--json` | | 否 | 以 JSON 格式输出（AI 友好） |

**全局参数**

| 参数 | 简写 | 必填 | 说明 |
|------|------|------|------|
| `--domain` | | 否 | 云效 API 域名，默认从配置文件读取 |
| `--org` | | 否 | 组织 ID，默认从配置文件读取 |
| `--token` | | 否 | 个人访问令牌，默认从配置文件读取 |

**示例**

```bash
# 交互式选择时间范围
yx-code summary

# 本周总结
yx-code summary --period week

# 本月总结
yx-code summary --period month

# 本季度总结
yx-code summary --period quarter

# 自定义时间范围
yx-code summary --start 2024-01-01 --end 2024-03-31

# 只统计未完成工作项
yx-code summary --period month --status undone

# JSON 格式输出（供 AI 整理展示）
yx-code summary --period month --json
```

**行为说明**

- 未指定时间参数时，交互式提示选择时间范围（本周/本月/本季度/本年/自定义）
- 时间范围按工作项**更新时间**筛选，能捕获到在该时间段内有实际进展的工作项（而非仅限于该时间段内创建的）
- 按需求、任务、缺陷三类分组展示，并附按项目汇总
- 用户 ID 首次运行时通过 token 自动获取并缓存

---

## 全局参数

所有子命令支持以下全局参数，用于覆盖配置文件：

```bash
yx-code <command> --domain <域名> --org <组织ID> --token <令牌>
```

| 参数 | 说明 |
|------|------|
| `--domain` | 云效 API 域名 |
| `--org` | 组织 ID |
| `--token` | 个人访问令牌 |

---

## 工作流示例

### 完整开发流程

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

# 7. 查看代码差异（可选，创建 MR 后会提示）
# 或单独运行
yx-code diff
```

### 查看工作项与生成总结

```bash
# 查看当前负责的未完成工作项
yx-code workitems

# 查看指定状态的工作项
yx-code workitems --status 开发中

# 生成本月工作总结
yx-code summary --period month

# 生成本周工作总结（JSON 格式，供 AI 进一步整理）
yx-code summary --period week --json

# 生成自定义时间段总结
yx-code summary --start 2024-03-01 --end 2024-03-31 --json
```

---

## 退出码

| 退出码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1 | 命令执行失败 |

---

## 故障排除

### 常见错误

**配置文件未找到**

```
读取配置失败: 配置文件不存在
```

解决方案：运行 `yx-code init` 初始化配置文件。

**令牌无效**

```
API 错误: Unauthorized
```

解决方案：检查 `.yunxiao.yaml` 中的 `token` 是否正确，或使用 `--token` 参数。

**分支不存在**

```
获取代码差异失败: branch not found
```

解决方案：确保目标分支和源分支名称正确，检查是否已推送到远程。

---

## 相关链接

- [阿里云云效官方文档](https://help.aliyun.com/zh/yunxiao/)
- [云效 OpenAPI 文档](https://help.aliyun.com/zh/yunxiao/developer-reference/)