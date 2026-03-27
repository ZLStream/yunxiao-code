package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"yx-code/internal/git"
)

const configFileName = ".yunxiao.yaml"

type Config struct {
	Domain         string `yaml:"domain"`
	OrganizationId string `yaml:"organization_id"`
	Token          string `yaml:"token"`
}

func findConfigFile() string {
	// 从当前目录向上查找 .yunxiao.yaml
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		p := filepath.Join(dir, configFileName)
		if _, err := os.Stat(p); err == nil {
			return p
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func Load() (*Config, error) {
	cfg := &Config{
		Domain: "openapi-rdc.aliyuncs.com",
	}

	// 从配置文件读取
	if path := findConfigFile(); path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			_ = yaml.Unmarshal(data, cfg)
		}
	}

	// 环境变量覆盖
	if v := os.Getenv("YUNXIAO_DOMAIN"); v != "" {
		cfg.Domain = v
	}
	if v := os.Getenv("YUNXIAO_ORGANIZATION_ID"); v != "" {
		cfg.OrganizationId = v
	}
	if v := os.Getenv("YUNXIAO_TOKEN"); v != "" {
		cfg.Token = v
	}

	// 如果 organization_id 仍未配置，尝试从 git remote URL 提取
	if cfg.OrganizationId == "" {
		if orgId, err := git.GetOrganizationId(); err == nil {
			cfg.OrganizationId = orgId
		}
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Domain == "" {
		return fmt.Errorf("domain 未配置")
	}
	if c.OrganizationId == "" {
		return fmt.Errorf("organization_id 未配置")
	}
	if c.Token == "" {
		return fmt.Errorf("token 未配置，请设置 YUNXIAO_TOKEN 环境变量或在 .yunxiao.yaml 中配置")
	}
	return nil
}

func Init() error {
	content := `# 云效合并请求 CLI 配置
domain: openapi-rdc.aliyuncs.com
organization_id: ""
token: ""
`
	if _, err := os.Stat(configFileName); err == nil {
		return fmt.Errorf("%s 已存在", configFileName)
	}
	if err := os.WriteFile(configFileName, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	fmt.Printf("已创建 %s，请填写配置信息\n", configFileName)
	return nil
}