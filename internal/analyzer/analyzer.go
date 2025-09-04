package analyzer

import (
	"github.com/yuhua2000/gocode-score/pkg/types"

	"golang.org/x/tools/go/packages"
)

// Analyzer 接口：每个具体分析器实现它
type Analyzer interface {
	Name() string
	Analyze(pkg *packages.Package) ([]types.Issue, error)
}
