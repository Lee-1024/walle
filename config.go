package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 对应 yaml 配置文件结构
type Config struct {
	Global struct {
		User           string `yaml:"user"`
		Password       string `yaml:"password"`
		PrivateKeyPath string `yaml:"private_key_path"`
		Port           int    `yaml:"port"`
	} `yaml:"global"`
	Hosts []HostInfo `yaml:"hosts"`
}

// HostInfo 单个主机连接信息
type HostInfo struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	PrivateKeyPath string `yaml:"private_key_path"`
	Alias          string `yaml:"alias"`
}

// LoadConfig 从指定路径加载并解析 YAML 配置文件
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	// 填充全局默认值
	if cfg.Global.Port == 0 {
		cfg.Global.Port = 22
	}
	if cfg.Global.User == "" {
		cfg.Global.User = "root"
	}

	for i := range cfg.Hosts {
		if cfg.Hosts[i].Port == 0 {
			cfg.Hosts[i].Port = cfg.Global.Port
		}
		if cfg.Hosts[i].User == "" {
			cfg.Hosts[i].User = cfg.Global.User
		}
		if cfg.Hosts[i].Password == "" && cfg.Hosts[i].PrivateKeyPath == "" {
			cfg.Hosts[i].Password = cfg.Global.Password
			cfg.Hosts[i].PrivateKeyPath = cfg.Global.PrivateKeyPath
		}
		if cfg.Hosts[i].Alias == "" {
			cfg.Hosts[i].Alias = cfg.Hosts[i].Host
		}
	}

	return &cfg, nil
}
