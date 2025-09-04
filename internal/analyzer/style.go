package analyzer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/yuhua2000/gocode-score/pkg/types"
)

// StyleAnalyzer: 检查 gofmt（使用 go/format）和导出符号注释
type StyleAnalyzer struct{}

func NewStyleAnalyzer() Analyzer { return &StyleAnalyzer{} }

func (s *StyleAnalyzer) Name() string { return "style" }

func (s *StyleAnalyzer) Analyze(pkg *packages.Package) ([]types.Issue, error) {
	var issues []types.Issue

	// 1) 检查文件格式是否与 go/format 格式化结果一致
	for _, fpath := range pkg.GoFiles {
		src, err := os.ReadFile(fpath)
		if err != nil {
			return nil, fmt.Errorf("read file %s: %w", fpath, err)
		}

		fmted, err := format.Source(src)
		if err != nil {
			// 如果格式化失败，我们记录为 issue
			issues = append(issues, types.Issue{
				Analyzer: s.Name(),
				File:     fpath,
				Line:     1,
				Col:      1,
				Message:  "go/format failed: " + err.Error(),
				Severity: "medium",
			})
			continue
		}

		if !bytes.Equal(bytes.TrimSpace(src), bytes.TrimSpace(fmted)) {
			issues = append(issues, types.Issue{
				Analyzer: s.Name(),
				File:     fpath,
				Line:     1,
				Col:      1,
				Message:  "file is not gofmt-formatted",
				Severity: "low",
			})
		}
	}

	// 2) 导出函数是否缺少注释
	for _, f := range pkg.Syntax {
		base := ""
		if f.Name != nil {
			base = f.Name.Name
		}
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Name.IsExported() {
					// 需要注释
					if d.Doc == nil || strings.TrimSpace(d.Doc.Text()) == "" {
						pos := pkg.Fset.Position(d.Pos())
						issues = append(issues, types.Issue{
							Analyzer: s.Name(),
							File:     filepath.Clean(pos.Filename),
							Line:     pos.Line,
							Col:      pos.Column,
							Message:  fmt.Sprintf("exported function %s in file %s is missing doc comment", d.Name.Name, base),
							Severity: "low",
						})
					}
				}
			}
		}
	}

	return issues, nil
}
