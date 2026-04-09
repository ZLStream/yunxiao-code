package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"yx-code/internal/api"
)

func newSummaryCmd() *cobra.Command {
	var (
		period       string
		startDate    string
		endDate      string
		statusFilter string
		outputJSON   bool
	)

	cmd := &cobra.Command{
		Use:   "summary",
		Short: "生成工作总结",
		Long:  "生成指定时间范围内的工作项总结（结构化、AI 友好格式）",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadAndOverrideConfig(cmd)
			if err != nil {
				return err
			}

			if err := resolveUserId(cfg); err != nil {
				return err
			}

			// 解析时间范围
			start, end, err := resolveDateRange(period, startDate, endDate)
			if err != nil {
				return err
			}

			fmt.Printf("统计时间: %s ~ %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))

			// 获取项目列表
			projects, err := fetchUserProjects(cfg)
			if err != nil {
				return err
			}

			if len(projects) == 0 {
				fmt.Println("未找到相关项目")
				return nil
			}

			fmt.Printf("共 %d 个项目，正在获取工作项...\n", len(projects))

			// 构建时间范围条件（毫秒时间戳）
			dateRangeJSON := fmt.Sprintf(`"%d","%d"`, start.UnixMilli(), end.UnixMilli())
			conditions := buildWorkitemConditions(cfg.UserId, dateRangeJSON)

			// 收集所有工作项
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
				summary := buildSummaryJSON(allItems, start, end, cfg.UserId)
				data, _ := json.MarshalIndent(summary, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			printSummaryReport(allItems, start, end)
			return nil
		},
	}

	cmd.Flags().StringVarP(&period, "period", "p", "", "时间周期: week/month/quarter/year（与 --start/--end 二选一）")
	cmd.Flags().StringVar(&startDate, "start", "", "开始日期，格式: 2006-01-02")
	cmd.Flags().StringVar(&endDate, "end", "", "结束日期，格式: 2006-01-02")
	cmd.Flags().StringVarP(&statusFilter, "status", "s", "all", "状态过滤: all（默认）/ undone（排除已完成）/ 中文状态名（如：开发中）")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "以 JSON 格式输出（AI 友好）")

	return cmd
}

// resolveDateRange 解析时间范围，未指定时交互提示用户
func resolveDateRange(period, startDate, endDate string) (time.Time, time.Time, error) {
	now := time.Now()

	// 优先使用 --start/--end
	if startDate != "" || endDate != "" {
		var start, end time.Time
		var err error

		if startDate != "" {
			start, err = time.ParseInLocation("2006-01-02", startDate, time.Local)
			if err != nil {
				return time.Time{}, time.Time{}, fmt.Errorf("开始日期格式错误，请使用 2006-01-02 格式")
			}
		} else {
			start = now.AddDate(0, -1, 0)
		}

		if endDate != "" {
			end, err = time.ParseInLocation("2006-01-02", endDate, time.Local)
			if err != nil {
				return time.Time{}, time.Time{}, fmt.Errorf("结束日期格式错误，请使用 2006-01-02 格式")
			}
			// 结束日期包含当天，设置到当天末尾
			end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		} else {
			end = now
		}

		return start, end, nil
	}

	// 使用 --period
	if period != "" {
		return periodToDateRange(period, now)
	}

	// 未指定时间，交互提示
	return promptForDateRange(now)
}

// periodToDateRange 将周期字符串转换为时间范围
func periodToDateRange(period string, now time.Time) (time.Time, time.Time, error) {
	switch strings.ToLower(period) {
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := now.AddDate(0, 0, -(weekday - 1))
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local)
		return start, now, nil

	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		return start, now, nil

	case "quarter":
		month := now.Month()
		quarterStart := time.Month(((int(month)-1)/3)*3 + 1)
		start := time.Date(now.Year(), quarterStart, 1, 0, 0, 0, 0, time.Local)
		return start, now, nil

	case "year":
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		return start, now, nil

	default:
		return time.Time{}, time.Time{}, fmt.Errorf("无效的时间周期: %s (可选: week/month/quarter/year)", period)
	}
}

// promptForDateRange 交互式提示用户选择时间范围
func promptForDateRange(now time.Time) (time.Time, time.Time, error) {
	fmt.Println("\n请选择统计时间范围:")
	fmt.Println("  1. 本周")
	fmt.Println("  2. 本月")
	fmt.Println("  3. 本季度")
	fmt.Println("  4. 本年")
	fmt.Println("  5. 自定义")
	fmt.Print("\n请输入选项 (1-5): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		return periodToDateRange("week", now)
	case "2":
		return periodToDateRange("month", now)
	case "3":
		return periodToDateRange("quarter", now)
	case "4":
		return periodToDateRange("year", now)
	case "5":
		fmt.Print("开始日期 (2006-01-02): ")
		startStr, _ := reader.ReadString('\n')
		startStr = strings.TrimSpace(startStr)

		fmt.Print("结束日期 (2006-01-02，留空为今天): ")
		endStr, _ := reader.ReadString('\n')
		endStr = strings.TrimSpace(endStr)

		start, err := time.ParseInLocation("2006-01-02", startStr, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("开始日期格式错误")
		}

		var end time.Time
		if endStr == "" {
			end = now
		} else {
			end, err = time.ParseInLocation("2006-01-02", endStr, time.Local)
			if err != nil {
				return time.Time{}, time.Time{}, fmt.Errorf("结束日期格式错误")
			}
			end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
		return start, end, nil

	default:
		return time.Time{}, time.Time{}, fmt.Errorf("无效选项: %s", input)
	}
}

// SummaryReport 工作总结结构（JSON 输出）
type SummaryReport struct {
	Period      string            `json:"period"`
	StartDate   string            `json:"start_date"`
	EndDate     string            `json:"end_date"`
	UserId      string            `json:"user_id"`
	TotalCount  int               `json:"total_count"`
	ByStatus    map[string]int    `json:"by_status"`
	ByProject   map[string]int    `json:"by_project"`
	ByCategory  map[string]int    `json:"by_category"`
	Items       []WorkitemSummary `json:"items"`
}

// WorkitemSummary 工作项摘要
type WorkitemSummary struct {
	Id           string `json:"id"`
	SerialNumber string `json:"serial_number"`
	Category     string `json:"category"`
	Subject      string `json:"subject"`
	Status       string `json:"status"`
	Priority     string `json:"priority"`
	Project      string `json:"project"`
	AssignedTo   string `json:"assigned_to"`
	Creator      string `json:"creator"`
	ParentId     string `json:"parent_id,omitempty"`
	CreatedAt    string `json:"created_at"`
	ModifiedAt   string `json:"modified_at"`
}

// buildSummaryJSON 构建 JSON 格式的总结
func buildSummaryJSON(items []api.Workitem, start, end time.Time, userId string) SummaryReport {
	byStatus := make(map[string]int)
	byProject := make(map[string]int)
	byCategory := make(map[string]int)
	var summaryItems []WorkitemSummary

	for _, item := range items {
		byStatus[item.Status.DisplayName]++
		byProject[item.Space.Name]++
		catName := categoryDisplayName(item.CategoryId)
		byCategory[catName]++

		modifiedAt := ""
		if item.GmtModified > 0 {
			modifiedAt = time.UnixMilli(item.GmtModified).Format("2006-01-02 15:04:05")
		}

		parentId := item.ParentId
		if parentId == "EMPTY_VALUE" {
			parentId = ""
		}

		summaryItems = append(summaryItems, WorkitemSummary{
			Id:           item.Id,
			SerialNumber: item.SerialNumber,
			Category:     catName,
			Subject:      item.Subject,
			Status:       item.Status.DisplayName,
			Priority:     extractPriority(item.CustomFieldValues),
			Project:      item.Space.Name,
			AssignedTo:   item.AssignedTo.Name,
			Creator:      item.Creator.Name,
			ParentId:     parentId,
			CreatedAt:    time.UnixMilli(item.GmtCreate).Format("2006-01-02 15:04:05"),
			ModifiedAt:   modifiedAt,
		})
	}

	return SummaryReport{
		Period:     fmt.Sprintf("%s ~ %s", start.Format("2006-01-02"), end.Format("2006-01-02")),
		StartDate:  start.Format("2006-01-02"),
		EndDate:    end.Format("2006-01-02"),
		UserId:     userId,
		TotalCount: len(items),
		ByStatus:   byStatus,
		ByProject:  byProject,
		ByCategory: byCategory,
		Items:      summaryItems,
	}
}

// printSummaryReport 打印可读的工作总结报告（按类别分组）
func printSummaryReport(items []api.Workitem, start, end time.Time) {
	fmt.Printf("\n========== 工作总结 ==========\n")
	fmt.Printf("统计周期: %s ~ %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
	fmt.Printf("工作项总数: %d\n", len(items))

	if len(items) == 0 {
		fmt.Println("该时间段内没有工作项")
		return
	}

	// 按类别分组
	categoryOrder := []string{"Req", "Task", "Bug"}
	grouped := make(map[string][]api.Workitem)
	for _, item := range items {
		grouped[item.CategoryId] = append(grouped[item.CategoryId], item)
	}

	for _, cat := range categoryOrder {
		catItems, ok := grouped[cat]
		if !ok || len(catItems) == 0 {
			continue
		}
		fmt.Printf("\n--- %s（%d 项）---\n", categoryDisplayName(cat), len(catItems))
		fmt.Printf("  %-12s %-36s %-10s %-6s %-8s %-16s\n",
			"编号", "标题", "状态", "优先级", "创建人", "所在项目")
		fmt.Printf("  %s\n", strings.Repeat("-", 95))
		for _, item := range catItems {
			subject := item.Subject
			if len([]rune(subject)) > 18 {
				subject = string([]rune(subject)[:18]) + "..."
			}
			priority := extractPriority(item.CustomFieldValues)
			fmt.Printf("  %-12s %-36s %-10s %-6s %-8s %-16s\n",
				item.SerialNumber,
				subject,
				item.Status.DisplayName,
				priority,
				item.Creator.Name,
				item.Space.Name,
			)
			if item.ParentId != "" && item.ParentId != "EMPTY_VALUE" {
				fmt.Printf("    父工作项: %s\n", item.ParentId)
			}
		}
	}

	// 按项目统计
	byProject := make(map[string]int)
	for _, item := range items {
		byProject[item.Space.Name]++
	}
	fmt.Println("\n--- 按项目汇总 ---")
	for project, count := range byProject {
		fmt.Printf("  %-30s %d 项\n", project, count)
	}

	fmt.Println("\n==============================")
}
