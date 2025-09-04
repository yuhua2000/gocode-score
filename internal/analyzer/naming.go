package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"

	"github.com/yuhua2000/gocode-score/pkg/types"
)

// NamingAnalyzer: 检查命名是否符合 Go 规范
type NamingAnalyzer struct{}

func NewNamingAnalyzer() Analyzer { return &NamingAnalyzer{} }

func (n *NamingAnalyzer) Name() string { return "naming" }

func (n *NamingAnalyzer) Analyze(pkg *packages.Package) ([]types.Issue, error) {
	var issues []types.Issue

	for _, f := range pkg.Syntax {
		// 检查包名
		if f.Name != nil {
			pkgName := f.Name.Name
			if pkgName != strings.ToLower(pkgName) {
				pos := pkg.Fset.Position(f.Name.Pos())
				issues = append(issues, types.Issue{
					Analyzer: n.Name(),
					File:     filepath.Clean(pos.Filename),
					Line:     pos.Line,
					Col:      pos.Column,
					Message:  fmt.Sprintf("package name %q should be all lowercase", pkgName),
					Severity: "low",
				})
			}
		}

		// 检查顶层声明
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						n.checkIdent(s.Name, pkg, &issues)
					case *ast.ValueSpec:
						for _, name := range s.Names {
							n.checkIdent(name, pkg, &issues)
						}
					}
				}
			case *ast.FuncDecl:
				n.checkIdent(d.Name, pkg, &issues)
			}
		}
	}

	return issues, nil
}

func (n *NamingAnalyzer) checkIdent(ident *ast.Ident, pkg *packages.Package, issues *[]types.Issue) {
	if ident == nil {
		return
	}
	name := ident.Name
	pos := pkg.Fset.Position(ident.Pos())

	// 忽略 _ 和 init
	if name == "_" || name == "init" {
		return
	}

	// 导出标识符：首字母必须大写
	if token.IsExported(name) {
		if !unicode.IsUpper(rune(name[0])) {
			*issues = append(*issues, types.Issue{
				Analyzer: n.Name(),
				File:     filepath.Clean(pos.Filename),
				Line:     pos.Line,
				Col:      pos.Column,
				Message:  fmt.Sprintf("exported identifier %q should start with uppercase letter", name),
				Severity: "low",
			})
		}
	} else {
		// 未导出标识符：首字母必须小写
		if !unicode.IsLower(rune(name[0])) {
			*issues = append(*issues, types.Issue{
				Analyzer: n.Name(),
				File:     filepath.Clean(pos.Filename),
				Line:     pos.Line,
				Col:      pos.Column,
				Message:  fmt.Sprintf("unexported identifier %q should start with lowercase letter", name),
				Severity: "low",
			})
		}
	}

	// 常量命名检查（假设 const 必须大写驼峰）
	if strings.ToUpper(name[:1]) == name[:1] && strings.Contains(name, "_") {
		*issues = append(*issues, types.Issue{
			Analyzer: n.Name(),
			File:     filepath.Clean(pos.Filename),
			Line:     pos.Line,
			Col:      pos.Column,
			Message:  fmt.Sprintf("constant %q should be in CamelCase, not snake_case", name),
			Severity: "low",
		})
	}
}
