package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"yx-code/internal/git"
)

const (
	configDir      = ".yunxiao"
	configFileName = "config.yaml"
)

type Config struct {
	Domain         string `yaml:"domain"`
	OrganizationId string `yaml:"organization_id"`
	Token          string `yaml:"token"`
	UserId         string `yaml:"user_id"`
}

func configFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, configDir, configFileName)
}

func Load() (*Config, error) {
	cfg := &Config{
		Domain: "openapi-rdc.aliyuncs.com",
	}

	if path := configFilePath(); path != "" {
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
	if v := os.Getenv("YUNXIAO_USER_ID"); v != "" {
		cfg.UserId = v
	}

	// 如果 organization_id 仍未配置，尝试从 git remote URL 提取
	if cfg.OrganizationId == "" {
		if orgId, err := git.GetOrganizationId(); err == nil {
			cfg.OrganizationId = orgId
		}
	}

	return cfg, nil
}

func (c *Config) Save() error {
	path := configFilePath()
	if path == "" {
		return fmt.Errorf("无法获取 home 目录")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	header := "# 云效 CLI 配置 (~/.yunxiao/config.yaml)\n"
	if err := os.WriteFile(path, append([]byte(header), data...), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	return nil
}

func (c *Config) Validate() error {
	if c.Domain == "" {
		return fmt.Errorf("domain 未配置")
	}
	if c.OrganizationId == "" {
		return fmt.Errorf("organization_id 未配置")
	}
	if c.Token == "" {
		return fmt.Errorf("token 未配置，请运行 yx-code init 或设置 YUNXIAO_TOKEN 环境变量")
	}
	return nil
}

func Init() error {
	path := configFilePath()
	if path == "" {
		return fmt.Errorf("无法获取 home 目录")
	}

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s 已存在", path)
	}

	cfg := &Config{
		Domain: "openapi-rdc.aliyuncs.com",
	}
	if err := cfg.Save(); err != nil {
		return err
	}
	fmt.Printf("已创建 %s，请填写配置信息\n", path)
	return nil
}
