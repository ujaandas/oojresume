package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	cfgPath  string
	cfgFile  string
	tmplPath string
	outDir   string
}

func loadAppCfg() Config {
	var cfg Config

	flag.StringVar(&cfg.cfgPath, "dir", ".", "set path to look for config file in")
	flag.StringVar(&cfg.cfgFile, "config", "resume.json", "config file name")
	flag.StringVar(&cfg.tmplPath, "tmpl", "tmpl", "set path to look for template files in")
	flag.StringVar(&cfg.outDir, "out", "out", "output directory for rendered resumes")

	flag.Parse()
	return cfg
}

func main() {
	cfg := loadAppCfg()

	resumes, err := getResumeConfig(cfg.cfgPath, cfg.cfgFile)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	tmpl, err := parseLatexTemplates(cfg.tmplPath)
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	if len(resumes) == 0 {
		log.Fatalf("no resumes in config")
	}

	if err := os.MkdirAll(cfg.outDir, 0o755); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}

	for i, r := range resumes {
		if err := validateResume(tmpl, "main.tex.tmpl", r); err != nil {
			log.Fatalf("resume %d validation failed: %v", i, err)
		}

		rendered, err := renderResume(tmpl, "main.tex.tmpl", r)
		if err != nil {
			log.Fatalf("resume %d render failed: %v", i, err)
		}

		outPath := filepath.Join(cfg.outDir, "main.tex")
		if err := os.WriteFile(outPath, []byte(rendered), 0o644); err != nil {
			log.Fatalf("failed to write %s: %v", outPath, err)
		}

		log.Printf("rendered resume %d to %s", i, outPath)
	}
}
