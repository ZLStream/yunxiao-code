---
name: yx-commands-summary
description: 生成云效工作项总结，支持按周/月/季度/年或自定义时间范围，同时汇总本地云效仓库的 git 提交记录
---

生成指定时间范围内的工作总结（含需求、任务、缺陷 + git 提交记录），适合写周报、月报或绩效总结。

**AI 生成工作总结的完整流程：**

## 第一步：确定时间范围

根据用户请求确定 `START_DATE` 和 `END_DATE`（格式 `YYYY-MM-DD`），例如本月：
- START_DATE = 本月第一天
- END_DATE = 今天

---

## 第二步：获取云效工作项

```bash
yx-code summary --start <START_DATE> --end <END_DATE> --json
```

也可使用预设周期：
```bash
yx-code summary --period week --json    # 本周
yx-code summary --period month --json   # 本月
yx-code summary --period quarter --json # 本季度
yx-code summary --period year --json    # 本年
```

**参数说明：**

| 参数 | 说明 | 可选值 | 默认值 |
|------|------|--------|--------|
| `--period` / `-p` | 时间周期 | week/month/quarter/year | 交互提示 |
| `--start` | 开始日期 | YYYY-MM-DD 格式 | - |
| `--end` | 结束日期 | YYYY-MM-DD 格式 | 今天 |
| `--status` / `-s` | 状态过滤 | all（默认）/ undone（排除已完成）/ 中文状态名 | all |
| `--json` | JSON 格式输出 | - | false |

---

## 第三步：扫描本地云效 git 仓库并收集提交记录

**3.1 获取 git 用户名**

```bash
git config --global user.name
```

**3.2 扫描所有 remote 指向 codeup.aliyun.com 的本地仓库**

在常见代码目录中查找，深度限制 5 层，并通过项目根目录 mtime 预过滤，跳过时间范围内未活动的仓库：

```bash
for gitdir in $(find ~/IdeaProjects ~/Projects ~/code ~/workspace ~/repos ~/src ~/dev ~/work ~/Documents/code 2>/dev/null -maxdepth 5 -name ".git" -type d); do
  repo=$(dirname "$gitdir")
  # 1. mtime 预过滤：项目根目录在时间范围内未被修改则跳过（纯文件系统操作，无需启动 git）
  if ! find "$repo" -maxdepth 0 -newermt "<START_DATE>" 2>/dev/null | grep -q .; then
    continue
  fi
  # 2. 检查 remote 是否为云效仓库
  if git -C "$repo" remote -v 2>/dev/null | grep -q "codeup.aliyun.com"; then
    echo "$repo"
  fi
done
```

> - 项目根目录 mtime 在任何文件编辑、IDE 打开等操作时均会更新，比 `.git/` 内部文件更灵敏
> - mtime 检查为纯文件系统操作，优先于需要启动 git 进程的 remote URL 检查
> - 某个目录不存在时 `2>/dev/null` 静默跳过，不影响其他目录扫描

**3.3 对每个仓库查询时间范围内的 commit**

```bash
git -C <repo_path> log \
  --after="<START_DATE> 00:00:00" \
  --before="<END_DATE> 23:59:59" \
  --author="<用户名>" \
  --format="%h|%s|%ai" \
  --no-merges
```

从 remote URL 提取仓库名作为展示标识：
```bash
git -C <repo_path> remote get-url origin
# git@codeup.aliyun.com:org/repo-name.git → 取最后一段去掉 .git → repo-name
```

---

## 第四步：整合展示

将云效工作项和 git 提交记录合并为以下结构：

```
## 工作总结（2024-03-01 ~ 2024-03-31）

### 云效工作项（N 项）

**总计：N 项**（需求 X 项 / 任务 Y 项 / 缺陷 Z 项）

#### 需求（X 项）
| 编号 | 标题 | 状态 | 优先级 | 所在项目 | 创建日期 |
|------|------|------|--------|----------|----------|
| XJHS-1898 | 太保：兼容团险网点人员 | 已上线 | 中 | 芯计划书 | 2024-03-10 |

#### 任务（Y 项）
...

#### 缺陷（Z 项）
...

#### 按项目汇总
| 项目 | 工作项数 |
|------|----------|
| 展业赋能 | 7 |

---

### Git 提交记录（共 N 次提交，涉及 R 个仓库）

#### 仓库：repo-a（M 次提交）
| 提交哈希 | 提交信息 | 日期 |
|----------|----------|------|
| a1b2c3d | feat: 新增XX功能 | 2024-03-15 |

#### 仓库：repo-b（K 次提交）
| 提交哈希 | 提交信息 | 日期 |
|----------|----------|------|
| e4f5g6h | fix: 修复登录问题 | 2024-03-20 |

---

### 综合统计
| 维度 | 数量 |
|------|------|
| 云效工作项 | X 项 |
| Git 提交 | Y 次 |
| 涉及云效项目 | Z 个 |
| 涉及代码仓库 | R 个 |
```

**展示规则：**
- 云效工作项按 `category` 字段分组：`Req`=需求、`Task`=任务、`Bug`=缺陷，无数据的类别跳过
- 标题超过 30 字时截断并加"..."
- 若 `parent_id` 非空，在标题后附注"（子项）"
- 优先级为空时显示"-"
- Git 提交按仓库分组，同一仓库内按日期倒序
- 未找到任何 codeup 仓库时，跳过 Git 部分
- yx-code 未安装时，跳过云效工作项部分
- 如需进一步分析，可请求："根据以上数据，帮我写一份适合汇报给领导的月度总结"
