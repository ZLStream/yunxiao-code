---
name: mr
description: 创建云效合并请求(Merge Request)
---

创建合并请求到目标分支。

**重要：默认目标分支为 develop，执行前需用户确认。**

## 前置条件

- 当前分支已推送到远程（如未推送，先执行推送）
- 工作区干净（无未提交变更）

> **注意**：认证信息（organization_id、token）由 CLI 自动获取，无需手动配置。

## 步骤

1. **获取分支状态和提交信息**：
```bash
git branch --show-current && git log --oneline -3
```

如果分支未推送到远程，先推送：
```bash
git push -u origin HEAD
```

2. **确定 MR 信息**：
   - 标题：优先使用用户提供的标题；如果用户未提供，从最近的 commit message 自动推导
   - 目标分支：**默认为 develop**，除非用户明确指定其他分支

3. **询问用户确认**：
使用 AskUserQuestion 询问用户：
- 问题："即将创建合并请求到分支 `<target_branch>`，是否继续？"
- 选项：Yes / No

4. **用户选择 Yes 后创建 MR**：
```bash
yx-code mr -m "<title>" [-d "<description>"] [-t "<target>"]
```

## 命令参数

```
yx-code mr [flags]

Flags:
  -m, --message string     MR 标题（必填）
  -d, --description string MR 描述（可选）
  -t, --target string      目标分支（默认: develop）
  -h, --help               帮助信息
```

## 示例

创建 MR 到默认分支（develop）：
```bash
yx-code mr -m "feat: 添加用户登录功能"
```

指定目标分支：
```bash
yx-code mr -m "feat: 添加用户登录功能" -t main
```

带描述创建：
```bash
yx-code mr -m "feat: 添加用户登录功能" -d "实现了登录、登出功能，包含单元测试" -t develop
```

## 创建后

创建成功后会返回：
- 合并请求 ID
- 合并请求 URL 链接

用户可在浏览器中访问链接进行代码审查和合并操作。