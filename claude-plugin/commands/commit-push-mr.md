---
name: commit-push-mr
description: 一键完成提交、推送、创建合并请求（完整工作流）
---

依次执行 commit、push、mr，完成从代码变更到合并请求的完整工作流。

**步骤：**

1. 查看当前变更状态：
```bash
git status && git diff --stat HEAD
```

2. 查看完整差异内容：
```bash
git diff HEAD
```

3. 分析变更后生成 commit message，执行提交：
```bash
yx-code commit -m "<生成的commit-message>"
```

4. 推送到远程：
```bash
yx-code push
```

5. 创建合并请求：
```bash
yx-code mr -m "<mr-title>"
```

**注意：**
- 如果用户提供了特定的 commit message，则使用用户提供的
- MR 标题默认使用 commit message
- 目标分支默认为 develop，可通过参数指定