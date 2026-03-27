---
name: push
description: 推送代码到云效远程仓库
---

推送代码到远程仓库。

**重要：执行前需用户确认。**

## 步骤

1. **获取当前分支信息**：
```bash
git branch --show-current
```

2. **询问用户确认**：
使用 AskUserQuestion 询问用户：
- 问题："即将推送代码到远程仓库，是否继续？"
- 选项：Yes / No

3. **用户选择 Yes 后执行推送**：
```bash
yx-code push
```

4. **如果推送失败且提示无上游分支**，自动设置上游并重试：
```bash
yx-code push -u origin HEAD
```