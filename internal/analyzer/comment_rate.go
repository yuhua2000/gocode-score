package analyzer

import (
	"fmt"
	"path/filepath"

	"golang.org/x/tools/go/packages"

	"github.com/yuhua2000/gocode-score/pkg/types"
)

// 定义分析器
type CommentRateAnalyzer struct {
	Threshold float64 // 最低注释率，例如 0.2 表示 20%
}

func NewCommentRateAnalyzer(threshold float64) Analyzer {
	return &CommentRateAnalyzer{
		Threshold: threshold,
	}
}

func (c *CommentRateAnalyzer) Name() string { return "comment_rate" }

func (c *CommentRateAnalyzer) Analyze(pkg *packages.Package) ([]types.Issue, error) {
	var issues []types.Issue

	for _, f := range pkg.Syntax {
		totalLines := f.End() - f.Pos()
		commentLines := 0

		// 统计注释行
		for _, cg := range f.Comments {
			for _, comment := range cg.List {
				start := pkg.Fset.Position(comment.Pos()).Line
				end := pkg.Fset.Position(comment.End()).Line
				commentLines += end - start + 1
			}
		}

		// 注释率计算
		rate := float64(commentLines) / float64(totalLines)
		if rate < c.Threshold {
			pos := pkg.Fset.Position(f.Pos())
			issues = append(issues, types.Issue{
				Analyzer: c.Name(),
				File:     filepath.Clean(pos.Filename),
				Line:     1,
				Col:      1,
				Message:  fmt.Sprintf("file comment rate %.2f%% is below threshold %.2f%%", rate*100, c.Threshold*100),
				Severity: "low",
			})
		}
	}

	return issues, nil
}
