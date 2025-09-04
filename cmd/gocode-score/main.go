package main

import (
	"flag"
	"log"
	"os"

	"github.com/yuhua2000/gocode-score/internal/engine"
)

func main() {
	var target string
	var format string
	var out string
	var configPath string

	flag.StringVar(&target, "target", "./...", "packages pattern to analyze (eg: ./... or ./pkg)")
	flag.StringVar(&format, "format", "text", "output format: text|json|md")
	flag.StringVar(&out, "out", "", "output file (default stdout)")
	flag.StringVar(&configPath, "config", "", "path to configuration file (default uses built-in config)")
	flag.Parse()

	runner, err := engine.NewRunner(configPath)
	if err != nil {
		log.Fatalf("failed to create runner: %v", err)
	}

	report, err := runner.Run(target)
	if err != nil {
		log.Fatalf("analysis failed: %v", err)
	}

	w := os.Stdout
	if out != "" {
		f, err := os.Create(out)
		if err != nil {
			log.Fatalf("cannot create output file: %v", err)
		}
		defer f.Close()
		w = f
	}

	switch format {
	case "json":
		if err := engine.WriteJSONReport(report, w); err != nil {
			log.Fatalf("write json failed: %v", err)
		}
	case "md":
		if err := engine.WriteMarkdownReport(report, w); err != nil {
			log.Fatalf("write md failed: %v", err)
		}
	default:
		if err := engine.WriteTextReport(report, w); err != nil {
			log.Fatalf("write text failed: %v", err)
		}
	}
}
