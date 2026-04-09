---
name: yx-commands
description: 阿里云云效(codeup)代码仓库工作流自动化。功能包括：智能提交代码、推送、创建MR(合并请求)、代码审查、分支差异对比、工作项查询、工作总结生成。触发场景：用户提到云效、codeup、阿里云代码仓库、yx-code工具；git地址以git@codeup.aliyun.com开头；需要提交代码并创建MR；说"提交合并"触发完整工作流；说"工作总结""周报""月报"触发总结功能。注意：github相关操作或纯git命令不要使用此技能。
tools: Bash, Read
---

# 云效代码仓库工作流

帮助用户在阿里云云效代码仓库平台进行开发工作流自动化。

> **认证信息由 CLI 自动获取**，无需手动配置 organization_id 和 token。

## 何时使用

当用户涉及以下场景时，主动使用此 skill：

1. **云效平台相关** - 提到云效、阿里云代码仓库、codeup、阿里云 DevOps、合并请求、MR、代码审查
2. **Git 操作** - 提交、推送
3. **yx-code 工具** - 提到云效 CLI
4. **复合工作流** - 用户说"提交合并"、"提交并合并"、"提交MR"时，按顺序执行：
   - 先执行 `/yx-commands-commit` 提交代码
   - 再执行 `/yx-commands-push` 推送到远程
   - 然后执行 `/yx-commands-mr` 创建合并请求

## 可用命令

| 命令 | 功能 | 典型用法 |
|------|------|----------|
| `/yx-commands-commit` | 智能提交 | 分析变更生成 commit |
| `/yx-commands-push` | 推送代码 | 推送到远程 |
| `/yx-commands-mr` | 创建 MR | 提交合并请求（默认 develop） |
| `/yx-commands-review` | 代码审查 | 查看分支差异 |
| `/yx-commands-commit-push-mr` | 完整工作流 | 一键完成提交、推送、创建MR |
| `/yx-commands-workitems` | 查询工作项 | 列出我负责/参与的工作项 |
| `/yx-commands-summary` | 工作总结 | 生成周报/月报/季报/年报 |

## 典型工作流

### 一键完成（推荐）
```
1. 编写代码...
2. /yx-commands-commit-push-mr  → 一键完成提交、推送、创建MR
```

### 分步执行
```
1. 编写代码...
2. /yx-commands-commit   → 智能生成 commit message 并提交
3. /yx-commands-push     → 推送到远程
4. /yx-commands-mr       → 创建合并请求
```

### 代码审查
```
/yx-commands-review -t develop -s feature-branch
```

### 工作总结
```
/yx-commands-summary          → 交互式选择时间范围
/yx-commands-summary --period month  → 本月总结
/yx-commands-workitems        → 查看当前进行中工作项
```

## 前置条件

- 已安装 `yx-code` CLI 工具
- `yx-code` 在系统 PATH 中可用
- 已运行 `yx-code init` 完成配置（`~/.yunxiao/config.yaml`）