package main

import (
	"github.com/caiocp/go-api/configs"
)

func main() {
	cfg, _ := configs.LoadConfig(".")
	println(cfg.DBHost)
}
