package engine

import (
	"fmt"
	"log"

	"golang.org/x/tools/go/packages"

	"github.com/yuhua2000/gocode-score/internal/analyzer"
	"github.com/yuhua2000/gocode-score/internal/config"
	"github.com/yuhua2000/gocode-score/pkg/types"
)

// Runner 负责调度所有分析器
type Runner struct {
	analyzers []analyzer.Analyzer
	cfg       *config.Config
}

func NewRunner(configPath string) (*Runner, error) {
	var cfg *config.Config
	var err error

	if configPath != "" {
		cfg, err = config.LoadFromFile(configPath)
		if err != nil {
			// 加载失败时使用默认配置
			log.Printf("failed to load config from %s: %v, using default config", configPath, err)
			return nil, err
		}
	} else {
		cfg = config.DefaultConfig()
	}

	analyzerDefs := []analyzer.Analyzer{
		analyzer.NewStyleAnalyzer(),
		analyzer.NewComplexityAnalyzer(),
		analyzer.NewNamingAnalyzer(),
		analyzer.NewCommentRateAnalyzer(cfg.MinCommentRate),
	}

	var analyzers []analyzer.Analyzer
	for _, a := range analyzerDefs {
		if enabled, ok := cfg.Enable[a.Name()]; ok && enabled {
			analyzers = append(analyzers, a)
		}
	}

	return &Runner{
		analyzers: analyzers,
		cfg:       cfg,
	}, nil
}

// Run 对给定的包模式运行分析，返回 Report
func (r *Runner) Run(pattern string) (*types.Report, error) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps}
	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return nil, fmt.Errorf("packages.Load failed: %w", err)
	}

	report := &types.Report{Target: pattern, Issues: []types.Issue{}, ByAnalyzer: map[string]int{}}

	for _, p := range pkgs {
		// 跳过空 package
		if p == nil || len(p.GoFiles) == 0 {
			continue
		}

		for _, f := range p.Syntax {
			start := f.Pos() // 文件开头位置
			end := f.End()   // 文件结束位置
			startPos := p.Fset.Position(start)
			endPos := p.Fset.Position(end)
			report.TotalLines += int64(endPos.Line - startPos.Line + 1)
		}

		for _, a := range r.analyzers {
			issues, err := a.Analyze(p)
			if err != nil {
				// 记录但继续分析
				fmt.Printf("analyzer %s failed on package %s: %v\n", a.Name(), p.PkgPath, err)
				continue
			}
			for _, it := range issues {
				report.Issues = append(report.Issues, it)
				report.ByAnalyzer[it.Analyzer]++
			}
		}
	}

	// 计算分数
	report.Score = CalculateScore(report, r.cfg)

	return report, nil
}
