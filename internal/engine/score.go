package engine

import (
	"github.com/yuhua2000/gocode-score/internal/config"
	"github.com/yuhua2000/gocode-score/pkg/types"
)

// CalculateScore 根据报告计算代码分数
//
// 打分规则：
// - score 范围为 0~100
// - 分数与文件总行数和 issue 数量挂钩：文件越大，相同数量的问题惩罚越低
// - 不同分析器可以通过 cfg.Weight 调整惩罚权重
//
// 参数：
// - r: *types.Report，包含分析出的 issue 列表以及总行数
// - cfg: *config.Config，包含各分析器的权重
//
// 返回值：
// - float64，最终得分
func CalculateScore(r *types.Report, cfg *config.Config) float64 {
	totalLines := float64(r.TotalLines) // 文件总行数
	if totalLines == 0 {
		return 100.0
	}

	var weightedIssues float64
	for _, it := range r.Issues {
		w := cfg.Weight[it.Analyzer]
		if w == 0 {
			w = 1.0
		}
		weightedIssues += w
	}

	// 按比例计算分数，issue越多，分数越低
	score := 100.0 * (1 - weightedIssues/totalLines)
	if score < 0 {
		score = 0
	}
	return score
}
