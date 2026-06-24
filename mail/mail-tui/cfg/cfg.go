package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host                  string `yaml:"host"`
	Port                  int    `yaml:"port"`
	MarkProcessedAfterRun bool   `yaml:"mark_processed_after_run"`
}

type LLMConfig struct {
	Endpoint    string  `yaml:"endpoint"`
	Model       string  `yaml:"model"`
	APIKey      string  `yaml:"api_key"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
	TopP        float64 `yaml:"top_p"`
}

type RelevanceConfig struct {
	Categories []string `yaml:"categories"`
}

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	LLM       LLMConfig       `yaml:"llm"`
	Relevance RelevanceConfig `yaml:"relevance"`
	CategoriesStr string
}

var (
	HOME, _ = os.UserHomeDir()
	DEFAULT_PATH = filepath.Join(HOME, ".config", "mtui", "config.yaml")
)

var Cfg *Config

func Init(path string) error {
	if path == "" {
		path = DEFAULT_PATH
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}
	if c.LLM.Endpoint == "" {
		c.LLM.Endpoint = "http://localhost:11434"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.LLM.MaxTokens == 0 {
		c.LLM.MaxTokens = 2048
	}
	c.CategoriesStr = strings.Join(c.Relevance.Categories, ", ")
	Cfg = &c
	return nil
}
