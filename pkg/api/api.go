package api

import (
	"github.com/yuhua2000/gocode-score/internal/engine"
	"github.com/yuhua2000/gocode-score/pkg/types"
)

// Analyze runs code analysis on the specified target path and returns the report.
// configPath can be empty to use the default configuration.
func Analyze(target string, configPath string) (*types.Report, error) {
	runner, err := engine.NewRunner(configPath)
	if err != nil {
		return nil, err
	}
	return runner.Run(target)
}
