package main

import (
	"flag"
	"fmt"
)

type Config struct {
	dir string
}

func loadAppCfg() Config {
	var cfg Config
	flag.StringVar(&cfg.dir, "dir", ".", "set directory to look for config file in")
	flag.Parse()
	return cfg
}

func main() {
	cfg := loadAppCfg()
	fmt.Println("config dir:", cfg.dir)
}
