package config

import (
	"fmt"
	"os"
	"time"

	"go.yaml.in/yaml/v2"
)

// Config defines the configuration for analyzers
type Config struct {
	MaxFileSize      int           // Maximum file size to read (bytes)
	MaxAnalysisDepth int           // Analysis depth
	Timeout          time.Duration // Analysis timeout

	Enable map[string]bool    // Enable flag for each analyzer
	Weight map[string]float64 // Weight of each analyzer

	MinCommentRate float64 // Minimum required comment rate (0~1)
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		MaxFileSize:      10 * 1024 * 1024, // 10 MB
		MaxAnalysisDepth: 5,
		Timeout:          30 * time.Second,
		Enable: map[string]bool{
			"style":        true,
			"complexity":   true,
			"naming":       true,
			"comment_rate": true,
		},
		Weight: map[string]float64{
			"style":        1.0,
			"complexity":   1.0,
			"naming":       1.0,
			"comment_rate": 1.0,
		},

		MinCommentRate: 0.3, // 默认最低注释率 30%
	}
}

// LoadFromFile 从 YAML 文件加载配置
func LoadFromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	err = yaml.NewDecoder(f).Decode(cfg)
	return cfg, err
}
