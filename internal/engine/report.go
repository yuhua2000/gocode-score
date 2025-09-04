package engine

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yuhua2000/gocode-score/pkg/types"
)

// 输出工具函数封装到 report.go
func WriteJSONReport(r *types.Report, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r)
}

func WriteTextReport(r *types.Report, w io.Writer) error {
	fmt.Fprintf(w, "Target: %s\n", r.Target)
	fmt.Fprintf(w, "Score: %.2f\n", r.Score)
	fmt.Fprintln(w, "Issues:")
	for _, issue := range r.Issues {
		fmt.Fprintf(w, "- [%s] %s: %s\n", issue.Severity, issue.Analyzer, issue.Message)
	}
	fmt.Fprintln(w, "By Analyzer:")
	for analyzer, count := range r.ByAnalyzer {
		fmt.Fprintf(w, "- %s: %d\n", analyzer, count)
	}
	return nil
}

func WriteMarkdownReport(r *types.Report, w io.Writer) error {
	fmt.Fprintf(w, "# Report for %s\n\n", r.Target)
	fmt.Fprintf(w, "**Score:** %.2f\n\n", r.Score)

	fmt.Fprintln(w, "## Issues")
	for _, issue := range r.Issues {
		fmt.Fprintf(w, "- **%s** (%s): %s\n", issue.Analyzer, issue.Severity, issue.Message)
	}

	fmt.Fprintln(w, "\n## By Analyzer")
	for analyzer, count := range r.ByAnalyzer {
		fmt.Fprintf(w, "- **%s**: %d\n", analyzer, count)
	}
	return nil
}
