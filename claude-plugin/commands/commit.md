---
name: commit
description: 分析代码变更并提交到云效仓库
---

分析当前的代码变更并提交。

**步骤：**

1. 一次性获取变更状态和差异概览：
```bash
git status && git diff --stat HEAD
```

2. 查看完整差异内容用于分析（包含已暂存和未暂存的所有变更）：
```bash
git diff HEAD
```

3. 分析变更后，根据变更内容生成合适的 commit message，格式遵循 Conventional Commits：
- `feat:` 新功能
- `fix:` 修复 bug
- `refactor:` 重构
- `docs:` 文档变更
- `style:` 代码格式
- `test:` 测试
- `chore:` 构建/工具

4. 执行提交：
```bash
yx-code commit -m "<生成的commit-message>"
```

如果用户提供了特定的 commit message，则使用用户提供的。