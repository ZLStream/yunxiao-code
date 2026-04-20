package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"yx-code/internal/config"
)

// UserInfo 表示当前用户信息（GetUserByToken 返回）
type UserInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	NickName string `json:"nickName"`
}

// GetCurrentUser 通过个人访问令牌获取当前用户信息
func GetCurrentUser(cfg *config.Config) (*UserInfo, error) {
	apiURL := fmt.Sprintf("https://%s/oapi/v1/platform/user", cfg.Domain)

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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("获取用户信息失败: HTTP %d, 响应: %s", resp.StatusCode, string(body))
	}

	var user UserInfo
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w\n响应: %s", err, string(body))
	}

	return &user, nil
}

// Project 表示云效项目信息
type Project struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	CustomCode string `json:"customCode"`
	Scope      string `json:"scope"`
}

// SearchProjectsRequest 表示搜索项目的请求体
type SearchProjectsRequest struct {
	Conditions      string `json:"conditions,omitempty"`
	ExtraConditions string `json:"extraConditions,omitempty"`
	OrderBy         string `json:"orderBy,omitempty"`
	Page            int    `json:"page"`
	PerPage         int    `json:"perPage"`
	Sort            string `json:"sort,omitempty"`
}

// WorkitemUser 表示工作项中的用户信息
type WorkitemUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// WorkitemStatus 表示工作项状态
type WorkitemStatus struct {
	DisplayName string `json:"displayName"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	NameEn      string `json:"nameEn"`
}

// WorkitemSpace 表示工作项所属项目
type WorkitemSpace struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// WorkitemType 表示工作项类型
type WorkitemType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// CustomFieldOption 表示自定义字段的单个选项值
type CustomFieldOption struct {
	DisplayValue string `json:"displayValue"`
	Identifier   string `json:"identifier"`
}

// CustomFieldValue 表示工作项的自定义字段值
type CustomFieldValue struct {
	FieldFormat string              `json:"fieldFormat"`
	FieldId     string              `json:"fieldId"`
	FieldName   string              `json:"fieldName"`
	Values      []CustomFieldOption `json:"values"`
}

// Workitem 表示单个工作项
type Workitem struct {
	Id                string             `json:"id"`
	Subject           string             `json:"subject"`
	SerialNumber      string             `json:"serialNumber"`
	CategoryId        string             `json:"categoryId"`
	Status            WorkitemStatus     `json:"status"`
	AssignedTo        WorkitemUser       `json:"assignedTo"`
	Creator           WorkitemUser       `json:"creator"`
	Participants      []WorkitemUser     `json:"participants"`
	Space             WorkitemSpace      `json:"space"`
	GmtCreate         int64              `json:"gmtCreate"`
	GmtModified       int64              `json:"gmtModified"`
	ParentId          string             `json:"parentId"`
	WorkitemType      WorkitemType       `json:"workitemType"`
	CustomFieldValues []CustomFieldValue `json:"customFieldValues"`
	UpdateStatusAt    int64              `json:"updateStatusAt"`
	Description       string             `json:"description"`
}

// SearchWorkitemsRequest 表示搜索工作项的请求体
type SearchWorkitemsRequest struct {
	Category   string `json:"category"`
	Conditions string `json:"conditions,omitempty"`
	SpaceId    string `json:"spaceId"`
	SpaceType  string `json:"spaceType"`
	OrderBy    string `json:"orderBy"`
	Sort       string `json:"sort"`
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
}

// SearchProjects 搜索项目列表（支持分页，自动获取全部）
// conditions 示例：我管理的项目 {"conditionGroups":[[{"className":"user","fieldIdentifier":"project.admin","format":"multiList","operator":"CONTAINS","value":["<userId>"]}]]}
func SearchProjects(cfg *config.Config, conditions string) ([]Project, error) {
	var allProjects []Project
	page := 1
	perPage := 200

	for {
		apiURL := fmt.Sprintf("https://%s/oapi/v1/projex/organizations/%s/projects:search",
			cfg.Domain, cfg.OrganizationId)

		reqBody := SearchProjectsRequest{
			Conditions: conditions,
			OrderBy:    "gmtCreate",
			Page:       page,
			PerPage:    perPage,
			Sort:       "desc",
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
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("读取响应失败: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("搜索项目失败: HTTP %d, 响应: %s", resp.StatusCode, string(body))
		}

		var projects []Project
		if err := json.Unmarshal(body, &projects); err != nil {
			return nil, fmt.Errorf("解析响应失败: %w\n响应: %s", err, string(body))
		}

		allProjects = append(allProjects, projects...)

		totalStr := resp.Header.Get("x-total")
		if totalStr != "" {
			total, _ := strconv.Atoi(totalStr)
			if page*perPage >= total {
				break
			}
		} else if len(projects) < perPage {
			break
		}
		page++
	}

	return allProjects, nil
}

// SearchWorkitems 搜索指定项目的工作项（支持分页，自动获取全部）
func SearchWorkitems(cfg *config.Config, spaceId, category, conditions string) ([]Workitem, error) {
	var allItems []Workitem
	page := 1
	perPage := 100

	for {
		apiURL := fmt.Sprintf("https://%s/oapi/v1/projex/organizations/%s/workitems:search",
			cfg.Domain, cfg.OrganizationId)

		reqBody := SearchWorkitemsRequest{
			Category:   category,
			Conditions: conditions,
			SpaceId:    spaceId,
			SpaceType:  "Project",
			OrderBy:    "gmtModified",
			Sort:       "desc",
			Page:       page,
			PerPage:    perPage,
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
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("读取响应失败: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("搜索工作项失败: HTTP %d, 响应: %s", resp.StatusCode, string(body))
		}

		var items []Workitem
		if err := json.Unmarshal(body, &items); err != nil {
			return nil, fmt.Errorf("解析响应失败: %w\n响应: %s", err, string(body))
		}

		allItems = append(allItems, items...)

		totalStr := resp.Header.Get("x-total")
		if totalStr != "" {
			total, _ := strconv.Atoi(totalStr)
			if page*perPage >= total {
				break
			}
		} else if len(items) < perPage {
			break
		}
		page++
	}

	return allItems, nil
}
