package main

import (
	"flag"
	"log"
)

type Config struct {
	cfgPath  string
	tmplPath string
}

func loadAppCfg() Config {
	var cfg Config

	flag.StringVar(&cfg.cfgPath, "dir", ".", "set path to look for config file in")
	flag.StringVar(&cfg.tmplPath, "tmpl", "tmpl", "set path to look for template files in")

	flag.Parse()
	return cfg
}

func main() {
	cfg := loadAppCfg()

	tmpl, err := parseLatexTemplates(cfg.tmplPath)
	if err != nil {
		log.Fatalf("failed to parse templates in %s: %s", cfg.tmplPath, err)
	}

	log.Printf("loaded %d templates", len(tmpl.Templates()))
}
