// gocode-score is a static code analysis and scoring tool for Go projects.
//
// # Overview
//
// gocode-score analyzes Go source code and produces a score (0–100)
// based on style, complexity, naming, and other rules.
// It is designed to help developers improve code quality.
//
// # Features
//
//   - Style checks (e.g., missing comments on exported functions)
//   - Complexity checks (e.g., cyclomatic complexity thresholds)
//   - Naming convention checks (e.g., package names lowercase, exported identifiers capitalized)
//   - Extensible analyzer framework
//   - Report generation in JSON, text, and Markdown
//
// # Usage
//
// From the command line:
//
//	gocode-score ./...
//
// Example with report output:
//
//	gocode-score -format=json ./...
//
// # Project Structure
//
//	internal/analyzer   → Built-in analyzers (style, complexity, naming, etc.)
//	internal/config     → Configuration handling
//	cmd/gocode-score    → CLI entry point
//
// # License
//
// gocode-score is released under the MIT License. See LICENSE file for details.
//
// # Contributing
//
// Contributions are welcome! Please open issues or pull requests on GitHub.
package main
