package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"yx-code/internal/api"
	"yx-code/internal/config"
)

var categoryNames = map[string]string{
	"Req":  "需求",
	"Task": "任务",
	"Bug":  "缺陷",
}

func categoryDisplayName(categoryId string) string {
	if name, ok := categoryNames[categoryId]; ok {
		return name
	}
	return categoryId
}

// extractPriority 从自定义字段中提取优先级显示名
func extractPriority(fields []api.CustomFieldValue) string {
	for _, f := range fields {
		if f.FieldName == "优先级" || strings.EqualFold(f.FieldName, "priority") {
			if len(f.Values) > 0 {
				return f.Values[0].DisplayValue
			}
		}
	}
	return ""
}

func newWorkitemsCmd() *cobra.Command {
	var (
		statusFilter string
		projectId    string
		outputJSON   bool
	)

	cmd := &cobra.Command{
		Use:   "workitems",
		Short: "列出工作项",
		Long:  "列出我负责的工作项（需求、任务、缺陷）",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadAndOverrideConfig(cmd)
			if err != nil {
				return err
			}

			if err := resolveUserId(cfg); err != nil {
				return err
			}

			projects, err := fetchUserProjects(cfg)
			if err != nil {
				return err
			}

			if len(projects) == 0 {
				fmt.Println("未找到相关项目")
				return nil
			}

			if projectId != "" {
				var filtered []api.Project
				for _, p := range projects {
					if p.Id == projectId || strings.EqualFold(p.Name, projectId) {
						filtered = append(filtered, p)
					}
				}
				if len(filtered) == 0 {
					return fmt.Errorf("未找到项目: %s", projectId)
				}
				projects = filtered
			}

			conditions := buildWorkitemConditions(cfg.UserId, "")

			var allItems []api.Workitem
			for _, proj := range projects {
				items, err := api.SearchWorkitems(cfg, proj.Id, "Req,Task,Bug", conditions)
				if err != nil {
					fmt.Printf("获取项目 %s 工作项失败: %v\n", proj.Name, err)
					continue
				}
				allItems = append(allItems, items...)
			}

			allItems = filterByStatus(allItems, statusFilter)

			if outputJSON {
				data, _ := json.MarshalIndent(allItems, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			printWorkitemsTable(allItems)
			return nil
		},
	}

	cmd.Flags().StringVarP(&statusFilter, "status", "s", "undone", "状态过滤: undone（默认，排除已完成）/ all / 中文状态名（如：开发中）")
	cmd.Flags().StringVarP(&projectId, "project", "p", "", "指定项目标识符或名称")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "以 JSON 格式输出")

	return cmd
}

// fetchUserProjects 获取当前用户参与的项目
func fetchUserProjects(cfg *config.Config) ([]api.Project, error) {
	cond := fmt.Sprintf(`{"conditionGroups":[[{"className":"user","fieldIdentifier":"users","format":"multiList","operator":"CONTAINS","value":["%s"]}]]}`, cfg.UserId)
	return api.SearchProjects(cfg, cond)
}

// buildWorkitemConditions 构建工作项查询条件 JSON（始终按 assignedTo 过滤）
func buildWorkitemConditions(userId, dateRange string) string {
	var conditions []string

	conditions = append(conditions, fmt.Sprintf(
		`{"fieldIdentifier":"assignedTo","operator":"CONTAINS","value":["%s"],"className":"user","format":"list"}`, userId))

	if dateRange != "" {
		conditions = append(conditions, fmt.Sprintf(
			`{"fieldIdentifier":"gmtCreate","operator":"BETWEEN","value":[%s],"className":"dateTime","format":"input"}`, dateRange))
	}

	condJSON := strings.Join(conditions, ",")
	return fmt.Sprintf(`{"conditionGroups":[[%s]]}`, condJSON)
}

// filterByStatus 客户端按状态名过滤工作项
func filterByStatus(items []api.Workitem, status string) []api.Workitem {
	switch status {
	case "all", "":
		return items
	case "undone":
		var filtered []api.Workitem
		for _, item := range items {
			if !isTerminalStatus(item.Status) {
				filtered = append(filtered, item)
			}
		}
		return filtered
	default:
		var filtered []api.Workitem
		for _, item := range items {
			if item.Status.Name == status || item.Status.DisplayName == status {
				filtered = append(filtered, item)
			}
		}
		return filtered
	}
}

// isTerminalStatus 判断状态是否为终态（已完成/已取消/已上线等）
func isTerminalStatus(s api.WorkitemStatus) bool {
	switch s.NameEn {
	case "Done", "Canceled", "Closed", "Rejected":
		return true
	}
	if strings.HasPrefix(s.Name, "已") || strings.HasPrefix(s.DisplayName, "已") {
		return true
	}
	return false
}

// printWorkitemsTable 以表格形式打印工作项完整信息
func printWorkitemsTable(items []api.Workitem) {
	if len(items) == 0 {
		fmt.Println("没有找到工作项")
		return
	}

	fmt.Printf("\n共 %d 个工作项:\n\n", len(items))

	// 按类别分组
	categoryOrder := []string{"Req", "Task", "Bug"}
	grouped := make(map[string][]api.Workitem)
	for _, item := range items {
		grouped[item.CategoryId] = append(grouped[item.CategoryId], item)
	}
	// 收集未知类别
	for _, item := range items {
		if _, known := categoryNames[item.CategoryId]; !known {
			grouped[item.CategoryId] = append(grouped[item.CategoryId], item)
		}
	}

	for _, cat := range categoryOrder {
		catItems, ok := grouped[cat]
		if !ok || len(catItems) == 0 {
			continue
		}
		fmt.Printf("【%s】%d 项\n", categoryDisplayName(cat), len(catItems))
		fmt.Printf("  %-12s %-36s %-10s %-6s %-8s %-16s %-16s\n",
			"编号", "标题", "状态", "优先级", "创建人", "所在项目", "父工作项")
		fmt.Printf("  %s\n", strings.Repeat("-", 110))
		for _, item := range catItems {
			subject := item.Subject
			if len([]rune(subject)) > 18 {
				subject = string([]rune(subject)[:18]) + "..."
			}
			priority := extractPriority(item.CustomFieldValues)
			parentId := item.ParentId
			if parentId == "EMPTY_VALUE" {
				parentId = ""
			} else if len(parentId) > 12 {
				parentId = parentId[:12] + "..."
			}
			fmt.Printf("  %-12s %-36s %-10s %-6s %-8s %-16s %-16s\n",
				item.SerialNumber,
				subject,
				item.Status.DisplayName,
				priority,
				item.Creator.Name,
				item.Space.Name,
				parentId,
			)
		}
		fmt.Println()
	}
}

// resolveUserId 自动获取并缓存用户 ID
func resolveUserId(cfg *config.Config) error {
	if cfg.UserId != "" {
		return nil
	}
	user, err := api.GetCurrentUser(cfg)
	if err != nil {
		return fmt.Errorf("获取当前用户信息失败: %w\n请运行 yx-code init 配置 token", err)
	}
	cfg.UserId = user.Id
	fmt.Printf("当前用户: %s (%s)\n", user.Name, user.Id)
	if err := cfg.Save(); err == nil {
		fmt.Println("已保存 user_id 到 ~/.yunxiao/config.yaml")
	}
	return nil
}
