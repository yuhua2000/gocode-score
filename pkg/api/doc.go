// Package api provides an external API for gocode-score.
//
// Using this package, other Go programs can analyze code in specified directories
// and obtain analysis reports without using the CLI.
//
// Example usage:
//
//	import (
//	    "log"
//	    "fmt"
//	    "github.com/yuhua2000/gocode-score/pkg/api"
//	)
//
//	// Run analysis on "./pkg" with optional config file path
//	report, err := api.Analyze("./pkg", "config.yaml")
//	if err != nil {
//	    log.Fatalf("analysis failed: %v", err)
//	}
//	fmt.Printf("Score: %.2f\n", report.Score)
package api
