package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"yx-code/internal/config"
)

type CreateMRRequest struct {
	SourceBranch    string `json:"sourceBranch"`
	TargetBranch    string `json:"targetBranch"`
	SourceProjectId int64  `json:"sourceProjectId"`
	TargetProjectId int64  `json:"targetProjectId"`
	Title           string `json:"title"`
	Description     string `json:"description,omitempty"`
}

type CreateMRResponse struct {
	LocalId      int64  `json:"localId"`
	DetailUrl    string `json:"detailUrl"`
	WebUrl       string `json:"webUrl"`
	Title        string `json:"title"`
	SourceBranch string `json:"sourceBranch"`
	TargetBranch string `json:"targetBranch"`
	// status 字段在成功时是字符串，在错误时是布尔值，使用 interface{} 兼容
	Status           interface{} `json:"status"`
	Code             int         `json:"code"`
	ErrorCode        string      `json:"errorCode"`
	ErrorMessage     string      `json:"errorMessage"`
	ErrorDescription string      `json:"errorDescription"`
	TraceId          string      `json:"traceId"`
}

type RepoInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// DiffFile 表示单个文件的 diff 信息
type DiffFile struct {
	Diff        string `json:"diff"`
	NewPath     string `json:"newPath"`
	OldPath     string `json:"oldPath"`
	NewFile     bool   `json:"newFile"`
	DeletedFile bool   `json:"deletedFile"`
	AddLines    int    `json:"addLines"`
	DelLines    int    `json:"delLines"`
}

// CommitInfo 表示提交信息
type CommitInfo struct {
	Id         string `json:"id"`
	ShortId    string `json:"shortId"`
	Title      string `json:"title"`
	Message    string `json:"message"`
	AuthorName string `json:"authorName"`
}

// CompareResponse 表示代码比较的响应
type CompareResponse struct {
	Commits []CommitInfo `json:"commits"`
	Diffs   []DiffFile   `json:"diffs"`
}

func GetProjectId(cfg *config.Config, repoName string) (int64, error) {
	apiURL := fmt.Sprintf("https://%s/oapi/v1/codeup/organizations/%s/repositories?search=%s",
		cfg.Domain, cfg.OrganizationId, repoName)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return 0, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("x-yunxiao-token", cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %w", err)
	}

	var repos []RepoInfo
	if err := json.Unmarshal(body, &repos); err != nil {
		return 0, fmt.Errorf("解析响应失败: %w\n响应: %s", err, string(body))
	}

	for _, repo := range repos {
		if repo.Name == repoName {
			return repo.Id, nil
		}
	}

	return 0, fmt.Errorf("未找到仓库: %s", repoName)
}

func CreateMergeRequest(cfg *config.Config, repositoryId int64, sourceBranch string, targetBranch string, title string, description string) (*CreateMRResponse, error) {
	apiURL := fmt.Sprintf("https://%s/oapi/v1/codeup/organizations/%s/repositories/%d/changeRequests",
		cfg.Domain, cfg.OrganizationId, repositoryId)

	reqBody := CreateMRRequest{
		SourceBranch:    sourceBranch,
		TargetBranch:    targetBranch,
		SourceProjectId: repositoryId,
		TargetProjectId: repositoryId,
		Title:           title,
		Description:     description,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-yunxiao-token", cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result CreateMRResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w\n响应: %s", err, string(body))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := result.ErrorDescription
		if errMsg == "" {
			errMsg = result.ErrorMessage
		}
		if errMsg == "" {
			errMsg = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("创建合并请求失败: %s", errMsg)
	}

	return &result, nil
}

// GetCompare 获取两个分支之间的代码差异
func GetCompare(cfg *config.Config, repositoryId int64, fromBranch, toBranch string) (*CompareResponse, error) {
	apiURL := fmt.Sprintf("https://%s/oapi/v1/codeup/organizations/%s/repositories/%d/compares?from=%s&to=%s&sourceType=branch&targetType=branch",
		cfg.Domain, cfg.OrganizationId, repositoryId, url.QueryEscape(fromBranch), url.QueryEscape(toBranch))

	fmt.Printf("🔍 请求 URL: %s\n", apiURL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("x-yunxiao-token", cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	fmt.Printf("📡 响应状态: %s\n", resp.Status)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("获取代码差异失败: HTTP %d, 响应: %s", resp.StatusCode, string(body))
	}

	var result CompareResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w\n响应: %s", err, string(body))
	}

	return &result, nil
}
