package git

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

func GetCurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("获取当前分支失败: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// GetRepoName 从 remote URL 提取仓库名（最后一段路径）
func GetRepoName() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", fmt.Errorf("获取 remote URL 失败: %w", err)
	}
	remoteURL := strings.TrimSpace(string(out))

	var path string
	if strings.HasPrefix(remoteURL, "git@") {
		idx := strings.Index(remoteURL, ":")
		if idx == -1 {
			return "", fmt.Errorf("无法解析 remote URL: %s", remoteURL)
		}
		path = remoteURL[idx+1:]
	} else {
		u, err := url.Parse(remoteURL)
		if err != nil {
			return "", fmt.Errorf("无法解析 remote URL: %w", err)
		}
		path = strings.TrimPrefix(u.Path, "/")
	}

	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("无法从 URL 中提取仓库名")
	}
	return parts[len(parts)-1], nil
}

// GetOrganizationId 从 remote URL 提取组织 ID（路径的第一段）
func GetOrganizationId() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", fmt.Errorf("获取 remote URL 失败: %w", err)
	}
	remoteURL := strings.TrimSpace(string(out))

	return ExtractOrganizationIdFromURL(remoteURL)
}

// ExtractOrganizationIdFromURL 从 git URL 提取组织 ID
func ExtractOrganizationIdFromURL(remoteURL string) (string, error) {
	var path string
	if strings.HasPrefix(remoteURL, "git@") {
		idx := strings.Index(remoteURL, ":")
		if idx == -1 {
			return "", fmt.Errorf("无法解析 remote URL: %s", remoteURL)
		}
		path = remoteURL[idx+1:]
	} else {
		u, err := url.Parse(remoteURL)
		if err != nil {
			return "", fmt.Errorf("无法解析 remote URL: %w", err)
		}
		path = strings.TrimPrefix(u.Path, "/")
	}

	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("无法从 URL 中提取组织 ID，格式应为 <org_id>/<repo_name>")
	}
	return parts[0], nil
}

func AddAndCommit(message string) error {
	addCmd := exec.Command("git", "add", ".")
	addCmd.Stdout = nil
	addCmd.Stderr = nil
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("git add 失败: %w", err)
	}

	commitCmd := exec.Command("git", "commit", "-m", message)
	out, err := commitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit 失败: %s", string(out))
	}
	fmt.Println(strings.TrimSpace(string(out)))
	return nil
}

func Clone(repoURL, destPath string) error {
	args := []string{"clone", repoURL}
	if destPath != "" {
		args = append(args, destPath)
	}
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone 失败: %s", string(out))
	}
	fmt.Println(strings.TrimSpace(string(out)))
	return nil
}

func DetectMainBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "branch", "-r")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取远程分支失败: %s", string(out))
	}
	branches := string(out)
	if strings.Contains(branches, "origin/main") {
		return "origin/main", nil
	}
	if strings.Contains(branches, "origin/master") {
		return "origin/master", nil
	}
	return "", fmt.Errorf("未检测到 origin/main 或 origin/master 分支")
}

func CheckoutNewBranch(repoPath, baseBranch, newBranch string) error {
	cmd := exec.Command("git", "-C", repoPath, "checkout", "-b", newBranch, baseBranch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("创建分支失败: %s", string(out))
	}
	fmt.Println(strings.TrimSpace(string(out)))
	return nil
}

func ExtractRepoNameFromURL(repoURL string) string {
	var path string
	if strings.HasPrefix(repoURL, "git@") {
		idx := strings.Index(repoURL, ":")
		if idx == -1 {
			return repoURL
		}
		path = repoURL[idx+1:]
	} else {
		u, err := url.Parse(repoURL)
		if err != nil {
			return repoURL
		}
		path = strings.TrimPrefix(u.Path, "/")
	}
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return repoURL
	}
	return parts[len(parts)-1]
}

func Push() error {
	branch, err := GetCurrentBranch()
	if err != nil {
		return err
	}

	// 受保护分支检查
	protectedBranches := []string{"main", "master", "develop"}
	for _, protected := range protectedBranches {
		if branch == protected {
			return fmt.Errorf("禁止直接推送到受保护分支: %s\n请创建 feature 分支并通过 Merge Request 提交代码", branch)
		}
	}

	cmd := exec.Command("git", "push", "origin", branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git push 失败: %s", string(out))
	}
	fmt.Println(strings.TrimSpace(string(out)))
	return nil
}