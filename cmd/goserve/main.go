package main

import (
	"goserve/cli"
	"goserve/server"
)

const version = "1.0.0"

func main() {
	cfg := cli.ParseArgs(version)
	if cfg == nil {
		return // help or version flag used
	}

	server.StartServer(cfg.RootPath, cfg.Port, version)
}
