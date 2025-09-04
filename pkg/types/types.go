package types

// Issue 表示分析器发现的问题
type Issue struct {
	Analyzer string `json:"analyzer"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Col      int    `json:"col"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// Report 是最终输出结果
type Report struct {
	Target     string         `json:"target"`
	Score      float64        `json:"score"`
	Issues     []Issue        `json:"issues"`
	ByAnalyzer map[string]int `json:"by_analyzer"`
	TotalLines int64          `json:"total_lines"`
}
