package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"golang.org/x/tools/go/packages"

	"github.com/yuhua2000/gocode-score/pkg/types"
)

// ComplexityAnalyzer: 检查函数过长和圈复杂度
type ComplexityAnalyzer struct {
	MaxFuncLines  int
	MaxComplexity int
}

func NewComplexityAnalyzer() Analyzer {
	return &ComplexityAnalyzer{MaxFuncLines: 50, MaxComplexity: 10}
}

func (c *ComplexityAnalyzer) Name() string { return "complexity" }

func (c *ComplexityAnalyzer) Analyze(pkg *packages.Package) ([]types.Issue, error) {
	var issues []types.Issue

	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Body == nil {
				continue
			}

			start := pkg.Fset.Position(fd.Pos()).Line
			end := pkg.Fset.Position(fd.End()).Line
			lines := end - start + 1
			if lines > c.MaxFuncLines {
				pos := pkg.Fset.Position(fd.Pos())
				issues = append(issues, types.Issue{
					Analyzer: c.Name(),
					File:     filepath.Clean(pos.Filename),
					Line:     pos.Line,
					Col:      pos.Column,
					Message:  fmt.Sprintf("function %s is too long (%d lines)", fd.Name.Name, lines),
					Severity: "medium",
				})
			}

			// 简单圈复杂度估算
			cc := 1
			ast.Inspect(fd.Body, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt, *ast.CaseClause:
					cc++
				case *ast.BinaryExpr:
					if x.Op == token.LAND || x.Op == token.LOR {
						cc++
					}
				}
				return true
			})

			if cc > c.MaxComplexity {
				pos := pkg.Fset.Position(fd.Pos())
				issues = append(issues, types.Issue{
					Analyzer: c.Name(),
					File:     filepath.Clean(pos.Filename),
					Line:     pos.Line,
					Col:      pos.Column,
					Message:  fmt.Sprintf("function %s has high cyclomatic complexity (%d)", fd.Name.Name, cc),
					Severity: "medium",
				})
			}
		}
	}

	return issues, nil
}
