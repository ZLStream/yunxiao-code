package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"yx-code/internal/api"
	"yx-code/internal/config"
	"yx-code/internal/git"
	"yx-code/internal/review"
)

// loadAndOverrideConfig 读取配置并用命令行参数覆盖
func loadAndOverrideConfig(cmd *cobra.Command) (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	if d, _ := cmd.Flags().GetString("domain"); d != "" {
		cfg.Domain = d
	}
	if o, _ := cmd.Flags().GetString("org"); o != "" {
		cfg.OrganizationId = o
	}
	if t, _ := cmd.Flags().GetString("token"); t != "" {
		cfg.Token = t
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "yx-code",
		Short: "云效 CLI 工具",
		Long:  "阿里云云效 CLI 工具，支持代码提交、合并请求、工作项管理和工作总结\n配置文件: ~/.yunxiao/config.yaml",
	}

	rootCmd.PersistentFlags().String("domain", "", "云效域名")
	rootCmd.PersistentFlags().String("org", "", "组织 ID")
	rootCmd.PersistentFlags().String("token", "", "个人访问令牌")

	// commit 子命令
	var commitMessage string
	commitCmd := &cobra.Command{
		Use:   "commit",
		Short: "提交代码",
		RunE: func(cmd *cobra.Command, args []string) error {
			if commitMessage == "" {
				return fmt.Errorf("请通过 -m 指定 commit message")
			}

			fmt.Println("--- Git Commit ---")
			if err := git.AddAndCommit(commitMessage); err != nil {
				return err
			}
			fmt.Println("✅ 提交成功!")
			return nil
		},
	}
	commitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "commit message（必填）")
	rootCmd.AddCommand(commitCmd)

	// push 子命令
	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "推送到远程仓库",
		RunE: func(cmd *cobra.Command, args []string) error {
			branch, err := git.GetCurrentBranch()
			if err != nil {
				return err
			}
			fmt.Printf("当前分支: %s\n", branch)

			fmt.Println("--- Git Push ---")
			if err := git.Push(); err != nil {
				return err
			}
			fmt.Println("✅ 推送成功!")
			return nil
		},
	}
	rootCmd.AddCommand(pushCmd)

	// mr 子命令
	var (
		mrTitle        string
		mrDescription  string
		mrTargetBranch string
	)
	mrCmd := &cobra.Command{
		Use:   "mr",
		Short: "创建合并请求",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mrTitle == "" {
				return fmt.Errorf("请通过 -m 指定 MR 标题")
			}

			cfg, err := loadAndOverrideConfig(cmd)
			if err != nil {
				return err
			}

			branch, err := git.GetCurrentBranch()
			if err != nil {
				return err
			}
			fmt.Printf("当前分支: %s\n", branch)

			repoName, err := git.GetRepoName()
			if err != nil {
				return err
			}
			fmt.Printf("仓库名称: %s\n", repoName)

			projectId, err := api.GetProjectId(cfg, repoName)
			if err != nil {
				return err
			}
			fmt.Printf("项目 ID: %d\n", projectId)

			// 设置目标分支，默认为 develop
			targetBranch := mrTargetBranch
			if targetBranch == "" {
				targetBranch = "develop"
			}

			fmt.Println("\n--- 创建合并请求 ---")
			fmt.Printf("标题: %s\n", mrTitle)
			fmt.Printf("源分支: %s → 目标分支: %s\n", branch, targetBranch)

			result, err := api.CreateMergeRequest(cfg, projectId, branch, targetBranch, mrTitle, mrDescription)
			if err != nil {
				return err
			}

			fmt.Printf("\n✅ 合并请求创建成功!\n")
			fmt.Printf("   ID: %d\n", result.LocalId)
			if result.DetailUrl != "" {
				fmt.Printf("   URL: %s\n", result.DetailUrl)
			}

			// 询问是否查看差异
			fmt.Print("\n是否查看代码差异? (y/N): ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))

			if input == "y" || input == "yes" {
				// 使用 API 响应中的分支信息，如果为空则使用默认值
				respTargetBranch := result.TargetBranch
				if respTargetBranch == "" {
					respTargetBranch = targetBranch
				}
				sourceBranch := result.SourceBranch
				if sourceBranch == "" {
					sourceBranch = branch // 使用当前分支
				}
				fmt.Printf("📋 比较分支: %s → %s\n", sourceBranch, respTargetBranch)
				if err := doDiff(cfg, projectId, respTargetBranch, sourceBranch); err != nil {
					fmt.Printf("\n查看差异失败: %v\n", err)
				}
			}

			return nil
		},
	}
	mrCmd.Flags().StringVarP(&mrTitle, "message", "m", "", "MR 标题（必填）")
	mrCmd.Flags().StringVarP(&mrDescription, "description", "d", "", "MR 描述")
	mrCmd.Flags().StringVarP(&mrTargetBranch, "target", "t", "develop", "目标分支（默认: develop）")
	rootCmd.AddCommand(mrCmd)

	// diff 子命令
	var (
		diffTargetBranch string
		diffSourceBranch string
	)
	diffCmd := &cobra.Command{
		Use:   "diff",
		Short: "查看代码差异",
		Long:  "查看指定分支间的代码差异统计",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadAndOverrideConfig(cmd)
			if err != nil {
				return err
			}

			repoName, err := git.GetRepoName()
			if err != nil {
				return err
			}
			fmt.Printf("仓库名称: %s\n", repoName)

			projectId, err := api.GetProjectId(cfg, repoName)
			if err != nil {
				return err
			}
			fmt.Printf("项目 ID: %d\n", projectId)

			// 获取当前分支作为默认源分支
			currentBranch, err := git.GetCurrentBranch()
			if err != nil {
				return err
			}

			targetBranch := diffTargetBranch
			if targetBranch == "" {
				targetBranch = "develop"
			}
			sourceBranch := diffSourceBranch
			if sourceBranch == "" {
				sourceBranch = currentBranch
			}

			fmt.Printf("📋 比较分支: %s → %s\n", sourceBranch, targetBranch)
			return doDiff(cfg, projectId, targetBranch, sourceBranch)
		},
	}
	diffCmd.Flags().StringVarP(&diffTargetBranch, "target", "t", "", "目标分支（默认: develop）")
	diffCmd.Flags().StringVarP(&diffSourceBranch, "source", "s", "", "源分支（默认: 当前分支）")
	rootCmd.AddCommand(diffCmd)

	// workitems 子命令
	rootCmd.AddCommand(newWorkitemsCmd())

	// summary 子命令
	rootCmd.AddCommand(newSummaryCmd())

	// init 子命令
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "初始化配置文件",
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.Init()
		},
	}
	rootCmd.AddCommand(initCmd)

	// clone 子命令
	var (
		cloneBranch string
		clonePath   string
	)
	cloneCmd := &cobra.Command{
		Use:   "clone <git-url>",
		Short: "克隆仓库",
		Long:  "根据 git 地址克隆代码。指定 -b 参数时，自动检测主分支（main/master）并从主分支创建新分支",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoURL := args[0]

			destPath := clonePath
			if destPath == "" {
				destPath = git.ExtractRepoNameFromURL(repoURL)
			}

			fmt.Printf("正在克隆仓库: %s → %s\n", repoURL, destPath)
			if err := git.Clone(repoURL, destPath); err != nil {
				return err
			}

			if cloneBranch != "" {
				mainBranch, err := git.DetectMainBranch(destPath)
				if err != nil {
					return err
				}
				fmt.Printf("检测到主分支: %s\n", mainBranch)

				if err := git.CheckoutNewBranch(destPath, mainBranch, cloneBranch); err != nil {
					return err
				}
				fmt.Printf("\n✅ 完成! 已在 %s 中创建并切换到分支: %s\n", destPath, cloneBranch)
			} else {
				fmt.Printf("\n✅ 克隆完成! 仓库位于: %s\n", destPath)
			}
			return nil
		},
	}
	cloneCmd.Flags().StringVarP(&cloneBranch, "branch", "b", "", "新分支名称（可选）")
	cloneCmd.Flags().StringVarP(&clonePath, "path", "p", "", "克隆目标路径（默认从 URL 提取仓库名）")
	rootCmd.AddCommand(cloneCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// doDiff 执行代码差异查看
func doDiff(cfg *config.Config, projectId int64, targetBranch, sourceBranch string) error {
	fmt.Println("\n--- 获取代码差异 ---")
	fmt.Printf("比较: %s...%s\n", targetBranch, sourceBranch)

	compare, err := api.GetCompare(cfg, projectId, targetBranch, sourceBranch)
	if err != nil {
		return fmt.Errorf("获取代码差异失败: %w", err)
	}

	// 转换 diff 格式
	diffs := make([]review.DiffFile, len(compare.Diffs))
	for i, d := range compare.Diffs {
		diffs[i] = review.DiffFile{
			Diff:        d.Diff,
			NewPath:     d.NewPath,
			OldPath:     d.OldPath,
			NewFile:     d.NewFile,
			DeletedFile: d.DeletedFile,
			AddLines:    d.AddLines,
			DelLines:    d.DelLines,
		}
	}

	reviewer := review.NewReviewer(diffs)
	files, addLines, delLines := reviewer.GetStats()

	fmt.Printf("变更统计: %d 个文件, +%d 行, -%d 行\n", files, addLines, delLines)

	if len(diffs) == 0 {
		fmt.Println("没有代码变更")
		return nil
	}

	fmt.Println("\n--- 代码差异 ---")
	result, err := reviewer.GetDiffResult()
	if err != nil {
		return fmt.Errorf("获取差异结果失败: %w", err)
	}

	fmt.Println(result)
	return nil
}