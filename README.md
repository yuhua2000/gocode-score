# gocode-score

`gocode-score` 是一个 **Go 代码质量评分工具**。  
它通过静态分析检测常见的可维护性和风格问题，并给出综合评分，帮助开发者提升代码质量。

## ✨ 特性

- 代码风格检查（命名规范、注释等）
- 复杂度分析
- 导出函数缺少文档注释检测
- 可扩展的分析器框架
- 多种报告输出格式（JSON / Text / Markdown）

## 🚀 安装

```bash
go install github.com/yuhua2000/gocode-score/cmd/gocode-score@latest
```

## 📦 使用方法

在项目根目录执行：

```bash
gocode-score ./...
```

输出示例：

```
Analyzing package: myproject/foo
[Style] exported function Foo is missing doc comment (foo.go:12)
[Complexity] function Bar is too complex (bar.go:34)

Final Score: 82.5 / 100
```

## ⚙️ 配置

默认配置文件位于 `config/config.go`，你可以修改权重和检测规则：

```go
cfg := config.DefaultConfig()
cfg.Weight["analyzer_name"] = 0.5
```

支持 YAML 配置文件。

## 📊 报告格式

生成不同格式的报告：

```bash
gocode-score -format json ./...   # JSON
gocode-score -format text ./...   # Text
gocode-score -format md ./...     # Markdown
```

## 🔌 扩展分析器

你可以通过实现 `Analyzer` 接口添加新的检测器：

```go
type Analyzer interface {
    Name() string
    Run(pkg *Package) ([]core.Issue, error)
}
```

并在 `runner.go` 中注册：

```go
analyzers := []analyzer.Analyzer{
    analyzer.NewStyleAnalyzer(),
    analyzer.NewComplexityAnalyzer(),
    analyzer.NewDocAnalyzer(),          // 新增
    analyzer.NewErrorHandlingAnalyzer(), // 新增
}
```

## 🤝 贡献

欢迎提交 Issue 或 PR，一起完善这个项目。

1. Fork 项目
2. 新建分支 (`git checkout -b feature-xxx`)
3. 提交修改 (`git commit -m 'add xxx'`)
4. 推送分支 (`git push origin feature-xxx`)
5. 发起 Pull Request

## 📜 许可证

[MIT License](LICENSE)
