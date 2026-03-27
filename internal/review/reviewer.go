package review

import (
	"fmt"
	"strings"
)

// DiffFile 表示单个文件的 diff
type DiffFile struct {
	Diff        string
	NewPath     string
	OldPath     string
	NewFile     bool
	DeletedFile bool
	AddLines    int
	DelLines    int
}

// Reviewer 代码差异查看器
type Reviewer struct {
	diffs []DiffFile
}

// NewReviewer 创建查看器
func NewReviewer(diffs []DiffFile) *Reviewer {
	return &Reviewer{diffs: diffs}
}

// GetStats 获取变更统计
func (r *Reviewer) GetStats() (files, addLines, delLines int) {
	files = len(r.diffs)
	for _, d := range r.diffs {
		// 如果 API 已提供行数，直接使用
		if d.AddLines > 0 || d.DelLines > 0 {
			addLines += d.AddLines
			delLines += d.DelLines
			continue
		}
		// 否则从 diff 内容解析
		for _, line := range strings.Split(d.Diff, "\n") {
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				addLines++
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				delLines++
			}
		}
	}
	return
}

// GetDiffResult 获取格式化的差异结果
func (r *Reviewer) GetDiffResult() (string, error) {
	return r.FormatDiff(), nil
}

// FormatDiff 格式化 diff 输出
func (r *Reviewer) FormatDiff() string {
	var buf strings.Builder

	for _, d := range r.diffs {
		if d.NewFile {
			buf.WriteString(fmt.Sprintf("+++ 新文件: %s (+%d 行)\n", d.NewPath, d.AddLines))
		} else if d.DeletedFile {
			buf.WriteString(fmt.Sprintf("--- 删除文件: %s (-%d 行)\n", d.OldPath, d.DelLines))
		} else {
			buf.WriteString(fmt.Sprintf("修改文件: %s (+%d/-%d 行)\n", d.NewPath, d.AddLines, d.DelLines))
		}
		buf.WriteString("```diff\n")
		buf.WriteString(d.Diff)
		if !strings.HasSuffix(d.Diff, "\n") {
			buf.WriteString("\n")
		}
		buf.WriteString("```\n\n")
	}

	return buf.String()
}